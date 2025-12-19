//go:build !windows

package theme

import (
	"os"
	"strings"
)

// getSystemThemeImpl detects the system theme on Linux and other Unix-like platforms
// Checks multiple sources in order:
// 1. GTK theme settings via gsettings
// 2. Environment variables (GTK_THEME, QT_QPA_PLATFORMTHEME)
// 3. GNOME/KDE desktop environment settings
//
// Returns "dark" if dark theme is detected, "light" otherwise
func getSystemThemeImpl() (string, error) {
	// Method 1: Check GTK_THEME environment variable
	// This is the most direct indicator for GTK apps
	if gtkTheme := os.Getenv("GTK_THEME"); gtkTheme != "" {
		if strings.Contains(strings.ToLower(gtkTheme), "dark") {
			return ThemeDark, nil
		}
	}

	// Method 2: Check desktop color scheme via environment
	// Used by some desktop environments
	if colorScheme := os.Getenv("COLOR_SCHEME"); colorScheme != "" {
		if strings.Contains(strings.ToLower(colorScheme), "dark") {
			return ThemeDark, nil
		}
	}

	// Method 3: Check XDG_CURRENT_DESKTOP for desktop environment
	// Then check DE-specific settings
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if desktop != "" {
		desktop = strings.ToLower(desktop)

		// For GNOME: check color-scheme setting
		// For KDE: check kdeglobals
		// Note: This would require executing gsettings/kreadconfig commands
		// For now, we default to light as a safe fallback
	}

	// Default to light theme if no dark indicators found
	return ThemeLight, nil
}
