package i18n

import (
	"fmt"
	"strings"
	"sync"
)

// Supported language codes
const (
	LangEnglish = "en"
	LangRussian = "ru"
)

// ValidLanguages is a whitelist of supported language codes
var ValidLanguages = []string{LangEnglish, LangRussian}

// Localizer manages localized string translations
// Thread-safe for concurrent access
type Localizer struct {
	// Current language code
	currentLang string

	// Strings map for the current language
	strings map[string]string

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// Global localizer instance
var globalLocalizer *Localizer
var localizerOnce sync.Once

// GetGlobalLocalizer returns the global localizer instance
// Creates it on first access with system default language
func GetGlobalLocalizer() *Localizer {
	localizerOnce.Do(func() {
		lang := DetectSystemLanguage()
		globalLocalizer = NewLocalizer(lang)
	})
	return globalLocalizer
}

// NewLocalizer creates a new localizer with the specified language
// Falls back to English if the language is not supported
//
// Security: Validates language code against whitelist
func NewLocalizer(langCode string) *Localizer {
	// Validate language code
	if !isValidLanguage(langCode) {
		langCode = LangEnglish
	}

	l := &Localizer{
		currentLang: langCode,
	}

	// Load strings for the language
	l.loadStrings(langCode)

	return l
}

// Get retrieves a localized string by key
// Returns the key itself if translation is not found (fallback behavior)
//
// Thread-safe with read lock
func (l *Localizer) Get(key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Lookup translation
	if value, ok := l.strings[key]; ok {
		return value
	}

	// Fallback: return key if not found
	return key
}

// SetLanguage changes the current language
// Reloads all strings for the new language
//
// Security: Validates language code against whitelist
// Thread-safe with write lock
func (l *Localizer) SetLanguage(langCode string) error {
	// Validate language code
	if !isValidLanguage(langCode) {
		return fmt.Errorf("invalid language code: %s (must be en or ru)", langCode)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.currentLang = langCode
	l.loadStrings(langCode)

	return nil
}

// GetCurrentLanguage returns the current language code
// Thread-safe with read lock
func (l *Localizer) GetCurrentLanguage() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.currentLang
}

// loadStrings loads the string map for the specified language
// Must be called with write lock held
func (l *Localizer) loadStrings(langCode string) {
	switch langCode {
	case LangEnglish:
		l.strings = GetEnglishStrings()
	case LangRussian:
		l.strings = GetRussianStrings()
	default:
		// Fallback to English
		l.strings = GetEnglishStrings()
	}
}

// isValidLanguage checks if the language code is supported
// Security: Prevents invalid language codes from being processed
func isValidLanguage(langCode string) bool {
	for _, valid := range ValidLanguages {
		if langCode == valid {
			return true
		}
	}
	return false
}

// DetectSystemLanguage attempts to detect the system's default language
// Returns language code (e.g., "en", "ru")
//
// Detection methods:
//   - Linux/Unix: Reads LANG, LC_ALL, or LANGUAGE environment variables
//   - Windows: Uses GetUserDefaultLocaleName API
//
// Defaults to English if detection fails
func DetectSystemLanguage() string {
	return detectSystemLanguageImpl()
}

// parseLangCode extracts the language code from a locale string
// Examples: "en_US.UTF-8" -> "en", "ru_RU" -> "ru", "en" -> "en"
func parseLangCode(locale string) string {
	// Handle empty string
	if locale == "" {
		return ""
	}

	// Convert to lowercase for consistency
	locale = strings.ToLower(locale)

	// Split by underscore or period to extract language code
	parts := strings.FieldsFunc(locale, func(r rune) bool {
		return r == '_' || r == '.' || r == '-'
	})

	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
