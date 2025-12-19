package windows

import (
	"fmt"
	"image/color"
	"log"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/dialogs"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/theme"
)

// ModernSettingsScreen represents the main settings menu with navigation cards
type ModernSettingsScreen struct {
	app AppInterface
}

// SettingsMenuCard represents a clickable menu card with hover effects
type SettingsMenuCard struct {
	widget.BaseWidget
	title       string
	description string
	accentColor color.Color
	onTapped    func()
	container   *fyne.Container
	background  *canvas.Rectangle
	hovered     bool
}

// SettingsContentCard represents a modern content card for settings pages
type SettingsContentCard struct {
	widget.BaseWidget
	title       string
	content     fyne.CanvasObject
	accentColor color.Color
	container   *fyne.Container
	background  *canvas.Rectangle
	hovered     bool
}

// NewModernSettingsScreen creates the main settings menu
func NewModernSettingsScreen(app AppInterface) fyne.CanvasObject {
	screen := &ModernSettingsScreen{
		app: app,
	}

	return screen.buildMenuLayout()
}

// buildMenuLayout creates the main settings menu with navigation cards
func (mss *ModernSettingsScreen) buildMenuLayout() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Header with Back button
	header := buildSettingsHeaderModern(SettingsHeaderConfig{
		Title:    loc.Get("settings.title"),
		Subtitle: loc.Get("settings.subtitle"),
		OnBack: func() {
			mss.app.ShowDashboard()
		},
	})

	// Color palette
	primaryColor := color.RGBA{R: 99, G: 102, B: 241, A: 255}
	infoColor := color.RGBA{R: 59, G: 130, B: 246, A: 255}
	warningColor := color.RGBA{R: 245, G: 158, B: 11, A: 255}
	successColor := color.RGBA{R: 16, G: 185, B: 129, A: 255}
	accentColor := color.RGBA{R: 168, G: 85, B: 247, A: 255}

	// Create menu cards
	generalCard := mss.createMenuCard(
		loc.Get("settings.general.title"),
		loc.Get("settings.general.desc"),
		primaryColor,
		func() { mss.showGeneralSettings() },
	)

	networkCard := mss.createMenuCard(
		loc.Get("settings.network.title"),
		loc.Get("settings.network.desc"),
		infoColor,
		func() { mss.showNetworkSettings() },
	)

	securityCard := mss.createMenuCard(
		loc.Get("settings.security.title"),
		loc.Get("settings.security.desc"),
		warningColor,
		func() { mss.showSecuritySettings() },
	)

	backupCard := mss.createMenuCard(
		loc.Get("settings.backup.title"),
		loc.Get("settings.backup.desc"),
		successColor,
		func() { mss.showBackupSettings() },
	)

	aboutCard := mss.createMenuCard(
		loc.Get("settings.about.title"),
		loc.Get("settings.about.desc"),
		accentColor,
		func() { mss.showAboutSettings() },
	)

	// Menu grid
	menuContent := container.NewVBox(
		generalCard,
		networkCard,
		securityCard,
		backupCard,
		aboutCard,
	)

	scrollContent := container.NewScroll(menuContent)

	return container.NewBorder(
		header,
		nil, nil, nil,
		scrollContent,
	)
}

// createMenuCard creates a clickable navigation card
func (mss *ModernSettingsScreen) createMenuCard(title, description string, accentColor color.Color, onTapped func()) fyne.CanvasObject {
	card := &SettingsMenuCard{
		title:       title,
		description: description,
		accentColor: accentColor,
		onTapped:    onTapped,
	}
	card.ExtendBaseWidget(card)
	return card
}

// showGeneralSettings displays the general settings page
func (mss *ModernSettingsScreen) showGeneralSettings() {
	content := mss.buildGeneralPage()
	mss.app.GetMainWindow().SetContent(content)
}

// showNetworkSettings displays the network settings page
func (mss *ModernSettingsScreen) showNetworkSettings() {
	content := mss.buildNetworkPage()
	mss.app.GetMainWindow().SetContent(content)
}

// showSecuritySettings displays the security settings page
func (mss *ModernSettingsScreen) showSecuritySettings() {
	content := mss.buildSecurityPage()
	mss.app.GetMainWindow().SetContent(content)
}

// showBackupSettings displays the backup settings page
func (mss *ModernSettingsScreen) showBackupSettings() {
	content := mss.buildBackupPage()
	mss.app.GetMainWindow().SetContent(content)
}

// showAboutSettings displays the about page
func (mss *ModernSettingsScreen) showAboutSettings() {
	content := mss.buildAboutPage()
	mss.app.GetMainWindow().SetContent(content)
}

// buildGeneralPage creates the general settings page with modern card layout
func (mss *ModernSettingsScreen) buildGeneralPage() fyne.CanvasObject {
	config := mss.app.GetConfig()
	loc := i18n.GetGlobalLocalizer()

	// Log current settings for debugging
	log.Printf("Building General Preferences page - Current theme: %s, Current language: %s",
		config.UIPreferences.Theme, config.UIPreferences.Language)

	// Header with Back button
	header := buildSettingsHeaderModern(SettingsHeaderConfig{
		Title:    loc.Get("settings.general.title"),
		Subtitle: loc.Get("settings.general.subtitle"),
		OnBack: func() {
			mss.app.ShowSettings()
		},
	})

	// Auto-start Card
	enabled, err := core.IsAutoStartEnabled()
	if err != nil {
		fmt.Printf("Warning: Failed to check auto-start status: %v\n", err)
		enabled = config.UIPreferences.AutoStart
	}

	autoStartCheck := widget.NewCheck(loc.Get("settings.autostart"), func(checked bool) {
		mss.handleAutoStartToggle(checked)
	})
	autoStartCheck.SetChecked(enabled)

	autoStartDesc := widget.NewLabel(loc.Get("settings.autostart.desc"))
	autoStartDesc.Wrapping = fyne.TextWrapWord
	autoStartDesc.TextStyle = fyne.TextStyle{Italic: true}

	autoStartContent := container.NewVBox(
		autoStartCheck,
		autoStartDesc,
	)

	autoStartCard := NewSettingsContentCard(
		loc.Get("settings.autostart.title"),
		autoStartContent,
		color.RGBA{R: 99, G: 102, B: 241, A: 255},
	)

	// Language Card
	languageLabel := widget.NewLabel(loc.Get("settings.language_label"))
	languageLabel.Wrapping = fyne.TextWrapWord

	currentLang := loc.Get("settings.language.select")
	if config.UIPreferences.Language == "ru" {
		currentLang = loc.Get("settings.language.select")
	}

	// Store original value to detect actual changes
	originalLang := config.UIPreferences.Language

	languageSelect := widget.NewSelect([]string{"English", "Русский (Russian)"}, func(selected string) {
		lang := "en"
		if selected == "Русский (Russian)" {
			lang = "ru"
		}

		// Only save and notify if value actually changed
		if lang == originalLang {
			log.Printf("Language selection unchanged (%s), skipping save", lang)
			return
		}

		log.Printf("Language changed from %s to %s, saving...", originalLang, lang)
		config.UIPreferences.Language = lang
		if err := config.Save(); err != nil {
			log.Printf("Failed to save language preference: %v", err)
			// Use old localizer since we haven't changed yet
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.error.title"), fmt.Sprintf(loc.Get("settings.error.save_failed"), err))
			return
		}

		// Apply language to global localizer
		localizer := i18n.GetGlobalLocalizer()
		if err := localizer.SetLanguage(lang); err != nil {
			log.Printf("Failed to apply language: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.error.language"), fmt.Sprintf(loc.Get("settings.error.apply_language"), err))
			return
		}

		// Update original value after successful save
		originalLang = lang
		log.Printf("Language preference saved and applied successfully")

		// Update system tray menu with new language
		mss.app.UpdateSystemTrayMenu()

		// Instant reload: Refresh the entire settings page with new language
		mss.app.GetMainWindow().SetContent(mss.buildGeneralPage())
		log.Printf("Settings page refreshed with new language: %s", lang)
	})
	languageSelect.SetSelected(currentLang)

	languageDesc := widget.NewLabel(loc.Get("settings.language.desc"))
	languageDesc.Wrapping = fyne.TextWrapWord
	languageDesc.TextStyle = fyne.TextStyle{Italic: true}

	languageContent := container.NewVBox(
		languageLabel,
		languageSelect,
		languageDesc,
	)

	languageCard := NewSettingsContentCard(
		loc.Get("settings.language.card"),
		languageContent,
		color.RGBA{R: 16, G: 185, B: 129, A: 255},
	)

	// Theme Card
	themeLabel := widget.NewLabel(loc.Get("settings.theme_label"))
	themeLabel.Wrapping = fyne.TextWrapWord

	currentTheme := loc.Get("settings.theme_system")
	switch config.UIPreferences.Theme {
	case "light":
		currentTheme = loc.Get("settings.theme_light")
	case "dark":
		currentTheme = loc.Get("settings.theme_dark")
	}

	// Store original value to detect actual changes
	originalTheme := config.UIPreferences.Theme

	themeOptions := []string{loc.Get("settings.theme_light"), loc.Get("settings.theme_dark"), loc.Get("settings.theme_system")}
	themeSelect := widget.NewSelect(themeOptions, func(selected string) {
		themeName := "system"
		// Compare with localized strings
		if selected == loc.Get("settings.theme_light") {
			themeName = "light"
		} else if selected == loc.Get("settings.theme_dark") {
			themeName = "dark"
		}

		// Only save and notify if value actually changed
		if themeName == originalTheme {
			log.Printf("Theme selection unchanged (%s), skipping save", themeName)
			return
		}

		log.Printf("Theme changed from %s to %s, saving...", originalTheme, themeName)
		config.UIPreferences.Theme = themeName
		if err := config.Save(); err != nil {
			log.Printf("Failed to save theme preference: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.error.title"), fmt.Sprintf(loc.Get("settings.error.save_failed"), err))
			return
		}

		if err := theme.ApplyTheme(mss.app.GetFyneApp(), themeName); err != nil {
			log.Printf("Failed to apply theme: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.error.theme"), fmt.Sprintf(loc.Get("settings.error.apply_theme"), err))
			return
		}

		// Update original value after successful save and apply
		originalTheme = themeName
		log.Printf("Theme preference saved and applied successfully")
	})
	themeSelect.SetSelected(currentTheme)

	themeDesc := widget.NewLabel(loc.Get("settings.theme.desc"))
	themeDesc.Wrapping = fyne.TextWrapWord
	themeDesc.TextStyle = fyne.TextStyle{Italic: true}

	themeContent := container.NewVBox(
		themeLabel,
		themeSelect,
		themeDesc,
	)

	themeCard := NewSettingsContentCard(
		loc.Get("settings.theme.card"),
		themeContent,
		color.RGBA{R: 168, G: 85, B: 247, A: 255},
	)

	// Main content
	mainContent := container.NewVBox(
		autoStartCard,
		languageCard,
		themeCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		header,
		nil, nil, nil,
		scrollContent,
	)
}

// buildNetworkPage creates the network settings page with modern card layout
func (mss *ModernSettingsScreen) buildNetworkPage() fyne.CanvasObject {
	config := mss.app.GetConfig()
	loc := i18n.GetGlobalLocalizer()

	// Header with Back button
	header := buildSettingsHeaderModern(SettingsHeaderConfig{
		Title:    loc.Get("settings.network.title"),
		Subtitle: loc.Get("settings.network.subtitle"),
		OnBack: func() {
			mss.app.ShowSettings()
		},
	})

	// Add Peer Card
	addPeerEntry := widget.NewEntry()
	addPeerEntry.SetPlaceHolder(loc.Get("settings.network.add_peer.placeholder"))

	addButton := widget.NewButton(loc.Get("settings.network.add_peer.button"), func() {
		address := addPeerEntry.Text

		if address == "" {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.network.error.invalid"), loc.Get("settings.network.error.empty"))
			return
		}

		config := mss.app.GetConfig()
		if err := config.AddPeer(address); err != nil {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.network.error.add_failed"), fmt.Sprintf(loc.Get("error.failed_to_add_peer"), err))
			return
		}

		// Save configuration
		if err := config.Save(); err != nil {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("error.failed_to_save_configuration"), err))
			return
		}

		// Clear entry and refresh page
		addPeerEntry.SetText("")
		mss.app.GetMainWindow().SetContent(mss.buildNetworkPage())

		dialogs.ShowSuccess(mss.app.GetMainWindow(), loc.Get("settings.network.success.added"), loc.Get("settings.network.success.added_msg"))
		log.Printf("Added peer: %s", address)
	})
	addButton.Importance = widget.HighImportance

	addPeerDesc := widget.NewLabel(loc.Get("settings.network.add_peer.desc"))
	addPeerDesc.Wrapping = fyne.TextWrapWord
	addPeerDesc.TextStyle = fyne.TextStyle{Italic: true}

	addPeerContent := container.NewVBox(
		widget.NewLabel(loc.Get("settings.network.add_peer.label")),
		addPeerEntry,
		addButton,
		addPeerDesc,
	)

	addPeerCard := NewSettingsContentCard(
		loc.Get("settings.network.add_peer.card"),
		addPeerContent,
		color.RGBA{R: 16, G: 185, B: 129, A: 255},
	)

	// Configured Peers List
	peersContainer := container.NewVBox()

	config = mss.app.GetConfig()
	serviceManager := mss.app.GetServiceManager()
	peerStats := serviceManager.GetPeerStats()

	if len(config.NetworkPeers) == 0 {
		emptyLabel := widget.NewLabel(loc.Get("settings.network.peers_list.empty"))
		emptyLabel.Alignment = fyne.TextAlignCenter
		emptyLabel.TextStyle = fyne.TextStyle{Italic: true}

		peersContainer.Objects = append(peersContainer.Objects, container.NewPadded(
			container.NewCenter(emptyLabel),
		))
	} else {
		for _, peer := range config.NetworkPeers {
			// Check if peer is connected
			isConnected := false
			for _, stat := range peerStats {
				if stat.Address == peer.Address && stat.Status {
					isConnected = true
					break
				}
			}

			peerCard := mss.buildPeerCard(peer, isConnected)
			peersContainer.Objects = append(peersContainer.Objects, peerCard)
		}
	}

	peersListLabel := widget.NewLabel(loc.Get("settings.network.peers_list.desc"))
	peersListLabel.Wrapping = fyne.TextWrapWord
	peersListLabel.TextStyle = fyne.TextStyle{Italic: true}

	peersListContent := container.NewVBox(
		peersListLabel,
		widget.NewSeparator(),
		peersContainer,
	)

	peersListCard := NewSettingsContentCard(
		loc.Get("settings.network.peers_list.title"),
		peersListContent,
		color.RGBA{R: 59, G: 130, B: 246, A: 255},
	)

	// Main content
	mainContent := container.NewVBox(
		addPeerCard,
		peersListCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		header,
		nil, nil, nil,
		scrollContent,
	)
}

// buildSecurityPage creates the security settings page with modern card layout
func (mss *ModernSettingsScreen) buildSecurityPage() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Header with Back button
	header := buildSettingsHeaderModern(SettingsHeaderConfig{
		Title:    loc.Get("settings.security.title"),
		Subtitle: loc.Get("settings.security.subtitle"),
		OnBack: func() {
			mss.app.ShowSettings()
		},
	})

	// Password Card
	changePasswordButton := widget.NewButton(loc.Get("settings.security.password.button"), func() {
		mss.showChangePasswordDialog()
	})
	changePasswordButton.Importance = widget.HighImportance

	passwordDesc := widget.NewLabel(loc.Get("settings.security.password.desc"))
	passwordDesc.Wrapping = fyne.TextWrapWord
	passwordDesc.TextStyle = fyne.TextStyle{Italic: true}

	passwordRequirements := widget.NewLabel(loc.Get("settings.security.password.requirements"))
	passwordRequirements.Wrapping = fyne.TextWrapWord

	passwordContent := container.NewVBox(
		changePasswordButton,
		passwordDesc,
		widget.NewSeparator(),
		passwordRequirements,
	)

	passwordCard := NewSettingsContentCard(
		loc.Get("settings.security.password.card"),
		passwordContent,
		color.RGBA{R: 245, G: 158, B: 11, A: 255},
	)

	// Regenerate Keys Card
	regenerateKeysButton := widget.NewButton(loc.Get("settings.security.regenerate_keys.button"), func() {
		mss.showRegenerateKeysDialog()
	})
	regenerateKeysButton.Importance = widget.DangerImportance

	regenerateKeysDesc := widget.NewLabel(loc.Get("settings.security.regenerate_keys.desc"))
	regenerateKeysDesc.Wrapping = fyne.TextWrapWord
	regenerateKeysDesc.TextStyle = fyne.TextStyle{Italic: true}

	regenerateKeysContent := container.NewVBox(
		regenerateKeysButton,
		regenerateKeysDesc,
	)

	regenerateKeysCard := NewSettingsContentCard(
		loc.Get("settings.security.regenerate_keys.card"),
		regenerateKeysContent,
		color.RGBA{R: 220, G: 38, B: 38, A: 255}, // Red accent for dangerous operation
	)

	// Main content
	mainContent := container.NewVBox(
		passwordCard,
		regenerateKeysCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		header,
		nil, nil, nil,
		scrollContent,
	)
}

// buildBackupPage creates the backup settings page with modern card layout
func (mss *ModernSettingsScreen) buildBackupPage() fyne.CanvasObject {
	config := mss.app.GetConfig()
	loc := i18n.GetGlobalLocalizer()

	// Header with Back button
	header := buildSettingsHeaderModern(SettingsHeaderConfig{
		Title:    loc.Get("settings.backup.title"),
		Subtitle: loc.Get("settings.backup.subtitle"),
		OnBack: func() {
			mss.app.ShowSettings()
		},
	})

	// Database Location Card
	dbLocationLabel := widget.NewLabel(loc.Get("settings.backup.db.label"))
	dbLocationLabel.TextStyle = fyne.TextStyle{Bold: true}

	dbValue := widget.NewLabel(config.ServiceSettings.DatabasePath)
	dbValue.TextStyle = fyne.TextStyle{Monospace: true}
	dbValue.Wrapping = fyne.TextWrapWord

	dbDesc := widget.NewLabel(loc.Get("settings.backup.db.desc"))
	dbDesc.Wrapping = fyne.TextWrapWord
	dbDesc.TextStyle = fyne.TextStyle{Italic: true}

	dbContent := container.NewVBox(
		dbLocationLabel,
		dbValue,
		dbDesc,
	)

	dbCard := NewSettingsContentCard(
		loc.Get("settings.backup.db.title"),
		dbContent,
		color.RGBA{R: 59, G: 130, B: 246, A: 255},
	)

	// Create Backup Card
	createBackupCard := mss.buildCreateBackupCard()

	// Restore Backup Card
	restoreBackupCard := mss.buildRestoreBackupCard()

	// Main content
	mainContent := container.NewVBox(
		dbCard,
		createBackupCard,
		restoreBackupCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		header,
		nil, nil, nil,
		scrollContent,
	)
}

// buildCreateBackupCard creates the backup creation card
func (mss *ModernSettingsScreen) buildCreateBackupCard() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	descLabel := widget.NewLabel(loc.Get("settings.backup.create.desc"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Password fields
	passwordLabel := widget.NewLabelWithStyle(loc.Get("settings.backup.create.password_label"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	backupPasswordEntry := widget.NewPasswordEntry()
	backupPasswordEntry.SetPlaceHolder(loc.Get("settings.backup.create.password_placeholder"))

	confirmLabel := widget.NewLabelWithStyle(loc.Get("settings.backup.create.confirm_label"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	backupPasswordConfirmEntry := widget.NewPasswordEntry()
	backupPasswordConfirmEntry.SetPlaceHolder(loc.Get("settings.backup.create.confirm_placeholder"))

	// Include database checkbox
	includeDBCheck := widget.NewCheck(loc.Get("settings.backup.create.include_db"), nil)
	includeDBCheck.SetChecked(true)

	dbCheckDesc := widget.NewLabel(loc.Get("settings.backup.create.include_db_desc"))
	dbCheckDesc.Wrapping = fyne.TextWrapWord
	dbCheckDesc.TextStyle = fyne.TextStyle{Italic: true}

	// Create button
	createButton := widget.NewButton(loc.Get("settings.backup.create.button"), func() {
		mss.createBackup(backupPasswordEntry.Text, backupPasswordConfirmEntry.Text, includeDBCheck.Checked, backupPasswordEntry, backupPasswordConfirmEntry)
	})
	createButton.Importance = widget.HighImportance

	// Requirements and Best Practices
	requirementsLabel := widget.NewLabel(loc.Get("settings.backup.create.requirements"))
	requirementsLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		descLabel,
		widget.NewSeparator(),
		passwordLabel,
		backupPasswordEntry,
		confirmLabel,
		backupPasswordConfirmEntry,
		widget.NewSeparator(),
		includeDBCheck,
		dbCheckDesc,
		widget.NewSeparator(),
		createButton,
		widget.NewSeparator(),
		requirementsLabel,
	)

	return NewSettingsContentCard(
		loc.Get("settings.backup.create.title"),
		content,
		color.RGBA{R: 16, G: 185, B: 129, A: 255},
	)
}

// buildRestoreBackupCard creates the restore backup card
func (mss *ModernSettingsScreen) buildRestoreBackupCard() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	descLabel := widget.NewLabel(loc.Get("settings.backup.restore.desc"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.TextStyle = fyne.TextStyle{Italic: true}

	// File selection
	fileLabel := widget.NewLabelWithStyle(loc.Get("settings.backup.restore.file_label"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	restoreFileLabel := widget.NewLabel(loc.Get("settings.backup.restore.file_empty"))
	restoreFileLabel.Wrapping = fyne.TextWrapWord
	restoreFileLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Store selected file path in closure
	var selectedRestoreFile string

	// Password field
	passwordLabel := widget.NewLabelWithStyle(loc.Get("settings.backup.restore.password_label"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	restorePasswordEntry := widget.NewPasswordEntry()
	restorePasswordEntry.SetPlaceHolder(loc.Get("settings.backup.restore.password_placeholder"))

	// Warning
	warningLabel := widget.NewLabel(loc.Get("backup.restore_warning"))
	warningLabel.Wrapping = fyne.TextWrapWord
	warningLabel.TextStyle = fyne.TextStyle{Bold: true}
	warningLabel.Importance = widget.WarningImportance

	// Restore button
	restoreButton := widget.NewButton(loc.Get("settings.backup.restore.button"), func() {
		mss.restoreBackup(selectedRestoreFile, restorePasswordEntry.Text, restorePasswordEntry, restoreFileLabel)
	})
	restoreButton.Importance = widget.DangerImportance
	restoreButton.Disable()

	chooseButton := widget.NewButton(loc.Get("settings.backup.restore.choose_button"), func() {
		mss.chooseBackupFile(&selectedRestoreFile, restoreFileLabel, restoreButton)
	})
	chooseButton.Importance = widget.MediumImportance

	// Instructions
	instructionsLabel := widget.NewLabel(loc.Get("settings.backup.restore.steps"))
	instructionsLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		descLabel,
		widget.NewSeparator(),
		fileLabel,
		restoreFileLabel,
		chooseButton,
		widget.NewSeparator(),
		passwordLabel,
		restorePasswordEntry,
		widget.NewSeparator(),
		warningLabel,
		restoreButton,
		widget.NewSeparator(),
		instructionsLabel,
	)

	return NewSettingsContentCard(
		loc.Get("settings.backup.restore.title"),
		content,
		color.RGBA{R: 245, G: 158, B: 11, A: 255},
	)
}

// createBackup creates an encrypted backup
func (mss *ModernSettingsScreen) createBackup(password, confirm string, includeDB bool, passwordEntry, confirmEntry *widget.Entry) {
	loc := i18n.GetGlobalLocalizer()

	// Validation
	if len(password) < 8 {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.invalid_password"), loc.Get("backup.password_too_short"))
		return
	}

	if password != confirm {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.password.error.mismatch_title"), loc.Get("backup.password_mismatch_msg"))
		return
	}

	// Show file save dialog
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.file_error"), fmt.Sprintf(loc.Get("backup.file_open_failed"), err))
			return
		}
		if writer == nil {
			return // User cancelled
		}
		defer writer.Close()

		// Create backup
		config := mss.app.GetConfig()

		backupData, err := core.CreateBackup(config, includeDB, password)
		if err != nil {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.backup_failed"), fmt.Sprintf(loc.Get("backup.backup_failed_msg"), err))
			return
		}

		// Write to file
		if _, err := writer.Write(backupData); err != nil {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.write_failed"), fmt.Sprintf(loc.Get("backup.write_failed_msg"), err))
			return
		}

		// Clear password fields
		passwordEntry.SetText("")
		confirmEntry.SetText("")

		log.Printf("Backup created successfully: %s", writer.URI().Path())
		dialogs.ShowSuccess(mss.app.GetMainWindow(), loc.Get("backup.backup_created"), loc.Get("backup.backup_created_msg"))
	}, mss.app.GetMainWindow())

	saveDialog.SetFileName("tyr_backup.tyrbackup")
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".tyrbackup"}))
	saveDialog.Show()
}

// chooseBackupFile shows a file picker for restore
func (mss *ModernSettingsScreen) chooseBackupFile(selectedFile *string, fileLabel *widget.Label, restoreButton *widget.Button) {
	loc := i18n.GetGlobalLocalizer()

	openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.file_error"), fmt.Sprintf(loc.Get("backup.file_open_failed"), err))
			return
		}
		if reader == nil {
			return // User cancelled
		}
		defer reader.Close()

		// Store file path
		*selectedFile = reader.URI().Path()
		fileLabel.SetText(reader.URI().Name())
		restoreButton.Enable()

		log.Printf("Selected backup file: %s", *selectedFile)
	}, mss.app.GetMainWindow())

	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".tyrbackup"}))
	openDialog.Show()
}

// restoreBackup restores from an encrypted backup
func (mss *ModernSettingsScreen) restoreBackup(selectedFile, password string, passwordEntry *widget.Entry, fileLabel *widget.Label) {
	loc := i18n.GetGlobalLocalizer()

	if selectedFile == "" {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.no_file_selected"), loc.Get("backup.select_file_first"))
		return
	}

	if password == "" {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.password_required"), loc.Get("backup.enter_password"))
		return
	}

	// Confirm restore
	dialogs.ShowConfirmation(mss.app.GetMainWindow(),
		loc.Get("backup.confirm_restore"),
		loc.Get("backup.confirm_restore_msg"),
		func(confirmed bool) {
			if confirmed {
				mss.performRestore(selectedFile, password, passwordEntry, fileLabel)
			}
		},
	)
}

// performRestore performs the actual restore operation
func (mss *ModernSettingsScreen) performRestore(selectedFile, password string, passwordEntry *widget.Entry, fileLabel *widget.Label) {
	loc := i18n.GetGlobalLocalizer()

	// Read backup file
	backupData, err := core.ReadBackupFile(selectedFile)
	if err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.read_failed"), fmt.Sprintf(loc.Get("backup.read_failed_msg"), err))
		return
	}

	// Restore backup
	restoredConfig, dbData, err := core.RestoreBackup(backupData, password)
	if err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.restore_failed"), fmt.Sprintf(loc.Get("backup.restore_failed_msg"), err))
		return
	}

	// Stop service if running (use SoftStop for clean peer disconnection)
	serviceManager := mss.app.GetServiceManager()
	if serviceManager.IsRunning() {
		if err := serviceManager.SoftStop(); err != nil {
			log.Printf("Warning: Failed to stop service: %v", err)
		}
	}

	// Save restored configuration
	if err := restoredConfig.Save(); err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("backup.save_failed_msg"), err))
		return
	}

	// Restore database if included
	if dbData != nil && len(dbData) > 0 {
		dbPath := restoredConfig.ServiceSettings.DatabasePath
		if err := core.WriteBackupFile(dbPath, dbData); err != nil {
			log.Printf("Warning: Failed to restore database: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.partial_restore"), loc.Get("backup.partial_restore_msg"))
			return
		}
	}

	// Clear UI
	passwordEntry.SetText("")
	fileLabel.SetText(loc.Get("settings.backup.restore.file_empty"))

	log.Println("Backup restored successfully")

	// Show success and offer to restart
	dialogs.ShowInfo(mss.app.GetMainWindow(),
		loc.Get("backup.restore_complete"),
		loc.Get("backup.restart_required"),
	)
}

// buildAboutPage creates the about page with modern card layout
func (mss *ModernSettingsScreen) buildAboutPage() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Header with Back button
	header := buildSettingsHeaderModern(SettingsHeaderConfig{
		Title:    loc.Get("settings.about.title"),
		Subtitle: loc.Get("settings.about.subtitle"),
		OnBack: func() {
			mss.app.ShowSettings()
		},
	})

	// Application Info Card
	appNameLabel := widget.NewRichTextFromMarkdown("### Tyr")
	taglineLabel := widget.NewLabel(loc.Get("settings.about.tagline"))
	taglineLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: true}

	versionLabel := widget.NewLabel(loc.Get("settings.about.version"))
	versionLabel.TextStyle = fyne.TextStyle{Monospace: true}

	descLabel := widget.NewLabel(loc.Get("settings.about.description"))
	descLabel.Wrapping = fyne.TextWrapWord

	appInfoContent := container.NewVBox(
		appNameLabel,
		taglineLabel,
		versionLabel,
		widget.NewSeparator(),
		descLabel,
	)

	appInfoCard := NewSettingsContentCard(
		loc.Get("settings.about.info_card"),
		appInfoContent,
		color.RGBA{R: 168, G: 85, B: 247, A: 255},
	)

	// License Card
	licenseLabel := widget.NewLabel(loc.Get("settings.about.license_label"))
	licenseLabel.TextStyle = fyne.TextStyle{Bold: true}

	copyrightLabel := widget.NewLabel(loc.Get("settings.about.copyright_label"))
	copyrightLabel.TextStyle = fyne.TextStyle{Italic: true}

	licenseDesc := widget.NewLabel(loc.Get("settings.about.license_desc"))
	licenseDesc.Wrapping = fyne.TextWrapWord

	githubButton := widget.NewButton(loc.Get("settings.about.github_button"), func() {
		parsedURL, err := url.Parse("https://github.com/JB-SelfCompany/Tyr-Desktop")
		if err == nil {
			fyne.CurrentApp().OpenURL(parsedURL)
		}
	})
	githubButton.Importance = widget.MediumImportance

	licenseContent := container.NewVBox(
		licenseLabel,
		copyrightLabel,
		widget.NewSeparator(),
		licenseDesc,
		githubButton,
	)

	licenseCard := NewSettingsContentCard(
		loc.Get("settings.about.license_card"),
		licenseContent,
		color.RGBA{R: 99, G: 102, B: 241, A: 255},
	)

	// Main content
	mainContent := container.NewVBox(
		appInfoCard,
		licenseCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		header,
		nil, nil, nil,
		scrollContent,
	)
}

// handleAutoStartToggle handles auto-start setting
func (mss *ModernSettingsScreen) handleAutoStartToggle(enabled bool) {
	config := mss.app.GetConfig()
	loc := i18n.GetGlobalLocalizer()

	// Check current autostart status to avoid unnecessary operations
	currentStatus, err := core.IsAutoStartEnabled()
	if err != nil {
		log.Printf("Warning: Failed to check current autostart status: %v", err)
		// Use config value as fallback
		currentStatus = config.UIPreferences.AutoStart
	}

	// Only proceed if the state actually changed
	if currentStatus == enabled {
		log.Printf("Autostart state unchanged (%v), skipping toggle", enabled)
		return
	}

	// Apply the change
	if enabled {
		err = core.EnableAutoStart()
	} else {
		err = core.DisableAutoStart()
	}

	if err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.error.title"), fmt.Sprintf(loc.Get("settings.error.save_failed"), err))
		return
	}

	config.UIPreferences.AutoStart = enabled
	if err := config.Save(); err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.error.title"), fmt.Sprintf(loc.Get("settings.error.save_failed"), err))
		return
	}

	status := loc.Get("settings.autostart_disabled")
	if enabled {
		status = loc.Get("settings.autostart_enabled")
	}
	dialogs.ShowSuccess(mss.app.GetMainWindow(), loc.Get("settings.autostart_updated"), status)
}

// showChangePasswordDialog shows password change dialog
func (mss *ModernSettingsScreen) showChangePasswordDialog() {
	loc := i18n.GetGlobalLocalizer()

	currentPasswordEntry := widget.NewPasswordEntry()
	currentPasswordEntry.SetPlaceHolder(loc.Get("password.current"))

	newPasswordEntry := widget.NewPasswordEntry()
	newPasswordEntry.SetPlaceHolder(loc.Get("password.new"))

	confirmPasswordEntry := widget.NewPasswordEntry()
	confirmPasswordEntry.SetPlaceHolder(loc.Get("password.confirm"))

	form := container.NewVBox(
		widget.NewLabelWithStyle(loc.Get("password.current")+":", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		currentPasswordEntry,
		widget.NewLabel(""),
		widget.NewLabelWithStyle(loc.Get("password.new")+":", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		newPasswordEntry,
		widget.NewLabel(""),
		widget.NewLabelWithStyle(loc.Get("password.confirm")+":", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		confirmPasswordEntry,
	)

	dialog.ShowCustomConfirm(loc.Get("settings.password.dialog.title"), loc.Get("settings.password.dialog.change_button"), loc.Get("action.cancel"),
		form,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			config := mss.app.GetConfig()
			currentPassword, err := config.GetPassword()
			if err != nil {
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.password.error.retrieve"), fmt.Sprintf(loc.Get("settings.password.error.retrieve_msg"), err))
				return
			}

			if currentPasswordEntry.Text != currentPassword {
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.password.error.incorrect_title"), loc.Get("settings.password.error.incorrect_msg"))
				return
			}

			if len(newPasswordEntry.Text) < 6 {
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("password.error.length"), loc.Get("password.error.length"))
				return
			}

			if newPasswordEntry.Text != confirmPasswordEntry.Text {
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.password.error.mismatch_title"), loc.Get("settings.password.error.mismatch_msg"))
				return
			}

			serviceManager := mss.app.GetServiceManager()
			if err := serviceManager.UpdatePassword(newPasswordEntry.Text); err != nil {
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.password.error.update_title"), fmt.Sprintf(loc.Get("settings.password.error.update_msg"), err))
				return
			}

			dialogs.ShowSuccess(mss.app.GetMainWindow(), loc.Get("settings.password.success.title"), loc.Get("settings.password.success.msg"))
		},
		mss.app.GetMainWindow(),
	)
}

// buildPeerCard creates a modern card for a peer
func (mss *ModernSettingsScreen) buildPeerCard(peer core.PeerConfig, isConnected bool) fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Determine accent color based on status
	accentColor := color.RGBA{R: 156, G: 163, B: 175, A: 255} // Gray for disabled
	if peer.Enabled {
		if isConnected {
			accentColor = color.RGBA{R: 16, G: 185, B: 129, A: 255} // Green for connected
		} else {
			accentColor = color.RGBA{R: 245, G: 158, B: 11, A: 255} // Orange for enabled but not connected
		}
	}

	// Accent bar
	accentBar := canvas.NewRectangle(accentColor)
	accentBar.SetMinSize(fyne.NewSize(4, 0))

	// Status badge
	statusText := loc.Get("widget.disabled")
	statusColor := color.RGBA{R: 156, G: 163, B: 175, A: 255}
	if peer.Enabled {
		if isConnected {
			statusText = loc.Get("peers.status.connected")
			statusColor = color.RGBA{R: 16, G: 185, B: 129, A: 255}
		} else {
			statusText = loc.Get("peers.enabled")
			statusColor = color.RGBA{R: 245, G: 158, B: 11, A: 255}
		}
	}

	statusBadge := canvas.NewRectangle(statusColor)
	statusBadge.SetMinSize(fyne.NewSize(100, 28))

	statusLabel := widget.NewLabel(statusText)
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}
	statusLabel.Alignment = fyne.TextAlignCenter

	statusContainer := container.NewMax(statusBadge, statusLabel)

	// Peer address
	addressLabel := widget.NewLabel(peer.Address)
	addressLabel.Wrapping = fyne.TextWrapWord
	addressLabel.TextStyle = fyne.TextStyle{Monospace: true}

	// Enable/Disable toggle
	enableCheck := widget.NewCheck(loc.Get("peers.enable_peer"), func(checked bool) {
		config := mss.app.GetConfig()
		if checked {
			if err := config.EnablePeer(peer.Address); err != nil {
				log.Printf("Failed to enable peer: %v", err)
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.peers.toggle.error.enable"), fmt.Sprintf(loc.Get("settings.peers.toggle.error.enable_msg"), err))
				return
			}
		} else {
			if err := config.DisablePeer(peer.Address); err != nil {
				log.Printf("Failed to disable peer: %v", err)
				dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.peers.toggle.error.disable"), fmt.Sprintf(loc.Get("settings.peers.toggle.error.disable_msg"), err))
				return
			}
		}

		// Save config and refresh page
		if err := config.Save(); err != nil {
			log.Printf("Failed to save config: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("backup.save_failed_msg"), err))
			return
		}

		// Refresh page to update UI
		mss.app.GetMainWindow().SetContent(mss.buildNetworkPage())
	})
	// IMPORTANT: Set checked state AFTER creating the widget to avoid triggering the callback
	enableCheck.Checked = peer.Enabled
	enableCheck.Refresh()

	// Delete button
	deleteButton := widget.NewButton(loc.Get("settings.peers.delete.button"), func() {
		mss.confirmDeletePeer(peer.Address)
	})
	deleteButton.Importance = widget.DangerImportance

	// Card content
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, nil, statusContainer, addressLabel),
		widget.NewSeparator(),
		enableCheck,
		deleteButton,
	)

	// Card with accent
	cardBody := container.NewPadded(cardContent)
	cardWithAccent := container.NewBorder(nil, nil, accentBar, nil, cardBody)

	return container.NewPadded(cardWithAccent)
}

// confirmDeletePeer shows a confirmation dialog before deleting a peer
func (mss *ModernSettingsScreen) confirmDeletePeer(address string) {
	loc := i18n.GetGlobalLocalizer()

	dialogs.ShowConfirmation(mss.app.GetMainWindow(),
		loc.Get("settings.peers.delete.title"),
		fmt.Sprintf("%s\n\n%s\n\n%s", loc.Get("peers.confirm_remove_msg"), address, loc.Get("backup.confirm_restore_msg")),
		func(confirmed bool) {
			if confirmed {
				mss.deletePeer(address)
			}
		},
	)
}

// deletePeer removes a peer from the configuration
func (mss *ModernSettingsScreen) deletePeer(address string) {
	loc := i18n.GetGlobalLocalizer()

	config := mss.app.GetConfig()
	if err := config.RemovePeer(address); err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("settings.peers.delete.error.title"), fmt.Sprintf(loc.Get("settings.peers.delete.error.msg"), err))
		return
	}

	// Save configuration
	if err := config.Save(); err != nil {
		dialogs.ShowError(mss.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("backup.save_failed_msg"), err))
		return
	}

	// Refresh page to update UI
	mss.app.GetMainWindow().SetContent(mss.buildNetworkPage())

	dialogs.ShowSuccess(mss.app.GetMainWindow(), loc.Get("settings.peers.delete.success"), loc.Get("settings.peers.delete.success_msg"))
	log.Printf("Deleted peer: %s", address)
}

// showRegenerateKeysDialog shows confirmation dialog for key regeneration
func (mss *ModernSettingsScreen) showRegenerateKeysDialog() {
	loc := i18n.GetGlobalLocalizer()

	dialogs.ShowConfirmation(mss.app.GetMainWindow(),
		loc.Get("settings.security.regenerate_keys.dialog.title"),
		loc.Get("settings.security.regenerate_keys.dialog.message"),
		func(confirmed bool) {
			if confirmed {
				mss.regenerateKeys()
			}
		},
	)
}

// regenerateKeys performs the key regeneration operation
func (mss *ModernSettingsScreen) regenerateKeys() {
	loc := i18n.GetGlobalLocalizer()
	config := mss.app.GetConfig()
	serviceManager := mss.app.GetServiceManager()

	// Check if service is running
	wasRunning := serviceManager.IsRunning()

	// Stop service if running and wait for it to fully stop
	if wasRunning {
		log.Println("Stopping service for key regeneration...")
		if err := serviceManager.SoftStop(); err != nil {
			log.Printf("Warning: Failed to stop service: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(),
				loc.Get("settings.security.regenerate_keys.error.title"),
				fmt.Sprintf(loc.Get("settings.security.regenerate_keys.error.msg"), err))
			return
		}

		// Wait for service to fully stop (up to 10 seconds)
		log.Println("Waiting for service to stop completely...")
		for i := 0; i < 50; i++ {
			if !serviceManager.IsRunning() {
				log.Println("Service stopped successfully")
				break
			}
			time.Sleep(200 * time.Millisecond)

			// Timeout after 10 seconds
			if i == 49 {
				log.Printf("Timeout waiting for service to stop")
				dialogs.ShowError(mss.app.GetMainWindow(),
					loc.Get("settings.security.regenerate_keys.error.title"),
					fmt.Sprintf(loc.Get("settings.security.regenerate_keys.error.msg"),
						fmt.Errorf("service did not stop within 10 seconds")))
				return
			}
		}
	}

	// Close the service to release the database file
	log.Println("Closing service to release database file...")
	if err := serviceManager.CloseService(); err != nil {
		log.Printf("Failed to close service: %v", err)
		dialogs.ShowError(mss.app.GetMainWindow(),
			loc.Get("settings.security.regenerate_keys.error.title"),
			fmt.Sprintf(loc.Get("settings.security.regenerate_keys.error.msg"), err))

		// Try to restart service if it was running
		if wasRunning {
			log.Println("Attempting to restart service after close failure...")
			if initErr := serviceManager.Initialize(); initErr != nil {
				log.Printf("Failed to reinitialize service: %v", initErr)
			} else if startErr := serviceManager.Start(); startErr != nil {
				log.Printf("Failed to restart service: %v", startErr)
			}
		}
		return
	}

	// Delete the database file
	log.Printf("Deleting database at: %s", config.ServiceSettings.DatabasePath)
	dbPath := config.ServiceSettings.DatabasePath
	if err := core.DeleteDatabase(dbPath); err != nil {
		log.Printf("Failed to delete database: %v", err)
		dialogs.ShowError(mss.app.GetMainWindow(),
			loc.Get("settings.security.regenerate_keys.error.title"),
			fmt.Sprintf(loc.Get("settings.security.regenerate_keys.error.msg"), err))

		// Try to restart service if it was running
		if wasRunning {
			log.Println("Attempting to restart service after database deletion failure...")
			if startErr := serviceManager.Start(); startErr != nil {
				log.Printf("Failed to restart service: %v", startErr)
			}
		}
		return
	}

	log.Println("Database deleted successfully, reinitializing service with new keys...")

	// Reinitialize service with new keys
	if err := serviceManager.Initialize(); err != nil {
		log.Printf("Failed to initialize service with new keys: %v", err)
		dialogs.ShowError(mss.app.GetMainWindow(),
			loc.Get("settings.security.regenerate_keys.error.title"),
			fmt.Sprintf(loc.Get("settings.security.regenerate_keys.error.msg"), err))
		return
	}

	// Restart service if it was running
	if wasRunning {
		log.Println("Restarting service with new keys...")
		if err := serviceManager.Start(); err != nil {
			log.Printf("Failed to restart service: %v", err)
			dialogs.ShowError(mss.app.GetMainWindow(),
				loc.Get("settings.security.regenerate_keys.error.title"),
				fmt.Sprintf(loc.Get("settings.security.regenerate_keys.error.msg"), err))
			return
		}
	}

	// Show success message
	dialogs.ShowInfo(mss.app.GetMainWindow(),
		loc.Get("settings.security.regenerate_keys.success.title"),
		loc.Get("settings.security.regenerate_keys.success.msg"),
	)

	log.Println("Keys regenerated successfully")
}

// NewSettingsContentCard creates a modern content card with dynamic sizing
func NewSettingsContentCard(title string, content fyne.CanvasObject, accentColor color.Color) *SettingsContentCard {
	card := &SettingsContentCard{
		title:       title,
		content:     content,
		accentColor: accentColor,
	}
	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer for SettingsContentCard
func (scc *SettingsContentCard) CreateRenderer() fyne.WidgetRenderer {
	// Background for hover effect
	scc.background = canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 0})
	scc.background.CornerRadius = 10

	// Accent bar (left side)
	accentBar := canvas.NewRectangle(scc.accentColor)
	accentBar.SetMinSize(fyne.NewSize(4, 0))

	// Title
	titleLabel := widget.NewLabelWithStyle(scc.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.Importance = widget.HighImportance

	// Header
	headerContent := container.NewVBox(titleLabel, widget.NewSeparator())

	// Card body
	cardBody := container.NewVBox(
		headerContent,
		container.NewPadded(scc.content),
	)

	// Combine accent + body
	cardWithAccent := container.NewBorder(nil, nil, accentBar, nil, cardBody)

	// Layer background behind card for visible hover effect
	cardWithHover := container.NewStack(scc.background, cardWithAccent)

	// Outer padding
	scc.container = container.NewPadded(cardWithHover)

	return widget.NewSimpleRenderer(scc.container)
}

// MouseIn handles mouse enter events for SettingsContentCard
func (scc *SettingsContentCard) MouseIn(*desktop.MouseEvent) {
	scc.hovered = true
	if scc.background != nil {
		// Subtle highlight with accent color
		r, g, b, _ := scc.accentColor.RGBA()
		scc.background.FillColor = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: 25, // Subtle alpha for content cards
		}
		scc.background.Refresh()
	}
}

// MouseOut handles mouse leave events for SettingsContentCard
func (scc *SettingsContentCard) MouseOut() {
	scc.hovered = false
	if scc.background != nil {
		// Transparent background when not hovered
		scc.background.FillColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}
		scc.background.Refresh()
	}
}

// MouseMoved handles mouse movement (required for desktop.Hoverable interface)
func (scc *SettingsContentCard) MouseMoved(*desktop.MouseEvent) {}

// Cursor returns the default cursor for content cards
func (scc *SettingsContentCard) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// CreateRenderer for SettingsMenuCard
func (smc *SettingsMenuCard) CreateRenderer() fyne.WidgetRenderer {
	// Background for hover effect (full card size)
	smc.background = canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 0})
	smc.background.CornerRadius = 10

	// Accent bar
	accentBar := canvas.NewRectangle(smc.accentColor)
	accentBar.SetMinSize(fyne.NewSize(4, 0))

	// Title
	titleLabel := widget.NewLabelWithStyle(smc.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.Importance = widget.HighImportance

	// Description
	descLabel := widget.NewLabel(smc.description)
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Content (removed arrow - it was causing rendering issues)
	contentBox := container.NewVBox(titleLabel, descLabel)

	// Card with accent
	cardBody := container.NewPadded(contentBox)
	cardWithAccent := container.NewBorder(nil, nil, accentBar, nil, cardBody)

	// Layer background behind card for visible hover effect
	cardWithHover := container.NewStack(smc.background, cardWithAccent)

	smc.container = container.NewPadded(cardWithHover)

	return widget.NewSimpleRenderer(smc.container)
}

// Tapped handles tap events
func (smc *SettingsMenuCard) Tapped(_ *fyne.PointEvent) {
	if smc.onTapped != nil {
		smc.onTapped()
	}
}

// TappedSecondary handles secondary tap (no-op)
func (smc *SettingsMenuCard) TappedSecondary(_ *fyne.PointEvent) {}

// MouseIn handles mouse enter events (hover effect)
func (smc *SettingsMenuCard) MouseIn(*desktop.MouseEvent) {
	smc.hovered = true
	if smc.background != nil {
		// Brighter overlay with accent color for clear hover feedback
		r, g, b, _ := smc.accentColor.RGBA()
		smc.background.FillColor = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: 40, // More visible alpha value
		}
		smc.background.Refresh()
	}
}

// MouseOut handles mouse leave events (hover effect)
func (smc *SettingsMenuCard) MouseOut() {
	smc.hovered = false
	if smc.background != nil {
		// Transparent background when not hovered
		smc.background.FillColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}
		smc.background.Refresh()
	}
}

// MouseMoved handles mouse movement (required for desktop.Hoverable interface)
func (smc *SettingsMenuCard) MouseMoved(*desktop.MouseEvent) {}

// Cursor returns the cursor for this widget
func (smc *SettingsMenuCard) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}
