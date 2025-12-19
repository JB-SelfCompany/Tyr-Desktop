package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/bindings/config"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/bindings/events"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/bindings/peerdiscovery"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/bindings/service"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/bindings/system"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/tray"
)

// App struct holds the application state and services
type App struct {
	// ctx is the Wails runtime context
	ctx context.Context

	// config holds the application configuration
	config *core.Config

	// serviceManager manages the yggmail service lifecycle
	serviceManager *core.ServiceManager

	// trayManager manages the system tray
	trayManager *tray.Manager

	// eventMonitorShutdown signals the event monitoring goroutine to stop
	eventMonitorShutdown chan struct{}

	// eventMonitorRunning tracks if event monitoring is already running
	eventMonitorRunning bool

	// statusMonitorShutdown signals the status monitoring goroutine to stop
	statusMonitorShutdown chan struct{}

	// statusMonitorRunning tracks if status monitoring is already running
	statusMonitorRunning bool

	// peerDiscoveryCtx is the context for peer discovery operations
	peerDiscoveryCtx context.Context

	// peerDiscoveryCancelFunc cancels the peer discovery context
	peerDiscoveryCancelFunc context.CancelFunc

	// allowQuit controls whether the application can actually quit
	allowQuit bool

	// startMinimized indicates if the app was started with --minimized flag
	startMinimized bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		eventMonitorShutdown:    make(chan struct{}),
		statusMonitorShutdown:   make(chan struct{}),
		peerDiscoveryCtx:        ctx,
		peerDiscoveryCancelFunc: cancel,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize configuration
	cfg, err := core.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}
	a.config = cfg

	// Initialize global localizer with config language to ensure tray uses correct language
	// This must be done before tray initialization in domReady
	if err := config.SetLanguage(a.config, a.config.UIPreferences.Language); err != nil {
		log.Printf("Failed to set initial language: %v", err)
	}

	// Initialize service manager
	if a.config.OnboardingComplete {
		sm, err := core.NewServiceManager(a.config)
		if err != nil {
			log.Printf("Failed to create service manager: %v", err)
		} else {
			a.serviceManager = sm

			if err := a.serviceManager.Initialize(); err != nil {
				log.Printf("Failed to initialize service: %v", err)
			} else if a.config.UIPreferences.AutoStart {
				if err := a.serviceManager.Start(); err != nil {
					log.Printf("Failed to auto-start service: %v", err)
				}
			}
		}
	}
}

// domReady is called after the frontend DOM is ready
func (a *App) domReady(ctx context.Context) {
	// Setup system tray
	a.setupTray()

	// Start event monitoring if service is available
	if a.serviceManager != nil && !a.eventMonitorRunning {
		a.eventMonitorRunning = true
		go a.startEventMonitoring()
	}

	// Start status monitoring for system tray updates
	if a.serviceManager != nil && !a.statusMonitorRunning {
		a.statusMonitorRunning = true
		go a.startStatusMonitoring()
	}
}

// beforeClose is called before the application window closes
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.allowQuit {
		a.actualShutdown()
		return false
	}

	tray.HideWindow(a.ctx)
	return true
}

// actualShutdown performs the actual shutdown operations
func (a *App) actualShutdown() {
	a.cancelPeerDiscoveryOperations()

	if a.eventMonitorShutdown != nil {
		select {
		case <-a.eventMonitorShutdown:
		default:
			close(a.eventMonitorShutdown)
		}
	}

	if a.statusMonitorShutdown != nil {
		select {
		case <-a.statusMonitorShutdown:
		default:
			close(a.statusMonitorShutdown)
		}
	}

	if a.serviceManager != nil && a.serviceManager.IsRunning() {
		if err := a.serviceManager.SoftStop(); err != nil {
			if err := a.serviceManager.Stop(); err != nil {
				log.Printf("Failed to stop service: %v", err)
			}
		}
	}

	if a.config != nil {
		if err := a.config.Save(); err != nil {
			log.Printf("Failed to save config: %v", err)
		}
	}
}

// shutdown is called when the application is terminating
func (a *App) shutdown(ctx context.Context) {
	a.cancelPeerDiscoveryOperations()

	if a.eventMonitorShutdown != nil {
		select {
		case <-a.eventMonitorShutdown:
		default:
			close(a.eventMonitorShutdown)
		}
	}

	if a.statusMonitorShutdown != nil {
		select {
		case <-a.statusMonitorShutdown:
		default:
			close(a.statusMonitorShutdown)
		}
	}

	if a.config != nil {
		if err := a.config.Save(); err != nil {
			log.Printf("Failed to save config: %v", err)
		}
	}

	if a.serviceManager != nil {
		if err := a.serviceManager.Shutdown(); err != nil {
			log.Printf("Failed to shutdown service manager: %v", err)
		}
	}
}

// OnStartupComplete is called after onboarding to initialize the service
func (a *App) OnStartupComplete() error {
	if a.serviceManager != nil {
		return fmt.Errorf("service manager already initialized")
	}

	sm, err := core.NewServiceManager(a.config)
	if err != nil {
		return fmt.Errorf("failed to create service manager: %w", err)
	}

	a.serviceManager = sm

	if err := a.serviceManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize service: %w", err)
	}

	if err := a.serviceManager.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	if !a.eventMonitorRunning {
		a.eventMonitorRunning = true
		go a.startEventMonitoring()
	}

	return nil
}

// IsOnboardingComplete returns whether the initial setup is complete
func (a *App) IsOnboardingComplete() bool {
	if a.config == nil {
		return false
	}
	return a.config.OnboardingComplete
}

// SetOnboardingComplete marks the onboarding as complete
func (a *App) SetOnboardingComplete() error {
	if a.config == nil {
		return fmt.Errorf("config not initialized")
	}

	a.config.OnboardingComplete = true
	a.config.UIPreferences.AutoStart = true

	if err := core.EnableAutoStart(); err != nil {
		log.Printf("Failed to enable system autostart: %v", err)
	}

	if err := a.config.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// ToggleWindowVisibility toggles between showing and hiding the window
func (a *App) ToggleWindowVisibility() {
	if a.ctx == nil {
		return
	}

	runtime.WindowShow(a.ctx)
	runtime.WindowUnminimise(a.ctx)
}

// getPeerDiscoveryContext returns a context for peer discovery operations
func (a *App) getPeerDiscoveryContext() context.Context {
	select {
	case <-a.peerDiscoveryCtx.Done():
		ctx, cancel := context.WithCancel(context.Background())
		a.peerDiscoveryCtx = ctx
		a.peerDiscoveryCancelFunc = cancel
		return ctx
	default:
		return a.peerDiscoveryCtx
	}
}

// startStatusMonitoring monitors service status changes and updates system tray
func (a *App) startStatusMonitoring() {
	if a.serviceManager == nil {
		a.statusMonitorRunning = false
		return
	}

	statusChan := a.serviceManager.GetStatusChannel()
	if statusChan == nil {
		a.statusMonitorRunning = false
		return
	}

	for {
		select {
		case <-a.statusMonitorShutdown:
			a.statusMonitorRunning = false
			return

		case _, ok := <-statusChan:
			if !ok {
				a.statusMonitorRunning = false
				return
			}
			a.UpdateSystemTrayStatus()
		}
	}
}

// cancelPeerDiscoveryOperations cancels all ongoing peer discovery operations
func (a *App) cancelPeerDiscoveryOperations() {
	if a.peerDiscoveryCancelFunc != nil {
		a.peerDiscoveryCancelFunc()
	}
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return version
}

// ==================== System Tray Methods ====================

// setupTray initializes the system tray
func (a *App) setupTray() {
	a.trayManager = tray.NewManager(
		a.ctx,
		a.config,
		a.serviceManager,
		a.showFromTray,
		a.showSettingsFromTray,
		a.quitFromTray,
	)
	a.trayManager.Setup()
}

// UpdateSystemTrayMenu updates the system tray menu
func (a *App) UpdateSystemTrayMenu() {
	if a.trayManager != nil {
		a.trayManager.UpdateMenu()
	}
}

// UpdateSystemTrayStatus updates the system tray status
func (a *App) UpdateSystemTrayStatus() {
	if a.trayManager != nil {
		a.trayManager.UpdateStatus()
	}
}

// showFromTray shows the window from system tray
func (a *App) showFromTray() {
	tray.ShowWindow(a.ctx)
}

// showSettingsFromTray shows settings from system tray
func (a *App) showSettingsFromTray() {
	tray.ShowSettingsWindow(a.ctx)
}

// quitFromTray quits the application from system tray
func (a *App) quitFromTray() {
	a.allowQuit = true
	tray.QuitApplication(a.ctx, a.actualShutdown)
}

// ==================== Event Monitoring ====================

// startEventMonitoring monitors backend events and forwards to frontend
func (a *App) startEventMonitoring() {
	events.StartEventMonitoring(
		a.serviceManager,
		func(eventName string, data interface{}) {
			if a.ctx != nil {
				runtime.EventsEmit(a.ctx, eventName, data)
			}
		},
		a.UpdateSystemTrayStatus,
		a.eventMonitorShutdown,
	)
	a.eventMonitorRunning = false
}

// ==================== Configuration Bindings ====================

// GetConfig returns the current application configuration
func (a *App) GetConfig() ConfigDTO {
	return config.GetConfigDTO(a.config)
}

// SaveConfig saves the entire configuration
func (a *App) SaveConfig(dto ConfigDTO) error {
	if err := config.UpdateConfigFromDTO(a.config, dto); err != nil {
		return err
	}
	a.UpdateSystemTrayMenu()
	return a.config.Save()
}

// AddPeer adds a new peer to the configuration
func (a *App) AddPeer(address string) error {
	if err := config.AddPeer(a.config, address); err != nil {
		return err
	}
	return a.config.Save()
}

// RemovePeer removes a peer from the configuration
func (a *App) RemovePeer(address string) error {
	if err := config.RemovePeer(a.config, address); err != nil {
		return err
	}
	return a.config.Save()
}

// EnablePeer enables a peer in the configuration
func (a *App) EnablePeer(address string) error {
	if err := config.EnablePeer(a.config, address); err != nil {
		return err
	}
	return a.config.Save()
}

// DisablePeer disables a peer in the configuration
func (a *App) DisablePeer(address string) error {
	if err := config.DisablePeer(a.config, address); err != nil {
		return err
	}
	return a.config.Save()
}

// SetPassword sets the yggmail password
func (a *App) SetPassword(password string) error {
	return config.SetPassword(a.config, a.serviceManager, password)
}

// ChangePassword changes the password after verifying the current password
func (a *App) ChangePassword(currentPassword, newPassword string) error {
	return config.ChangePassword(a.config, a.serviceManager, currentPassword, newPassword)
}

// RegenerateKeys regenerates Yggdrasil keys (WARNING: deletes all mail data)
func (a *App) RegenerateKeys(password string) error {
	return config.RegenerateKeys(a.config, a.serviceManager, password)
}

// SetLanguage sets the UI language
func (a *App) SetLanguage(language string) error {
	if err := config.SetLanguage(a.config, language); err != nil {
		return err
	}
	a.UpdateSystemTrayMenu()
	return a.config.Save()
}

// SetTheme sets the UI theme
func (a *App) SetTheme(theme string) error {
	if err := config.SetTheme(a.config, theme); err != nil {
		return err
	}
	return a.config.Save()
}

// SetAutoStart sets whether the service should start on system boot
func (a *App) SetAutoStart(enabled bool) error {
	if err := config.SetAutoStart(a.config, enabled); err != nil {
		return err
	}
	return a.config.Save()
}

// GetDefaultPeers returns a list of recommended default peers
func (a *App) GetDefaultPeers() []string {
	return config.GetDefaultPeers()
}

// ==================== Service Bindings ====================

// InitializeService initializes the yggmail service
func (a *App) InitializeService() error {
	return service.InitializeService(a.serviceManager)
}

// StartService starts the yggmail service
func (a *App) StartService() error {
	shouldStartMonitoring, err := service.StartService(a.serviceManager)
	if err != nil {
		return err
	}

	if shouldStartMonitoring && !a.eventMonitorRunning {
		a.eventMonitorRunning = true
		go a.startEventMonitoring()
	}

	a.UpdateSystemTrayStatus()
	return nil
}

// StopService stops the yggmail service
func (a *App) StopService() error {
	a.cancelPeerDiscoveryOperations()

	if err := service.StopService(a.serviceManager); err != nil {
		return err
	}

	a.UpdateSystemTrayStatus()
	return nil
}

// RestartService restarts the yggmail service
func (a *App) RestartService() error {
	a.cancelPeerDiscoveryOperations()

	if err := service.RestartService(a.serviceManager); err != nil {
		return err
	}

	a.UpdateSystemTrayStatus()
	return nil
}

// GetServiceStatus returns the current service status
func (a *App) GetServiceStatus() ServiceStatusDTO {
	return service.GetServiceStatusDTO(a.serviceManager, a.config)
}

// GetPeerStats returns statistics for all configured peers
func (a *App) GetPeerStats() []PeerInfoDTO {
	return service.GetPeerStatsDTO(a.serviceManager, a.config)
}

// HotReloadPeers reloads the peer list without stopping the service
func (a *App) HotReloadPeers() error {
	return service.HotReloadPeers(a.serviceManager, a.config)
}

// GetMailAddress returns the current yggmail address
func (a *App) GetMailAddress() string {
	return service.GetMailAddress(a.serviceManager)
}

// IsServiceRunning returns whether the service is currently running
func (a *App) IsServiceRunning() bool {
	return service.IsServiceRunning(a.serviceManager)
}

// ==================== System Bindings ====================

// ShowOpenFileDialog shows a file open dialog and returns the selected file path
func (a *App) ShowOpenFileDialog(title string) (string, error) {
	return system.ShowOpenFileDialog(a.ctx, title)
}

// ShowSaveFileDialog shows a file save dialog and returns the selected file path
func (a *App) ShowSaveFileDialog(title string, defaultFilename string) (string, error) {
	return system.ShowSaveFileDialog(a.ctx, title, defaultFilename)
}

// ShowOpenDirectoryDialog shows a directory selection dialog
func (a *App) ShowOpenDirectoryDialog(title string) (string, error) {
	return system.ShowOpenDirectoryDialog(a.ctx, title)
}

// ShowMessageDialog shows a message dialog with OK button
func (a *App) ShowMessageDialog(title string, message string) error {
	return system.ShowMessageDialog(a.ctx, title, message)
}

// ShowErrorDialog shows an error dialog
func (a *App) ShowErrorDialog(title string, message string) error {
	return system.ShowErrorDialog(a.ctx, title, message)
}

// ShowQuestionDialog shows a question dialog with Yes/No buttons
func (a *App) ShowQuestionDialog(title string, message string) (bool, error) {
	return system.ShowQuestionDialog(a.ctx, title, message)
}

// CopyToClipboard copies text to the system clipboard
func (a *App) CopyToClipboard(text string) error {
	return system.CopyToClipboard(a.ctx, text)
}

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) error {
	return system.OpenURL(a.ctx, url)
}

// CreateBackup creates an encrypted backup of the configuration and optionally database
func (a *App) CreateBackup(options BackupOptionsDTO) (ResultDTO, error) {
	return system.CreateBackup(a.ctx, a.config, options)
}

// RestoreBackup restores configuration and optionally database from an encrypted backup
func (a *App) RestoreBackup(options RestoreOptionsDTO) (ResultDTO, error) {
	restoredConfig, result, err := system.RestoreBackup(a.ctx, options)
	if err != nil {
		return result, err
	}

	// Update app's config reference if restore was successful
	if restoredConfig != nil {
		a.config = restoredConfig
	}

	return result, nil
}

// QuitApplication gracefully quits the application
func (a *App) QuitApplication() {
	system.QuitApplication(a.ctx)
}

// HideWindow hides the application window
func (a *App) HideWindow() {
	system.HideWindow(a.ctx)
}

// ShowWindow shows the application window
func (a *App) ShowWindow() {
	system.ShowWindow(a.ctx)
}

// MinimizeWindow minimizes the application window
func (a *App) MinimizeWindow() {
	system.MinimizeWindow(a.ctx)
}

// MaximizeWindow maximizes the application window
func (a *App) MaximizeWindow() {
	system.MaximizeWindow(a.ctx)
}

// ToggleMaximize toggles between maximized and normal window state
func (a *App) ToggleMaximize() {
	system.ToggleMaximize(a.ctx)
}

// GetSystemTheme returns the current system theme preference
func (a *App) GetSystemTheme() (string, error) {
	return system.GetSystemTheme()
}

// GetSystemLanguage returns the system's default language code
func (a *App) GetSystemLanguage() string {
	return system.GetSystemLanguage()
}

// OpenDeltaChat opens DeltaChat with auto-configured account
func (a *App) OpenDeltaChat() error {
	return system.OpenDeltaChat(a.ctx, a.config, a.serviceManager)
}

// ==================== Peer Discovery Bindings ====================

// FindAvailablePeers discovers available Yggdrasil peers
func (a *App) FindAvailablePeers(protocols string, region string, maxRTTMs int) (*core.PeerDiscoveryResult, error) {
	return peerdiscovery.FindAvailablePeers(
		a.ctx,
		a.getPeerDiscoveryContext(),
		a.config,
		protocols,
		region,
		maxRTTMs,
	)
}

// GetCachedDiscoveredPeers returns cached discovered peers
func (a *App) GetCachedDiscoveredPeers() []core.DiscoveredPeer {
	return peerdiscovery.GetCachedDiscoveredPeers(a.config)
}

// ClearCachedDiscoveredPeers clears the cached discovered peers
func (a *App) ClearCachedDiscoveredPeers() error {
	return peerdiscovery.ClearCachedDiscoveredPeers(a.config)
}

// GetAvailableRegions returns a list of all available peer regions
func (a *App) GetAvailableRegions() ([]string, error) {
	return peerdiscovery.GetAvailableRegions(a.getPeerDiscoveryContext())
}

// CheckCustomPeers checks a list of user-provided peer URIs
func (a *App) CheckCustomPeers(peerURIs []string) ([]core.DiscoveredPeer, error) {
	return peerdiscovery.CheckCustomPeers(a.getPeerDiscoveryContext(), peerURIs)
}

// AddDiscoveredPeer adds a discovered peer to the configuration
func (a *App) AddDiscoveredPeer(peer core.DiscoveredPeer) error {
	return peerdiscovery.AddDiscoveredPeer(a.config, peer)
}

// AddDiscoveredPeers adds multiple discovered peers to the configuration
func (a *App) AddDiscoveredPeers(peers []core.DiscoveredPeer) error {
	return peerdiscovery.AddDiscoveredPeers(a.config, peers)
}

// GetPeerDiscoverySystemInfo returns system information for debugging
func (a *App) GetPeerDiscoverySystemInfo() map[string]interface{} {
	return peerdiscovery.GetPeerDiscoverySystemInfo()
}
