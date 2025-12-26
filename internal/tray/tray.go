package tray

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/energye/systray"
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
	mutex           sync.Mutex
	initialized     bool
	mStatus         *systray.MenuItem
	mShow           *systray.MenuItem
	mSettings       *systray.MenuItem
	mQuit           *systray.MenuItem
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
	}
}

// Setup initializes and configures the system tray
func (m *Manager) Setup() {
	go systray.Run(m.onTrayReady, m.onTrayExit)
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

	// Create menu items once
	m.mStatus = systray.AddMenuItem("Service: Stopped", "")
	m.mStatus.Disable()

	systray.AddSeparator()

	m.mShow = systray.AddMenuItem(localizer.Get("app.show"), localizer.Get("app.show"))
	m.mSettings = systray.AddMenuItem(localizer.Get("app.settings"), localizer.Get("app.settings"))

	systray.AddSeparator()

	m.mQuit = systray.AddMenuItem(localizer.Get("app.quit"), localizer.Get("app.quit"))

	// Set up menu item click handlers using callback functions
	m.mShow.Click(func() {
		log.Println("Tray: Show window clicked")
		if m.onShowCallback != nil {
			m.onShowCallback()
		}
	})

	m.mSettings.Click(func() {
		log.Println("Tray: Settings clicked")
		if m.onSettingsCallback != nil {
			m.onSettingsCallback()
		}
	})

	m.mQuit.Click(func() {
		log.Println("Tray: Quit clicked")
		if m.onQuitCallback != nil {
			m.onQuitCallback()
		}
	})

	// Set double-click handler to show window
	systray.SetOnDClick(func(menu systray.IMenu) {
		log.Println("Tray: Double-click detected, showing window")
		if m.onShowCallback != nil {
			m.onShowCallback()
		}
	})

	m.initialized = true

	// Update menu with current status
	m.updateMenu()

	log.Println("System tray initialized successfully")
}

// onTrayExit is called when the tray is exiting
func (m *Manager) onTrayExit() {
	log.Println("System tray exiting")
}

// updateMenu updates system tray menu with current service status
func (m *Manager) updateMenu() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

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

	log.Printf("System tray menu updated: %s", statusInfo)
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

// ShowWindow shows the window from system tray with robust recovery
func ShowWindow(ctx context.Context) {
	if ctx == nil {
		log.Println("ShowWindow: context is nil, cannot show window")
		return
	}

	log.Println("ShowWindow: attempting to show window from tray")

	// Multi-step approach to ensure window is shown reliably on Windows
	// This fixes issues where window becomes unresponsive after being minimized for a long time

	// Step 1: Unminimize first (important for Windows)
	runtime.WindowUnminimise(ctx)
	log.Println("ShowWindow: window unminimized")

	// Step 2: Show the window
	runtime.WindowShow(ctx)
	log.Println("ShowWindow: window shown")

	// Step 3: Set window to always on top temporarily to force it to foreground
	// This is crucial on Windows where the window might be hidden behind other windows
	runtime.WindowSetAlwaysOnTop(ctx, true)
	log.Println("ShowWindow: window set to always on top (temporary)")

	// Step 4: Center the window
	runtime.WindowCenter(ctx)
	log.Println("ShowWindow: window centered")

	// Step 5: Remove always on top after a short delay
	// This ensures the window has time to appear before we remove the flag
	go func() {
		// Small delay to ensure window is fully visible
		// Using a goroutine so we don't block the tray click handler
		runtime.EventsEmit(ctx, "window:showing", nil)

		// Wait 200ms for the window manager to bring the window forward
		// This delay is crucial on Windows to ensure the window fully appears
		// before we remove the always-on-top flag
		time.Sleep(200 * time.Millisecond)

		runtime.WindowSetAlwaysOnTop(ctx, false)
		log.Println("ShowWindow: always on top removed, window should now be visible")
	}()
}

// ShowSettingsWindow shows the settings window from system tray
func ShowSettingsWindow(ctx context.Context) {
	if ctx == nil {
		log.Println("ShowSettingsWindow: context is nil, cannot show window")
		return
	}

	log.Println("ShowSettingsWindow: attempting to show settings from tray")

	// Use the same robust window showing approach as ShowWindow
	runtime.WindowUnminimise(ctx)
	runtime.WindowShow(ctx)
	runtime.WindowSetAlwaysOnTop(ctx, true)
	runtime.WindowCenter(ctx)

	log.Println("ShowSettingsWindow: window shown, navigating to settings")

	// Navigate to settings and remove always-on-top after delay
	go func() {
		runtime.EventsEmit(ctx, "window:showing", nil)

		// Small delay to ensure window is visible before navigation
		time.Sleep(100 * time.Millisecond)

		// Navigate to settings using direct path manipulation
		script := `
			if (window.location.pathname !== '/settings') {
				window.history.pushState({}, '', '/settings');
				window.dispatchEvent(new PopStateEvent('popstate'));
			}
		`
		runtime.WindowExecJS(ctx, script)

		// Wait a bit more before removing always-on-top
		time.Sleep(100 * time.Millisecond)
		runtime.WindowSetAlwaysOnTop(ctx, false)
		log.Println("ShowSettingsWindow: navigation complete, always on top removed")
	}()
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

	// Quit systray (this will exit the tray goroutine)
	systray.Quit()

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
