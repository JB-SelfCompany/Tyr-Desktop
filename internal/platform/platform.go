// Package platform provides platform-specific utilities and portable data storage.
// All application data is stored in a "data" subdirectory next to the executable,
// making the application fully portable.
package platform

import (
	"os"
	"path/filepath"
	"runtime"
)

// OS represents the operating system type
type OS string

const (
	Windows OS = "windows"
	Darwin  OS = "darwin"
	Linux   OS = "linux"
	Unknown OS = "unknown"
)

// Info contains information about the current platform
type Info struct {
	OS       OS
	Arch     string
	HomeDir  string
	DataDir  string
	CacheDir string
	LogDir   string
}

// Current returns information about the current platform
func Current() *Info {
	info := &Info{
		OS:      GetOS(),
		Arch:    runtime.GOARCH,
		HomeDir: getHomeDir(),
	}

	info.DataDir = GetDataDir()
	info.CacheDir = GetCacheDir()
	info.LogDir = GetLogDir()

	return info
}

// GetOS returns the current operating system
func GetOS() OS {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	default:
		return Unknown
	}
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsDarwin returns true if running on macOS
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// GetExecutableDir returns the directory where the executable is located.
// Used for portable mode - all config/data files are stored next to the executable.
func GetExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to current working directory
		cwd, _ := os.Getwd()
		return cwd
	}
	// Resolve symlinks to get the real path
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return filepath.Dir(exePath)
	}
	return filepath.Dir(realPath)
}

// GetDataDir returns the data directory for the application.
// PORTABLE MODE: All files are stored in "data" subdirectory next to the executable.
// This makes the application fully portable - just copy the folder to another location.
func GetDataDir() string {
	return filepath.Join(GetExecutableDir(), "data")
}

// GetCacheDir returns the cache directory for the application.
// PORTABLE MODE: cache is stored in "data/cache" subdirectory next to the executable.
func GetCacheDir() string {
	return filepath.Join(GetDataDir(), "cache")
}

// GetLogDir returns the log directory for the application.
// PORTABLE MODE: logs are stored in "data/logs" subdirectory next to the executable.
func GetLogDir() string {
	return filepath.Join(GetDataDir(), "logs")
}

// EnsureDirectories creates all necessary application directories
func EnsureDirectories() error {
	dirs := []string{
		GetDataDir(),
		GetCacheDir(),
		GetLogDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	return filepath.Join(GetDataDir(), "config.toml")
}

// GetDatabasePath returns the path to the yggmail database file
func GetDatabasePath() string {
	return filepath.Join(GetDataDir(), "yggmail.db")
}

// GetLegacyConfigDirs returns the list of legacy configuration directories
// that may contain data from previous non-portable versions.
// Windows: %APPDATA%\Tyr
// Linux: ~/.config/tyr
func GetLegacyConfigDirs() []string {
	var dirs []string

	switch GetOS() {
	case Windows:
		// Windows: %APPDATA%\Tyr
		appData := os.Getenv("APPDATA")
		if appData != "" {
			dirs = append(dirs, filepath.Join(appData, "Tyr"))
		}
	case Linux, Darwin:
		// Linux/macOS: ~/.config/tyr
		home := getHomeDir()
		if home != "" {
			dirs = append(dirs, filepath.Join(home, ".config", "tyr"))
		}
	}

	return dirs
}
