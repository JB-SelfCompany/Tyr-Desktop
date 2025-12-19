//go:build windows

package core

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const (
	// Windows registry key path for Run entries
	registryRunKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
)

// enableAutoStartWindows enables autostart on Windows by creating a registry entry
// Registry location: HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
// Value name: Tyr
// Value data: Full path to executable with --minimized flag
func enableAutoStartWindows() error {
	// Get executable path
	exePath, err := getExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Validate executable path for security
	if err := validateExecutablePath(exePath); err != nil {
		return fmt.Errorf("invalid executable path: %w", err)
	}

	// Add --minimized flag to start in system tray
	// Quote path to handle spaces in path
	autostartCommand := fmt.Sprintf("\"%s\" --minimized", exePath)

	// Open registry key for writing
	// HKEY_CURRENT_USER is used for per-user autostart (no admin rights required)
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		registryRunKeyPath,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Set string value with executable path and --minimized flag
	// This creates or updates the Tyr entry
	if err := key.SetStringValue(AutostartAppName, autostartCommand); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

// disableAutoStartWindows disables autostart on Windows by removing the registry entry
func disableAutoStartWindows() error {
	// Open registry key for writing
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		registryRunKeyPath,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Delete the registry value
	// Returns no error if value doesn't exist
	if err := key.DeleteValue(AutostartAppName); err != nil {
		// Check if error is "value not found" - this is not an error for disable
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

// isAutoStartEnabledWindows checks if autostart is enabled on Windows
// Returns true if the registry entry exists
func isAutoStartEnabledWindows() (bool, error) {
	// Open registry key for reading
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		registryRunKeyPath,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Try to read the value
	value, _, err := key.GetStringValue(AutostartAppName)
	if err != nil {
		// If value doesn't exist, autostart is disabled
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("failed to read registry value: %w", err)
	}

	// Value exists - check if it's non-empty
	return value != "", nil
}

// Stub implementations for Linux functions (not used on Windows)
// These are needed for compilation when building on Windows

func enableAutoStartLinux() error {
	return fmt.Errorf("Linux autostart not available on Windows")
}

func disableAutoStartLinux() error {
	return fmt.Errorf("Linux autostart not available on Windows")
}

func isAutoStartEnabledLinux() (bool, error) {
	return false, fmt.Errorf("Linux autostart not available on Windows")
}
