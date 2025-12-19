package windows

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/logger"
)

// LogsScreen displays real-time application logs
type LogsScreen struct {
	app AppInterface

	// UI components
	logText   *widget.RichText
	logScroll *container.Scroll
	logLines  []string
	maxLines  int
	mu        sync.Mutex
	lastCount int

	// Update ticker
	ticker *time.Ticker
	stopCh chan struct{}
}

// NewLogsScreen creates a new logs screen with real-time updates
func NewLogsScreen(app AppInterface) fyne.CanvasObject {
	screen := &LogsScreen{
		app:      app,
		logLines: []string{},
		maxLines: 500,
		stopCh:   make(chan struct{}),
	}

	return screen.buildLayout()
}

// buildLayout creates the logs screen layout
func (ls *LogsScreen) buildLayout() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Back button - compact with icon and text
	backButton := widget.NewButtonWithIcon(loc.Get("action.back"), theme.NavigateBackIcon(), func() {
		ls.stopUpdates()
		ls.app.ShowDashboard()
	})
	backButton.Importance = widget.LowImportance

	// Create a container for back button aligned to the left
	backButtonContainer := container.NewHBox(backButton)

	// Title and subtitle
	titleLabel := widget.NewRichTextFromMarkdown("# " + loc.Get("logs.title"))
	subtitleLabel := widget.NewLabel(loc.Get("logs.subtitle"))
	subtitleLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Copy all logs button
	copyButton := widget.NewButton(loc.Get("logs.copy_all"), func() {
		ls.copyAllLogs()
	})
	copyButton.Importance = widget.MediumImportance

	// Clear logs button
	clearButton := widget.NewButton(loc.Get("logs.clear_logs"), func() {
		ls.clearLogs()
	})
	clearButton.Importance = widget.MediumImportance

	controlBar := container.NewHBox(
		copyButton,
		clearButton,
	)

	header := container.NewVBox(
		backButtonContainer,
		titleLabel,
		subtitleLabel,
		controlBar,
		widget.NewSeparator(),
	)

	// Log display with RichText (read-only, supports text selection and copy)
	ls.logText = widget.NewRichText()
	ls.logText.Wrapping = fyne.TextWrapWord // Enable word wrapping for long lines

	// Initialize with existing logs
	ls.updateLogs()

	// Log container with scroll - save reference for auto-scrolling
	ls.logScroll = container.NewScroll(ls.logText)

	// Start live updates
	ls.startUpdates()

	// Return layout with header at top and logs filling the rest
	return container.NewBorder(
		container.NewPadded(header),
		nil, nil, nil,
		ls.logScroll,
	)
}

// isAtBottom checks if the scroll is at or near the bottom
func (ls *LogsScreen) isAtBottom() bool {
	if ls.logScroll == nil {
		return true
	}

	// Get the current scroll position
	currentOffset := ls.logScroll.Offset.Y

	// Use MinSize() instead of Size() to handle pre-render cases
	// Size() can return 0 before widget is rendered
	contentHeight := ls.logText.MinSize().Height
	containerHeight := ls.logScroll.Size().Height
	if containerHeight == 0 {
		containerHeight = ls.logScroll.MinSize().Height
	}

	// Calculate maximum scroll offset
	maxOffset := contentHeight - containerHeight
	if maxOffset < 0 {
		maxOffset = 0
	}

	// Consider "at bottom" if within 100 pixels of the bottom
	// or if content fits entirely in view
	return currentOffset >= maxOffset-100 || contentHeight <= containerHeight
}

// updateLogs refreshes the log display
func (ls *LogsScreen) updateLogs() {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	entries := logger.GetLogEntries()

	// Check if there are new entries
	if len(entries) == ls.lastCount {
		return // No new logs, skip update
	}
	ls.lastCount = len(entries)

	// Build plain text with emoji indicators
	var builder strings.Builder

	for _, entry := range entries {
		timestamp := entry.Timestamp.Format("15:04:05")
		levelEmoji := ls.getLevelEmoji(entry.Level)

		// Format: 🔵 [15:04:05] INFO: message
		builder.WriteString(fmt.Sprintf("%s [%s] %s: %s\n", levelEmoji, timestamp, entry.Level, entry.Message))
	}

	// Update RichText content
	// UI update must happen on the Fyne main thread
	text := builder.String()
	fyne.Do(func() {
		// Check if user was at the bottom before update
		wasAtBottom := ls.isAtBottom()

		// Save current scroll position before text update
		savedOffset := ls.logScroll.Offset

		// Update the RichText content with plain text segments
		// RichText doesn't have a cursor, so it won't auto-scroll like Entry
		if text == "" {
			loc := i18n.GetGlobalLocalizer()
			text = loc.Get("logs.no_logs")
		}
		ls.logText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Monospace: true},
				},
			},
		}
		ls.logText.Refresh()

		// Small delay to let the widget update its size
		time.AfterFunc(10*time.Millisecond, func() {
			fyne.Do(func() {
				if wasAtBottom {
					// User was at bottom, scroll to new bottom
					ls.scrollToBottom()
				} else {
					// User was reading older logs, keep their position
					ls.logScroll.Offset = savedOffset
					ls.logScroll.Refresh()
				}
			})
		})
	})
}

// scrollToBottom scrolls the log view to the bottom
func (ls *LogsScreen) scrollToBottom() {
	if ls.logScroll == nil {
		return
	}

	// Use MinSize() to handle pre-render cases properly
	// Size() can return 0 before widget is rendered
	contentHeight := ls.logText.MinSize().Height
	containerHeight := ls.logScroll.Size().Height
	if containerHeight == 0 {
		containerHeight = ls.logScroll.MinSize().Height
	}

	// Calculate the maximum scroll offset to reach the bottom
	maxOffset := contentHeight - containerHeight

	if maxOffset > 0 {
		ls.logScroll.Offset = fyne.NewPos(0, maxOffset)
	} else {
		ls.logScroll.Offset = fyne.NewPos(0, 0)
	}
	ls.logScroll.Refresh()
}

// getLevelEmoji returns a colored emoji for log levels
func (ls *LogsScreen) getLevelEmoji(level string) string {
	switch strings.ToUpper(level) {
	case "ERROR":
		return "🔴"
	case "WARN", "WARNING":
		return "🟡"
	case "INFO":
		return "🔵"
	case "DEBUG":
		return "🟢"
	default:
		return "⚪"
	}
}

// clearLogs clears all logs
func (ls *LogsScreen) clearLogs() {
	logger.ClearLogs()
	ls.updateLogs()
}

// copyAllLogs copies all logs to clipboard
func (ls *LogsScreen) copyAllLogs() {
	entries := logger.GetLogEntries()
	if len(entries) == 0 {
		return
	}

	// Build plain text version for clipboard
	var builder strings.Builder
	for _, entry := range entries {
		timestamp := entry.Timestamp.Format("15:04:05")
		levelEmoji := ls.getLevelEmoji(entry.Level)
		builder.WriteString(fmt.Sprintf("%s [%s] %s: %s\n", levelEmoji, timestamp, entry.Level, entry.Message))
	}

	// Copy to clipboard
	clipboard := ls.app.GetMainWindow().Clipboard()
	clipboard.SetContent(builder.String())
}

// startUpdates starts periodic log updates
func (ls *LogsScreen) startUpdates() {
	ls.ticker = time.NewTicker(500 * time.Millisecond) // Update every 500ms

	go func() {
		for {
			select {
			case <-ls.ticker.C:
				ls.updateLogs()
			case <-ls.stopCh:
				return
			}
		}
	}()
}

// stopUpdates stops the periodic updates
func (ls *LogsScreen) stopUpdates() {
	if ls.ticker != nil {
		ls.ticker.Stop()
		close(ls.stopCh)
	}
}
