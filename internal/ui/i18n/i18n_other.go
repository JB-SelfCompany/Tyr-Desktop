//go:build !windows

package i18n

import (
	"os"
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
			if isValidLanguage(langCode) {
				return langCode
			}
		}
	}

	// Default to English
	return LangEnglish
}
