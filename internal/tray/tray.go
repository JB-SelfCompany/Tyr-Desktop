package tray

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)


// Manager manages the system tray functionality
type Manager struct {
	ctx            context.Context
	config         *core.Config
	serviceManager *core.ServiceManager

	// Callbacks
	onShowCallback     func()
	onSettingsCallback func()
	onQuitCallback     func()

	// Menu items
	mutex       sync.Mutex
	initialized bool
	mStatus     *systray.MenuItem
	mShow       *systray.MenuItem
	mSettings   *systray.MenuItem
	mQuit       *systray.MenuItem

	// Shutdown channel to stop click handler goroutines
	shutdownCh chan struct{}
}

// NewManager creates a new tray manager
func NewManager(
	ctx context.Context,
	cfg *core.Config,
	sm *core.ServiceManager,
	onShow, onSettings, onQuit func(),
) *Manager {
	return &Manager{
		ctx:                ctx,
		config:             cfg,
		serviceManager:     sm,
		onShowCallback:     onShow,
		onSettingsCallback: onSettings,
		onQuitCallback:     onQuit,
		shutdownCh:         make(chan struct{}),
	}
}

// Setup initializes and configures the system tray
// Uses systray.Run in a goroutine for proper integration with Wails
// Based on fyne.io/systray example and community recommendations
func (m *Manager) Setup() {
	// Run systray in a dedicated goroutine
	// This prevents blocking Wails' main event loop while allowing systray to have its own message pump
	// See: https://github.com/fyne-io/systray/blob/master/example/main.go
	go func() {
		log.Println("Starting system tray...")
		systray.Run(m.onTrayReady, m.onTrayExit)
		log.Println("System tray has exited")
	}()
}

// onTrayReady is called when the system tray is ready
func (m *Manager) onTrayReady() {
	// Set tray icon using platform-specific icon data
	// Windows uses ICO format (icon_windows.go)
	// Linux uses PNG format (icon_linux.go)
	systray.SetIcon(GetIconData())
	systray.SetTitle("Tyr")

	// Get localizer for translations
	localizer := i18n.GetGlobalLocalizer()
	systray.SetTooltip(localizer.Get("systray.description"))

	// Create menu items once - use localized initial status
	initialStatus := fmt.Sprintf("%s: %s", localizer.Get("systray.service_status"), localizer.Get("dashboard.status.stopped"))
	m.mStatus = systray.AddMenuItem(initialStatus, "")
	m.mStatus.Disable()

	systray.AddSeparator()

	m.mShow = systray.AddMenuItem(localizer.Get("app.show"), localizer.Get("app.show"))
	m.mSettings = systray.AddMenuItem(localizer.Get("app.settings"), localizer.Get("app.settings"))

	systray.AddSeparator()

	m.mQuit = systray.AddMenuItem(localizer.Get("app.quit"), localizer.Get("app.quit"))

	// Set up menu item click handlers using channels (fyne.io/systray API)
	// Each handler runs in its own goroutine and listens to the ClickedCh channel
	go m.handleShowClicks()
	go m.handleSettingsClicks()
	go m.handleQuitClicks()

	// REMOVED: SetOnTapped double-click implementation
	// Reason: SetOnTapped can cause event loop blocking and unresponsiveness on Windows
	// especially after the window has been hidden for a while.
	// Users should use the "Show" menu item instead - this is more reliable.

	m.initialized = true

	// Update menu with current status
	m.updateMenu()

	log.Println("System tray initialized successfully")
}

// handleShowClicks listens for clicks on the Show menu item
func (m *Manager) handleShowClicks() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("handleShowClicks: recovered from panic: %v", r)
		}
	}()

	for {
		select {
		case <-m.shutdownCh:
			return
		case <-m.mShow.ClickedCh:
			log.Println("Tray: Show window clicked")

			// Wrap callback in recovery to prevent panics from crashing the handler
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("handleShowClicks callback: recovered from panic: %v", r)
					}
				}()

				if m.onShowCallback != nil {
					m.onShowCallback()
				}
			}()
		}
	}
}

// handleSettingsClicks listens for clicks on the Settings menu item
func (m *Manager) handleSettingsClicks() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("handleSettingsClicks: recovered from panic: %v", r)
		}
	}()

	for {
		select {
		case <-m.shutdownCh:
			return
		case <-m.mSettings.ClickedCh:
			log.Println("Tray: Settings clicked")

			// Wrap callback in recovery to prevent panics from crashing the handler
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("handleSettingsClicks callback: recovered from panic: %v", r)
					}
				}()

				if m.onSettingsCallback != nil {
					m.onSettingsCallback()
				}
			}()
		}
	}
}

// handleQuitClicks listens for clicks on the Quit menu item
func (m *Manager) handleQuitClicks() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("handleQuitClicks: recovered from panic: %v", r)
		}
	}()

	for {
		select {
		case <-m.shutdownCh:
			return
		case <-m.mQuit.ClickedCh:
			log.Println("Tray: Quit clicked")

			// Wrap callback in recovery to prevent panics from crashing the handler
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("handleQuitClicks callback: recovered from panic: %v", r)
					}
				}()

				if m.onQuitCallback != nil {
					m.onQuitCallback()
				}
			}()
		}
	}
}

// onTrayExit is called when the tray is exiting
func (m *Manager) onTrayExit() {
	log.Println("System tray exiting")
}

// Cleanup properly cleans up the system tray on application exit
// This must be called before the application quits to prevent resource leaks
func (m *Manager) Cleanup() {
	log.Println("Cleaning up system tray...")

	// Signal all click handler goroutines to stop
	select {
	case <-m.shutdownCh:
		// Already closed
	default:
		close(m.shutdownCh)
	}

	// Quit the systray - this will trigger onTrayExit callback
	// and cleanly shutdown the systray event loop
	systray.Quit()

	log.Println("System tray cleanup complete")
}

// updateMenuInternal updates system tray menu with current service status
// MUST be called with mutex already locked or when mutex is not needed
func (m *Manager) updateMenuInternal() {
	if !m.initialized {
		log.Println("Tray not initialized yet, skipping update")
		return
	}

	if m.serviceManager == nil {
		log.Println("Service manager not initialized, skipping tray update")
		return
	}

	// Get localizer for translations
	localizer := i18n.GetGlobalLocalizer()

	// Get service status
	status := m.serviceManager.GetStatus()
	statusKey := "dashboard.status." + strings.ToLower(status.String())
	statusStr := localizer.Get(statusKey)
	statusInfo := fmt.Sprintf("%s: %s", localizer.Get("systray.service_status"), statusStr)

	// Update menu item titles
	m.mStatus.SetTitle(statusInfo)
	m.mShow.SetTitle(localizer.Get("app.show"))
	m.mSettings.SetTitle(localizer.Get("app.settings"))
	m.mQuit.SetTitle(localizer.Get("app.quit"))

	// Update tooltip
	systray.SetTooltip(localizer.Get("systray.description"))
}

// updateMenu updates system tray menu with current service status
func (m *Manager) updateMenu() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.updateMenuInternal()
}

// UpdateMenu updates the system tray menu with current localization
// Call this method after changing the application language
func (m *Manager) UpdateMenu() {
	m.updateMenu()
	log.Println("System tray menu updated with new language")
}

// UpdateStatus updates the system tray menu with current service status
// Call this when service status changes (start/stop/error)
func (m *Manager) UpdateStatus() {
	m.updateMenu()
}

// SetServiceManager updates the service manager reference
// This should be called when service manager is initialized after tray creation (e.g., after onboarding)
func (m *Manager) SetServiceManager(sm *core.ServiceManager) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.serviceManager = sm

	// Update menu immediately with new service manager status
	if m.initialized && sm != nil {
		m.updateMenuInternal()
	}
}

// ShowWindow shows the window from system tray with robust recovery
func ShowWindow(ctx context.Context) {
	if ctx == nil {
		log.Println("ShowWindow: context is nil, cannot show window")
		return
	}

	log.Println("ShowWindow: attempting to show window from tray")

	// Simplified approach for showing window on Windows
	// Previous complex approach with goroutines and delays could cause race conditions
	// and event loop blocking, leading to unresponsiveness

	// Use defer with recover to prevent panics from blocking the tray
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ShowWindow: recovered from panic: %v", r)
		}
	}()

	// Step 1: Unminimize first (important for Windows)
	runtime.WindowUnminimise(ctx)

	// Step 2: Show the window
	runtime.WindowShow(ctx)

	// Step 3: Bring to front using WindowSetAlwaysOnTop temporarily
	// IMPORTANT: We do this synchronously now, no goroutines with delays
	runtime.WindowSetAlwaysOnTop(ctx, true)
	runtime.WindowSetAlwaysOnTop(ctx, false)

	// Step 4: Center the window
	runtime.WindowCenter(ctx)

	// Emit event for frontend
	runtime.EventsEmit(ctx, "window:showing", nil)

	log.Println("ShowWindow: window shown successfully")
}

// ShowSettingsWindow shows the settings window from system tray
func ShowSettingsWindow(ctx context.Context) {
	if ctx == nil {
		log.Println("ShowSettingsWindow: context is nil, cannot show window")
		return
	}

	log.Println("ShowSettingsWindow: attempting to show settings from tray")

	// Use defer with recover to prevent panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ShowSettingsWindow: recovered from panic: %v", r)
		}
	}()

	// Simplified approach - show window first
	runtime.WindowUnminimise(ctx)
	runtime.WindowShow(ctx)
	runtime.WindowSetAlwaysOnTop(ctx, true)
	runtime.WindowSetAlwaysOnTop(ctx, false)
	runtime.WindowCenter(ctx)

	// Navigate to settings using direct path manipulation
	script := `
		if (window.location.pathname !== '/settings') {
			window.history.pushState({}, '', '/settings');
			window.dispatchEvent(new PopStateEvent('popstate'));
		}
	`
	runtime.WindowExecJS(ctx, script)

	// Emit event for frontend
	runtime.EventsEmit(ctx, "window:showing", nil)

	log.Println("ShowSettingsWindow: window shown and navigated to settings")
}

// QuitApplication quits the application from system tray
func QuitApplication(ctx context.Context, performShutdown func()) {
	if ctx == nil {
		return
	}

	// Perform cleanup
	if performShutdown != nil {
		performShutdown()
	}

	// Note: systray cleanup is handled by Manager.Cleanup()
	// which is called from performShutdown (app.actualShutdown)

	// Quit the application
	runtime.Quit(ctx)
}

// HideWindow hides the window to system tray
func HideWindow(ctx context.Context) {
	if ctx == nil {
		return
	}

	runtime.WindowHide(ctx)
}
