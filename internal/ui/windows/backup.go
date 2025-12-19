package windows

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/dialogs"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)

// ModernBackupScreen represents the modern backup/restore screen
type ModernBackupScreen struct {
	app AppInterface

	// Backup UI
	backupPasswordEntry        *widget.Entry
	backupPasswordConfirmEntry *widget.Entry
	includeDBCheck             *widget.Check

	// Restore UI
	restoreFileLabel     *widget.Label
	restorePasswordEntry *widget.Entry
	restoreButton        *widget.Button
	selectedRestoreFile  string
}

// NewModernBackupScreen creates a new modern backup/restore screen
func NewModernBackupScreen(app AppInterface) fyne.CanvasObject {
	screen := &ModernBackupScreen{
		app: app,
	}

	return screen.buildLayout()
}

// buildLayout creates the modern backup/restore screen layout
func (mbs *ModernBackupScreen) buildLayout() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Header
	backButton := widget.NewButton(loc.Get("action.back"), func() {
		// Navigate back to Backup settings page
		content := &ModernSettingsScreen{app: mbs.app}
		mbs.app.GetMainWindow().SetContent(content.buildBackupPage())
	})
	backButton.Importance = widget.LowImportance

	titleLabel := widget.NewRichTextFromMarkdown("# " + loc.Get("backup.title"))
	subtitleLabel := widget.NewLabel(loc.Get("backup.protect_subtitle"))
	subtitleLabel.TextStyle = fyne.TextStyle{Italic: true}

	header := container.NewVBox(
		backButton,
		titleLabel,
		subtitleLabel,
		widget.NewSeparator(),
	)

	// Create backup card
	createBackupCard := mbs.buildCreateBackupCard()

	// Restore backup card
	restoreBackupCard := mbs.buildRestoreBackupCard()

	// Main content
	mainContent := container.NewVBox(
		createBackupCard,
		restoreBackupCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		container.NewPadded(header),
		nil, nil, nil,
		scrollContent,
	)
}

// buildCreateBackupCard creates the backup creation card
func (mbs *ModernBackupScreen) buildCreateBackupCard() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	descLabel := widget.NewLabel(loc.Get("backup.create_description"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Password fields
	passwordLabel := widget.NewLabelWithStyle(loc.Get("backup.backup_password"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mbs.backupPasswordEntry = widget.NewPasswordEntry()
	mbs.backupPasswordEntry.SetPlaceHolder(loc.Get("backup.password"))

	confirmLabel := widget.NewLabelWithStyle(loc.Get("backup.confirm_password_label"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mbs.backupPasswordConfirmEntry = widget.NewPasswordEntry()
	mbs.backupPasswordConfirmEntry.SetPlaceHolder(loc.Get("backup.confirm_password"))

	// Include database checkbox
	mbs.includeDBCheck = widget.NewCheck(loc.Get("backup.include_db"), nil)
	mbs.includeDBCheck.SetChecked(true)

	dbCheckDesc := widget.NewLabel(loc.Get("backup.include_db_desc"))
	dbCheckDesc.Wrapping = fyne.TextWrapWord
	dbCheckDesc.TextStyle = fyne.TextStyle{Italic: true}

	// Create button
	createButton := widget.NewButton(loc.Get("backup.create_file"), func() {
		mbs.createBackup()
	})
	createButton.Importance = widget.HighImportance

	// Requirements and Best Practices
	requirementsLabel := widget.NewLabel(loc.Get("backup.requirements"))
	requirementsLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		descLabel,
		widget.NewSeparator(),
		passwordLabel,
		mbs.backupPasswordEntry,
		confirmLabel,
		mbs.backupPasswordConfirmEntry,
		widget.NewSeparator(),
		mbs.includeDBCheck,
		dbCheckDesc,
		widget.NewSeparator(),
		createButton,
		widget.NewSeparator(),
		requirementsLabel,
	)

	return NewSettingsContentCard(
		loc.Get("backup.create_backup"),
		content,
		color.RGBA{R: 16, G: 185, B: 129, A: 255},
	)
}

// buildRestoreBackupCard creates the restore backup card
func (mbs *ModernBackupScreen) buildRestoreBackupCard() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	descLabel := widget.NewLabel(loc.Get("backup.restore_description"))
	descLabel.Wrapping = fyne.TextWrapWord
	descLabel.TextStyle = fyne.TextStyle{Italic: true}

	// File selection
	fileLabel := widget.NewLabelWithStyle(loc.Get("backup.selected_file"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mbs.restoreFileLabel = widget.NewLabel(loc.Get("backup.no_file_selected"))
	mbs.restoreFileLabel.Wrapping = fyne.TextWrapWord
	mbs.restoreFileLabel.TextStyle = fyne.TextStyle{Italic: true}

	chooseButton := widget.NewButton(loc.Get("backup.choose_backup"), func() {
		mbs.chooseBackupFile()
	})
	chooseButton.Importance = widget.MediumImportance

	// Password field
	passwordLabel := widget.NewLabelWithStyle(loc.Get("backup.backup_password"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mbs.restorePasswordEntry = widget.NewPasswordEntry()
	mbs.restorePasswordEntry.SetPlaceHolder(loc.Get("backup.restore_password"))

	// Warning
	warningLabel := widget.NewLabel(loc.Get("backup.restore_warning"))
	warningLabel.Wrapping = fyne.TextWrapWord
	warningLabel.TextStyle = fyne.TextStyle{Bold: true}
	warningLabel.Importance = widget.WarningImportance

	// Restore button
	mbs.restoreButton = widget.NewButton(loc.Get("backup.restore_button"), func() {
		mbs.restoreBackup()
	})
	mbs.restoreButton.Importance = widget.DangerImportance
	mbs.restoreButton.Disable()

	// Instructions
	instructionsLabel := widget.NewLabel(loc.Get("backup.restore_steps"))
	instructionsLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		descLabel,
		widget.NewSeparator(),
		fileLabel,
		mbs.restoreFileLabel,
		chooseButton,
		widget.NewSeparator(),
		passwordLabel,
		mbs.restorePasswordEntry,
		widget.NewSeparator(),
		warningLabel,
		mbs.restoreButton,
		widget.NewSeparator(),
		instructionsLabel,
	)

	return NewSettingsContentCard(
		loc.Get("backup.restore_backup"),
		content,
		color.RGBA{R: 245, G: 158, B: 11, A: 255},
	)
}

// buildInfoCard creates the information card
func (mbs *ModernBackupScreen) buildInfoCard() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	infoTitle := widget.NewLabel(loc.Get("backup.info.why_title"))
	infoTitle.TextStyle = fyne.TextStyle{Bold: true}

	infoPoints := []string{
		loc.Get("backup.info.point1"),
		loc.Get("backup.info.point2"),
		loc.Get("backup.info.point3"),
		loc.Get("backup.info.point4"),
	}

	infoList := container.NewVBox()
	for _, point := range infoPoints {
		pointLabel := widget.NewLabel(point)
		pointLabel.Wrapping = fyne.TextWrapWord
		infoList.Add(pointLabel)
	}

	bestPracticesTitle := widget.NewLabel(loc.Get("backup.info.best_practices"))
	bestPracticesTitle.TextStyle = fyne.TextStyle{Bold: true}

	bestPractices := []string{
		loc.Get("backup.info.practice1"),
		loc.Get("backup.info.practice2"),
		loc.Get("backup.info.practice3"),
		loc.Get("backup.info.practice4"),
		loc.Get("backup.info.practice5"),
	}

	bestPracticesList := container.NewVBox()
	for _, practice := range bestPractices {
		practiceLabel := widget.NewLabel(practice)
		practiceLabel.Wrapping = fyne.TextWrapWord
		bestPracticesList.Add(practiceLabel)
	}

	content := container.NewVBox(
		infoTitle,
		widget.NewSeparator(),
		infoList,
		widget.NewSeparator(),
		bestPracticesTitle,
		widget.NewSeparator(),
		bestPracticesList,
	)

	return NewSettingsContentCard(
		loc.Get("backup.info_card"),
		content,
		color.RGBA{R: 59, G: 130, B: 246, A: 255},
	)
}

// createBackup creates an encrypted backup
func (mbs *ModernBackupScreen) createBackup() {
	loc := i18n.GetGlobalLocalizer()

	password := mbs.backupPasswordEntry.Text
	confirm := mbs.backupPasswordConfirmEntry.Text

	// Validation
	if len(password) < 8 {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.invalid_password"), loc.Get("backup.password_too_short"))
		return
	}

	if password != confirm {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.password_mismatch"), loc.Get("backup.password_mismatch_msg"))
		return
	}

	// Show file save dialog
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.file_error"), fmt.Sprintf(loc.Get("backup.file_open_failed"), err))
			return
		}
		if writer == nil {
			return // User cancelled
		}
		defer writer.Close()

		// Create backup
		config := mbs.app.GetConfig()
		includeDB := mbs.includeDBCheck.Checked

		backupData, err := core.CreateBackup(config, includeDB, password)
		if err != nil {
			dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.backup_failed"), fmt.Sprintf(loc.Get("backup.backup_failed_msg"), err))
			return
		}

		// Write to file
		if _, err := writer.Write(backupData); err != nil {
			dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.write_failed"), fmt.Sprintf(loc.Get("backup.write_failed_msg"), err))
			return
		}

		// Clear password fields
		mbs.backupPasswordEntry.SetText("")
		mbs.backupPasswordConfirmEntry.SetText("")

		log.Printf("Backup created successfully: %s", writer.URI().Path())
		dialogs.ShowSuccess(mbs.app.GetMainWindow(), loc.Get("backup.backup_created"), loc.Get("backup.backup_created_msg"))
	}, mbs.app.GetMainWindow())

	saveDialog.SetFileName("tyr_backup.tyrbackup")
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".tyrbackup"}))
	saveDialog.Show()
}

// chooseBackupFile shows a file picker for restore
func (mbs *ModernBackupScreen) chooseBackupFile() {
	loc := i18n.GetGlobalLocalizer()

	openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.file_error"), fmt.Sprintf(loc.Get("backup.file_open_failed"), err))
			return
		}
		if reader == nil {
			return // User cancelled
		}
		defer reader.Close()

		// Store file path
		mbs.selectedRestoreFile = reader.URI().Path()
		mbs.restoreFileLabel.SetText(reader.URI().Name())
		mbs.restoreButton.Enable()

		log.Printf("Selected backup file: %s", mbs.selectedRestoreFile)
	}, mbs.app.GetMainWindow())

	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".tyrbackup"}))
	openDialog.Show()
}

// restoreBackup restores from an encrypted backup
func (mbs *ModernBackupScreen) restoreBackup() {
	loc := i18n.GetGlobalLocalizer()

	if mbs.selectedRestoreFile == "" {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.no_file_selected"), loc.Get("backup.select_file_first"))
		return
	}

	password := mbs.restorePasswordEntry.Text
	if password == "" {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.password_required"), loc.Get("backup.enter_password"))
		return
	}

	// Confirm restore
	dialogs.ShowConfirmation(mbs.app.GetMainWindow(),
		loc.Get("backup.confirm_restore"),
		loc.Get("backup.confirm_restore_msg"),
		func(confirmed bool) {
			if confirmed {
				mbs.performRestore(password)
			}
		},
	)
}

// performRestore performs the actual restore operation
func (mbs *ModernBackupScreen) performRestore(password string) {
	loc := i18n.GetGlobalLocalizer()

	// Read backup file
	backupData, err := core.ReadBackupFile(mbs.selectedRestoreFile)
	if err != nil {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.read_failed"), fmt.Sprintf(loc.Get("backup.read_failed_msg"), err))
		return
	}

	// Restore backup
	restoredConfig, dbData, err := core.RestoreBackup(backupData, password)
	if err != nil {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.restore_failed"), fmt.Sprintf(loc.Get("backup.restore_failed_msg"), err))
		return
	}

	// Stop service if running (use SoftStop for clean peer disconnection)
	serviceManager := mbs.app.GetServiceManager()
	if serviceManager.IsRunning() {
		if err := serviceManager.SoftStop(); err != nil {
			log.Printf("Warning: Failed to stop service: %v", err)
		}
	}

	// Save restored configuration
	if err := restoredConfig.Save(); err != nil {
		dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("backup.save_failed_msg"), err))
		return
	}

	// Restore database if included
	if dbData != nil && len(dbData) > 0 {
		dbPath := restoredConfig.ServiceSettings.DatabasePath
		if err := core.WriteBackupFile(dbPath, dbData); err != nil {
			log.Printf("Warning: Failed to restore database: %v", err)
			dialogs.ShowError(mbs.app.GetMainWindow(), loc.Get("backup.partial_restore"), loc.Get("backup.partial_restore_msg"))
			return
		}
	}

	// Clear UI
	mbs.restorePasswordEntry.SetText("")
	mbs.restoreFileLabel.SetText(loc.Get("backup.no_file_selected"))
	mbs.selectedRestoreFile = ""
	mbs.restoreButton.Disable()

	log.Println("Backup restored successfully")

	// Show success and offer to restart
	dialogs.ShowInfo(mbs.app.GetMainWindow(),
		loc.Get("backup.restore_complete"),
		loc.Get("backup.restart_required"),
	)
}
