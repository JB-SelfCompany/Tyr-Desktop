package tray

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/getlantern/systray"
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

	m.initialized = true

	// Update menu with current status
	m.updateMenu()

	// Handle menu actions in goroutine
	go m.handleTrayActions()

	log.Println("System tray initialized successfully")
}

// onTrayExit is called when the tray is exiting
func (m *Manager) onTrayExit() {
	log.Println("System tray exiting")
}

// handleTrayActions handles clicks on tray menu items
func (m *Manager) handleTrayActions() {
	for {
		select {
		case <-m.mShow.ClickedCh:
			log.Println("Tray: Show window clicked")
			if m.onShowCallback != nil {
				m.onShowCallback()
			}
		case <-m.mSettings.ClickedCh:
			log.Println("Tray: Settings clicked")
			if m.onSettingsCallback != nil {
				m.onSettingsCallback()
			}
		case <-m.mQuit.ClickedCh:
			log.Println("Tray: Quit clicked")
			if m.onQuitCallback != nil {
				m.onQuitCallback()
			}
			return
		}
	}
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

// ShowWindow shows the window from system tray
func ShowWindow(ctx context.Context) {
	if ctx == nil {
		return
	}

	runtime.WindowShow(ctx)
	runtime.WindowUnminimise(ctx)
	runtime.WindowCenter(ctx)
}

// ShowSettingsWindow shows the settings window from system tray
func ShowSettingsWindow(ctx context.Context) {
	if ctx == nil {
		return
	}

	// Show window first
	runtime.WindowShow(ctx)
	runtime.WindowUnminimise(ctx)
	runtime.WindowCenter(ctx)

	// Small delay to ensure window is visible before navigation
	// Then navigate to settings using direct path manipulation
	script := `
		setTimeout(() => {
			if (window.location.pathname !== '/settings') {
				window.history.pushState({}, '', '/settings');
				window.dispatchEvent(new PopStateEvent('popstate'));
			}
		}, 100);
	`
	runtime.WindowExecJS(ctx, script)
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
