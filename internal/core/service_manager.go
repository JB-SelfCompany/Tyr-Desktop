package core

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/autoconfig"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/yggmail"
)

// ServiceManager coordinates the lifecycle of the Yggmail service
// Integrates configuration, status monitoring, and event forwarding
// All methods are thread-safe and can be called from multiple goroutines
type ServiceManager struct {
	// Core yggmail service
	yggmailService *yggmail.Service

	// Autoconfiguration server
	autoconfigServer *autoconfig.Server

	// Configuration reference
	config *Config

	// Status monitoring
	statusChan chan yggmail.ServiceStatus
	eventChans *yggmail.EventChannels

	// State management
	mu       sync.RWMutex
	running  bool
	shutdown bool

	// Shutdown coordination
	stopChan chan struct{}
	wg       sync.WaitGroup

	// Auto-restart configuration
	autoRestart      bool
	maxRestartDelay  time.Duration
	restartAttempts  int
	maxRestartCount  int
	restartResetTime time.Duration
	lastRestart      time.Time
}

// ServiceManagerOptions contains optional configuration for the service manager
type ServiceManagerOptions struct {
	// AutoRestart enables automatic service restart on errors
	AutoRestart bool

	// MaxRestartDelay is the maximum delay between restart attempts
	MaxRestartDelay time.Duration

	// MaxRestartCount is the maximum number of restart attempts before giving up
	MaxRestartCount int

	// RestartResetTime is the time after which restart counter is reset
	RestartResetTime time.Duration
}

// DefaultServiceManagerOptions returns default options for the service manager
func DefaultServiceManagerOptions() ServiceManagerOptions {
	return ServiceManagerOptions{
		AutoRestart:      false, // Disabled by default for user control
		MaxRestartDelay:  5 * time.Minute,
		MaxRestartCount:  5,
		RestartResetTime: 1 * time.Hour,
	}
}

// NewServiceManager creates a new service manager with the given configuration
// Does not start the service - call Initialize() and Start() explicitly
func NewServiceManager(config *Config) (*ServiceManager, error) {
	return NewServiceManagerWithOptions(config, DefaultServiceManagerOptions())
}

// NewServiceManagerWithOptions creates a new service manager with custom options
func NewServiceManagerWithOptions(config *Config, options ServiceManagerOptions) (*ServiceManager, error) {
	// Validate configuration
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Validate service settings
	if config.ServiceSettings.DatabasePath == "" {
		return nil, fmt.Errorf("database path cannot be empty")
	}
	if config.ServiceSettings.SMTPAddress == "" {
		return nil, fmt.Errorf("SMTP address cannot be empty")
	}
	if config.ServiceSettings.IMAPAddress == "" {
		return nil, fmt.Errorf("IMAP address cannot be empty")
	}

	sm := &ServiceManager{
		config:           config,
		statusChan:       make(chan yggmail.ServiceStatus, 10),
		stopChan:         make(chan struct{}),
		autoRestart:      options.AutoRestart,
		maxRestartDelay:  options.MaxRestartDelay,
		maxRestartCount:  options.MaxRestartCount,
		restartResetTime: options.RestartResetTime,
	}

	return sm, nil
}

// Initialize sets up the yggmail service with configuration
// Must be called before Start()
// Thread-safe with write lock
func (sm *ServiceManager) Initialize() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.running {
		return fmt.Errorf("service is already running, stop it first")
	}

	if sm.yggmailService != nil {
		// Service already initialized, close it first
		if err := sm.yggmailService.Close(); err != nil {
			return fmt.Errorf("failed to close existing service: %w", err)
		}
	}

	// Stop existing autoconfig server if running
	// This prevents "bind: address already in use" errors
	if sm.autoconfigServer != nil && sm.autoconfigServer.IsRunning() {
		if err := sm.autoconfigServer.Stop(); err != nil {
			log.Printf("Warning: failed to stop existing autoconfig server: %v", err)
		}
	}

	// Create new yggmail service
	service, err := yggmail.New(
		sm.config.ServiceSettings.DatabasePath,
		sm.config.ServiceSettings.SMTPAddress,
		sm.config.ServiceSettings.IMAPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create yggmail service: %w", err)
	}

	// Initialize the service
	if err := service.Initialize(); err != nil {
		service.Close()
		return fmt.Errorf("failed to initialize yggmail service: %w", err)
	}

	// Set password from keyring only on first initialization
	// Note: During onboarding, SetPassword is called BEFORE serviceManager exists,
	// so we need to set it here on first run. On subsequent restarts, the password
	// persists in yggmail database, so we skip this to avoid unnecessary updates.
	if !sm.config.ServiceSettings.PasswordInitialized {
		password, err := sm.config.GetPassword()
		if err != nil {
			service.Close()
			return fmt.Errorf("failed to retrieve password: %w", err)
		}

		if password != "" {
			if err := service.SetPassword(password); err != nil {
				service.Close()
				return fmt.Errorf("failed to set password: %w", err)
			}

			// Mark password as initialized and save config
			sm.config.ServiceSettings.PasswordInitialized = true
			if err := sm.config.Save(); err != nil {
				log.Printf("Warning: failed to save password_initialized flag: %v", err)
			}
		}
	}

	// Set maximum message size limit from configuration
	// This limits the size of individual messages that can be received
	maxSizeMB := sm.config.ServiceSettings.MaxMessageSizeMB
	if maxSizeMB == 0 {
		maxSizeMB = 10 // Fallback if not set in config
	}
	if err := service.SetMaxMessageSizeMB(maxSizeMB); err != nil {
		log.Printf("Warning: failed to set max message size: %v", err)
	} else {
		log.Printf("Maximum message size set to %d MB", maxSizeMB)
	}

	sm.yggmailService = service
	sm.eventChans = service.GetEventChannels()

	// Initialize and start autoconfiguration server
	if err := sm.startAutoconfigServer(); err != nil {
		service.Close()
		return fmt.Errorf("failed to start autoconfig server: %w", err)
	}

	// Start status monitoring goroutine
	sm.wg.Add(1)
	go sm.monitorStatus()

	// Note: Log monitoring is handled by App.startEventMonitoring() in events.go
	// to avoid duplicate consumption of log events from the same channel

	return nil
}

// Start starts the yggmail service with configured enabled peers
// Service must be initialized first
// Thread-safe with write lock
func (sm *ServiceManager) Start() error {
	sm.mu.Lock()

	if sm.running {
		sm.mu.Unlock()
		return fmt.Errorf("service is already running")
	}

	if sm.yggmailService == nil {
		sm.mu.Unlock()
		return fmt.Errorf("service not initialized, call Initialize() first")
	}

	// Get enabled peers from configuration
	peers := sm.config.GetEnabledPeers()

	// Allow starting with no peers - service can run locally without network peers
	// This is useful for testing or when user wants to configure peers after start
	if len(peers) == 0 {
		log.Println("Starting service with no peers (running in local-only mode)")
	} else {
		log.Printf("Starting service with %d enabled peer(s)", len(peers))
	}

	sm.mu.Unlock()

	// Start the service (unlocked to prevent deadlock with status updates)
	if err := sm.yggmailService.Start(peers); err != nil {
		return fmt.Errorf("failed to start yggmail service: %w", err)
	}

	sm.mu.Lock()
	sm.running = true
	sm.mu.Unlock()

	// Send initial status update (non-blocking to avoid panic on closed channel)
	select {
	case sm.statusChan <- yggmail.StatusRunning:
	default:
		// Channel closed or full, skip update
	}

	return nil
}

// Stop gracefully stops the yggmail service
// Also stops autoconfig server to release port 8080
// Thread-safe with write lock
func (sm *ServiceManager) Stop() error {
	sm.mu.Lock()

	if !sm.running {
		sm.mu.Unlock()
		return fmt.Errorf("service is not running")
	}

	if sm.yggmailService == nil {
		sm.mu.Unlock()
		return fmt.Errorf("service not initialized")
	}

	// Check if already shutting down
	if sm.shutdown {
		sm.mu.Unlock()
		return fmt.Errorf("service manager is shutting down")
	}

	sm.mu.Unlock()

	// Stop the service
	if err := sm.yggmailService.Stop(); err != nil {
		return fmt.Errorf("failed to stop yggmail service: %w", err)
	}

	// Stop autoconfig server to release port 8080
	// This allows clean restart without port conflicts
	sm.mu.Lock()
	if sm.autoconfigServer != nil && sm.autoconfigServer.IsRunning() {
		if err := sm.autoconfigServer.Stop(); err != nil {
			log.Printf("Warning: failed to stop autoconfig server: %v", err)
		}
	}
	sm.running = false
	sm.mu.Unlock()

	// Send status update (non-blocking to avoid panic on closed channel)
	select {
	case sm.statusChan <- yggmail.StatusStopped:
	default:
		// Channel closed or full, skip update
	}

	return nil
}

// SoftStop performs a graceful shutdown by first disconnecting all peers cleanly
func (sm *ServiceManager) SoftStop() error {
	sm.mu.Lock()

	if !sm.running {
		sm.mu.Unlock()
		return fmt.Errorf("service is not running")
	}

	if sm.yggmailService == nil {
		sm.mu.Unlock()
		return fmt.Errorf("service not initialized")
	}

	// Check if already shutting down
	if sm.shutdown {
		sm.mu.Unlock()
		return fmt.Errorf("service manager is shutting down")
	}

	sm.mu.Unlock()

	// Use SoftStop from yggmail service
	if err := sm.yggmailService.SoftStop(); err != nil {
		return fmt.Errorf("failed to soft stop yggmail service: %w", err)
	}

	// Stop autoconfig server to release port 8080
	sm.mu.Lock()
	if sm.autoconfigServer != nil && sm.autoconfigServer.IsRunning() {
		if err := sm.autoconfigServer.Stop(); err != nil {
			log.Printf("Warning: failed to stop autoconfig server: %v", err)
		}
	}
	sm.running = false
	sm.mu.Unlock()

	// Send status update (non-blocking to avoid panic on closed channel)
	select {
	case sm.statusChan <- yggmail.StatusStopped:
	default:
		// Channel closed or full, skip update
	}

	return nil
}

// Restart stops and then starts the service
// Useful for applying configuration changes
// Thread-safe
func (sm *ServiceManager) Restart() error {
	// Stop if running
	sm.mu.RLock()
	wasRunning := sm.running
	sm.mu.RUnlock()

	if wasRunning {
		if err := sm.Stop(); err != nil {
			return fmt.Errorf("failed to stop service during restart: %w", err)
		}

		// Brief delay to allow clean shutdown
		time.Sleep(500 * time.Millisecond)
	}

	// Reinitialize with current configuration
	// Initialize() will handle stopping old autoconfig server if needed
	if err := sm.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize service during restart: %w", err)
	}

	// Start service
	if err := sm.Start(); err != nil {
		return fmt.Errorf("failed to start service during restart: %w", err)
	}

	return nil
}

// CloseService closes the yggmail service and releases the database file
// Unlike Shutdown, this allows re-initialization afterwards
// Must be called when service is stopped
// Thread-safe with write lock
func (sm *ServiceManager) CloseService() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Stop autoconfig server if running
	if sm.autoconfigServer != nil && sm.autoconfigServer.IsRunning() {
		if err := sm.autoconfigServer.Stop(); err != nil {
			log.Printf("Warning: failed to stop autoconfig server: %v", err)
		}
		sm.autoconfigServer = nil
	}

	// Close yggmail service to release database
	if sm.yggmailService != nil {
		if err := sm.yggmailService.Close(); err != nil {
			return fmt.Errorf("failed to close service: %w", err)
		}
		sm.yggmailService = nil
	}

	return nil
}

// Shutdown performs a complete shutdown of the service manager
// Stops monitoring, closes service, and waits for all goroutines
// Uses SoftStop to prevent ErrClosed errors in logs
// Thread-safe and idempotent - safe to call multiple times
func (sm *ServiceManager) Shutdown() error {
	sm.mu.Lock()

	// Check if already shutdown
	if sm.shutdown {
		sm.mu.Unlock()
		log.Println("ServiceManager already shutdown, skipping")
		return nil
	}

	// Mark as shutdown
	sm.shutdown = true

	// Stop service if running (use SoftStop for clean peer disconnection)
	if sm.running && sm.yggmailService != nil {
		sm.mu.Unlock()
		if err := sm.yggmailService.SoftStop(); err != nil {
			log.Printf("Warning: failed to soft stop service during shutdown: %v", err)
			// Try normal stop as fallback
			if err := sm.yggmailService.Stop(); err != nil {
				return fmt.Errorf("failed to stop service during shutdown: %w", err)
			}
		}
		sm.mu.Lock()
		sm.running = false
	}

	sm.mu.Unlock()

	// Signal monitoring goroutine to stop
	close(sm.stopChan)

	// Wait for all goroutines to complete
	sm.wg.Wait()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Stop autoconfiguration server
	if sm.autoconfigServer != nil {
		if err := sm.autoconfigServer.Stop(); err != nil {
			return fmt.Errorf("failed to stop autoconfig server: %w", err)
		}
		sm.autoconfigServer = nil
	}

	// Close service
	if sm.yggmailService != nil {
		if err := sm.yggmailService.Close(); err != nil {
			return fmt.Errorf("failed to close service: %w", err)
		}
		sm.yggmailService = nil
	}

	// Close channels
	close(sm.statusChan)

	return nil
}

// GetStatus returns the current service status
// Thread-safe with read lock
func (sm *ServiceManager) GetStatus() yggmail.ServiceStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.yggmailService == nil {
		return yggmail.StatusStopped
	}

	return sm.yggmailService.GetStatus()
}

// IsRunning returns true if the service is currently running
// Thread-safe with read lock
func (sm *ServiceManager) IsRunning() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.running
}

// GetEventChannels returns the event channels for UI updates
// Returns nil if service is not initialized
// Thread-safe with read lock
func (sm *ServiceManager) GetEventChannels() *yggmail.EventChannels {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.eventChans
}

// GetMailAddress returns the email address for this node
// Returns empty string if service is not initialized
// Thread-safe with read lock
func (sm *ServiceManager) GetMailAddress() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.yggmailService == nil {
		return ""
	}

	return sm.yggmailService.GetMailAddress()
}

// GetPublicKey returns the hex-encoded public key
// Returns empty string if service is not initialized
// Thread-safe with read lock
func (sm *ServiceManager) GetPublicKey() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.yggmailService == nil {
		return ""
	}

	return sm.yggmailService.GetPublicKey()
}

// GetPeerStats returns statistics for all connected peers
// Returns empty slice if service is not running
// Thread-safe with read lock
func (sm *ServiceManager) GetPeerStats() []yggmail.PeerInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.yggmailService == nil {
		return []yggmail.PeerInfo{}
	}

	return sm.yggmailService.GetPeerStats()
}

// UpdatePassword updates the password and saves it to keyring
// Service must be initialized for the password to take effect
// Thread-safe
func (sm *ServiceManager) UpdatePassword(password string) error {
	// Save to keyring first
	if err := sm.config.SetPassword(password); err != nil {
		return fmt.Errorf("failed to save password to keyring: %w", err)
	}

	// Update service if initialized
	sm.mu.RLock()
	service := sm.yggmailService
	sm.mu.RUnlock()

	if service != nil {
		if err := service.SetPassword(password); err != nil {
			return fmt.Errorf("failed to update service password: %w", err)
		}

		// Mark password as initialized after successful update
		if !sm.config.ServiceSettings.PasswordInitialized {
			sm.config.ServiceSettings.PasswordInitialized = true
			if err := sm.config.Save(); err != nil {
				log.Printf("Warning: failed to save password_initialized flag: %v", err)
			}
		}
	}

	return nil
}

// UpdatePeers updates the peer list without restarting the service
// Only works if service is running
// Thread-safe
func (sm *ServiceManager) UpdatePeers(peers []string) error {
	sm.mu.RLock()
	service := sm.yggmailService
	running := sm.running
	sm.mu.RUnlock()

	if !running || service == nil {
		return fmt.Errorf("service must be running to update peers")
	}

	if err := service.UpdatePeers(peers); err != nil {
		return fmt.Errorf("failed to update peers: %w", err)
	}

	return nil
}

// HotReloadPeers updates the peer list without restarting the service
// This is an alias for UpdatePeers for API compatibility with Android version
// Uses Yggdrasil Core's AddPeer/RemovePeer for live updates without reconnection
// Thread-safe
func (sm *ServiceManager) HotReloadPeers(peers []string) error {
	log.Printf("Hot reloading peers with %d enabled peers", len(peers))
	return sm.UpdatePeers(peers)
}

// HotReloadMaxMessageSize updates the maximum message size without restarting the service
// Applies the new limit to the running yggmail service without disconnecting peers
// Thread-safe
func (sm *ServiceManager) HotReloadMaxMessageSize(sizeMB int64) error {
	sm.mu.RLock()
	service := sm.yggmailService
	running := sm.running
	sm.mu.RUnlock()

	if !running || service == nil {
		return fmt.Errorf("service must be running to update max message size")
	}

	log.Printf("Hot reloading max message size to %d MB", sizeMB)
	if err := service.SetMaxMessageSizeMB(sizeMB); err != nil {
		return fmt.Errorf("failed to set max message size: %w", err)
	}

	log.Printf("Max message size updated successfully to %d MB", sizeMB)
	return nil
}

// GetStatusChannel returns a channel that receives status updates
// Buffered channel with capacity of 10
func (sm *ServiceManager) GetStatusChannel() <-chan yggmail.ServiceStatus {
	return sm.statusChan
}

// SetAutoRestart enables or disables automatic service restart on errors
// Thread-safe with write lock
func (sm *ServiceManager) SetAutoRestart(enable bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.autoRestart = enable
}

// monitorStatus monitors service status and handles automatic restarts
// Runs in a background goroutine
func (sm *ServiceManager) monitorStatus() {
	defer sm.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sm.stopChan:
			return

		case <-ticker.C:
			sm.mu.RLock()
			service := sm.yggmailService
			running := sm.running
			autoRestart := sm.autoRestart
			sm.mu.RUnlock()

			if service == nil {
				continue
			}

			status := service.GetStatus()

			// Check for error state
			if status == yggmail.StatusError && running && autoRestart {
				sm.handleServiceError()
			}

			// Send status update
			select {
			case sm.statusChan <- status:
			default:
				// Channel full, skip update
			}
		}
	}
}

// handleServiceError attempts to restart the service after an error
// Implements exponential backoff with maximum retry count
func (sm *ServiceManager) handleServiceError() {
	sm.mu.Lock()

	// Check if we should reset restart counter
	if time.Since(sm.lastRestart) > sm.restartResetTime {
		sm.restartAttempts = 0
	}

	// Check if we've exceeded max restart attempts
	if sm.restartAttempts >= sm.maxRestartCount {
		sm.mu.Unlock()
		// Give up and let user manually restart
		return
	}

	sm.restartAttempts++
	sm.lastRestart = time.Now()

	// Calculate backoff delay (exponential with maximum)
	delay := time.Duration(sm.restartAttempts) * time.Second
	if delay > sm.maxRestartDelay {
		delay = sm.maxRestartDelay
	}

	sm.mu.Unlock()

	// Wait before restart attempt
	time.Sleep(delay)

	// Attempt restart
	if err := sm.Restart(); err != nil {
		// Restart failed, will try again on next monitoring cycle
		select {
		case sm.statusChan <- yggmail.StatusError:
		default:
		}
	}
}

// startAutoconfigServer initializes and starts the autoconfiguration HTTP server
func (sm *ServiceManager) startAutoconfigServer() error {
	// Parse SMTP and IMAP addresses
	smtpHost, smtpPort := autoconfig.ParseSMTPAddress(sm.config.ServiceSettings.SMTPAddress)
	imapHost, imapPort := autoconfig.ParseIMAPAddress(sm.config.ServiceSettings.IMAPAddress)

	// Create autoconfig server
	server, err := autoconfig.NewServer(autoconfig.ServerConfig{
		MailDomain:  "yggmail",
		SMTPHost:    smtpHost,
		SMTPPort:    smtpPort,
		IMAPHost:    imapHost,
		IMAPPort:    imapPort,
		ListenAddr:  "127.0.0.1:8080",
		DisplayName: "Yggmail",
		ShortName:   "Yggmail",
	})
	if err != nil {
		return fmt.Errorf("failed to create autoconfig server: %w", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start autoconfig server: %w", err)
	}

	sm.autoconfigServer = server
	return nil
}

// GetAutoconfigURL returns the URL for the autoconfiguration server
// Returns empty string if autoconfig server is not running
func (sm *ServiceManager) GetAutoconfigURL() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.autoconfigServer == nil || !sm.autoconfigServer.IsRunning() {
		return ""
	}

	return "http://" + sm.autoconfigServer.GetListenAddr()
}

// IsAutoconfigRunning returns true if the autoconfig server is running
func (sm *ServiceManager) IsAutoconfigRunning() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.autoconfigServer != nil && sm.autoconfigServer.IsRunning()
}

// GetMaxMessageSizeMB returns the current maximum message size in megabytes
// Thread-safe with read lock
func (sm *ServiceManager) GetMaxMessageSizeMB() (int64, error) {
	sm.mu.RLock()
	service := sm.yggmailService
	sm.mu.RUnlock()

	if service == nil {
		return 0, fmt.Errorf("service not initialized")
	}

	return service.GetMaxMessageSizeMB()
}

// CheckRecipientMessageSizeLimit checks if recipient can accept a message of given size
// This should be called BEFORE sending in 1-on-1 chats to avoid wasting bandwidth
// For group chats, skip this check - send to all, those with capacity will accept
//
// Parameters:
//   - recipientEmail: Full email address (e.g., "abc123...@yggmail")
//   - messageSizeBytes: Size of message to send in bytes
//
// Returns MessageSizeLimitCheckResult with CanSend=true if message size is acceptable, false otherwise
// Thread-safe with read lock
func (sm *ServiceManager) CheckRecipientMessageSizeLimit(recipientEmail string, messageSizeBytes int64) (*yggmail.MessageSizeLimitCheckResult, error) {
	sm.mu.RLock()
	service := sm.yggmailService
	running := sm.running
	sm.mu.RUnlock()

	if service == nil {
		return nil, fmt.Errorf("service not initialized")
	}

	if !running {
		return nil, fmt.Errorf("service must be running to check recipient message size limit")
	}

	return service.CheckRecipientMessageSizeLimit(recipientEmail, messageSizeBytes)
}

// Note: monitorLogs was removed to prevent duplicate consumption of log events
// Log events are now exclusively handled by App.startEventMonitoring() in events.go
// This ensures all logs are properly forwarded to the frontend via Wails runtime events
