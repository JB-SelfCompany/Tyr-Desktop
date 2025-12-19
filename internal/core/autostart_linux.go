//go:build !windows

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Desktop file name for autostart
	desktopFileName = "tyr.desktop"
)

// Follows XDG Base Directory specification: ~/.config/autostart/
func getAutostartDir() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Check for XDG_CONFIG_HOME environment variable
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		// Use default: ~/.config
		configHome = filepath.Join(homeDir, ".config")
	}

	autostartDir := filepath.Join(configHome, "autostart")
	return autostartDir, nil
}

// getDesktopFilePath returns the full path to the .desktop file
func getDesktopFilePath() (string, error) {
	autostartDir, err := getAutostartDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(autostartDir, desktopFileName), nil
}

// createDesktopFileContent generates the content for the .desktop file
// Follows freedesktop.org Desktop Entry Specification
func createDesktopFileContent(execPath string) (string, error) {
	// Validate executable path for security
	if err := validateExecutablePath(execPath); err != nil {
		return "", fmt.Errorf("invalid executable path: %w", err)
	}

	// Sanitize exec path - escape any special characters
	// This prevents command injection attacks
	sanitizedExec := sanitizeDesktopExecPath(execPath)

	// Add --minimized flag to start in system tray
	execWithArgs := fmt.Sprintf("%s --minimized", sanitizedExec)

	// Create .desktop file content
	// Format follows freedesktop.org specification
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Comment=Secure P2P Email over Yggdrasil Network
Exec=%s
Icon=mail-client
Terminal=false
Hidden=false
X-GNOME-Autostart-enabled=true
StartupNotify=false
Categories=Network;Email;
`,
		AutostartAppName,
		execWithArgs,
	)

	return content, nil
}

// sanitizeDesktopExecPath sanitizes the executable path for use in .desktop file
// Escapes special characters that could be interpreted as shell commands
func sanitizeDesktopExecPath(path string) string {
	// Characters that need escaping in .desktop Exec field
	// Based on freedesktop.org specification
	replacer := strings.NewReplacer(
		`\`, `\\`,   // Backslash
		`"`, `\"`,   // Double quote
		"`", "\\`",  // Backtick
		"$", "\\$",  // Dollar sign (variable expansion)
	)

	return replacer.Replace(path)
}

// enableAutoStartLinux enables autostart on Linux by creating a .desktop file
// Location: ~/.config/autostart/tyr.desktop
func enableAutoStartLinux() error {
	// Get executable path
	exePath, err := getExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get autostart directory
	autostartDir, err := getAutostartDir()
	if err != nil {
		return fmt.Errorf("failed to get autostart directory: %w", err)
	}

	// Create autostart directory if it doesn't exist
	// Use 0700 permissions (user-only access) for security
	if err := os.MkdirAll(autostartDir, 0700); err != nil {
		return fmt.Errorf("failed to create autostart directory: %w", err)
	}

	// Generate .desktop file content
	content, err := createDesktopFileContent(exePath)
	if err != nil {
		return fmt.Errorf("failed to create desktop file content: %w", err)
	}

	// Get desktop file path
	desktopFilePath, err := getDesktopFilePath()
	if err != nil {
		return fmt.Errorf("failed to get desktop file path: %w", err)
	}

	// Write .desktop file
	// Use 0644 permissions (readable by all, writable by owner)
	if err := os.WriteFile(desktopFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	return nil
}

// disableAutoStartLinux disables autostart on Linux by removing the .desktop file
func disableAutoStartLinux() error {
	// Get desktop file path
	desktopFilePath, err := getDesktopFilePath()
	if err != nil {
		return fmt.Errorf("failed to get desktop file path: %w", err)
	}

	// Remove the .desktop file
	// os.Remove returns nil if file doesn't exist
	if err := os.Remove(desktopFilePath); err != nil {
		// Check if error is "file not found" - this is not an error for disable
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}

	return nil
}

// isAutoStartEnabledLinux checks if autostart is enabled on Linux
// Returns true if the .desktop file exists
func isAutoStartEnabledLinux() (bool, error) {
	// Get desktop file path
	desktopFilePath, err := getDesktopFilePath()
	if err != nil {
		return false, fmt.Errorf("failed to get desktop file path: %w", err)
	}

	// Check if file exists
	_, err = os.Stat(desktopFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check desktop file: %w", err)
	}

	return true, nil
}

// Stub implementations for Windows functions (not used on Linux)
// These are needed for compilation when building on Linux

func enableAutoStartWindows() error {
	return fmt.Errorf("Windows autostart not available on Linux")
}

func disableAutoStartWindows() error {
	return fmt.Errorf("Windows autostart not available on Linux")
}

func isAutoStartEnabledWindows() (bool, error) {
	return false, fmt.Errorf("Windows autostart not available on Linux")
}
