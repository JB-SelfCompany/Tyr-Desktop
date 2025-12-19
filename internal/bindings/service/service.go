package service

import (
	"fmt"
	"log"

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
func RestartService(sm *core.ServiceManager) error {
	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	log.Println("Restarting service...")

	// Stop service gracefully
	if err := sm.SoftStop(); err != nil {
		log.Printf("Warning: failed to soft stop service: %v, trying normal stop", err)
		if err := sm.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
	}

	// Reinitialize and start
	if err := sm.Initialize(); err != nil {
		return fmt.Errorf("failed to reinitialize service: %w", err)
	}

	if err := sm.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
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

	peers := cfg.GetEnabledPeers()
	if len(peers) == 0 {
		return fmt.Errorf("no enabled peers configured")
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
