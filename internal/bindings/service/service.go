package service

import (
	"fmt"
	"log"
	"time"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/models"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/yggmail"
)

// InitializeService initializes the yggmail service
func InitializeService(sm *core.ServiceManager) error {
	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	if err := sm.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize service: %w", err)
	}

	log.Println("Service initialized successfully")
	return nil
}

// StartService starts the yggmail service
// Automatically initializes if not already initialized
// Returns true if event monitoring should be started
func StartService(sm *core.ServiceManager) (shouldStartMonitoring bool, err error) {
	if sm == nil {
		return false, fmt.Errorf("service manager not initialized")
	}

	if sm.IsRunning() {
		return false, fmt.Errorf("service is already running")
	}

	// Check if initialized
	if sm.GetEventChannels() == nil {
		log.Println("Service not initialized, initializing before start...")
		if err := sm.Initialize(); err != nil {
			return false, fmt.Errorf("failed to initialize service: %w", err)
		}
		log.Println("Service initialized successfully")
		shouldStartMonitoring = true
	}

	if err := sm.Start(); err != nil {
		return false, fmt.Errorf("failed to start service: %w", err)
	}

	log.Println("Service started successfully")
	return shouldStartMonitoring, nil
}

// StopService stops the yggmail service gracefully
func StopService(sm *core.ServiceManager) error {
	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	// Use SoftStop to gracefully disconnect peers first
	if err := sm.SoftStop(); err != nil {
		log.Printf("Warning: failed to soft stop service: %v, trying normal stop", err)
		if err := sm.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
	}

	log.Println("Service stopped successfully")
	return nil
}

// RestartService restarts the yggmail service
// Implements robust restart with proper waiting for complete shutdown
// Based on Android implementation pattern
func RestartService(sm *core.ServiceManager) error {
	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	log.Println("Restarting service...")

	// Step 1: Stop service gracefully
	if sm.IsRunning() {
		log.Println("Stopping service...")
		if err := sm.SoftStop(); err != nil {
			log.Printf("Warning: failed to soft stop service: %v, trying normal stop", err)
			if err := sm.Stop(); err != nil {
				return fmt.Errorf("failed to stop service: %w", err)
			}
		}

		// Step 2: Wait for service to fully stop (up to 10 seconds)
		// This is critical to prevent "address already in use" errors
		log.Println("Waiting for service to fully stop...")
		for i := 0; i < 50; i++ {
			if !sm.IsRunning() {
				log.Printf("Service stopped after %d ms", i*200)
				break
			}
			time.Sleep(200 * time.Millisecond)

			// Timeout after 10 seconds
			if i == 49 {
				return fmt.Errorf("service did not stop within 10 seconds")
			}
		}

		// Step 3: Additional brief delay to ensure ports are released
		// The native library may need extra time to release resources
		time.Sleep(500 * time.Millisecond)
	}

	// Step 4: Reinitialize service
	log.Println("Reinitializing service...")
	if err := sm.Initialize(); err != nil {
		return fmt.Errorf("failed to reinitialize service: %w", err)
	}

	// Step 5: Start service
	log.Println("Starting service...")
	if err := sm.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	// Step 6: Verify service started successfully
	time.Sleep(1 * time.Second)
	if !sm.IsRunning() {
		return fmt.Errorf("service failed to start after restart")
	}

	log.Println("Service restarted successfully")
	return nil
}

// GetServiceStatusDTO returns the current service status
func GetServiceStatusDTO(sm *core.ServiceManager, cfg *core.Config) models.ServiceStatusDTO {
	status := models.ServiceStatusDTO{
		Status:       "Stopped",
		Running:      false,
		MailAddress:  "",
		SMTPAddress:  "",
		IMAPAddress:  "",
		DatabasePath: "",
		ErrorMessage: "",
	}

	if sm == nil {
		return status
	}

	serviceStatus := sm.GetStatus()
	status.Status = serviceStatus.String()
	status.Running = (serviceStatus == yggmail.StatusRunning)

	if cfg != nil {
		status.SMTPAddress = cfg.ServiceSettings.SMTPAddress
		status.IMAPAddress = cfg.ServiceSettings.IMAPAddress
		status.DatabasePath = cfg.ServiceSettings.DatabasePath

		if status.Running {
			status.MailAddress = sm.GetMailAddress()
		}
	}

	return status
}

// GetPeerStatsDTO returns statistics for all configured peers
func GetPeerStatsDTO(sm *core.ServiceManager, cfg *core.Config) []models.PeerInfoDTO {
	if sm == nil || cfg == nil {
		return []models.PeerInfoDTO{}
	}

	peerStats := sm.GetPeerStats()
	peerStatsMap := make(map[string]*yggmail.PeerInfo)
	for i := range peerStats {
		peerStatsMap[peerStats[i].Address] = &peerStats[i]
	}

	var result []models.PeerInfoDTO
	for _, peerCfg := range cfg.NetworkPeers {
		dto := models.PeerInfoDTO{
			Address:   peerCfg.Address,
			Enabled:   peerCfg.Enabled,
			Connected: false,
			Latency:   0,
			Uptime:    0,
			RXBytes:   0,
			TXBytes:   0,
			RXRate:    0,
			TXRate:    0,
			LastError: "",
		}

		if stats, exists := peerStatsMap[peerCfg.Address]; exists {
			dto.Connected = stats.Status
			dto.Latency = stats.Latency
			dto.Uptime = stats.Uptime
			dto.RXBytes = stats.RXBytes
			dto.TXBytes = stats.TXBytes
			dto.RXRate = stats.RXRate
			dto.TXRate = stats.TXRate
			dto.LastError = stats.LastError
		}

		result = append(result, dto)
	}

	return result
}

// HotReloadPeers reloads the peer list without stopping the service
func HotReloadPeers(sm *core.ServiceManager, cfg *core.Config) error {
	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	if !sm.IsRunning() {
		return fmt.Errorf("service is not running")
	}

	// Get enabled peers from config
	peers := cfg.GetEnabledPeers()

	// Allow empty peer list - this will disconnect from all peers
	// This is useful when user wants to disable all peers temporarily
	if len(peers) == 0 {
		log.Println("Hot-reloading with empty peer list (will disconnect from all peers)")
	} else {
		log.Printf("Hot-reloading with %d enabled peer(s)", len(peers))
	}

	if err := sm.HotReloadPeers(peers); err != nil {
		return fmt.Errorf("failed to hot-reload peers: %w", err)
	}

	log.Println("Peers hot-reloaded successfully")
	return nil
}

// GetMailAddress returns the current yggmail address
func GetMailAddress(sm *core.ServiceManager) string {
	if sm == nil {
		return ""
	}
	return sm.GetMailAddress()
}

// IsServiceRunning returns whether the service is running
func IsServiceRunning(sm *core.ServiceManager) bool {
	if sm == nil {
		return false
	}
	return sm.IsRunning()
}
