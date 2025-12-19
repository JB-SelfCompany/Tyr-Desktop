//go:build windows

package theme

import (
	"golang.org/x/sys/windows/registry"
)

// getSystemThemeImpl reads the Windows registry to determine system theme
// Reads HKCU\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize\AppsUseLightTheme
//
// Security: Uses registry.OpenKey with READ access only
// Returns "dark" if AppsUseLightTheme is 0, "light" otherwise or on error
func getSystemThemeImpl() (string, error) {
	// Open registry key with read-only access
	// Security: KEY_READ ensures we cannot modify registry
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		// Registry key not found or access denied - default to light
		return ThemeLight, err
	}
	defer key.Close()

	// Read the AppsUseLightTheme value
	// 0 = dark theme, 1 = light theme
	value, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		// Value not found - default to light
		return ThemeLight, err
	}

	// Interpret the value
	if value == 0 {
		return ThemeDark, nil
	}
	return ThemeLight, nil
}
