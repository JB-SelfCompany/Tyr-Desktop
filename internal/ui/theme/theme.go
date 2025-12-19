package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Theme variant constants
const (
	ThemeLight = "light"
	ThemeDark  = "dark"
)

type ModernTheme struct {
	variant fyne.ThemeVariant
}

// NewModernTheme creates a new modern theme instance
func NewModernTheme(variant fyne.ThemeVariant) fyne.Theme {
	return &ModernTheme{variant: variant}
}

// Color returns theme colors
func (mt *ModernTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Use the theme's variant if specified, otherwise use provided variant
	v := mt.variant
	if v == 0 {
		v = variant
	}

	// Modern color palette
	switch name {
	case theme.ColorNamePrimary:
		if v == theme.VariantLight {
			// Vibrant Indigo
			return color.RGBA{R: 99, G: 102, B: 241, A: 255}
		}
		// Brighter Indigo for dark mode
		return color.RGBA{R: 129, G: 140, B: 248, A: 255}

	case theme.ColorNameForeground:
		if v == theme.VariantLight {
			// Near-black with slight warmth
			return color.RGBA{R: 31, G: 41, B: 55, A: 255}
		}
		// Soft white for dark mode
		return color.RGBA{R: 243, G: 244, B: 246, A: 255}

	case theme.ColorNameBackground:
		if v == theme.VariantLight {
			// Pure white with slight warmth
			return color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		// Deep dark with slight blue tint (2025 dark mode trend)
		return color.RGBA{R: 17, G: 24, B: 39, A: 255}

	case theme.ColorNameButton:
		if v == theme.VariantLight {
			// Vibrant Indigo
			return color.RGBA{R: 99, G: 102, B: 241, A: 255}
		}
		return color.RGBA{R: 129, G: 140, B: 248, A: 255}

	case theme.ColorNameHover:
		if v == theme.VariantLight {
			// Light Indigo
			return color.RGBA{R: 224, G: 231, B: 255, A: 255}
		}
		// Darker Indigo for dark mode
		return color.RGBA{R: 67, G: 56, B: 202, A: 255}

	case theme.ColorNameDisabled:
		if v == theme.VariantLight {
			return color.RGBA{R: 209, G: 213, B: 219, A: 255}
		}
		return color.RGBA{R: 75, G: 85, B: 99, A: 255}

	case theme.ColorNameSuccess:
		// Vibrant Green 
		if v == theme.VariantLight {
			return color.RGBA{R: 16, G: 185, B: 129, A: 255}
		}
		return color.RGBA{R: 52, G: 211, B: 153, A: 255}

	case theme.ColorNameWarning:
		// Bold Amber
		if v == theme.VariantLight {
			return color.RGBA{R: 245, G: 158, B: 11, A: 255}
		}
		return color.RGBA{R: 251, G: 191, B: 36, A: 255}

	case theme.ColorNameError:
		// Vibrant Red 
		if v == theme.VariantLight {
			return color.RGBA{R: 239, G: 68, B: 68, A: 255}
		}
		return color.RGBA{R: 248, G: 113, B: 113, A: 255}

	case theme.ColorNameShadow:
		if v == theme.VariantLight {
			return color.RGBA{R: 0, G: 0, B: 0, A: 30}
		}
		return color.RGBA{R: 0, G: 0, B: 0, A: 80}

	case theme.ColorNameInputBackground:
		if v == theme.VariantLight {
			return color.RGBA{R: 249, G: 250, B: 251, A: 255}
		}
		return color.RGBA{R: 31, G: 41, B: 55, A: 255}

	case theme.ColorNamePlaceHolder:
		if v == theme.VariantLight {
			return color.RGBA{R: 156, G: 163, B: 175, A: 255}
		}
		return color.RGBA{R: 107, G: 114, B: 128, A: 255}

	case theme.ColorNameScrollBar:
		if v == theme.VariantLight {
			return color.RGBA{R: 209, G: 213, B: 219, A: 255}
		}
		return color.RGBA{R: 75, G: 85, B: 99, A: 255}

	case theme.ColorNameSeparator:
		if v == theme.VariantLight {
			return color.RGBA{R: 229, G: 231, B: 235, A: 255}
		}
		return color.RGBA{R: 55, G: 65, B: 81, A: 255}

	default:
		// Fallback to default theme
		return theme.DefaultTheme().Color(name, v)
	}
}

// Font returns modern typography settings
func (mt *ModernTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Use system fonts for better cross-platform support
	// Fyne will use system defaults which are already modern
	return theme.DefaultTheme().Font(style)
}

// Icon returns theme icons
func (mt *ModernTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns spacing and sizing with modern proportions
func (mt *ModernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		// Increased padding for better spacing 
		return 8

	case theme.SizeNameInlineIcon:
		return 24

	case theme.SizeNameScrollBar:
		// Thinner scroll bars (modern trend)
		return 12

	case theme.SizeNameScrollBarSmall:
		return 8

	case theme.SizeNameSeparatorThickness:
		return 1

	case theme.SizeNameText:
		// Slightly larger base text (accessibility)
		return 15

	case theme.SizeNameHeadingText:
		// Bold heading sizes (modern typography)
		return 24

	case theme.SizeNameSubHeadingText:
		return 18

	case theme.SizeNameCaptionText:
		return 12

	case theme.SizeNameInputBorder:
		return 2

	default:
		return theme.DefaultTheme().Size(name)
	}
}

// ApplyModernTheme applies the modern theme to the app
func ApplyModernTheme(app fyne.App, variant fyne.ThemeVariant) {
	modernTheme := NewModernTheme(variant)
	app.Settings().SetTheme(modernTheme)
}

// GetSystemTheme returns the system theme ("light" or "dark")
// This is a wrapper around the platform-specific implementation
func GetSystemTheme() (string, error) {
	return getSystemThemeImpl()
}

// ApplyTheme applies a theme to the app based on the theme name
// Supports "light", "dark", "auto", or "system"
func ApplyTheme(app fyne.App, themeName string) error {
	var variant fyne.ThemeVariant

	switch themeName {
	case ThemeLight: // "light"
		variant = theme.VariantLight
	case ThemeDark: // "dark"
		variant = theme.VariantDark
	case "auto", "system":
		// Detect system theme
		systemTheme, err := getSystemThemeImpl()
		if err != nil {
			// Default to light on error
			variant = theme.VariantLight
		} else if systemTheme == ThemeDark {
			variant = theme.VariantDark
		} else {
			variant = theme.VariantLight
		}
	default:
		// Default to light theme
		variant = theme.VariantLight
	}

	ApplyModernTheme(app, variant)
	return nil
}
