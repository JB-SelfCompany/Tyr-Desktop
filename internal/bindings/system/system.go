package system

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	goruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/models"
	uitheme "github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/theme"
)

// Dialog Functions

// ShowOpenFileDialog shows a file open dialog and returns the selected file path
func ShowOpenFileDialog(ctx context.Context, title string) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: title,
	})

	return selection, err
}

// ShowSaveFileDialog shows a file save dialog and returns the selected file path
func ShowSaveFileDialog(ctx context.Context, title string, defaultFilename string) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	selection, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
	})

	return selection, err
}

// ShowOpenDirectoryDialog shows a directory selection dialog
func ShowOpenDirectoryDialog(ctx context.Context, title string) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context not initialized")
	}

	selection, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title: title,
	})

	return selection, err
}

// ShowMessageDialog shows a message dialog with OK button
func ShowMessageDialog(ctx context.Context, title string, message string) error {
	if ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	_, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})

	return err
}

// ShowErrorDialog shows an error dialog
func ShowErrorDialog(ctx context.Context, title string, message string) error {
	if ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	_, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   title,
		Message: message,
	})

	return err
}

// ShowQuestionDialog shows a question dialog with Yes/No buttons
// Returns true if user clicked Yes, false otherwise
func ShowQuestionDialog(ctx context.Context, title string, message string) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("context not initialized")
	}

	result, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   title,
		Message: message,
		Buttons: []string{"Yes", "No"},
	})

	return result == "Yes", err
}

// Clipboard and Browser Functions

// CopyToClipboard copies text to the system clipboard
func CopyToClipboard(ctx context.Context, text string) error {
	if ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	return runtime.ClipboardSetText(ctx, text)
}

// OpenURL opens a URL in the default browser
func OpenURL(ctx context.Context, url string) error {
	if ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	runtime.BrowserOpenURL(ctx, url)
	return nil
}

// Backup and Restore Functions

// CreateBackup creates an encrypted backup of the configuration and optionally database
func CreateBackup(ctx context.Context, cfg *core.Config, options models.BackupOptionsDTO) (models.ResultDTO, error) {
	if cfg == nil {
		return models.ResultDTO{Success: false, Message: "Config not initialized"}, nil
	}

	// Validate backup path
	if options.BackupPath == "" {
		return models.ResultDTO{Success: false, Message: "Backup path cannot be empty"}, nil
	}

	// Validate password
	if options.Password == "" {
		return models.ResultDTO{Success: false, Message: "Password cannot be empty"}, nil
	}

	// Create backup data
	log.Printf("Creating backup with includeDatabase=%v", options.IncludeDatabase)
	runtime.EventsEmit(ctx, "backup:progress", map[string]interface{}{"progress": 30, "message": "Creating backup data..."})
	backupData, err := core.CreateBackup(cfg, options.IncludeDatabase, options.Password)
	if err != nil {
		log.Printf("ERROR: Failed to create backup: %v", err)
		return models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to create backup: %v", err)}, nil
	}
	log.Printf("Backup data created successfully (%d bytes)", len(backupData))

	// Write backup to file
	runtime.EventsEmit(ctx, "backup:progress", map[string]interface{}{"progress": 80, "message": "Writing backup file..."})
	if err := core.WriteBackupFile(options.BackupPath, backupData); err != nil {
		return models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to write backup file: %v", err)}, nil
	}

	runtime.EventsEmit(ctx, "backup:progress", map[string]interface{}{"progress": 100, "message": "Backup completed successfully!"})
	return models.ResultDTO{Success: true, Message: "Backup created successfully", Data: options.BackupPath}, nil
}

// RestoreBackup restores configuration and optionally database from an encrypted backup
// Returns the restored config and result
func RestoreBackup(ctx context.Context, options models.RestoreOptionsDTO) (*core.Config, models.ResultDTO, error) {
	// CRITICAL: Load CURRENT config to get current database path
	// We want to restore database to CURRENT location, not to location from backup
	currentConfig, err := core.Load()
	currentDatabasePath := ""
	if err == nil {
		currentDatabasePath = currentConfig.ServiceSettings.DatabasePath
		log.Printf("[RestoreBackup] Current database path from config: %s", currentDatabasePath)
	} else {
		log.Printf("[RestoreBackup] Warning: Failed to load current config: %v - will use path from backup", err)
	}

	// Show open file dialog if path not provided
	runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 10, "message": "Selecting backup file..."})
	backupPath := options.BackupPath
	if backupPath == "" {
		var err error
		backupPath, err = ShowOpenFileDialog(ctx, "Select Backup File")
		if err != nil {
			return nil, models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to select file: %v", err)}, nil
		}
		if backupPath == "" {
			return nil, models.ResultDTO{Success: false, Message: "No file selected"}, nil
		}
	}

	// Validate password
	if options.Password == "" {
		return nil, models.ResultDTO{Success: false, Message: "Password cannot be empty"}, nil
	}

	// Read backup file
	runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 30, "message": "Reading backup file..."})
	backupData, err := core.ReadBackupFile(backupPath)
	if err != nil {
		return nil, models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to read backup file: %v", err)}, nil
	}

	// Restore backup
	runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 50, "message": "Decrypting backup..."})
	restoredConfig, dbData, err := core.RestoreBackup(backupData, options.Password)
	if err != nil {
		return nil, models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to restore backup: %v", err)}, nil
	}

	// CRITICAL: Override database path with CURRENT path if available
	// This ensures database is restored to user's current location, not backup location
	if currentDatabasePath != "" && restoredConfig.ServiceSettings.DatabasePath != currentDatabasePath {
		log.Printf("[RestoreBackup] Overriding database path from backup (%s) with current path (%s)",
			restoredConfig.ServiceSettings.DatabasePath, currentDatabasePath)
		restoredConfig.ServiceSettings.DatabasePath = currentDatabasePath
	}

	// IMPORTANT: Mark onboarding as complete BEFORE saving
	// This ensures the user doesn't see onboarding screen again
	restoredConfig.OnboardingComplete = true

	// CRITICAL: If database is included, restore it BEFORE saving config
	// This ensures database path in restoredConfig is used for restoration
	if dbData != nil && len(dbData) > 0 {
		log.Printf("Restoring database (%d bytes) to path: %s", len(dbData), restoredConfig.ServiceSettings.DatabasePath)
		runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 70, "message": "Restoring database..."})
		if err := core.RestoreDatabase(restoredConfig, dbData); err != nil {
			log.Printf("ERROR: Failed to restore database: %v", err)
			return restoredConfig, models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to restore database: %v", err)}, nil
		}
		log.Println("Database restored successfully")
	} else {
		log.Printf("WARNING: No database data to restore (dbData nil=%v, len=%d)", dbData == nil, len(dbData))
	}

	// Save config to disk AFTER database restoration
	runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 85, "message": "Saving configuration..."})
	if err := restoredConfig.Save(); err != nil {
		return nil, models.ResultDTO{Success: false, Message: fmt.Sprintf("Failed to save config: %v", err)}, nil
	}

	// CRITICAL: Reload config from disk to ensure in-memory state matches disk
	// This is necessary because window.location.reload() doesn't restart the Go backend
	runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 90, "message": "Reloading configuration..."})
	reloadedConfig, err := core.Load()
	if err != nil {
		log.Printf("Failed to reload config from disk: %v", err)
		return restoredConfig, models.ResultDTO{Success: true, Message: "Backup restored successfully"}, nil
	}

	log.Printf("Config reloaded from disk - DatabasePath: %s", reloadedConfig.ServiceSettings.DatabasePath)

	runtime.EventsEmit(ctx, "restore:progress", map[string]interface{}{"progress": 100, "message": "Restore completed successfully!"})

	// Emit config:restored event to notify frontend that configuration was restored
	// Frontend should reload its config state without full page reload to avoid tray issues
	runtime.EventsEmit(ctx, "config:restored", map[string]interface{}{"success": true})

	return reloadedConfig, models.ResultDTO{Success: true, Message: "Backup restored successfully"}, nil
}

// Window Control Functions

// QuitApplication gracefully quits the application
func QuitApplication(ctx context.Context) {
	if ctx == nil {
		return
	}
	runtime.Quit(ctx)
}

// HideWindow hides the application window (useful for system tray)
func HideWindow(ctx context.Context) {
	if ctx == nil {
		return
	}
	runtime.WindowHide(ctx)
}

// ShowWindow shows the application window
func ShowWindow(ctx context.Context) {
	if ctx == nil {
		return
	}
	runtime.WindowShow(ctx)
}

// MinimizeWindow minimizes the application window
func MinimizeWindow(ctx context.Context) {
	if ctx == nil {
		return
	}
	runtime.WindowMinimise(ctx)
}

// MaximizeWindow maximizes the application window
func MaximizeWindow(ctx context.Context) {
	if ctx == nil {
		return
	}
	runtime.WindowMaximise(ctx)
}

// ToggleMaximize toggles between maximized and normal window state
func ToggleMaximize(ctx context.Context) {
	if ctx == nil {
		return
	}
	runtime.WindowToggleMaximise(ctx)
}

// System Information Functions

// GetSystemTheme returns the current system theme preference ("light" or "dark")
// Uses platform-specific detection (Windows registry, etc.)
func GetSystemTheme() (string, error) {
	return uitheme.GetSystemTheme()
}

// GetSystemLanguage returns the system's default language code ("en" or "ru")
// Uses platform-specific detection (Windows API or environment variables)
func GetSystemLanguage() string {
	return detectSystemLanguage()
}

// detectSystemLanguage is a helper that detects system language
// Uses platform-specific methods (Windows API on Windows, env vars on Linux)
func detectSystemLanguage() string {
	return detectSystemLanguageImpl()
}

// detectSystemLanguageImpl is implemented in platform-specific files:
// - system_windows.go for Windows
// - system_other.go for Linux/Unix

// DeltaChat Integration Functions

// OpenDeltaChat opens DeltaChat with auto-configured account using dclogin:// URL
func OpenDeltaChat(ctx context.Context, cfg *core.Config, sm *core.ServiceManager) error {
	if sm == nil {
		return fmt.Errorf("service manager not initialized")
	}

	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	// Get mail address
	mailAddress := sm.GetMailAddress()
	if mailAddress == "" {
		return fmt.Errorf("mail address not available - service may not be initialized")
	}

	// Get password
	password, err := cfg.GetPassword()
	if err != nil || password == "" {
		return fmt.Errorf("failed to retrieve password: %w", err)
	}

	// Parse SMTP and IMAP addresses
	smtpHost, smtpPort := parseAddress(cfg.ServiceSettings.SMTPAddress)
	imapHost, imapPort := parseAddress(cfg.ServiceSettings.IMAPAddress)

	// Generate dclogin:// URL with full IMAP/SMTP configuration
	dcloginURL := generateDCLoginURL(mailAddress, password, imapHost, imapPort, smtpHost, smtpPort)

	// Try to open DeltaChat directly with the dclogin:// URL
	if err := openDCLoginURL(dcloginURL); err != nil {
		// If opening fails, copy URL to clipboard as fallback
		log.Printf("Failed to open DeltaChat automatically: %v", err)
		log.Println("Copying dclogin URL to clipboard as fallback")

		if ctx != nil {
			if clipErr := runtime.ClipboardSetText(ctx, dcloginURL); clipErr != nil {
				return fmt.Errorf("failed to open DeltaChat and failed to copy URL to clipboard: %w", err)
			}
		}

		// Return error indicating that automatic opening failed but clipboard copy succeeded
		// The frontend should show this as a warning/info message, not an error
		return fmt.Errorf("could not open DeltaChat automatically - dclogin URL has been copied to clipboard. Please paste it in DeltaChat manually")
	}

	return nil
}

// Helper Functions for DeltaChat Integration

// generateDCLoginURL creates a dclogin:// URL for DeltaChat auto-configuration
// Format: dclogin://user@host/?v=1&p=password&ih=imaphost&ip=imapport&is=plain&ic=3&sh=smtphost&sp=smtpport&ss=plain&sc=3
func generateDCLoginURL(email, password, imapHost, imapPort, smtpHost, smtpPort string) string {
	// Build query parameters
	params := url.Values{}
	params.Set("v", "1")       // Version
	params.Set("p", password)  // Password
	params.Set("ih", imapHost) // IMAP hostname
	params.Set("ip", imapPort) // IMAP port
	params.Set("is", "plain")  // IMAP security (plain = no encryption)
	params.Set("ic", "3")      // IMAP certificate checks (3 = accept invalid, for localhost)
	params.Set("sh", smtpHost) // SMTP hostname
	params.Set("sp", smtpPort) // SMTP port
	params.Set("ss", "plain")  // SMTP security (plain = no encryption)
	params.Set("sc", "3")      // SMTP certificate checks (3 = accept invalid, for localhost)

	// Build dclogin URL
	// Format: dclogin://user@host/?parameters
	return fmt.Sprintf("dclogin://%s/?%s", email, params.Encode())
}

// parseAddress splits an address string into host and port
// Example: "127.0.0.1:1025" -> ("127.0.0.1", "1025")
func parseAddress(addr string) (host, port string) {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return addr, ""
}

// openDCLoginURL opens a dclogin:// URL with proper escaping for Windows
func openDCLoginURL(dcloginURL string) error {
	switch goruntime.GOOS {
	case "windows":
		// On Windows, cmd /c start requires special handling:
		// The first argument after "start" is the window title (can be empty "")
		// Then comes the URL which must be in quotes if it contains special chars
		// Use rundll32 instead to avoid cmd parsing issues
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", dcloginURL).Start()
	case "darwin":
		return exec.Command("open", dcloginURL).Start()
	default: // "linux", "freebsd", "openbsd", "netbsd"
		return exec.Command("xdg-open", dcloginURL).Start()
	}
}
