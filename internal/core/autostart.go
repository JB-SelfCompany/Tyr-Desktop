package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// AutostartAppName is the application name used for autostart configuration
const AutostartAppName = "Tyr"

// EnableAutoStart enables the application to start automatically on system boot
// Platform-specific implementation:
//   - Windows: Creates registry key in HKCU\Software\Microsoft\Windows\CurrentVersion\Run
//   - Linux: Creates .desktop file in ~/.config/autostart/
// Returns error if operation fails or platform is unsupported
func EnableAutoStart() error {
	switch runtime.GOOS {
	case "windows":
		return enableAutoStartWindows()
	case "linux":
		return enableAutoStartLinux()
	default:
		return fmt.Errorf("autostart not supported on platform: %s", runtime.GOOS)
	}
}

// DisableAutoStart disables the application from starting automatically on system boot
// Platform-specific implementation:
//   - Windows: Removes registry key from HKCU\Software\Microsoft\Windows\CurrentVersion\Run
//   - Linux: Removes .desktop file from ~/.config/autostart/
// Returns error if operation fails or platform is unsupported
func DisableAutoStart() error {
	switch runtime.GOOS {
	case "windows":
		return disableAutoStartWindows()
	case "linux":
		return disableAutoStartLinux()
	default:
		return fmt.Errorf("autostart not supported on platform: %s", runtime.GOOS)
	}
}

// IsAutoStartEnabled checks if autostart is currently enabled
// Platform-specific implementation:
//   - Windows: Checks for registry key existence
//   - Linux: Checks for .desktop file existence
// Returns (enabled, error)
func IsAutoStartEnabled() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return isAutoStartEnabledWindows()
	case "linux":
		return isAutoStartEnabledLinux()
	default:
		return false, fmt.Errorf("autostart not supported on platform: %s", runtime.GOOS)
	}
}

// getExecutablePath returns the absolute path to the current executable
// Uses os.Executable() for security - never hardcoded paths
// Resolves symlinks to get the real path
func getExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks to get the real executable path
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		// If symlink resolution fails, use the original path
		return exePath, nil
	}

	return realPath, nil
}

// validateExecutablePath performs security validation on the executable path
// Prevents path traversal and injection attacks
func validateExecutablePath(path string) error {
	// Check for empty path
	if path == "" {
		return fmt.Errorf("executable path is empty")
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("executable path does not exist: %w", err)
	}

	// Check if it's a regular file (not directory or device)
	if !info.Mode().IsRegular() {
		return fmt.Errorf("executable path is not a regular file")
	}

	// Verify path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("executable path must be absolute")
	}

	return nil
}
