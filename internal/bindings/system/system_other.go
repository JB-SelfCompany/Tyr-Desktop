//go:build !windows

package system

import (
	"os"
	"strings"
)

// detectSystemLanguageImpl detects the system language on Linux and other Unix-like platforms
// Checks environment variables: LANG, LC_ALL, LANGUAGE
// Returns "en" or "ru" based on the locale, defaults to "en"
func detectSystemLanguageImpl() string {
	// Check environment variables in order of preference
	langVars := []string{"LANG", "LC_ALL", "LANGUAGE"}

	for _, envVar := range langVars {
		langValue := os.Getenv(envVar)
		if langValue != "" {
			// Parse language code from locale string
			// Examples: "en_US.UTF-8" -> "en", "ru_RU" -> "ru"
			langCode := parseLangCode(langValue)
			if langCode == "en" || langCode == "ru" {
				return langCode
			}
		}
	}

	// Default to English
	return "en"
}

// parseLangCode extracts the language code from a locale string
// Examples: "en_US.UTF-8" -> "en", "ru_RU" -> "ru", "en" -> "en"
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
