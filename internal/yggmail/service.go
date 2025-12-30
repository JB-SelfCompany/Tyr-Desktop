package yggmail

import (
	"fmt"
	"sync"
	"time"

	"github.com/JB-SelfCompany/yggmail/mobile"
)

// Service wraps the yggmail library with lifecycle management and event handling
// All methods are thread-safe and can be called from multiple goroutines
type Service struct {
	// Core yggmail service
	yggmailService *mobile.YggmailService

	// Service configuration
	dbPath   string
	smtpAddr string
	imapAddr string

	// State management
	mu     sync.RWMutex
	status ServiceStatus

	// Event channels for UI communication
	events *EventChannels

	// Shutdown coordination
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// New creates a new Service instance
// dbPath: absolute path to SQLite database file (will be created if not exists)
// smtpAddr: SMTP server listen address (e.g., "127.0.0.1:1025")
// imapAddr: IMAP server listen address (e.g., "127.0.0.1:1143")
func New(dbPath, smtpAddr, imapAddr string) (*Service, error) {
	// Validate inputs
	if dbPath == "" {
		return nil, fmt.Errorf("database path cannot be empty")
	}
	if smtpAddr == "" {
		smtpAddr = "127.0.0.1:1025"
	}
	if imapAddr == "" {
		imapAddr = "127.0.0.1:1143"
	}

	// Create yggmail service instance
	yggService, err := mobile.NewYggmailService(dbPath, smtpAddr, imapAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create yggmail service: %w", err)
	}

	service := &Service{
		yggmailService: yggService,
		dbPath:         dbPath,
		smtpAddr:       smtpAddr,
		imapAddr:       imapAddr,
		status:         StatusStopped,
		events:         NewEventChannels(),
		stopChan:       make(chan struct{}),
	}

	// Set up callbacks to bridge yggmail events to our channels
	service.setupCallbacks()

	return service, nil
}

// Initialize initializes the service and creates/loads keys
// This must be called before Start()
func (s *Service) Initialize() error {
	s.mu.Lock()
	if s.status != StatusStopped {
		s.mu.Unlock()
		return fmt.Errorf("service must be stopped before initialization")
	}
	s.status = StatusStarting
	s.mu.Unlock()

	// Initialize yggmail service
	if err := s.yggmailService.Initialize(); err != nil {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		return fmt.Errorf("failed to initialize yggmail: %w", err)
	}

	s.mu.Lock()
	s.status = StatusStopped
	s.mu.Unlock()

	return nil
}

// Start starts the Yggmail service with network connectivity
// peers: list of static peer URIs (e.g., ["tls://example.com:12345"])
func (s *Service) Start(peers []string) error {
	s.mu.Lock()
	if s.status == StatusRunning || s.status == StatusStarting {
		s.mu.Unlock()
		return fmt.Errorf("service already running or starting")
	}
	s.status = StatusStarting
	s.mu.Unlock()

	// Validate peer list
	if len(peers) == 0 {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		return fmt.Errorf("must provide at least one peer")
	}

	// Convert peer slice to comma-separated string for yggmail API
	peerStr := ""
	for i, peer := range peers {
		if i > 0 {
			peerStr += ","
		}
		peerStr += peer
	}

	// Start yggmail service
	if err := s.yggmailService.Start(peerStr); err != nil {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		return fmt.Errorf("failed to start yggmail: %w", err)
	}

	s.mu.Lock()
	s.status = StatusRunning
	s.mu.Unlock()

	return nil
}

// Stop gracefully stops the Yggmail service
// This will close all connections and wait for background tasks to complete
func (s *Service) Stop() error {
	s.mu.Lock()
	if s.status == StatusStopped || s.status == StatusStopping {
		s.mu.Unlock()
		return fmt.Errorf("service already stopped or stopping")
	}
	s.status = StatusStopping
	s.mu.Unlock()

	// Stop yggmail service
	if err := s.yggmailService.Stop(); err != nil {
		s.mu.Lock()
		s.status = StatusError
		s.mu.Unlock()
		return fmt.Errorf("failed to stop yggmail: %w", err)
	}

	s.mu.Lock()
	s.status = StatusStopped
	s.mu.Unlock()

	return nil
}

// SoftStop performs a graceful shutdown by first disconnecting all peers cleanly
// This method prevents "ErrClosed" errors in logs by:
// 1. Disconnecting all peers using UpdatePeers with empty list
// 2. Waiting for graceful disconnection to complete
// 3. Then performing normal service shutdown
// Recommended to use instead of Stop() for clean shutdown
func (s *Service) SoftStop() error {
	s.mu.RLock()
	status := s.status
	s.mu.RUnlock()

	if status == StatusStopped || status == StatusStopping {
		return fmt.Errorf("service already stopped or stopping")
	}

	// First, gracefully disconnect all peers by updating to empty peer list
	// This uses Yggdrasil Core's RemovePeer for clean disconnection
	if status == StatusRunning {
		if err := s.yggmailService.UpdatePeers(""); err != nil {
			// Log error but continue with stop anyway
			fmt.Printf("Warning: failed to disconnect peers gracefully: %v\n", err)
		}

		// Give time for graceful disconnection to complete
		time.Sleep(500 * time.Millisecond)
	}

	// Now perform normal stop
	return s.Stop()
}

// Close releases all resources
// Service must be stopped before calling Close
func (s *Service) Close() error {
	s.mu.Lock()
	if s.status == StatusRunning || s.status == StatusStarting {
		s.mu.Unlock()
		return fmt.Errorf("service must be stopped before closing")
	}
	s.mu.Unlock()

	// Close yggmail service
	if err := s.yggmailService.Close(); err != nil {
		return fmt.Errorf("failed to close yggmail: %w", err)
	}

	// Close event channels
	s.events.Close()

	return nil
}

// GetStatus returns the current service status
func (s *Service) GetStatus() ServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// GetMailAddress returns the email address for this node
// Returns empty string if service is not initialized
func (s *Service) GetMailAddress() string {
	return s.yggmailService.GetMailAddress()
}

// GetPublicKey returns the hex-encoded public key
// Returns empty string if service is not initialized
func (s *Service) GetPublicKey() string {
	return s.yggmailService.GetPublicKey()
}

// SetPassword sets the IMAP/SMTP authentication password
// The password is hashed using bcrypt before storage
func (s *Service) SetPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if err := s.yggmailService.SetPassword(password); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

// GetPeerStats returns statistics for all connected peers
func (s *Service) GetPeerStats() []PeerInfo {
	s.mu.RLock()
	status := s.status
	s.mu.RUnlock()

	// Return empty slice if not running
	if status != StatusRunning {
		return []PeerInfo{}
	}

	// Get peer connections from yggmail
	peerConns := s.yggmailService.GetPeerConnections()
	if peerConns == nil {
		return []PeerInfo{}
	}

	// Convert to our PeerInfo type
	peers := make([]PeerInfo, 0, len(peerConns))
	for _, pc := range peerConns {
		peers = append(peers, PeerInfo{
			Address:   pc.URI,
			Latency:   pc.LatencyMs,
			Status:    pc.Up,
			Inbound:   pc.Inbound,
			Key:       pc.Key,
			Uptime:    pc.Uptime,
			RXBytes:   pc.RXBytes,
			TXBytes:   pc.TXBytes,
			RXRate:    pc.RXRate,
			TXRate:    pc.TXRate,
			LastError: pc.LastError,
		})
	}

	return peers
}

// GetSMTPAddress returns the local SMTP server address
func (s *Service) GetSMTPAddress() string {
	return s.smtpAddr
}

// GetIMAPAddress returns the local IMAP server address
func (s *Service) GetIMAPAddress() string {
	return s.imapAddr
}

// GetEventChannels returns the event channels for receiving service events
// The returned channels should be monitored in a separate goroutine
func (s *Service) GetEventChannels() *EventChannels {
	return s.events
}

// UpdatePeers updates the peer configuration without restarting
// This allows adding/removing peers dynamically
func (s *Service) UpdatePeers(peers []string) error {
	s.mu.RLock()
	status := s.status
	s.mu.RUnlock()

	if status != StatusRunning {
		return fmt.Errorf("service must be running to update peers")
	}

	// Convert peer slice to comma-separated string
	peerStr := ""
	for i, peer := range peers {
		if i > 0 {
			peerStr += ","
		}
		peerStr += peer
	}

	if err := s.yggmailService.UpdatePeers(peerStr); err != nil {
		return fmt.Errorf("failed to update peers: %w", err)
	}

	return nil
}

// SetMaxMessageSizeMB sets the maximum message size in megabytes
// This limits the size of individual messages that can be received
func (s *Service) SetMaxMessageSizeMB(megabytes int64) error {
	if megabytes < 0 {
		return fmt.Errorf("max message size cannot be negative")
	}

	if err := s.yggmailService.SetMaxMessageSizeMB(megabytes); err != nil {
		return fmt.Errorf("failed to set max message size: %w", err)
	}

	return nil
}

// GetMaxMessageSizeMB returns the current maximum message size in megabytes
func (s *Service) GetMaxMessageSizeMB() (int64, error) {
	sizeMB, err := s.yggmailService.GetMaxMessageSizeMB()
	if err != nil {
		return 0, fmt.Errorf("failed to get max message size: %w", err)
	}
	return sizeMB, nil
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
func (s *Service) CheckRecipientMessageSizeLimit(recipientEmail string, messageSizeBytes int64) (*MessageSizeLimitCheckResult, error) {
	s.mu.RLock()
	status := s.status
	s.mu.RUnlock()

	if status != StatusRunning {
		return nil, fmt.Errorf("service must be running to check recipient message size limit")
	}

	// Call yggmail library's CheckRecipientMessageSizeLimit
	result, err := s.yggmailService.CheckRecipientMessageSizeLimit(recipientEmail, messageSizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to check recipient message size limit: %w", err)
	}

	// Convert from mobile library result to our type
	return &MessageSizeLimitCheckResult{
		CanSend:       result.CanSend,
		ErrorMessage:  result.ErrorMessage,
		RecipientAddr: result.RecipientAddr,
		MessageSizeMB: result.MessageSizeMB,
	}, nil
}

// setupCallbacks configures the yggmail callbacks to forward events to our channels
func (s *Service) setupCallbacks() {
	// Set log callback
	s.yggmailService.SetLogCallback(&logCallback{service: s})

	// Set mail callback
	s.yggmailService.SetMailCallback(&mailCallback{service: s})

	// Set connection callback
	s.yggmailService.SetConnectionCallback(&connectionCallback{service: s})
}

// logCallback implements mobile.LogCallback interface
type logCallback struct {
	service *Service
}

// OnLog is called when yggmail generates a log message
func (lc *logCallback) OnLog(level, tag, message string) {
	// Check if channels are closed before sending
	if lc.service.events.IsClosed() {
		return
	}

	// Non-blocking send to log channel
	select {
	case lc.service.events.Log <- LogEvent{
		Timestamp: time.Now(),
		Level:     level,
		Tag:       tag,
		Message:   message,
	}:
	default:
		// Channel full, drop event to prevent blocking
	}
}

// mailCallback implements mobile.MailCallback interface
type mailCallback struct {
	service *Service
}

// OnNewMail is called when new mail is received
func (mc *mailCallback) OnNewMail(mailbox, from, subject string, mailID int) {
	// Check if channels are closed before sending
	if mc.service.events.IsClosed() {
		return
	}

	select {
	case mc.service.events.Mail <- MailEvent{
		Timestamp: time.Now(),
		Type:      "new_mail",
		Mailbox:   mailbox,
		From:      from,
		Subject:   subject,
		MailID:    mailID,
	}:
	default:
		// Channel full, drop event
	}
}

// OnMailSent is called when mail is successfully sent
func (mc *mailCallback) OnMailSent(to, subject string) {
	// Check if channels are closed before sending
	if mc.service.events.IsClosed() {
		return
	}

	select {
	case mc.service.events.Mail <- MailEvent{
		Timestamp: time.Now(),
		Type:      "sent",
		To:        to,
		Subject:   subject,
	}:
	default:
		// Channel full, drop event
	}
}

// OnMailError is called when mail sending fails
func (mc *mailCallback) OnMailError(to, subject, errorMsg string) {
	// Check if channels are closed before sending
	if mc.service.events.IsClosed() {
		return
	}

	select {
	case mc.service.events.Mail <- MailEvent{
		Timestamp:    time.Now(),
		Type:         "error",
		To:           to,
		Subject:      subject,
		ErrorMessage: errorMsg,
	}:
	default:
		// Channel full, drop event
	}
}

// connectionCallback implements mobile.ConnectionCallback interface
type connectionCallback struct {
	service *Service
}

// OnConnected is called when a peer connection is established
func (cc *connectionCallback) OnConnected(peer string) {
	// Check if channels are closed before sending
	if cc.service.events.IsClosed() {
		return
	}

	select {
	case cc.service.events.Connection <- ConnectionEvent{
		Timestamp: time.Now(),
		Type:      "connected",
		Peer:      peer,
	}:
	default:
		// Channel full, drop event
	}
}

// OnDisconnected is called when a peer connection is lost
func (cc *connectionCallback) OnDisconnected(peer string) {
	// Check if channels are closed before sending
	if cc.service.events.IsClosed() {
		return
	}

	select {
	case cc.service.events.Connection <- ConnectionEvent{
		Timestamp: time.Now(),
		Type:      "disconnected",
		Peer:      peer,
	}:
	default:
		// Channel full, drop event
	}
}

// OnConnectionError is called when a connection error occurs
func (cc *connectionCallback) OnConnectionError(peer, errorMsg string) {
	// Check if channels are closed before sending
	if cc.service.events.IsClosed() {
		return
	}

	select {
	case cc.service.events.Connection <- ConnectionEvent{
		Timestamp:    time.Now(),
		Type:         "error",
		Peer:         peer,
		ErrorMessage: errorMsg,
	}:
	default:
		// Channel full, drop event
	}
}
