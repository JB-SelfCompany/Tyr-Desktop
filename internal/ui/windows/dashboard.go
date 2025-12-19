package windows

import (
	"fmt"
	"image/color"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/dialogs"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/widgets"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/yggmail"
)

// ModernDashboardScreen represents the modern redesigned main status screen
// Modern design trends: Bento Grid layout, bold colors, micro-animations,
// enhanced accessibility, and mobile-first responsive design
type ModernDashboardScreen struct {
	app AppInterface

	// UI components
	mailAddressCard     *ModernCard
	serverInfoCard      *ModernCard
	peersCard           *ModernCard
	quickActionsCard    *ModernCard
	mailAddressLabel    *widget.RichText
	copyButton          *widget.Button
	deltaChatButton     *widget.Button
	serverInfoContainer *fyne.Container
	peersContainer      *fyne.Container
	startStopButton     *widget.Button

	// Update ticker
	ticker *time.Ticker
	stopCh chan struct{}
}

// ModernCard represents a modern card container with shadow effect and rounded corners
type ModernCard struct {
	widget.BaseWidget
	content   fyne.CanvasObject
	title     string
	accentColor color.Color
	container *fyne.Container
}

// NewModernCard creates a new modern card with title and accent color
func NewModernCard(title string, content fyne.CanvasObject, accentColor color.Color) *ModernCard {
	card := &ModernCard{
		content:     content,
		title:       title,
		accentColor: accentColor,
	}
	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer creates the card renderer with modern styling
func (mc *ModernCard) CreateRenderer() fyne.WidgetRenderer {
	// Accent bar (left side colored bar)
	accentBar := canvas.NewRectangle(mc.accentColor)
	accentBar.SetMinSize(fyne.NewSize(4, 0))

	// Title with bold, larger text
	titleLabel := widget.NewLabel(mc.title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Importance = widget.HighImportance

	// Card body with proper padding
	cardBody := container.NewBorder(
		container.NewVBox(titleLabel, widget.NewSeparator()), // Top: Header
		nil, // Bottom
		nil, // Left
		nil, // Right
		container.NewPadded(mc.content), // Center: Content with padding
	)

	// Combine accent bar + body
	mc.container = container.NewBorder(
		nil, nil,
		accentBar, nil,
		cardBody,
	)

	return widget.NewSimpleRenderer(mc.container)
}

// NewModernDashboardScreen creates a new modern dashboard screen
func NewModernDashboardScreen(app AppInterface) fyne.CanvasObject {
	screen := &ModernDashboardScreen{
		app:    app,
		stopCh: make(chan struct{}),
	}

	screen.setupUI()
	screen.startUpdates()

	return screen.buildLayout()
}

// setupUI initializes all UI components with modern design
func (mds *ModernDashboardScreen) setupUI() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()

	// Mail address display with rich text formatting
	mailAddress := serviceManager.GetMailAddress()
	if mailAddress == "" {
		mailAddress = loc.Get("dashboard.service_not_initialized")
	}

	mds.mailAddressLabel = widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**", mailAddress))
	mds.mailAddressLabel.Wrapping = fyne.TextWrapBreak

	// Modern gradient buttons
	mds.copyButton = widget.NewButton(fmt.Sprintf("📋 %s", loc.Get("dashboard.copy_address")), func() {
		mds.copyMailAddress()
	})
	mds.copyButton.Importance = widget.HighImportance

	mds.deltaChatButton = widget.NewButton(fmt.Sprintf("💬 %s", loc.Get("dashboard.setup_deltachat")), func() {
		mds.openDeltaChat()
	})
	mds.deltaChatButton.Importance = widget.HighImportance

	// Server info
	mds.serverInfoContainer = mds.buildServerInfo()

	// Peers container
	mds.peersContainer = container.NewVBox()
	mds.updatePeersList()

	// Start/Stop button with accent styling
	mds.startStopButton = widget.NewButton(fmt.Sprintf("▶ %s", loc.Get("dashboard.start_service")), func() {
		mds.toggleService()
	})
	mds.startStopButton.Importance = widget.HighImportance
	mds.updateStartStopButton()
}

// buildLayout creates the modern responsive layout
func (mds *ModernDashboardScreen) buildLayout() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Color palette - vibrant 2025 colors
	primaryColor := color.RGBA{R: 99, G: 102, B: 241, A: 255}    // Indigo
	successColor := color.RGBA{R: 16, G: 185, B: 129, A: 255}    // Green
	infoColor := color.RGBA{R: 59, G: 130, B: 246, A: 255}       // Blue
	warningColor := color.RGBA{R: 245, G: 158, B: 11, A: 255}    // Orange

	// Settings button (top-left toolbar)
	settingsButton := widget.NewButton(fmt.Sprintf("⚙ %s", loc.Get("dashboard.settings")), func() {
		mds.app.ShowSettings()
	})
	settingsButton.Importance = widget.LowImportance

	// Logs button
	logsButton := widget.NewButton(fmt.Sprintf("📝 %s", loc.Get("dashboard.logs")), func() {
		mds.app.ShowLogs()
	})
	logsButton.Importance = widget.LowImportance

	// Top toolbar with settings and logs buttons
	toolbar := container.NewHBox(
		settingsButton,
		logsButton,
	)

	// Service Control Card (Full Width, Primary Position)
	serviceControlContent := container.NewVBox(
		mds.startStopButton,
	)
	serviceControlCard := NewModernCard(loc.Get("dashboard.service_control"), serviceControlContent, primaryColor)

	// Mail Address Card (Full Width)
	mailContent := container.NewVBox(
		mds.mailAddressLabel,
		container.NewGridWithColumns(2,
			mds.copyButton,
			mds.deltaChatButton,
		),
	)
	mds.mailAddressCard = NewModernCard(loc.Get("dashboard.your_mail_address"), mailContent, successColor)

	// Server Info Card
	mds.serverInfoCard = NewModernCard(loc.Get("dashboard.server_configuration"), mds.serverInfoContainer, infoColor)

	// Peers Card (Dynamic height) - adapts to content size
	mds.peersCard = NewModernCard(loc.Get("dashboard.connected_peers"), mds.peersContainer, warningColor)

	// Adaptive layout - optimized for compact display
	// Reorganized: Service Control first, then Mail, then Server Info, then Peers
	mainContent := container.NewVBox(
		// Row 1: Service Control (full width, primary position)
		serviceControlCard,

		// Row 2: Mail Address (full width)
		mds.mailAddressCard,

		// Row 3: Server Info (full width)
		mds.serverInfoCard,

		// Row 4: Peers (full width, dynamic height)
		mds.peersCard,
	)

	// Wrap in scroll container for full window adaptability
	scrollContent := container.NewScroll(mainContent)

	// Combine toolbar and main content
	finalLayout := container.NewBorder(
		toolbar, // Top: Settings toolbar
		nil,     // Bottom
		nil,     // Left
		nil,     // Right
		scrollContent, // Center: Main content
	)

	return finalLayout
}

// buildServerInfo creates the server information content
func (mds *ModernDashboardScreen) buildServerInfo() *fyne.Container {
	loc := i18n.GetGlobalLocalizer()
	config := mds.app.GetConfig()

	smtpLabel := widget.NewLabel(fmt.Sprintf("📧 %s: %s", loc.Get("dashboard.smtp_address"), config.ServiceSettings.SMTPAddress))
	smtpLabel.Wrapping = fyne.TextWrapWord

	imapLabel := widget.NewLabel(fmt.Sprintf("📬 %s: %s", loc.Get("dashboard.imap_address"), config.ServiceSettings.IMAPAddress))
	imapLabel.Wrapping = fyne.TextWrapWord

	dbLabel := widget.NewLabel(fmt.Sprintf("💾 %s: %s", loc.Get("settings.backup.db.title"), config.ServiceSettings.DatabasePath))
	dbLabel.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		smtpLabel,
		imapLabel,
		dbLabel,
	)
}

// updatePeersList updates the peers display with modern peer cards
func (mds *ModernDashboardScreen) updatePeersList() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()
	peers := serviceManager.GetPeerStats()

	// Clear existing peers
	mds.peersContainer.Objects = []fyne.CanvasObject{}

	if len(peers) == 0 {
		emptyIcon := widget.NewLabel("🌐")
		emptyText := widget.NewLabel(loc.Get("widget.no_peers"))
		emptyText.Alignment = fyne.TextAlignCenter

		emptyState := container.NewVBox(
			emptyIcon,
			emptyText,
		)
		mds.peersContainer.Objects = append(mds.peersContainer.Objects, emptyState)
	} else {
		for _, peer := range peers {
			peerCard := widgets.NewModernPeerCard(&peer)
			mds.peersContainer.Objects = append(mds.peersContainer.Objects, peerCard)
		}
	}

	mds.peersContainer.Refresh()
}

// updateStartStopButton updates the start/stop button state with icon
func (mds *ModernDashboardScreen) updateStartStopButton() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()

	if serviceManager.IsRunning() {
		mds.startStopButton.SetText(fmt.Sprintf("⏹ %s", loc.Get("dashboard.stop_service")))
		mds.startStopButton.Importance = widget.DangerImportance
	} else {
		mds.startStopButton.SetText(fmt.Sprintf("▶ %s", loc.Get("dashboard.start_service")))
		mds.startStopButton.Importance = widget.HighImportance
	}
}

// toggleService starts or stops the service
func (mds *ModernDashboardScreen) toggleService() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()

	if serviceManager.IsRunning() {
		// Stop service using SoftStop for clean peer disconnection
		mds.startStopButton.Disable()
		go func() {
			if err := serviceManager.SoftStop(); err != nil {
				log.Printf("Failed to stop service: %v", err)
				dialogs.ShowError(mds.app.GetMainWindow(), "", err.Error())
			}
			// Always enable button and update its state regardless of success/failure
			fyne.Do(func() {
				mds.startStopButton.Enable()
				mds.updateStartStopButton()
				mds.startStopButton.Refresh()
				// Update system tray status
				mds.app.UpdateSystemTrayStatus()
			})
		}()
	} else {
		// Start service
		mds.startStopButton.Disable()
		go func() {
			// Initialize if needed
			if serviceManager.GetStatus() == yggmail.StatusStopped {
				if err := serviceManager.Initialize(); err != nil {
					log.Printf("Failed to initialize service: %v", err)
					dialogs.ShowError(
						mds.app.GetMainWindow(),
						"",
						fmt.Sprintf("%s: %v", loc.Get("dashboard.error_initialize_service"), err),
					)
					fyne.Do(func() {
						mds.startStopButton.Enable()
						mds.updateStartStopButton()
						mds.startStopButton.Refresh()
					})
					return
				}
			}

			// Start service
			if err := serviceManager.Start(); err != nil {
				log.Printf("Failed to start service: %v", err)
				dialogs.ShowError(
					mds.app.GetMainWindow(),
					"",
					fmt.Sprintf("%s: %v", loc.Get("dashboard.error_start_service"), err),
				)
			}
			// Always enable button and update its state regardless of success/failure
			fyne.Do(func() {
				mds.startStopButton.Enable()
				mds.updateStartStopButton()
				mds.startStopButton.Refresh()
				// Update system tray status
				mds.app.UpdateSystemTrayStatus()
			})
		}()
	}
}

// copyMailAddress copies the mail address to clipboard
func (mds *ModernDashboardScreen) copyMailAddress() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()
	mailAddress := serviceManager.GetMailAddress()

	if mailAddress == "" {
		dialogs.ShowInfo(mds.app.GetMainWindow(),
			loc.Get("dashboard.no_address"),
			loc.Get("dashboard.no_address_message"),
		)
		return
	}

	// Copy to clipboard
	clipboard := mds.app.GetMainWindow().Clipboard()
	clipboard.SetContent(mailAddress)

	dialogs.ShowInfo(mds.app.GetMainWindow(),
		loc.Get("dashboard.copied"),
		loc.Get("dashboard.copied_message"),
	)
}

// openDeltaChat opens DeltaChat with auto-configured account
func (mds *ModernDashboardScreen) openDeltaChat() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()
	config := mds.app.GetConfig()

	// Get mail address
	mailAddress := serviceManager.GetMailAddress()
	if mailAddress == "" {
		dialogs.ShowError(
			mds.app.GetMainWindow(),
			"",
			loc.Get("dashboard.no_address_message"),
		)
		return
	}

	// Get password
	password, err := config.GetPassword()
	if err != nil || password == "" {
		dialogs.ShowError(
			mds.app.GetMainWindow(),
			"",
			fmt.Sprintf("%s: %v", loc.Get("dashboard.error_retrieve_password"), err),
		)
		return
	}

	// Parse SMTP and IMAP addresses
	smtpHost, smtpPort := parseAddress(config.ServiceSettings.SMTPAddress)
	imapHost, imapPort := parseAddress(config.ServiceSettings.IMAPAddress)

	// Generate dclogin:// URL with full IMAP/SMTP configuration
	dcloginURL := generateDCLoginURL(mailAddress, password, imapHost, imapPort, smtpHost, smtpPort)

	// Log the URL (for debugging)
	log.Printf("Opening DeltaChat with URL: %s", dcloginURL)

	// Open DeltaChat directly
	if err := openDCLoginURL(dcloginURL); err != nil {
		log.Printf("Failed to open DeltaChat URL: %v", err)
		// Copy to clipboard as fallback
		clipboard := mds.app.GetMainWindow().Clipboard()
		clipboard.SetContent(dcloginURL)
		dialogs.ShowError(
			mds.app.GetMainWindow(),
			"",
			loc.Get("dashboard.error_open_deltachat"),
		)
	}
}

// startUpdates starts periodic UI updates
func (mds *ModernDashboardScreen) startUpdates() {
	mds.ticker = time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-mds.ticker.C:
				mds.updateUI()
			case <-mds.stopCh:
				return
			}
		}
	}()
}

// updateUI updates all dynamic UI elements
func (mds *ModernDashboardScreen) updateUI() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mds.app.GetServiceManager()

	// All UI updates must happen on the Fyne main thread
	fyne.Do(func() {
		// Update mail address
		mailAddress := serviceManager.GetMailAddress()
		if mailAddress == "" {
			mailAddress = loc.Get("dashboard.service_not_initialized")
		}
		mds.mailAddressLabel.ParseMarkdown(fmt.Sprintf("**%s**", mailAddress))

		// Update peers list
		mds.updatePeersList()

		// Update button state
		mds.updateStartStopButton()

		// Note: System tray is now updated reactively via status channel monitoring
		// See app.startStatusMonitoring() - no need to call UpdateSystemTrayStatus() here
	})
}

// StopUpdates stops the periodic updates
func (mds *ModernDashboardScreen) StopUpdates() {
	if mds.ticker != nil {
		mds.ticker.Stop()
		close(mds.stopCh)
	}
}

// openURL opens a URL in the default browser or application handler
func openURL(urlStr string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", urlStr}
	case "darwin":
		cmd = "open"
		args = []string{urlStr}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{urlStr}
	}

	return exec.Command(cmd, args...).Start()
}

// openDCLoginURL opens a dclogin:// URL with proper escaping for Windows
func openDCLoginURL(dcloginURL string) error {
	switch runtime.GOOS {
	case "windows":
		// On Windows, cmd /c start requires special handling:
		// The first argument after "start" is the window title (can be empty "")
		// Then comes the URL which must be in quotes if it contains special chars
		// Use rundll32 instead to avoid cmd parsing issues
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", dcloginURL).Start()
	case "darwin":
		return exec.Command("open", dcloginURL).Start()
	default: // "linux", "freebsd", "openbsd", "netbsd"
		return exec.Command("xdg-open", dcloginURL).Start()
	}
}

// generateDCLoginURL creates a dclogin:// URL for DeltaChat auto-configuration
// Format: dclogin://user@host/?v=1&p=password&ih=imaphost&ip=imapport&is=plain&ic=3&sh=smtphost&sp=smtpport&ss=plain&sc=3
func generateDCLoginURL(email, password, imapHost, imapPort, smtpHost, smtpPort string) string {
	// Build query parameters
	params := url.Values{}
	params.Set("v", "1")       // Version
	params.Set("p", password)  // Password
	params.Set("ih", imapHost) // IMAP hostname
	params.Set("ip", imapPort) // IMAP port
	params.Set("is", "plain")  // IMAP security (plain = no encryption)
	params.Set("ic", "3")      // IMAP certificate checks (3 = accept invalid, for localhost)
	params.Set("sh", smtpHost) // SMTP hostname
	params.Set("sp", smtpPort) // SMTP port
	params.Set("ss", "plain")  // SMTP security (plain = no encryption)
	params.Set("sc", "3")      // SMTP certificate checks (3 = accept invalid, for localhost)

	// Build dclogin URL
	// Format: dclogin://user@host/?parameters
	return fmt.Sprintf("dclogin://%s/?%s", email, params.Encode())
}

// parseAddress splits an address string into host and port
// Example: "127.0.0.1:1025" -> ("127.0.0.1", "1025")
func parseAddress(addr string) (host, port string) {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return addr, ""
}
