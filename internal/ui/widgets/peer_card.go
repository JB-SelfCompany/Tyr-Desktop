package widgets

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/yggmail"
)

// ModernPeerCard is a modern redesigned widget for displaying peer information
// Modern design: tall card layout, vibrant colors, micro-interactions,
// enhanced data visualization with progress bars
type ModernPeerCard struct {
	widget.BaseWidget

	peer *yggmail.PeerInfo

	// UI components
	statusIndicator  *canvas.Circle
	addressLabel     *widget.RichText
	directionBadge   *ModernBadge
	latencyLabel     *widget.Label
	uptimeLabel      *widget.Label
	rxLabel          *widget.Label
	txLabel          *widget.Label
	errorLabel       *widget.Label
	container        *fyne.Container
}

// ModernBadge represents a small colored badge for labels
type ModernBadge struct {
	widget.BaseWidget
	text      string
	bgColor   color.Color
	textColor color.Color
	container *fyne.Container
}

// NewModernBadge creates a new modern badge
func NewModernBadge(text string, bgColor, textColor color.Color) *ModernBadge {
	badge := &ModernBadge{
		text:      text,
		bgColor:   bgColor,
		textColor: textColor,
	}
	badge.ExtendBaseWidget(badge)
	return badge
}

// CreateRenderer creates the badge renderer
func (mb *ModernBadge) CreateRenderer() fyne.WidgetRenderer {
	// Background
	bg := canvas.NewRectangle(mb.bgColor)

	// Text
	label := widget.NewLabel(mb.text)
	label.TextStyle = fyne.TextStyle{Bold: true}

	// Container
	content := container.NewPadded(
		container.NewCenter(label),
	)

	mb.container = container.NewMax(bg, content)
	mb.container.Resize(fyne.NewSize(80, 24))

	return widget.NewSimpleRenderer(mb.container)
}

// NewModernPeerCard creates a new modern peer card widget
func NewModernPeerCard(peer *yggmail.PeerInfo) fyne.CanvasObject {
	if peer == nil {
		return widget.NewLabel("Invalid peer data")
	}

	card := &ModernPeerCard{
		peer: peer,
	}
	card.ExtendBaseWidget(card)

	return card.build()
}

// build constructs the modern tall card layout
func (mpc *ModernPeerCard) build() fyne.CanvasObject {
	// Status indicator with glow effect
	mpc.statusIndicator = canvas.NewCircle(mpc.getStatusColor())
	mpc.statusIndicator.Resize(fyne.NewSize(24, 24))
	mpc.statusIndicator.StrokeWidth = 0

	// Peer address with bold typography and proper wrapping
	addressText := mpc.getSanitizedAddress()
	mpc.addressLabel = widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**", addressText))
	mpc.addressLabel.Wrapping = fyne.TextWrapWord

	// Direction badge with modern colors
	directionText := "IN"
	badgeColor := color.RGBA{R: 59, G: 130, B: 246, A: 255} // Blue for inbound
	if !mpc.peer.Inbound {
		directionText = "OUT"
		badgeColor = color.RGBA{R: 16, G: 185, B: 129, A: 255} // Green for outbound
	}
	mpc.directionBadge = NewModernBadge(directionText, badgeColor, color.White)

	// Header: Status + Address (full width with wrapping)
	header := container.NewBorder(
		nil, nil,
		container.NewHBox(mpc.statusIndicator),
		mpc.directionBadge,
		mpc.addressLabel,
	)

	// Network stats with icons - compact inline layout
	mpc.latencyLabel = widget.NewLabel("⚡ " + mpc.getLatencyText())
	mpc.uptimeLabel = widget.NewLabel("⏱ " + mpc.getUptimeText())

	statsGrid := container.NewGridWithColumns(2,
		mpc.latencyLabel,
		mpc.uptimeLabel,
	)

	// Data transfer - compact text-only display (no progress bars for smaller size)
	rxBytes, rxRate := mpc.peer.RXBytes, mpc.peer.RXRate
	txBytes, txRate := mpc.peer.TXBytes, mpc.peer.TXRate

	mpc.rxLabel = widget.NewLabel(mpc.getDataTransferText("📥", rxBytes, rxRate))
	mpc.rxLabel.Wrapping = fyne.TextWrapWord

	mpc.txLabel = widget.NewLabel(mpc.getDataTransferText("📤", txBytes, txRate))
	mpc.txLabel.Wrapping = fyne.TextWrapWord

	dataTransferGrid := container.NewGridWithColumns(2,
		mpc.rxLabel,
		mpc.txLabel,
	)

	// Build main card content
	cardContent := container.NewVBox(
		header,
		widget.NewSeparator(),
		statsGrid,
		dataTransferGrid,
	)

	// Add error message if present
	if mpc.peer.LastError != "" {
		errorIcon := widget.NewLabel("⚠")
		mpc.errorLabel = widget.NewLabel(mpc.getSanitizedError())
		mpc.errorLabel.Wrapping = fyne.TextWrapWord
		mpc.errorLabel.Importance = widget.DangerImportance

		errorRow := container.NewHBox(errorIcon, mpc.errorLabel)
		cardContent.Add(widget.NewSeparator())
		cardContent.Add(errorRow)
	}

	// Card background with accent border
	accentBar := canvas.NewRectangle(mpc.getAccentColor())
	accentBar.SetMinSize(fyne.NewSize(4, 0))

	cardBody := container.NewPadded(cardContent)

	cardWithAccent := container.NewBorder(
		nil, nil,
		accentBar, nil,
		cardBody,
	)

	// Outer container with shadow effect (padding)
	mpc.container = container.NewPadded(cardWithAccent)

	return mpc.container
}

// getStatusColor returns the color for the status indicator
func (mpc *ModernPeerCard) getStatusColor() color.Color {
	if mpc.peer.Status {
		// Connected - vibrant green 
		return color.RGBA{R: 16, G: 185, B: 129, A: 255}
	}
	// Disconnected - bold red 
	return color.RGBA{R: 239, G: 68, B: 68, A: 255}
}

// getAccentColor returns the accent color for the card border
func (mpc *ModernPeerCard) getAccentColor() color.Color {
	if mpc.peer.Status {
		return color.RGBA{R: 16, G: 185, B: 129, A: 255} // Green
	}
	return color.RGBA{R: 239, G: 68, B: 68, A: 255} // Red
}

// getSanitizedAddress returns the peer address with sanitization
func (mpc *ModernPeerCard) getSanitizedAddress() string {
	// Return full address - wrapping will handle long addresses
	return mpc.peer.Address
}

// getLatencyText returns formatted latency text
func (mpc *ModernPeerCard) getLatencyText() string {
	if mpc.peer.Latency <= 0 {
		return "N/A"
	}

	latency := mpc.peer.Latency
	if latency < 50 {
		return fmt.Sprintf("%d ms (Excellent)", latency)
	} else if latency < 150 {
		return fmt.Sprintf("%d ms (Good)", latency)
	} else {
		return fmt.Sprintf("%d ms (High)", latency)
	}
}

// getUptimeText returns formatted uptime text
func (mpc *ModernPeerCard) getUptimeText() string {
	uptime := mpc.peer.Uptime

	if uptime <= 0 {
		return "N/A"
	}

	if uptime < 60 {
		return fmt.Sprintf("%ds", uptime)
	} else if uptime < 3600 {
		return fmt.Sprintf("%dm", uptime/60)
	} else if uptime < 86400 {
		hours := uptime / 3600
		minutes := (uptime % 3600) / 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		days := uptime / 86400
		hours := (uptime % 86400) / 3600
		return fmt.Sprintf("%dd %dh", days, hours)
	}
}

// getDataTransferText returns formatted data transfer text
func (mpc *ModernPeerCard) getDataTransferText(direction string, bytes, rate int64) string {
	if bytes < 0 {
		bytes = 0
	}
	if rate < 0 {
		rate = 0
	}

	bytesStr := formatBytes(bytes)
	rateStr := formatBytes(rate) + "/s"

	return fmt.Sprintf("%s: %s (%s)", direction, bytesStr, rateStr)
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	} else if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
	}
}

// getSanitizedError returns sanitized error message
func (mpc *ModernPeerCard) getSanitizedError() string {
	errorMsg := mpc.peer.LastError

	const maxLength = 150
	if len(errorMsg) > maxLength {
		return errorMsg[:maxLength] + "..."
	}

	return errorMsg
}

// CreateRenderer implements the fyne.Widget interface
func (mpc *ModernPeerCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(mpc.container)
}

// Update refreshes the peer card with new data
func (mpc *ModernPeerCard) Update(peer *yggmail.PeerInfo) {
	if peer == nil {
		return
	}

	mpc.peer = peer

	// Update all components
	if mpc.statusIndicator != nil {
		mpc.statusIndicator.FillColor = mpc.getStatusColor()
		mpc.statusIndicator.Refresh()
	}

	if mpc.addressLabel != nil {
		addressText := mpc.getSanitizedAddress()
		mpc.addressLabel.ParseMarkdown(fmt.Sprintf("**%s**", addressText))
	}

	if mpc.latencyLabel != nil {
		mpc.latencyLabel.SetText("⚡ " + mpc.getLatencyText())
	}

	if mpc.uptimeLabel != nil {
		mpc.uptimeLabel.SetText("⏱ " + mpc.getUptimeText())
	}

	// Update data transfer labels
	if mpc.rxLabel != nil {
		mpc.rxLabel.SetText(mpc.getDataTransferText("📥", peer.RXBytes, peer.RXRate))
	}

	if mpc.txLabel != nil {
		mpc.txLabel.SetText(mpc.getDataTransferText("📤", peer.TXBytes, peer.TXRate))
	}

	mpc.Refresh()
}
