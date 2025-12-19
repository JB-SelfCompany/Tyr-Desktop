package dialogs

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)

// Modern 2025 color scheme with vibrant, eye-catching colors
var (
	// ErrorColor - Modern red with glassmorphic accent
	ErrorColor = color.RGBA{R: 239, G: 68, B: 68, A: 255}
	ErrorGradientColor = color.RGBA{R: 239, G: 68, B: 68, A: 180}

	// SuccessColor - Modern green with glassmorphic accent
	SuccessColor = color.RGBA{R: 16, G: 185, B: 129, A: 255}
	SuccessGradientColor = color.RGBA{R: 16, G: 185, B: 129, A: 180}

	// WarningColor - Modern orange with glassmorphic accent
	WarningColor = color.RGBA{R: 245, G: 158, B: 11, A: 255}
	WarningGradientColor = color.RGBA{R: 245, G: 158, B: 11, A: 180}

	// InfoColor - Modern blue with glassmorphic accent
	InfoColor = color.RGBA{R: 59, G: 130, B: 246, A: 255}
	InfoGradientColor = color.RGBA{R: 59, G: 130, B: 246, A: 180}
)

// createModernDialogContent creates a compact modern styled dialog content
// Layout: [Accent Bar] [Message Text] arranged horizontally
func createModernDialogContent(icon string, iconColor, accentColor color.Color, message string) fyne.CanvasObject {
	// Vertical accent bar on the left (thin colored strip)
	accentBar := canvas.NewRectangle(accentColor)
	accentBar.SetMinSize(fyne.NewSize(4, 60))

	// Message label - aligned left, wraps text
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	// Final layout with accent bar on the left
	return container.NewBorder(
		nil, nil,
		accentBar,
		nil,
		container.NewPadded(messageLabel),
	)
}

// ShowError displays a modern compact error dialog
// Blocks until the user dismisses the dialog
func ShowError(window fyne.Window, title, message string) {
	loc := i18n.GetGlobalLocalizer()

	content := createModernDialogContent(
		"✗",
		ErrorColor,
		ErrorGradientColor,
		message,
	)

	// Use localized title if empty
	if title == "" {
		title = loc.Get("dialog.error")
	}

	d := dialog.NewCustom(title, loc.Get("action.ok"), content, window)
	d.Resize(fyne.NewSize(450, 120))
	d.Show()
}

// ShowSuccess displays a modern compact success notification
// Auto-dismisses after 3 seconds
func ShowSuccess(window fyne.Window, title, message string) func() {
	loc := i18n.GetGlobalLocalizer()

	content := createModernDialogContent(
		"✓",
		SuccessColor,
		SuccessGradientColor,
		message,
	)

	// Use localized title if empty
	if title == "" {
		title = loc.Get("dialog.success")
	}

	d := dialog.NewCustom(title, loc.Get("action.ok"), content, window)
	d.Resize(fyne.NewSize(450, 120))
	d.Show()

	// Auto-dismiss after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		fyne.Do(func() {
			d.Hide()
		})
	}()

	return func() {
		d.Hide()
	}
}

// ShowWarning displays a modern compact warning dialog
// Blocks until the user dismisses the dialog
func ShowWarning(window fyne.Window, title, message string) {
	loc := i18n.GetGlobalLocalizer()

	content := createModernDialogContent(
		"⚠",
		WarningColor,
		WarningGradientColor,
		message,
	)

	// Use localized title if empty
	if title == "" {
		title = loc.Get("dialog.warning")
	}

	d := dialog.NewCustom(title, loc.Get("action.ok"), content, window)
	d.Resize(fyne.NewSize(450, 120))
	d.Show()
}

// ShowInfo displays a modern compact information dialog
// Blocks until the user dismisses the dialog
func ShowInfo(window fyne.Window, title, message string) {
	loc := i18n.GetGlobalLocalizer()

	content := createModernDialogContent(
		"ℹ",
		InfoColor,
		InfoGradientColor,
		message,
	)

	// Use localized title if empty
	if title == "" {
		title = loc.Get("dialog.info")
	}

	d := dialog.NewCustom(title, loc.Get("action.ok"), content, window)
	d.Resize(fyne.NewSize(450, 120))
	d.Show()
}

// ShowConfirmation displays a modern compact confirmation dialog
// Calls the callback with true if user confirms, false if user cancels
func ShowConfirmation(window fyne.Window, title, message string, callback func(bool)) {
	loc := i18n.GetGlobalLocalizer()

	// Use warning icon instead of question mark for better UX
	content := createModernDialogContent(
		"⚠",
		WarningColor,
		WarningGradientColor,
		message,
	)

	// Use localized title if empty
	if title == "" {
		title = loc.Get("dialog.confirm")
	}

	d := dialog.NewCustomConfirm(title, loc.Get("action.yes"), loc.Get("action.no"), content, callback, window)
	d.Resize(fyne.NewSize(450, 120))
	d.Show()
}

// ShowErrorWithDetails displays a modern error dialog with expandable details
func ShowErrorWithDetails(window fyne.Window, title, message, details string) {
	loc := i18n.GetGlobalLocalizer()

	// Vertical accent bar
	accentBar := canvas.NewRectangle(ErrorGradientColor)
	accentBar.SetMinSize(fyne.NewSize(4, 0))

	// Message label
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord

	// Details section
	detailsLabel := widget.NewLabel(details)
	detailsLabel.Wrapping = fyne.TextWrapWord

	// Accordion for details
	accordion := widget.NewAccordion(
		widget.NewAccordionItem("Details", container.NewScroll(detailsLabel)),
	)

	// Vertical layout for error with details
	contentBox := container.NewVBox(
		container.NewPadded(messageLabel),
		widget.NewSeparator(),
		accordion,
	)

	// Add accent bar
	content := container.NewBorder(nil, nil, accentBar, nil, contentBox)

	// Use localized title if empty
	if title == "" {
		title = loc.Get("dialog.error")
	}

	d := dialog.NewCustom(title, loc.Get("action.ok"), content, window)
	d.Resize(fyne.NewSize(500, 300))
	d.Show()
}

// ShowProgressDialog displays a modern compact progress dialog
func ShowProgressDialog(window fyne.Window, initialMessage string) (updateMessage func(string), dismiss func()) {
	loc := i18n.GetGlobalLocalizer()

	// Vertical accent bar
	accentBar := canvas.NewRectangle(InfoGradientColor)
	accentBar.SetMinSize(fyne.NewSize(4, 60))

	// Progress bar
	progress := widget.NewProgressBarInfinite()

	// Message label
	messageLabel := widget.NewLabel(initialMessage)
	messageLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Vertical layout
	contentBox := container.NewVBox(
		messageLabel,
		progress,
	)

	content := container.NewBorder(nil, nil, accentBar, nil, container.NewPadded(contentBox))

	d := dialog.NewCustomWithoutButtons(loc.Get("dialog.please_wait"), content, window)
	d.Resize(fyne.NewSize(400, 100))
	d.Show()

	updateMessage = func(msg string) {
		messageLabel.SetText(msg)
	}

	dismiss = func() {
		d.Hide()
	}

	return updateMessage, dismiss
}

// ShowQuickNotification displays a brief compact notification
// Auto-dismisses after 2 seconds
func ShowQuickNotification(window fyne.Window, message string) {
	// Vertical accent bar
	accentBar := canvas.NewRectangle(SuccessGradientColor)
	accentBar.SetMinSize(fyne.NewSize(4, 40))

	// Message label
	messageLabel := widget.NewLabel(message)

	// Horizontal layout
	content := container.NewBorder(
		nil, nil,
		accentBar,
		nil,
		container.NewPadded(messageLabel),
	)

	d := dialog.NewCustomWithoutButtons("", content, window)
	d.Resize(fyne.NewSize(350, 60))
	d.Show()

	// Auto-dismiss after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		fyne.Do(func() {
			d.Hide()
		})
	}()
}
