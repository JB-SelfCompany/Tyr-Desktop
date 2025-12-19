//go:build windows

package i18n

import (
	"syscall"
	"unsafe"
)

const (
	LOCALE_NAME_MAX_LENGTH = 85
)

var (
	kernel32                  = syscall.NewLazyDLL("kernel32.dll")
	getUserDefaultLocaleName  = kernel32.NewProc("GetUserDefaultLocaleName")
)

// detectSystemLanguageImpl detects the system language on Windows
// Uses Windows GetUserDefaultLocaleName API to get the user's locale
// Returns "en" or "ru" based on the locale, defaults to "en"
func detectSystemLanguageImpl() string {
	// Get user default locale name
	// This returns locale names like "en-US", "ru-RU", etc.
	var localeName [LOCALE_NAME_MAX_LENGTH]uint16

	ret, _, _ := getUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&localeName[0])),
		uintptr(LOCALE_NAME_MAX_LENGTH),
	)

	if ret == 0 {
		// Failed to get locale, default to English
		return LangEnglish
	}

	// Convert UTF-16 to string
	localeStr := syscall.UTF16ToString(localeName[:])

	// Parse language code from locale
	langCode := parseLangCode(localeStr)

	// Validate and return
	if isValidLanguage(langCode) {
		return langCode
	}

	return LangEnglish
}
