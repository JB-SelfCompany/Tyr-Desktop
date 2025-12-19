package ui

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/windows"
)

// App represents the main application coordinator
// Manages window navigation, system tray, and application lifecycle
type App struct {
	// Fyne application instance
	fyneApp fyne.App

	// Main application window
	mainWindow fyne.Window

	// Core components
	config         *core.Config
	serviceManager *core.ServiceManager

	// Current screen tracking
	currentScreen string

	// System tray (desktop-specific)
	systray desktop.App

	// System tray state cache
	lastTrayStatus string

	// Status monitoring control
	stopStatusMonitoring chan struct{}
}

// Screen names
const (
	ScreenOnboarding = "onboarding"
	ScreenDashboard  = "dashboard"
	ScreenSettings   = "settings"
	ScreenPeers      = "peers"
	ScreenBackup     = "backup"
	ScreenLogs       = "logs"
)

// NewApp creates a new application instance
// Initializes UI components and system tray
func NewApp(fyneApp fyne.App, config *core.Config, serviceManager *core.ServiceManager) *App {
	app := &App{
		fyneApp:              fyneApp,
		config:               config,
		serviceManager:       serviceManager,
		stopStatusMonitoring: make(chan struct{}),
	}

	// Initialize localization with configured language
	localizer := i18n.GetGlobalLocalizer()
	if config.UIPreferences.Language != "" {
		if err := localizer.SetLanguage(config.UIPreferences.Language); err != nil {
			log.Printf("Warning: Failed to set language to %s: %v", config.UIPreferences.Language, err)
		} else {
			log.Printf("Localization initialized with language: %s", config.UIPreferences.Language)
		}
	}

	// Create main window
	app.mainWindow = fyneApp.NewWindow("Tyr")

	// Restore window state from config
	app.restoreWindowState()

	// Set up window close handler (minimize to tray instead of quit)
	app.mainWindow.SetCloseIntercept(func() {
		app.saveWindowState()
		app.mainWindow.Hide()
	})

	// Set up system tray if available
	if desk, ok := fyneApp.(desktop.App); ok {
		app.systray = desk
		app.setupSystemTray()
	}

	// Set up keyboard shortcuts
	app.SetupKeyboardShortcuts()

	return app
}

// Run starts the application
// Shows onboarding or dashboard based on configuration
// If minimized is true, starts in system tray without showing window
func (a *App) Run(minimized bool) {
	// Start monitoring service status for reactive tray updates
	a.startStatusMonitoring()

	// Determine initial screen
	if !a.config.OnboardingComplete {
		a.ShowOnboarding()
	} else {
		a.ShowDashboard()

		// Auto-start service if minimized (autostart mode) or if autostart is enabled in preferences
		if minimized || a.config.UIPreferences.AutoStart {
			log.Println("Auto-starting service...")

			// First initialize the service (required before Start)
			if err := a.serviceManager.Initialize(); err != nil {
				log.Printf("Warning: Failed to initialize service for auto-start: %v", err)
			} else if err := a.serviceManager.Start(); err != nil {
				log.Printf("Warning: Failed to auto-start service: %v", err)
				// Don't show error dialog on startup - user can manually start if needed
			} else {
				log.Println("Service auto-started successfully")
				// Status will be updated automatically via status monitoring
			}
		}
	}

	// If minimized mode, hide window and run in system tray
	// Otherwise show main window normally
	if minimized && a.config.OnboardingComplete {
		log.Println("Running in minimized mode (system tray)")
		// Don't show window - just run the app
		a.mainWindow.Hide()
		a.fyneApp.Run()
	} else {
		// Show main window
		a.mainWindow.ShowAndRun()
	}
}

// ShowOnboarding displays the onboarding wizard
func (a *App) ShowOnboarding() {
	a.currentScreen = ScreenOnboarding
	a.mainWindow.SetTitle("Tyr - Setup")

	content := windows.NewModernOnboardingScreen(a)
	a.mainWindow.SetContent(content)
}

// ShowDashboard displays the main dashboard screen
func (a *App) ShowDashboard() {
	a.currentScreen = ScreenDashboard
	a.mainWindow.SetTitle("Tyr")

	content := windows.NewModernDashboardScreen(a)
	a.mainWindow.SetContent(content)

	// Show window if hidden
	a.mainWindow.Show()
}

// ShowSettings displays the settings screen
func (a *App) ShowSettings() {
	a.currentScreen = ScreenSettings
	a.mainWindow.SetTitle("Tyr - Settings")

	content := windows.NewModernSettingsScreen(a)
	a.mainWindow.SetContent(content)

	// Show window if hidden
	a.mainWindow.Show()
}

// ShowPeers displays the peer management screen
func (a *App) ShowPeers() {
	a.currentScreen = ScreenPeers
	a.mainWindow.SetTitle("Tyr - Network Peers")

	content := windows.NewModernPeersScreen(a)
	a.mainWindow.SetContent(content)

	// Show window if hidden
	a.mainWindow.Show()
}

// ShowBackup displays the backup/restore screen
func (a *App) ShowBackup() {
	a.currentScreen = ScreenBackup
	a.mainWindow.SetTitle("Tyr - Backup & Restore")

	content := windows.NewModernBackupScreen(a)
	a.mainWindow.SetContent(content)

	// Show window if hidden
	a.mainWindow.Show()
}

// ShowLogs displays the logs screen with real-time logging
func (a *App) ShowLogs() {
	a.currentScreen = ScreenLogs
	a.mainWindow.SetTitle("Tyr - System Logs")

	content := windows.NewLogsScreen(a)
	a.mainWindow.SetContent(content)

	// Show window if hidden
	a.mainWindow.Show()
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *core.Config {
	return a.config
}

// GetServiceManager returns the service manager
func (a *App) GetServiceManager() *core.ServiceManager {
	return a.serviceManager
}

// GetMainWindow returns the main window
func (a *App) GetMainWindow() fyne.Window {
	return a.mainWindow
}

// GetFyneApp returns the Fyne application instance
func (a *App) GetFyneApp() fyne.App {
	return a.fyneApp
}

// setupSystemTray configures the system tray icon and menu
func (a *App) setupSystemTray() {
	if a.systray == nil {
		return
	}

	// Update tray menu with current status
	a.updateTrayMenu()

	log.Println("System tray configured")
}

// updateTrayMenu updates system tray menu with current service status
func (a *App) updateTrayMenu() {
	if a.systray == nil {
		return
	}

	// Get localizer for translations
	localizer := i18n.GetGlobalLocalizer()

	// Get service status
	status := a.serviceManager.GetStatus()
	statusKey := "dashboard.status." + strings.ToLower(status.String())
	statusStr := localizer.Get(statusKey)
	statusInfo := fmt.Sprintf("%s: %s", localizer.Get("systray.service_status"), statusStr)

	// Check if status changed - avoid unnecessary menu recreation
	if a.lastTrayStatus == statusInfo {
		return
	}

	// Update cache
	a.lastTrayStatus = statusInfo

	// Create system tray menu with status
	menu := fyne.NewMenu("Tyr",
		// Service status (non-clickable info item)
		fyne.NewMenuItem(statusInfo, func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(localizer.Get("app.show"), func() {
			a.mainWindow.Show()
			a.mainWindow.RequestFocus()
		}),
		fyne.NewMenuItem(localizer.Get("app.settings"), func() {
			a.ShowSettings()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(localizer.Get("app.quit"), func() {
			a.Quit()
		}),
	)

	a.systray.SetSystemTrayMenu(menu)
	log.Printf("System tray menu updated: %s", statusInfo)
}

// UpdateSystemTrayMenu updates the system tray menu with current localization
// Call this method after changing the application language
func (a *App) UpdateSystemTrayMenu() {
	// Reset cache to force menu recreation with new translations
	a.lastTrayStatus = ""
	a.updateTrayMenu()
	log.Println("System tray menu updated with new language")
}

// UpdateSystemTrayStatus updates the system tray menu with current service status
// Call this when service status changes (start/stop/error)
func (a *App) UpdateSystemTrayStatus() {
	a.updateTrayMenu()
}

// Quit performs application cleanup and exits
func (a *App) Quit() {
	log.Println("Quitting application...")

	// Stop status monitoring
	close(a.stopStatusMonitoring)

	// Save window state before quitting
	a.saveWindowState()

	// Stop service if running (use SoftStop for clean peer disconnection)
	if a.serviceManager.IsRunning() {
		log.Println("Stopping service...")
		if err := a.serviceManager.SoftStop(); err != nil {
			log.Printf("Error stopping service: %v", err)
		}
	}

	// Shutdown service manager
	if err := a.serviceManager.Shutdown(); err != nil {
		log.Printf("Error shutting down service manager: %v", err)
	}

	// Quit Fyne application
	a.fyneApp.Quit()
}

// startStatusMonitoring starts a background goroutine that monitors service status
// and reactively updates the system tray menu when status changes
func (a *App) startStatusMonitoring() {
	// Get status channel from service manager
	statusChan := a.serviceManager.GetStatusChannel()

	go func() {
		log.Println("Status monitoring started for system tray updates")
		for {
			select {
			case <-a.stopStatusMonitoring:
				log.Println("Status monitoring stopped")
				return

			case status, ok := <-statusChan:
				if !ok {
					// Channel closed
					log.Println("Status channel closed, stopping monitoring")
					return
				}

				// Update system tray menu on Fyne main thread
				fyne.Do(func() {
					a.updateTrayMenu()
					log.Printf("System tray updated with status: %s", status.String())
				})
			}
		}
	}()
}

// ShowError displays an error dialog
func (a *App) ShowError(title string, err error) {
	log.Printf("Error [%s]: %v", title, err)

	if a.mainWindow != nil {
		// Use Fyne dialog to show error
		fyne.Do(func() {
			dialog := fyne.NewErrorDialog(title, err, a.mainWindow)
			dialog.Show()
		})
	}
}

// ShowInfo displays an information dialog
func (a *App) ShowInfo(title, message string) {
	log.Printf("Info [%s]: %s", title, message)

	if a.mainWindow != nil {
		// Use Fyne dialog to show information
		fyne.Do(func() {
			dialog := fyne.NewInformationDialog(title, message, a.mainWindow)
			dialog.Show()
		})
	}
}

// restoreWindowState restores the window size and position from config
// Always uses default compact size for consistent user experience
func (a *App) restoreWindowState() {
	// Always use default compact size (ignore saved size for consistency)
	width := core.DefaultWindowWidth
	height := core.DefaultWindowHeight

	a.mainWindow.Resize(fyne.NewSize(float32(width), float32(height)))

	// Always center on screen
	a.mainWindow.CenterOnScreen()

	log.Printf("Window state set to default: %dx%d (centered)", width, height)
}

// saveWindowState saves the current window size and position to config
// Called when window closes or application quits
func (a *App) saveWindowState() {
	// Get current window size
	size := a.mainWindow.Canvas().Size()
	width := int(size.Width)
	height := int(size.Height)

	// Note: Fyne v2 doesn't provide window position API
	// We save size but keep position as -1 (centered)
	x := -1
	y := -1

	// Save to config
	if err := a.config.SaveWindowState(width, height, x, y); err != nil {
		log.Printf("Error saving window state: %v", err)
		return
	}

	// Persist config to disk
	if err := a.config.Save(); err != nil {
		log.Printf("Error persisting config: %v", err)
	}

	log.Printf("Window state saved: %dx%d", width, height)
}

