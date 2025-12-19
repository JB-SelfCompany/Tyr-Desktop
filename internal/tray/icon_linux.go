//go:build linux

package tray

// GetIconData returns the platform-specific tray icon data
// On Linux, system tray also works better with smaller icons (32x32 pixels)
func GetIconData() []byte {
	return ResizeIcon()
}
