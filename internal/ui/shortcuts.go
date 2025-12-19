package ui

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// SetupKeyboardShortcuts configures global keyboard shortcuts for the application
// Platform-aware: Uses Ctrl on Windows/Linux, Cmd on macOS
//
// Shortcuts:
//   - Ctrl/Cmd+Q: Quit application
//   - Ctrl/Cmd+S: Open Settings
//   - Ctrl/Cmd+P: Open Peers management
//   - Ctrl/Cmd+B: Open Backup screen
//   - Ctrl/Cmd+,: Open Settings (alternative)
//   - F5: Refresh dashboard
//   - Escape: Close current dialog (handled by Fyne automatically)
func (a *App) SetupKeyboardShortcuts() {
	window := a.mainWindow
	canvas := window.Canvas()

	// Determine the appropriate modifier key based on platform
	// macOS uses Command (Cmd), Windows/Linux use Control (Ctrl)
	modifier := fyne.KeyModifierControl
	if runtime.GOOS == "darwin" {
		modifier = fyne.KeyModifierSuper // Super = Cmd on macOS
	}

	// Ctrl/Cmd+Q - Quit application
	quitShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: modifier,
	}
	canvas.AddShortcut(quitShortcut, func(shortcut fyne.Shortcut) {
		a.Quit()
	})

	// Ctrl/Cmd+S - Open Settings
	settingsShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: modifier,
	}
	canvas.AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		a.ShowSettings()
	})

	// Ctrl/Cmd+P - Open Peers management
	peersShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: modifier,
	}
	canvas.AddShortcut(peersShortcut, func(shortcut fyne.Shortcut) {
		a.ShowPeers()
	})

	// Ctrl/Cmd+B - Open Backup screen
	backupShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyB,
		Modifier: modifier,
	}
	canvas.AddShortcut(backupShortcut, func(shortcut fyne.Shortcut) {
		a.ShowBackup()
	})

	// Ctrl/Cmd+, - Open Settings (alternative, common on macOS)
	settingsAltShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyComma,
		Modifier: modifier,
	}
	canvas.AddShortcut(settingsAltShortcut, func(shortcut fyne.Shortcut) {
		a.ShowSettings()
	})

	// F5 - Refresh dashboard
	refreshShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyF5,
		Modifier: 0, // No modifier key
	}
	canvas.AddShortcut(refreshShortcut, func(shortcut fyne.Shortcut) {
		// Only refresh if we're on the dashboard screen
		if a.currentScreen == ScreenDashboard {
			a.ShowDashboard()
		}
	})

	// Ctrl/Cmd+R - Refresh (alternative)
	refreshAltShortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyR,
		Modifier: modifier,
	}
	canvas.AddShortcut(refreshAltShortcut, func(shortcut fyne.Shortcut) {
		// Refresh current screen by re-showing it
		switch a.currentScreen {
		case ScreenDashboard:
			a.ShowDashboard()
		case ScreenSettings:
			a.ShowSettings()
		case ScreenPeers:
			a.ShowPeers()
		case ScreenBackup:
			a.ShowBackup()
		}
	})

	// Note: Escape key for closing dialogs is handled automatically by Fyne
	// No need to implement it manually
}
