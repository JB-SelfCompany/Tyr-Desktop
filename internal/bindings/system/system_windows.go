//go:build windows

package system

import (
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultLocaleName = kernel32.NewProc("GetUserDefaultLocaleName")
)

// detectSystemLanguageImpl detects the system language on Windows
// Uses Windows GetUserDefaultLocaleName API to get the user's locale
// Returns "en" or "ru" based on the locale, defaults to "en"
func detectSystemLanguageImpl() string {
	// Get user default locale name
	// This returns locale names like "en-US", "ru-RU", etc.
	var localeName [85]uint16 // LOCALE_NAME_MAX_LENGTH is 85
	ret, _, _ := procGetUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&localeName[0])),
		uintptr(len(localeName)),
	)

	if ret == 0 {
		// Failed to get locale, default to English
		return "en"
	}

	// Convert UTF-16 to string
	localeStr := syscall.UTF16ToString(localeName[:])

	// Parse language code from locale
	langCode := parseLangCode(localeStr)

	// Validate and return (only "en" or "ru" are supported)
	if langCode == "en" || langCode == "ru" {
		return langCode
	}

	return "en"
}

// parseLangCode extracts the language code from a locale string
// Examples: "en-US" -> "en", "ru-RU" -> "ru", "en" -> "en"
func parseLangCode(locale string) string {
	if locale == "" {
		return ""
	}

	locale = strings.ToLower(locale)

	// Split by underscore, hyphen, or period to extract language code
	for _, sep := range []rune{'_', '.', '-'} {
		if idx := strings.IndexRune(locale, sep); idx > 0 {
			return locale[:idx]
		}
	}

	return locale
}
