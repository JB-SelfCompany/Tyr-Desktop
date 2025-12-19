package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)

// SettingsHeaderConfig holds configuration for settings page headers
type SettingsHeaderConfig struct {
	Title    string
	Subtitle string
	OnBack   func()
}

// SettingsHeader represents a header with back and apply buttons
type SettingsHeader struct {
	backText          string
	onBack            func()
	onApply           func()
	hasUnsavedChanges *bool
	app               AppInterface
	container         *fyne.Container
}

// NewSettingsHeader creates a new settings header with back and apply buttons
func NewSettingsHeader(backText string, onBack func(), onApply func(), hasUnsavedChanges *bool, app AppInterface) *SettingsHeader {
	header := &SettingsHeader{
		backText:          backText,
		onBack:            onBack,
		onApply:           onApply,
		hasUnsavedChanges: hasUnsavedChanges,
		app:               app,
	}
	header.buildContainer()
	return header
}

// GetContainer returns the header's container
func (sh *SettingsHeader) GetContainer() fyne.CanvasObject {
	return sh.container
}

// buildContainer creates the header container
func (sh *SettingsHeader) buildContainer() {
	loc := i18n.GetGlobalLocalizer()

	// Back button
	backButton := widget.NewButton(sh.backText, sh.onBack)
	backButton.Importance = widget.LowImportance

	// Apply button
	applyButton := widget.NewButton(loc.Get("action.apply"), sh.onApply)
	applyButton.Importance = widget.HighImportance

	// Create header with buttons
	sh.container = container.NewBorder(
		nil, nil,
		backButton,
		applyButton,
		container.NewCenter(widget.NewLabel("")),
	)
}

// UpdateApplyButton updates the apply button state based on whether there are unsaved changes
func (sh *SettingsHeader) UpdateApplyButton(enabled bool) {
	// This is a placeholder implementation
	// In a real implementation, you would enable/disable the button based on the enabled parameter
	// For now, we'll just refresh the container
	sh.container.Refresh()
}

// buildSettingsHeaderModern creates a modern header with Back button above title
func buildSettingsHeaderModern(config SettingsHeaderConfig) fyne.CanvasObject {
	// Get localized strings
	loc := i18n.GetGlobalLocalizer()

	// Back button - compact with icon and text
	backButton := widget.NewButtonWithIcon(loc.Get("action.back"), theme.NavigateBackIcon(), config.OnBack)
	backButton.Importance = widget.LowImportance

	// Create a container for back button aligned to the left
	backButtonContainer := container.NewHBox(backButton)

	// Title and subtitle
	titleLabel := widget.NewRichTextFromMarkdown("# " + config.Title)
	subtitleLabel := widget.NewLabel(config.Subtitle)
	subtitleLabel.TextStyle = fyne.TextStyle{Italic: true}
	subtitleLabel.Wrapping = fyne.TextWrapWord

	// Vertical layout: Back button at top, then title, then subtitle
	header := container.NewVBox(
		backButtonContainer,
		titleLabel,
		subtitleLabel,
		widget.NewSeparator(),
	)

	return container.NewPadded(header)
}
