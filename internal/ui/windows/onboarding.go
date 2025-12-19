package windows

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/dialogs"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)

// ModernOnboardingScreen represents the modern onboarding wizard
// Features: Clean design, smooth transitions, restore from backup, modern color palette
type ModernOnboardingScreen struct {
	app AppInterface

	// Current step (0-4: welcome, setup-mode, password/restore, peers, complete)
	currentStep int

	// Setup mode: "new" or "restore"
	setupMode string

	// UI elements
	stepIndicator   *fyne.Container
	contentArea     *fyne.Container
	navigationBar   *fyne.Container
	backButton      *widget.Button
	nextButton      *widget.Button
	finishButton    *widget.Button
	skipButton      *widget.Button

	// Step data - New Setup
	passwordEntry        *widget.Entry
	passwordConfirmEntry *widget.Entry
	peersEntry           *widget.Entry

	// Step data - Restore from Backup
	restoreFileLabel     *widget.Label
	restorePasswordEntry *widget.Entry
	selectedRestoreFile  string
}

// NewModernOnboardingScreen creates a new modern onboarding screen
func NewModernOnboardingScreen(app AppInterface) fyne.CanvasObject {
	screen := &ModernOnboardingScreen{
		app:         app,
		currentStep: 0,
		setupMode:   "", // Not yet chosen
	}

	screen.setupUI()
	screen.showStep(0)

	return screen.buildLayout()
}

// setupUI initializes all UI components
func (mos *ModernOnboardingScreen) setupUI() {
	loc := i18n.GetGlobalLocalizer()

	// Navigation buttons (без стрелок, чтобы избежать проблем с символами)
	mos.backButton = widget.NewButton(loc.Get("action.back"), func() {
		mos.goBack()
	})
	mos.backButton.Importance = widget.LowImportance
	mos.backButton.Disable()

	mos.nextButton = widget.NewButton(loc.Get("action.next"), func() {
		mos.goNext()
	})
	mos.nextButton.Importance = widget.HighImportance

	mos.skipButton = widget.NewButton("Skip", func() {
		mos.skipStep()
	})
	mos.skipButton.Importance = widget.LowImportance
	mos.skipButton.Hide()

	mos.finishButton = widget.NewButton(loc.Get("action.finish"), func() {
		mos.finish()
	})
	mos.finishButton.Importance = widget.HighImportance
	mos.finishButton.Hide()

	// Step indicator (will be built dynamically)
	mos.stepIndicator = container.NewHBox()

	// Content area
	mos.contentArea = container.NewVBox()

	// Navigation bar
	mos.navigationBar = container.NewBorder(
		nil, nil,
		mos.backButton,
		container.NewHBox(mos.skipButton, mos.nextButton, mos.finishButton),
		widget.NewLabel(""), // Spacer
	)
}

// buildLayout creates the main layout with modern styling
func (mos *ModernOnboardingScreen) buildLayout() fyne.CanvasObject {
	// Main content with padding - stretch to full width but add margins
	scrollContent := container.NewScroll(mos.contentArea)
	scrollContent.SetMinSize(fyne.NewSize(700, 300))

	mainWithPadding := container.NewPadded(scrollContent)

	// Footer with navigation
	footer := container.NewPadded(mos.navigationBar)

	// Combine all sections - no header
	return container.NewBorder(
		nil,     // Top (no header)
		footer,  // Bottom
		nil,     // Left
		nil,     // Right
		mainWithPadding, // Center
	)
}

// buildStepIndicator creates the step indicator dots
func (mos *ModernOnboardingScreen) buildStepIndicator(totalSteps int) *fyne.Container {
	dots := make([]fyne.CanvasObject, totalSteps)

	for i := 0; i < totalSteps; i++ {
		var dotColor color.Color
		var size float32 = 10

		if i == mos.currentStep {
			// Active step - blue, larger
			dotColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
			size = 12
		} else if i < mos.currentStep {
			// Completed step - white semi-transparent
			dotColor = color.RGBA{R: 255, G: 255, B: 255, A: 180}
		} else {
			// Future step - white very transparent
			dotColor = color.RGBA{R: 255, G: 255, B: 255, A: 100}
		}

		circle := canvas.NewCircle(dotColor)
		circle.Resize(fyne.NewSize(size, size))

		// Wrap in container to control size
		dotContainer := container.NewWithoutLayout(circle)
		dotContainer.Resize(fyne.NewSize(size+4, size+4))

		dots[i] = dotContainer
	}

	return container.NewHBox(dots...)
}

// showStep displays the content for a specific step
func (mos *ModernOnboardingScreen) showStep(step int) {
	mos.currentStep = step

	// Update button visibility
	if step == 0 {
		mos.backButton.Disable()
	} else {
		mos.backButton.Enable()
	}

	// Determine content based on setup mode
	switch step {
	case 0:
		mos.showWelcomeStep()
		mos.stepIndicator.Objects = mos.buildStepIndicator(3).Objects
	case 1:
		mos.showSetupModeStep()
		mos.stepIndicator.Objects = mos.buildStepIndicator(3).Objects
	case 2:
		if mos.setupMode == "restore" {
			mos.showRestoreStep()
			mos.stepIndicator.Objects = mos.buildStepIndicator(3).Objects
			mos.nextButton.Hide()
			mos.finishButton.Show()
		} else {
			mos.showPasswordStep()
			mos.stepIndicator.Objects = mos.buildStepIndicator(5).Objects
		}
	case 3:
		mos.showPeerStep()
		mos.stepIndicator.Objects = mos.buildStepIndicator(5).Objects
		mos.skipButton.Show()
	case 4:
		mos.showCompleteStep()
		mos.stepIndicator.Objects = mos.buildStepIndicator(5).Objects
		mos.nextButton.Hide()
		mos.skipButton.Hide()
		mos.finishButton.Show()
	}

	mos.stepIndicator.Refresh()
	mos.contentArea.Refresh()
}

// showWelcomeStep shows the welcome screen with modern design
func (mos *ModernOnboardingScreen) showWelcomeStep() {
	loc := i18n.GetGlobalLocalizer()

	// Main title with large size
	titleLabel := widget.NewLabel(loc.Get("onboarding.welcome_message"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.Importance = widget.HighImportance

	// Subtitle
	subtitleLabel := widget.NewLabel(loc.Get("onboarding.welcome_subtitle"))
	subtitleLabel.Alignment = fyne.TextAlignCenter
	subtitleLabel.Wrapping = fyne.TextWrapWord

	// Feature highlights with colored backgrounds
	features := []struct {
		title       string
		description string
		bgColor     color.Color
	}{
		{
			loc.Get("onboarding.feature_encrypted"),
			loc.Get("onboarding.feature_encrypted_desc"),
			color.RGBA{R: 16, G: 185, B: 129, A: 255}, // Green
		},
		{
			loc.Get("onboarding.feature_decentralized"),
			loc.Get("onboarding.feature_decentralized_desc"),
			color.RGBA{R: 59, G: 130, B: 246, A: 255}, // Blue
		},
		{
			loc.Get("onboarding.feature_compatible"),
			loc.Get("onboarding.feature_compatible_desc"),
			color.RGBA{R: 147, G: 51, B: 234, A: 255}, // Purple
		},
		{
			loc.Get("onboarding.feature_p2p"),
			loc.Get("onboarding.feature_p2p_desc"),
			color.RGBA{R: 245, G: 158, B: 11, A: 255}, // Orange
		},
	}

	featureGrid := container.NewGridWithColumns(2)
	for _, f := range features {
		card := mos.buildHighlightCard(f.title, f.description, f.bgColor)
		featureGrid.Add(card)
	}

	// Call to action
	ctaLabel := widget.NewLabel(loc.Get("onboarding.cta_setup"))
	ctaLabel.Alignment = fyne.TextAlignCenter
	ctaLabel.Wrapping = fyne.TextWrapWord
	ctaLabel.TextStyle = fyne.TextStyle{Italic: true}

	mos.contentArea.Objects = []fyne.CanvasObject{
		container.NewVBox(
			titleLabel,
			subtitleLabel,
			featureGrid,
			ctaLabel,
		),
	}
}

// buildHighlightCard creates an attractive feature card with colored accent
func (mos *ModernOnboardingScreen) buildHighlightCard(title, description string, accentColor color.Color) fyne.CanvasObject {
	// Title with bold text
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Wrapping = fyne.TextWrapWord

	// Description text
	descLabel := widget.NewLabel(description)
	descLabel.Wrapping = fyne.TextWrapWord

	// Colored accent bar at top
	accentBar := canvas.NewRectangle(accentColor)
	accentBar.SetMinSize(fyne.NewSize(0, 4))

	content := container.NewVBox(
		titleLabel,
		descLabel,
	)

	cardContent := container.NewBorder(
		accentBar, // Top accent
		nil, nil, nil,
		container.NewPadded(content),
	)

	return NewSettingsContentCard(
		"",
		cardContent,
		color.RGBA{R: 248, G: 250, B: 252, A: 255}, // Very light gray background
	)
}

// showSetupModeStep shows the choice between new setup and restore
func (mos *ModernOnboardingScreen) showSetupModeStep() {
	loc := i18n.GetGlobalLocalizer()

	titleLabel := widget.NewLabel(loc.Get("onboarding.setup_mode"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.Importance = widget.HighImportance

	descLabel := widget.NewLabel(loc.Get("onboarding.setup_mode_desc"))
	descLabel.Alignment = fyne.TextAlignCenter
	descLabel.Wrapping = fyne.TextWrapWord

	// New Setup Option
	newSetupTitle := widget.NewLabel(loc.Get("onboarding.new_setup"))
	newSetupTitle.TextStyle = fyne.TextStyle{Bold: true}
	newSetupTitle.Alignment = fyne.TextAlignCenter

	newSetupDesc := widget.NewLabel(loc.Get("onboarding.new_setup_desc"))
	newSetupDesc.Wrapping = fyne.TextWrapWord
	newSetupDesc.Alignment = fyne.TextAlignCenter

	newSetupButton := widget.NewButton(loc.Get("onboarding.start_new_setup"), func() {
		mos.setupMode = "new"
		mos.nextButton.Enable()
		mos.goNext()
	})
	newSetupButton.Importance = widget.HighImportance

	newSetupContent := container.NewVBox(
		newSetupTitle,
		newSetupDesc,
		container.NewCenter(newSetupButton),
	)

	newSetupCard := NewSettingsContentCard(
		"",
		newSetupContent,
		color.RGBA{R: 16, G: 185, B: 129, A: 255}, // Green accent
	)

	// Restore from Backup Option
	restoreTitle := widget.NewLabel(loc.Get("onboarding.restore_backup"))
	restoreTitle.TextStyle = fyne.TextStyle{Bold: true}
	restoreTitle.Alignment = fyne.TextAlignCenter

	restoreDesc := widget.NewLabel(loc.Get("onboarding.restore_backup_desc"))
	restoreDesc.Wrapping = fyne.TextWrapWord
	restoreDesc.Alignment = fyne.TextAlignCenter

	restoreButton := widget.NewButton(loc.Get("onboarding.restore_from_backup"), func() {
		mos.setupMode = "restore"
		mos.nextButton.Enable()
		mos.goNext()
	})
	restoreButton.Importance = widget.MediumImportance

	restoreContent := container.NewVBox(
		restoreTitle,
		restoreDesc,
		container.NewCenter(restoreButton),
	)

	restoreCard := NewSettingsContentCard(
		"",
		restoreContent,
		color.RGBA{R: 245, G: 158, B: 11, A: 255}, // Orange accent
	)

	// Disable default next button - mode must be chosen via buttons
	mos.nextButton.Disable()

	// Equal-height grid layout
	cardsGrid := container.NewGridWithColumns(2,
		container.NewMax(newSetupCard),
		container.NewMax(restoreCard),
	)

	mos.contentArea.Objects = []fyne.CanvasObject{
		container.NewVBox(
			titleLabel,
			descLabel,
			cardsGrid,
		),
	}
}

// showPasswordStep shows the password setup screen
func (mos *ModernOnboardingScreen) showPasswordStep() {
	loc := i18n.GetGlobalLocalizer()

	titleLabel := widget.NewLabel(loc.Get("onboarding.password.title"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.Importance = widget.HighImportance

	descLabel := widget.NewLabel(loc.Get("onboarding.password.message"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.Alignment = fyne.TextAlignCenter

	// Password strength indicator
	strengthLabel := widget.NewLabel(fmt.Sprintf(loc.Get("onboarding.password_strength"), loc.Get("onboarding.password_strength_not_set")))
	strengthLabel.Alignment = fyne.TextAlignCenter

	// Password entry fields (without labels, use only placeholders)
	mos.passwordEntry = widget.NewPasswordEntry()
	mos.passwordEntry.SetPlaceHolder(loc.Get("onboarding.password.enter"))
	mos.passwordEntry.OnChanged = func(text string) {
		strength := mos.calculatePasswordStrength(text)
		strengthLabel.SetText(fmt.Sprintf(loc.Get("onboarding.password_strength"), strength))
		strengthLabel.Refresh()
	}

	mos.passwordConfirmEntry = widget.NewPasswordEntry()
	mos.passwordConfirmEntry.SetPlaceHolder(loc.Get("onboarding.password.confirm"))

	// Compact security tips - only 2 most important
	tips := []string{
		loc.Get("onboarding.tip_12chars"),
		loc.Get("onboarding.tip_mix"),
	}

	tipsList := container.NewVBox()
	for _, tip := range tips {
		bullet := widget.NewLabel("• " + tip)
		bullet.Wrapping = fyne.TextWrapWord
		tipsList.Add(bullet)
	}

	tipsCard := NewSettingsContentCard(
		"",
		tipsList,
		color.RGBA{R: 59, G: 130, B: 246, A: 255}, // Blue accent
	)

	form := container.NewVBox(
		mos.passwordEntry,
		strengthLabel,
		mos.passwordConfirmEntry,
	)

	formCard := NewSettingsContentCard(
		"",
		form,
		color.RGBA{R: 248, G: 250, B: 252, A: 255}, // Light background
	)

	mos.contentArea.Objects = []fyne.CanvasObject{
		container.NewVBox(
			titleLabel,
			descLabel,
			formCard,
			tipsCard,
		),
	}

	// Re-enable next button
	mos.nextButton.Enable()
	mos.nextButton.Show()
}

// showRestoreStep shows the restore from backup screen
func (mos *ModernOnboardingScreen) showRestoreStep() {
	loc := i18n.GetGlobalLocalizer()

	titleLabel := widget.NewLabel(loc.Get("onboarding.restore_backup"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.Importance = widget.HighImportance

	descLabel := widget.NewLabel(loc.Get("backup.restore_description"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.Alignment = fyne.TextAlignCenter

	// File selection
	fileLabel := widget.NewLabel(loc.Get("onboarding.restore_file"))
	fileLabel.TextStyle = fyne.TextStyle{Bold: true}

	mos.restoreFileLabel = widget.NewLabel(loc.Get("onboarding.restore_no_file"))
	mos.restoreFileLabel.Wrapping = fyne.TextWrapWord
	mos.restoreFileLabel.Alignment = fyne.TextAlignCenter

	chooseButton := widget.NewButton(loc.Get("onboarding.restore_choose"), func() {
		mos.chooseBackupFile()
	})
	chooseButton.Importance = widget.MediumImportance

	fileSection := container.NewVBox(
		fileLabel,
		mos.restoreFileLabel,
		container.NewCenter(chooseButton),
	)

	fileCard := NewSettingsContentCard(
		"",
		fileSection,
		color.RGBA{R: 248, G: 250, B: 252, A: 255},
	)

	// Password field
	passwordLabel := widget.NewLabel(loc.Get("onboarding.restore_password"))
	passwordLabel.TextStyle = fyne.TextStyle{Bold: true}

	mos.restorePasswordEntry = widget.NewPasswordEntry()
	mos.restorePasswordEntry.SetPlaceHolder(loc.Get("onboarding.restore_enter_password"))

	passwordSection := container.NewVBox(
		passwordLabel,
		mos.restorePasswordEntry,
	)

	passwordCard := NewSettingsContentCard(
		"",
		passwordSection,
		color.RGBA{R: 248, G: 250, B: 252, A: 255},
	)

	// Important notice
	warningLabel := widget.NewLabel(loc.Get("onboarding.restore_warning"))
	warningLabel.Wrapping = fyne.TextWrapWord
	warningLabel.Alignment = fyne.TextAlignCenter
	warningLabel.TextStyle = fyne.TextStyle{Bold: true}

	warningCard := NewSettingsContentCard(
		"",
		container.NewVBox(warningLabel),
		color.RGBA{R: 245, G: 158, B: 11, A: 255}, // Orange accent for warning
	)

	mos.contentArea.Objects = []fyne.CanvasObject{
		container.NewVBox(
			titleLabel,
			descLabel,
			fileCard,
			passwordCard,
			warningCard,
		),
	}
}

// showPeerStep shows the peer configuration screen
func (mos *ModernOnboardingScreen) showPeerStep() {
	loc := i18n.GetGlobalLocalizer()

	titleLabel := widget.NewLabel(loc.Get("onboarding.network_config"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.Importance = widget.HighImportance

	descLabel := widget.NewLabel(loc.Get("onboarding.network_config_desc"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.Alignment = fyne.TextAlignCenter

	// Get default peers
	defaultPeers := strings.Join(core.DefaultPeers, "\n")

	// Peers text area
	peersLabel := widget.NewLabel(loc.Get("onboarding.network_peers"))
	peersLabel.TextStyle = fyne.TextStyle{Bold: true}

	mos.peersEntry = widget.NewMultiLineEntry()
	mos.peersEntry.SetPlaceHolder(loc.Get("onboarding.peers_placeholder_multi"))
	mos.peersEntry.SetText(defaultPeers)
	mos.peersEntry.SetMinRowsVisible(5)

	peersSection := container.NewVBox(
		peersLabel,
		mos.peersEntry,
	)

	peersCard := NewSettingsContentCard(
		"",
		peersSection,
		color.RGBA{R: 248, G: 250, B: 252, A: 255},
	)

	// Info about peers
	infoLabel := widget.NewLabel(loc.Get("onboarding.peers_info"))
	infoLabel.Wrapping = fyne.TextWrapWord
	infoLabel.Alignment = fyne.TextAlignCenter

	infoCard := NewSettingsContentCard(
		"",
		container.NewVBox(infoLabel),
		color.RGBA{R: 99, G: 102, B: 241, A: 255}, // Indigo accent
	)

	mos.contentArea.Objects = []fyne.CanvasObject{
		container.NewVBox(
			titleLabel,
			descLabel,
			peersCard,
			infoCard,
		),
	}
}

// showCompleteStep shows the completion screen
func (mos *ModernOnboardingScreen) showCompleteStep() {
	loc := i18n.GetGlobalLocalizer()

	titleLabel := widget.NewLabel(loc.Get("onboarding.complete.title"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.Importance = widget.HighImportance

	descLabel := widget.NewLabel(loc.Get("onboarding.complete_subtitle"))
	descLabel.Alignment = fyne.TextAlignCenter
	descLabel.Wrapping = fyne.TextWrapWord

	// Success message in green card
	successMessage := widget.NewLabel(loc.Get("onboarding.complete_ready"))
	successMessage.Alignment = fyne.TextAlignCenter
	successMessage.TextStyle = fyne.TextStyle{Bold: true}

	successCard := NewSettingsContentCard(
		"",
		container.NewVBox(successMessage),
		color.RGBA{R: 16, G: 185, B: 129, A: 255}, // Green
	)

	// Next steps
	nextStepsTitle := widget.NewLabel(loc.Get("onboarding.whats_next"))
	nextStepsTitle.TextStyle = fyne.TextStyle{Bold: true}

	steps := []string{
		loc.Get("onboarding.next_view_address"),
		loc.Get("onboarding.next_connect_client"),
		loc.Get("onboarding.next_manage_peers"),
		loc.Get("onboarding.next_create_backups"),
	}

	stepsList := container.NewVBox()
	for i, s := range steps {
		stepLabel := widget.NewLabel(fmt.Sprintf("%d. %s", i+1, s))
		stepLabel.Wrapping = fyne.TextWrapWord
		stepsList.Add(stepLabel)
	}

	stepsCard := NewSettingsContentCard(
		"",
		container.NewVBox(nextStepsTitle, stepsList),
		color.RGBA{R: 248, G: 250, B: 252, A: 255},
	)

	mos.contentArea.Objects = []fyne.CanvasObject{
		container.NewVBox(
			titleLabel,
			descLabel,
			successCard,
			stepsCard,
		),
	}
}

// calculatePasswordStrength calculates password strength
func (mos *ModernOnboardingScreen) calculatePasswordStrength(password string) string {
	loc := i18n.GetGlobalLocalizer()

	if len(password) == 0 {
		return loc.Get("onboarding.password_strength_not_set")
	}
	if len(password) < 6 {
		return loc.Get("onboarding.password_strength_weak")
	}
	if len(password) < 8 {
		return loc.Get("onboarding.password_strength_weak")
	}
	if len(password) < 12 {
		return loc.Get("onboarding.password_strength_moderate")
	}

	// Check for variety
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false

	for _, ch := range password {
		if ch >= 'A' && ch <= 'Z' {
			hasUpper = true
		} else if ch >= 'a' && ch <= 'z' {
			hasLower = true
		} else if ch >= '0' && ch <= '9' {
			hasDigit = true
		} else {
			hasSymbol = true
		}
	}

	variety := 0
	if hasUpper {
		variety++
	}
	if hasLower {
		variety++
	}
	if hasDigit {
		variety++
	}
	if hasSymbol {
		variety++
	}

	if len(password) >= 16 && variety >= 3 {
		return loc.Get("onboarding.password_strength_very_strong")
	}
	if len(password) >= 12 && variety >= 3 {
		return loc.Get("onboarding.password_strength_strong")
	}
	return loc.Get("onboarding.password_strength_moderate")
}

// goBack navigates to the previous step
func (mos *ModernOnboardingScreen) goBack() {
	if mos.currentStep > 0 {
		// Handle special navigation for restore flow
		if mos.setupMode == "restore" && mos.currentStep == 2 {
			// Go back to mode selection and reset mode
			mos.setupMode = ""
			mos.showStep(1)
		} else {
			mos.showStep(mos.currentStep - 1)
		}
	}
}

// goNext validates current step and navigates to next
func (mos *ModernOnboardingScreen) goNext() {
	// Validate current step
	if !mos.validateCurrentStep() {
		return
	}

	// Move to next step based on mode
	if mos.setupMode == "restore" && mos.currentStep == 1 {
		// Skip directly to restore step
		mos.showStep(2)
	} else if mos.setupMode == "new" && mos.currentStep == 1 {
		// Go to password step for new setup
		mos.showStep(2)
	} else if mos.currentStep < 4 {
		mos.showStep(mos.currentStep + 1)
	}
}

// skipStep skips the current optional step
func (mos *ModernOnboardingScreen) skipStep() {
	// Only peer step is skippable
	if mos.currentStep == 3 {
		mos.showStep(4)
	}
}

// validateCurrentStep validates the current step data
func (mos *ModernOnboardingScreen) validateCurrentStep() bool {
	loc := i18n.GetGlobalLocalizer()

	switch mos.currentStep {
	case 0:
		// Welcome step - no validation needed
		return true

	case 1:
		// Setup mode step - mode must be chosen (handled by buttons)
		return mos.setupMode != ""

	case 2:
		if mos.setupMode == "new" {
			// Password step - validate password
			password := mos.passwordEntry.Text
			confirm := mos.passwordConfirmEntry.Text

			if len(password) < 6 {
				dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("password.error.length"), loc.Get("onboarding.error.password"))
				return false
			}

			if password != confirm {
				dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("password.error.mismatch"), loc.Get("onboarding.error.mismatch"))
				return false
			}
		} else if mos.setupMode == "restore" {
			// Restore step - validate file and password
			if mos.selectedRestoreFile == "" {
				dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.no_file_selected"), loc.Get("backup.select_file_first"))
				return false
			}

			if mos.restorePasswordEntry.Text == "" {
				dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.password_required"), loc.Get("backup.enter_password"))
				return false
			}
		}
		return true

	case 3:
		// Peer step - validate peer addresses
		peerText := strings.TrimSpace(mos.peersEntry.Text)
		if peerText == "" {
			dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("peers.error.add"), loc.Get("onboarding.error.peers"))
			return false
		}

		// Basic validation of peer format
		peers := strings.Split(peerText, "\n")
		for _, peer := range peers {
			peer = strings.TrimSpace(peer)
			if peer == "" {
				continue
			}

			if !strings.HasPrefix(peer, "tcp://") &&
				!strings.HasPrefix(peer, "tls://") &&
				!strings.HasPrefix(peer, "quic://") {
				dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("peers.error.invalid"),
					fmt.Sprintf(loc.Get("onboarding.error_invalid_peer_format"), peer))
				return false
			}
		}

		return true

	case 4:
		// Complete step - no validation needed
		return true
	}

	return true
}

// chooseBackupFile shows a file picker for restore
func (mos *ModernOnboardingScreen) chooseBackupFile() {
	loc := i18n.GetGlobalLocalizer()

	openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.file_error"), fmt.Sprintf(loc.Get("backup.file_open_failed"), err))
			return
		}
		if reader == nil {
			return // User cancelled
		}
		defer reader.Close()

		// Store file path
		mos.selectedRestoreFile = reader.URI().Path()
		mos.restoreFileLabel.SetText(reader.URI().Name())

		log.Printf("Selected backup file: %s", mos.selectedRestoreFile)
	}, mos.app.GetMainWindow())

	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".tyrbackup"}))
	openDialog.Show()
}

// finish completes the onboarding process
func (mos *ModernOnboardingScreen) finish() {
	if mos.setupMode == "restore" {
		mos.finishRestore()
	} else {
		mos.finishNewSetup()
	}
}

// finishNewSetup completes the new setup flow
func (mos *ModernOnboardingScreen) finishNewSetup() {
	loc := i18n.GetGlobalLocalizer()
	config := mos.app.GetConfig()

	// Save password
	password := mos.passwordEntry.Text
	if err := config.SetPassword(password); err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("password.error.update"), fmt.Sprintf(loc.Get("error.failed_to_save_password"), err))
		return
	}

	// Clear existing peers and add new ones
	config.NetworkPeers = []core.PeerConfig{}

	peerText := strings.TrimSpace(mos.peersEntry.Text)
	peers := strings.Split(peerText, "\n")
	for _, peer := range peers {
		peer = strings.TrimSpace(peer)
		if peer == "" {
			continue
		}

		if err := config.AddPeer(peer); err != nil {
			log.Printf("Warning: Failed to add peer %s: %v", peer, err)
			continue
		}
	}

	// Mark onboarding as complete
	config.OnboardingComplete = true

	// Enable autostart by default after first run
	config.UIPreferences.AutoStart = true
	if err := core.EnableAutoStart(); err != nil {
		log.Printf("Warning: Failed to enable autostart: %v", err)
		// Continue anyway - not critical
	} else {
		log.Println("Autostart enabled by default after first run")
	}

	// Save configuration
	if err := config.Save(); err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("error.config"), fmt.Sprintf(loc.Get("error.failed_to_save_configuration"), err))
		return
	}

	log.Println("Onboarding completed successfully")

	// Initialize and start service manager
	serviceManager := mos.app.GetServiceManager()
	if err := serviceManager.Initialize(); err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("onboarding.error.init"), fmt.Sprintf(loc.Get("error.failed_to_initialize_service"), err))
		return
	}

	// Auto-start the service
	if err := serviceManager.Start(); err != nil {
		log.Printf("Warning: Failed to auto-start service: %v", err)
		// Don't show error - user can start manually from dashboard
	}

	// Navigate to dashboard
	mos.app.ShowDashboard()

	// Show success message
	dialogs.ShowSuccess(mos.app.GetMainWindow(), loc.Get("onboarding.welcome_message"), loc.Get("onboarding.setup_complete_msg"))
}

// finishRestore completes the restore flow
func (mos *ModernOnboardingScreen) finishRestore() {
	loc := i18n.GetGlobalLocalizer()
	password := mos.restorePasswordEntry.Text

	// Read backup file
	backupData, err := core.ReadBackupFile(mos.selectedRestoreFile)
	if err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.read_failed"), fmt.Sprintf(loc.Get("backup.read_failed_msg"), err))
		return
	}

	// Restore backup
	restoredConfig, dbData, err := core.RestoreBackup(backupData, password)
	if err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.restore_failed"), fmt.Sprintf(loc.Get("backup.restore_failed_msg"), err))
		return
	}

	// Mark onboarding as complete
	restoredConfig.OnboardingComplete = true

	// Enable autostart by default after first run (unless already configured in backup)
	if !restoredConfig.UIPreferences.AutoStart {
		restoredConfig.UIPreferences.AutoStart = true
		if err := core.EnableAutoStart(); err != nil {
			log.Printf("Warning: Failed to enable autostart: %v", err)
			// Continue anyway - not critical
		} else {
			log.Println("Autostart enabled by default after restore")
		}
	}

	// Save restored configuration
	if err := restoredConfig.Save(); err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("backup.save_failed_msg"), err))
		return
	}

	// Restore database if included
	if dbData != nil && len(dbData) > 0 {
		dbPath := restoredConfig.ServiceSettings.DatabasePath
		if err := core.WriteBackupFile(dbPath, dbData); err != nil {
			log.Printf("Warning: Failed to restore database: %v", err)
			dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("backup.partial_restore"), loc.Get("backup.partial_restore_msg"))
			return
		}
	}

	log.Println("Backup restored successfully during onboarding")

	// Reload the config in the app
	appConfig := mos.app.GetConfig()
	*appConfig = *restoredConfig

	// Initialize and start service manager
	serviceManager := mos.app.GetServiceManager()
	if err := serviceManager.Initialize(); err != nil {
		dialogs.ShowError(mos.app.GetMainWindow(), loc.Get("onboarding.error.init"), fmt.Sprintf(loc.Get("error.failed_to_initialize_service"), err))
		return
	}

	// Auto-start the service
	if err := serviceManager.Start(); err != nil {
		log.Printf("Warning: Failed to auto-start service: %v", err)
		// Don't show error - user can start manually
	}

	// Navigate to dashboard
	mos.app.ShowDashboard()

	// Show success message
	dialogs.ShowSuccess(mos.app.GetMainWindow(), loc.Get("backup.restore_complete"), loc.Get("onboarding.restore_complete_msg"))
}
