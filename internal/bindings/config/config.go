package config

import (
	"fmt"
	"log"
	"time"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/models"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)

// GetConfigDTO converts core.Config to ConfigDTO
func GetConfigDTO(cfg *core.Config) models.ConfigDTO {
	if cfg == nil {
		return models.ConfigDTO{}
	}

	// Build peer list
	peers := make([]models.PeerConfigDTO, len(cfg.NetworkPeers))
	for i, peer := range cfg.NetworkPeers {
		peers[i] = models.PeerConfigDTO{
			Address: peer.Address,
			Enabled: peer.Enabled,
		}
	}

	return models.ConfigDTO{
		OnboardingComplete: cfg.OnboardingComplete,
		Peers:              peers,
		Language:           cfg.UIPreferences.Language,
		Theme:              cfg.UIPreferences.Theme,
		AutoStart:          cfg.UIPreferences.AutoStart,
		SMTPAddress:        cfg.ServiceSettings.SMTPAddress,
		IMAPAddress:        cfg.ServiceSettings.IMAPAddress,
	}
}

// UpdateConfigFromDTO updates core.Config from ConfigDTO
func UpdateConfigFromDTO(cfg *core.Config, dto models.ConfigDTO) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Update peers
	cfg.NetworkPeers = make([]core.PeerConfig, len(dto.Peers))
	for i, peer := range dto.Peers {
		cfg.NetworkPeers[i] = core.PeerConfig{
			Address: peer.Address,
			Enabled: peer.Enabled,
		}
	}

	// Update UI preferences
	cfg.UIPreferences.Language = dto.Language
	cfg.UIPreferences.Theme = dto.Theme
	cfg.UIPreferences.AutoStart = dto.AutoStart

	// Update service settings
	cfg.ServiceSettings.SMTPAddress = dto.SMTPAddress
	cfg.ServiceSettings.IMAPAddress = dto.IMAPAddress

	return nil
}

// AddPeer adds a new peer to the configuration
func AddPeer(cfg *core.Config, address string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Use thread-safe method from config
	return cfg.AddPeer(address)
}

// RemovePeer removes a peer from the configuration
func RemovePeer(cfg *core.Config, address string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Use thread-safe method from config
	return cfg.RemovePeer(address)
}

// EnablePeer enables a peer in the configuration
func EnablePeer(cfg *core.Config, address string) error {
	return setPeerEnabled(cfg, address, true)
}

// DisablePeer disables a peer in the configuration
func DisablePeer(cfg *core.Config, address string) error {
	return setPeerEnabled(cfg, address, false)
}

// setPeerEnabled sets the enabled status of a peer
func setPeerEnabled(cfg *core.Config, address string, enabled bool) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Use thread-safe methods from config
	var err error
	if enabled {
		err = cfg.EnablePeer(address)
	} else {
		err = cfg.DisablePeer(address)
	}

	return err
}

// SetPassword sets the yggmail password
// Password is stored securely in the OS keyring
// Note: During onboarding, serviceManager doesn't exist yet.
// The password will be set in yggmail database on first Initialize() call.
func SetPassword(cfg *core.Config, sm *core.ServiceManager, password string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// If service manager is already initialized, use UpdatePassword
	// (saves to keyring AND updates service)
	if sm != nil {
		if err := sm.UpdatePassword(password); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}
	} else {
		// During onboarding, serviceManager doesn't exist yet
		// Just save to keyring - will be set in yggmail on first Initialize()
		if err := cfg.SetPassword(password); err != nil {
			return fmt.Errorf("failed to set password: %w", err)
		}
	}

	return nil
}

// SetLanguage sets the UI language and updates the i18n localizer
func SetLanguage(cfg *core.Config, language string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Validate language
	if language != "en" && language != "ru" {
		return fmt.Errorf("invalid language: %s (must be 'en' or 'ru')", language)
	}

	cfg.UIPreferences.Language = language

	// Update global localizer with new language
	localizer := i18n.GetGlobalLocalizer()
	if err := localizer.SetLanguage(language); err != nil {
		log.Printf("Failed to set localizer language: %v", err)
	}

	return nil
}

// SetTheme sets the UI theme
func SetTheme(cfg *core.Config, theme string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Validate theme
	if theme != "light" && theme != "dark" && theme != "system" {
		return fmt.Errorf("invalid theme: %s (must be 'light', 'dark', or 'system')", theme)
	}

	cfg.UIPreferences.Theme = theme
	return nil
}

// SetAutoStart sets whether the service should start on system boot
func SetAutoStart(cfg *core.Config, enabled bool) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	cfg.UIPreferences.AutoStart = enabled

	// Apply autostart setting to OS
	if enabled {
		if err := core.EnableAutoStart(); err != nil {
			log.Printf("Failed to enable autostart: %v", err)
		}
	} else {
		if err := core.DisableAutoStart(); err != nil {
			log.Printf("Failed to disable autostart: %v", err)
		}
	}

	return nil
}

// GetDefaultPeers returns a list of recommended default peers
func GetDefaultPeers() []string {
	return core.DefaultPeers
}

// ChangePassword changes the password after verifying the current password
// Returns error if current password is incorrect or if password update fails
func ChangePassword(cfg *core.Config, sm *core.ServiceManager, currentPassword, newPassword string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	if currentPassword == "" {
		return fmt.Errorf("current password cannot be empty")
	}

	if newPassword == "" {
		return fmt.Errorf("new password cannot be empty")
	}

	if len(newPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters")
	}

	// Verify current password by comparing with stored password
	storedPassword, err := cfg.GetPassword()
	if err != nil {
		return fmt.Errorf("failed to retrieve stored password: %w", err)
	}

	if storedPassword != currentPassword {
		return fmt.Errorf("current password is incorrect")
	}

	// Update password using service manager
	if sm != nil {
		if err := sm.UpdatePassword(newPassword); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}
	} else {
		// Fallback if service manager is not available
		if err := cfg.SetPassword(newPassword); err != nil {
			return fmt.Errorf("failed to set password: %w", err)
		}
	}

	return nil
}

// RegenerateKeys regenerates Yggdrasil keys by deleting the database and reinitializing
// WARNING: This will delete ALL mail data and change the email address
// Requires password verification for security
func RegenerateKeys(cfg *core.Config, sm *core.ServiceManager, password string) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Verify password before allowing destructive operation
	storedPassword, err := cfg.GetPassword()
	if err != nil {
		return fmt.Errorf("failed to retrieve stored password: %w", err)
	}

	if storedPassword != password {
		return fmt.Errorf("password is incorrect")
	}

	// Stop service if running
	wasRunning := sm.IsRunning()
	if wasRunning {
		if err := sm.SoftStop(); err != nil {
			if err := sm.Stop(); err != nil {
				return fmt.Errorf("failed to stop service: %w", err)
			}
		}

		// Wait for service to fully stop (up to 10 seconds)
		for i := 0; i < 50; i++ {
			if !sm.IsRunning() {
				break
			}
			time.Sleep(200 * time.Millisecond)

			// Timeout after 10 seconds
			if i == 49 {
				return fmt.Errorf("service did not stop within 10 seconds")
			}
		}
	}

	// Close the service to release the database file
	if err := sm.CloseService(); err != nil {
		// Try to restart service if it was running
		if wasRunning {
			if initErr := sm.Initialize(); initErr == nil {
				sm.Start()
			}
		}
		return fmt.Errorf("failed to close service: %w", err)
	}

	// Delete database to force key regeneration
	dbPath := cfg.ServiceSettings.DatabasePath
	if err := core.DeleteDatabase(dbPath); err != nil {
		// Try to restart service if it was running
		if wasRunning {
			if initErr := sm.Initialize(); initErr == nil {
				sm.Start()
			}
		}
		return fmt.Errorf("failed to delete database: %w", err)
	}

	// Mark password as not initialized to force re-setting on next Initialize()
	cfg.ServiceSettings.PasswordInitialized = false
	if err := cfg.Save(); err != nil {
		log.Printf("Failed to save config: %v", err)
	}

	// Reinitialize service (this creates new database with new keys)
	if err := sm.Initialize(); err != nil {
		return fmt.Errorf("failed to reinitialize service: %w", err)
	}

	// Restart service if it was running
	if wasRunning {
		if err := sm.Start(); err != nil {
			return fmt.Errorf("failed to restart service: %w", err)
		}
	}

	return nil
}
