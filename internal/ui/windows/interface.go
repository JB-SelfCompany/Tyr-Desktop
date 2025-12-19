package windows

import (
	"fyne.io/fyne/v2"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
)

// AppInterface defines the interface that windows need from the main application
// This allows windows to interact with the app without creating circular dependencies
type AppInterface interface {
	// GetConfig returns the application configuration
	GetConfig() *core.Config

	// GetServiceManager returns the service manager
	GetServiceManager() *core.ServiceManager

	// GetMainWindow returns the main application window
	GetMainWindow() fyne.Window

	// GetFyneApp returns the Fyne application instance
	GetFyneApp() fyne.App

	// ShowDashboard displays the dashboard screen
	ShowDashboard()

	// ShowSettings displays the settings screen
	ShowSettings()

	// ShowLogs displays the logs screen
	ShowLogs()

	// UpdateSystemTrayMenu updates the system tray menu with current localization
	UpdateSystemTrayMenu()

	// UpdateSystemTrayStatus updates the system tray menu with current service status
	UpdateSystemTrayStatus()
}
