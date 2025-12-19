//go:build windows

package tray

import (
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/resources"
)

// GetIconData returns the platform-specific tray icon data
// On Windows, systray requires ICO format (not PNG)
func GetIconData() []byte {
	// Use ICO file - required by Windows systray API
	return resources.ResourceTyrIco.StaticContent
}
