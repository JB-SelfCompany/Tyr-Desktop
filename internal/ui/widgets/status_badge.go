package widgets

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/yggmail"
)

// ModernStatusBadge displays the service status with animated visual effects
// Modern design : pulsing animation, gradient backgrounds, bold typography
type ModernStatusBadge struct {
	widget.BaseWidget

	status          yggmail.ServiceStatus
	outerCircle     *canvas.Circle
	innerCircle     *canvas.Circle
	label           *widget.RichText
	pulseTimer      *time.Ticker
	pulsePhase      float64
	container       *fyne.Container
}

// NewModernStatusBadge creates a new modern animated status badge
func NewModernStatusBadge(status yggmail.ServiceStatus) *ModernStatusBadge {
	badge := &ModernStatusBadge{
		status:      status,
		pulsePhase:  0,
	}

	badge.ExtendBaseWidget(badge)
	badge.startAnimation()

	return badge
}

// SetStatus updates the displayed status
func (msb *ModernStatusBadge) SetStatus(status yggmail.ServiceStatus) {
	if msb.status != status {
		msb.status = status
		msb.updateUI()
		msb.Refresh()
	}
}

// GetStatus returns the current status
func (msb *ModernStatusBadge) GetStatus() yggmail.ServiceStatus {
	return msb.status
}

// CreateRenderer creates the renderer with modern animated design
func (msb *ModernStatusBadge) CreateRenderer() fyne.WidgetRenderer {
	// Outer circle (pulsing effect) - fixed size, SMALLER for compact display
	msb.outerCircle = canvas.NewCircle(msb.getStatusColor())
	msb.outerCircle.Resize(fyne.NewSize(50, 50))
	msb.outerCircle.StrokeWidth = 0

	// Inner circle (solid core) - fixed size, SMALLER for compact display
	msb.innerCircle = canvas.NewCircle(msb.getStatusColor())
	msb.innerCircle.Resize(fyne.NewSize(40, 40))
	msb.innerCircle.StrokeWidth = 0

	// Status label with medium bold typography (not huge headline)
	statusText := msb.getStatusText()
	msb.label = widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**", statusText))

	// Icon based on status - larger font for visibility
	iconLabel := widget.NewLabel(msb.getStatusIcon())
	iconLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Layout: Horizontal for compact display
	badgeContent := container.NewHBox(
		container.NewMax(
			msb.outerCircle,
			container.NewCenter(msb.innerCircle),
		),
		container.NewVBox(
			iconLabel,
			msb.label,
		),
	)

	msb.container = container.NewCenter(badgeContent)

	return widget.NewSimpleRenderer(msb.container)
}

// updateUI updates the visual elements based on current status
func (msb *ModernStatusBadge) updateUI() {
	if msb.outerCircle != nil {
		msb.outerCircle.FillColor = msb.getStatusColorWithAlpha(0.3)
		msb.outerCircle.Refresh()
	}

	if msb.innerCircle != nil {
		msb.innerCircle.FillColor = msb.getStatusColor()
		msb.innerCircle.Refresh()
	}

	if msb.label != nil {
		statusText := msb.getStatusText()
		msb.label.ParseMarkdown(fmt.Sprintf("# **%s**", statusText))
		msb.label.Refresh()
	}
}

// startAnimation starts the pulsing animation for active states
func (msb *ModernStatusBadge) startAnimation() {
	msb.pulseTimer = time.NewTicker(50 * time.Millisecond)

	go func() {
		for range msb.pulseTimer.C {
			// Only animate for running/starting/stopping states
			if msb.status == yggmail.StatusRunning ||
				msb.status == yggmail.StatusStarting ||
				msb.status == yggmail.StatusStopping {

				msb.pulsePhase += 0.1
				if msb.pulsePhase > 2*math.Pi {
					msb.pulsePhase = 0
				}

				// Only animate opacity, NOT size (to prevent layout shifts)
				if msb.outerCircle != nil {
					// Fade effect - smooth pulsing opacity
					alpha := 0.15 + 0.15*math.Sin(msb.pulsePhase)
					fyne.Do(func() {
						msb.outerCircle.FillColor = msb.getStatusColorWithAlpha(alpha)
						msb.outerCircle.Refresh()
					})
				}
			}
		}
	}()
}

// StopAnimation stops the pulsing animation
func (msb *ModernStatusBadge) StopAnimation() {
	if msb.pulseTimer != nil {
		msb.pulseTimer.Stop()
	}
}

// getStatusColor returns the solid color for the current status
func (msb *ModernStatusBadge) getStatusColor() color.Color {
	switch msb.status {
	case yggmail.StatusRunning:
		// Vibrant green
		return color.RGBA{R: 16, G: 185, B: 129, A: 255}
	case yggmail.StatusStopped:
		// Modern gray
		return color.RGBA{R: 107, G: 114, B: 128, A: 255}
	case yggmail.StatusError:
		// Bold red
		return color.RGBA{R: 239, G: 68, B: 68, A: 255}
	case yggmail.StatusStarting:
		// Vibrant blue 
		return color.RGBA{R: 59, G: 130, B: 246, A: 255}
	case yggmail.StatusStopping:
		// Warm orange 
		return color.RGBA{R: 245, G: 158, B: 11, A: 255}
	default:
		return color.RGBA{R: 107, G: 114, B: 128, A: 255}
	}
}

// getStatusColorWithAlpha returns the status color with custom alpha
func (msb *ModernStatusBadge) getStatusColorWithAlpha(alpha float64) color.Color {
	baseColor := msb.getStatusColor()
	r, g, b, _ := baseColor.RGBA()

	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(alpha * 255),
	}
}

// getStatusText returns bold, modern status text
func (msb *ModernStatusBadge) getStatusText() string {
	switch msb.status {
	case yggmail.StatusRunning:
		return "ONLINE"
	case yggmail.StatusStopped:
		return "OFFLINE"
	case yggmail.StatusError:
		return "ERROR"
	case yggmail.StatusStarting:
		return "STARTING"
	case yggmail.StatusStopping:
		return "STOPPING"
	default:
		return "UNKNOWN"
	}
}

// getStatusIcon returns emoji icon for status
func (msb *ModernStatusBadge) getStatusIcon() string {
	switch msb.status {
	case yggmail.StatusRunning:
		return "✓"
	case yggmail.StatusStopped:
		return "○"
	case yggmail.StatusError:
		return "✕"
	case yggmail.StatusStarting:
		return "↻"
	case yggmail.StatusStopping:
		return "⏸"
	default:
		return "?"
	}
}

