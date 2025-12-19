package windows

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/dialogs"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/ui/i18n"
)

// ModernPeersScreen represents the modern peer management screen
type ModernPeersScreen struct {
	app AppInterface

	// UI components
	peersContainer *fyne.Container
	addPeerEntry   *widget.Entry
	header         *SettingsHeader

	// State tracking
	hasUnsavedChanges bool
	wasServiceRunning bool
}

// PeerCard represents a modern peer card widget
type PeerCard struct {
	widget.BaseWidget
	peer          core.PeerConfig
	isConnected   bool
	accentColor   color.Color
	onToggle      func(bool)
	onDelete      func()
	container     *fyne.Container
}

// NewModernPeersScreen creates a new modern peers management screen
func NewModernPeersScreen(app AppInterface) fyne.CanvasObject {
	screen := &ModernPeersScreen{
		app:               app,
		peersContainer:    container.NewVBox(),
		hasUnsavedChanges: false,
		wasServiceRunning: app.GetServiceManager().IsRunning(),
	}

	return screen.buildLayout()
}

// buildLayout creates the modern peers screen layout
func (mps *ModernPeersScreen) buildLayout() fyne.CanvasObject {
	loc := i18n.GetGlobalLocalizer()

	// Create header with back and apply buttons
	mps.header = NewSettingsHeader(
		"← "+loc.Get("peers.back"),
		func() {
			// Navigate back to Network settings
			content := &ModernSettingsScreen{app: mps.app}
			mps.app.GetMainWindow().SetContent(content.buildNetworkPage())
		},
		func() {
			// Apply changes
			mps.applyPeerChanges()
		},
		&mps.hasUnsavedChanges,
		mps.app,
	)

	// Title and subtitle
	titleLabel := widget.NewRichTextFromMarkdown("# " + loc.Get("peers.network_peers"))
	subtitleLabel := widget.NewLabel(loc.Get("peers.manage_subtitle"))
	subtitleLabel.TextStyle = fyne.TextStyle{Italic: true}

	header := container.NewVBox(
		mps.header.GetContainer(),
		widget.NewSeparator(),
		titleLabel,
		subtitleLabel,
		widget.NewSeparator(),
	)

	// Add peer card
	mps.addPeerEntry = widget.NewEntry()
	mps.addPeerEntry.SetPlaceHolder(loc.Get("peers.enter_address"))

	addButton := widget.NewButton(loc.Get("peers.add_peer"), func() {
		mps.addPeer()
	})
	addButton.Importance = widget.HighImportance

	addPeerDesc := widget.NewLabel(loc.Get("peers.add_description"))
	addPeerDesc.Wrapping = fyne.TextWrapWord
	addPeerDesc.TextStyle = fyne.TextStyle{Italic: true}

	addPeerContent := container.NewVBox(
		widget.NewLabel(loc.Get("peers.new_peer_address")),
		mps.addPeerEntry,
		addButton,
		addPeerDesc,
	)

	addPeerCard := NewSettingsContentCard(
		loc.Get("peers.add_new_peer"),
		addPeerContent,
		color.RGBA{R: 16, G: 185, B: 129, A: 255},
	)

	// Update peers list
	mps.updatePeersList()

	// Peers list card
	peersListLabel := widget.NewLabel(loc.Get("peers.manage_description"))
	peersListLabel.Wrapping = fyne.TextWrapWord
	peersListLabel.TextStyle = fyne.TextStyle{Italic: true}

	peersListContent := container.NewVBox(
		peersListLabel,
		widget.NewSeparator(),
		mps.peersContainer,
	)

	peersListCard := NewSettingsContentCard(
		loc.Get("peers.configured_peers"),
		peersListContent,
		color.RGBA{R: 59, G: 130, B: 246, A: 255},
	)

	// Info card
	infoTitle := widget.NewLabel(loc.Get("peers.about_what"))
	infoTitle.TextStyle = fyne.TextStyle{Bold: true}

	infoPoints := []string{
		loc.Get("peers.about_point1"),
		loc.Get("peers.about_point2"),
		loc.Get("peers.about_point3"),
		loc.Get("peers.about_point4"),
		loc.Get("peers.about_point5"),
	}

	infoList := container.NewVBox()
	for _, point := range infoPoints {
		pointLabel := widget.NewLabel(point)
		pointLabel.Wrapping = fyne.TextWrapWord
		infoList.Add(pointLabel)
	}

	infoContent := container.NewVBox(
		infoTitle,
		widget.NewSeparator(),
		infoList,
	)

	infoCard := NewSettingsContentCard(
		loc.Get("peers.about_title"),
		infoContent,
		color.RGBA{R: 168, G: 85, B: 247, A: 255},
	)

	// Main content
	mainContent := container.NewVBox(
		addPeerCard,
		peersListCard,
		infoCard,
	)

	scrollContent := container.NewScroll(mainContent)

	return container.NewBorder(
		container.NewPadded(header),
		nil,
		nil, nil,
		scrollContent,
	)
}

// updatePeersList updates the peers list display
func (mps *ModernPeersScreen) updatePeersList() {
	loc := i18n.GetGlobalLocalizer()
	config := mps.app.GetConfig()
	serviceManager := mps.app.GetServiceManager()
	peerStats := serviceManager.GetPeerStats()

	mps.peersContainer.Objects = []fyne.CanvasObject{}

	if len(config.NetworkPeers) == 0 {
		emptyLabel := widget.NewLabel(loc.Get("peers.no_configured"))
		emptyLabel.Alignment = fyne.TextAlignCenter
		emptyLabel.TextStyle = fyne.TextStyle{Italic: true}

		emptyCard := container.NewPadded(
			container.NewCenter(emptyLabel),
		)

		mps.peersContainer.Objects = append(mps.peersContainer.Objects, emptyCard)
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

			peerCard := mps.buildPeerCard(peer, isConnected)
			mps.peersContainer.Objects = append(mps.peersContainer.Objects, peerCard)
		}
	}

	mps.peersContainer.Refresh()
}

// buildPeerCard creates a modern card for a peer
func (mps *ModernPeersScreen) buildPeerCard(peer core.PeerConfig, isConnected bool) fyne.CanvasObject {
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
		config := mps.app.GetConfig()
		if checked {
			if err := config.EnablePeer(peer.Address); err != nil {
				log.Printf("Failed to enable peer: %v", err)
				dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.edit_failed"), fmt.Sprintf(loc.Get("error.failed_to_enable_peer"), err))
				return
			}
		} else {
			if err := config.DisablePeer(peer.Address); err != nil {
				log.Printf("Failed to disable peer: %v", err)
				dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.edit_failed"), fmt.Sprintf(loc.Get("error.failed_to_disable_peer"), err))
				return
			}
		}

		// Save config immediately so it persists
		if err := config.Save(); err != nil {
			log.Printf("Failed to save config: %v", err)
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("error.failed_to_save"), err))
			return
		}

		// Mark as unsaved for service update
		log.Printf("[PeersScreen] Peer toggled - setting hasUnsavedChanges = true")
		mps.hasUnsavedChanges = true
		mps.updateApplyButtonVisibility()

		log.Printf("Peer %s %s - changes need to be applied", peer.Address, map[bool]string{true: "enabled", false: "disabled"}[checked])
	})
	// IMPORTANT: Set checked state AFTER creating the widget to avoid triggering the callback
	enableCheck.Checked = peer.Enabled
	enableCheck.Refresh()

	// Edit button
	editButton := widget.NewButton(loc.Get("peers.edit"), func() {
		mps.showEditPeerDialog(peer.Address)
	})
	editButton.Importance = widget.MediumImportance

	// Delete button
	deleteButton := widget.NewButton(loc.Get("action.delete"), func() {
		mps.confirmDeletePeer(peer.Address)
	})
	deleteButton.Importance = widget.DangerImportance

	// Card content
	cardContent := container.NewVBox(
		container.NewBorder(nil, nil, nil, statusContainer, addressLabel),
		widget.NewSeparator(),
		enableCheck,
		container.NewGridWithColumns(2, editButton, deleteButton),
	)

	// Card with accent
	cardBody := container.NewPadded(cardContent)
	cardWithAccent := container.NewBorder(nil, nil, accentBar, nil, cardBody)

	return container.NewPadded(cardWithAccent)
}

// addPeer adds a new peer to the configuration
func (mps *ModernPeersScreen) addPeer() {
	loc := i18n.GetGlobalLocalizer()
	address := mps.addPeerEntry.Text

	if address == "" {
		dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("error.invalid_input"), loc.Get("peers.error.invalid"))
		return
	}

	config := mps.app.GetConfig()
	if err := config.AddPeer(address); err != nil {
		dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.error.add"), fmt.Sprintf(loc.Get("error.failed_to_add_peer"), err))
		return
	}

	// Save configuration
	if err := config.Save(); err != nil {
		dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("error.failed_to_save_configuration"), err))
		return
	}

	// Clear entry and update list
	mps.addPeerEntry.SetText("")
	log.Printf("[PeersScreen] Peer added - setting hasUnsavedChanges = true")
	mps.hasUnsavedChanges = true
	mps.updateApplyButtonVisibility()
	mps.updatePeersList()

	dialogs.ShowSuccess(mps.app.GetMainWindow(), loc.Get("peers.peer_added"), loc.Get("peers.peer_added"))
	log.Printf("Added peer: %s", address)
}

// confirmDeletePeer shows a confirmation dialog before deleting a peer
func (mps *ModernPeersScreen) confirmDeletePeer(address string) {
	loc := i18n.GetGlobalLocalizer()
	dialogs.ShowConfirmation(mps.app.GetMainWindow(),
		loc.Get("peers.confirm_remove"),
		fmt.Sprintf("%s\n\n%s", loc.Get("peers.confirm_remove_msg"), address),
		func(confirmed bool) {
			if confirmed {
				mps.deletePeer(address)
			}
		},
	)
}

// deletePeer removes a peer from the configuration
func (mps *ModernPeersScreen) deletePeer(address string) {
	loc := i18n.GetGlobalLocalizer()
	config := mps.app.GetConfig()
	if err := config.RemovePeer(address); err != nil {
		dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.error.remove"), fmt.Sprintf(loc.Get("error.failed_to_delete_peer"), err))
		return
	}

	// Save configuration
	if err := config.Save(); err != nil {
		dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("error.failed_to_save_configuration"), err))
		return
	}

	log.Printf("[PeersScreen] Peer deleted - setting hasUnsavedChanges = true")
	mps.hasUnsavedChanges = true
	mps.updateApplyButtonVisibility()
	mps.updatePeersList()
	dialogs.ShowSuccess(mps.app.GetMainWindow(), loc.Get("peers.peer_removed"), loc.Get("peers.peer_removed"))
	log.Printf("Deleted peer: %s", address)
}

// showEditPeerDialog shows a dialog to edit an existing peer
func (mps *ModernPeersScreen) showEditPeerDialog(oldAddress string) {
	loc := i18n.GetGlobalLocalizer()

	// Create entry for new address
	newAddressEntry := widget.NewEntry()
	newAddressEntry.SetPlaceHolder(loc.Get("peers.enter_address"))
	newAddressEntry.SetText(oldAddress)

	// Create form
	formContent := container.NewVBox(
		widget.NewLabel(loc.Get("peers.edit_address")),
		newAddressEntry,
	)

	// Create dialog
	dlg := dialog.NewCustomConfirm(loc.Get("peers.edit_peer"), loc.Get("action.save"), loc.Get("action.cancel"), formContent, func(save bool) {
		if !save {
			return
		}

		newAddress := newAddressEntry.Text
		if newAddress == "" {
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("error.invalid_input"), loc.Get("peers.error.invalid"))
			return
		}

		if newAddress == oldAddress {
			// No changes made
			return
		}

		config := mps.app.GetConfig()

		// Check if new address already exists (and it's not the same peer)
		for _, peer := range config.NetworkPeers {
			if peer.Address == newAddress {
				dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.duplicate_peer"), loc.Get("peers.duplicate_exists"))
				return
			}
		}

		// Find the old peer to preserve its enabled state
		var wasEnabled bool = true
		for _, peer := range config.NetworkPeers {
			if peer.Address == oldAddress {
				wasEnabled = peer.Enabled
				break
			}
		}

		// Remove old peer
		if err := config.RemovePeer(oldAddress); err != nil {
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.edit_failed"), fmt.Sprintf(loc.Get("peers.remove_old_failed"), err))
			return
		}

		// Add new peer with preserved enabled state
		if err := config.AddPeer(newAddress); err != nil {
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.edit_failed"), fmt.Sprintf(loc.Get("peers.add_new_failed"), err))
			// Try to restore old peer
			config.AddPeer(oldAddress)
			return
		}

		// Set enabled state
		if wasEnabled {
			config.EnablePeer(newAddress)
		} else {
			config.DisablePeer(newAddress)
		}

		// Save configuration
		if err := config.Save(); err != nil {
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("backup.save_failed"), fmt.Sprintf(loc.Get("error.failed_to_save_configuration"), err))
			return
		}

		log.Printf("[PeersScreen] Peer edited - setting hasUnsavedChanges = true")
		mps.hasUnsavedChanges = true
		mps.updateApplyButtonVisibility()
		mps.updatePeersList()

		dialogs.ShowSuccess(mps.app.GetMainWindow(), loc.Get("peers.peer_updated"), loc.Get("peers.peer_updated_msg"))
		log.Printf("Edited peer: %s -> %s", oldAddress, newAddress)
	}, mps.app.GetMainWindow())

	dlg.Resize(fyne.NewSize(500, 150))
	dlg.Show()
}

// updateApplyButtonVisibility enables/disables the apply button and changes its color
func (mps *ModernPeersScreen) updateApplyButtonVisibility() {
	log.Printf("[PeersScreen] updateApplyButtonVisibility called - hasUnsavedChanges: %v", mps.hasUnsavedChanges)
	mps.header.UpdateApplyButton(mps.hasUnsavedChanges)
}

// applyPeerChanges applies the peer configuration changes to the running service
func (mps *ModernPeersScreen) applyPeerChanges() {
	loc := i18n.GetGlobalLocalizer()
	serviceManager := mps.app.GetServiceManager()

	if !mps.wasServiceRunning {
		// Service was not running, just close the activity
		mps.hasUnsavedChanges = false
		mps.updateApplyButtonVisibility()
		dialogs.ShowSuccess(mps.app.GetMainWindow(), loc.Get("peers.changes_applied"), loc.Get("peers.config_saved"))
		return
	}

	// If service is running, use hot reload to update peers without restart
	if serviceManager.IsRunning() {
		config := mps.app.GetConfig()
		enabledPeers := config.GetEnabledPeers()

		if len(enabledPeers) == 0 {
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.error.none_enabled"), loc.Get("peers.error.none_enabled"))
			return
		}

		// Hot reload peers using the service's UpdatePeers method
		if err := serviceManager.HotReloadPeers(enabledPeers); err != nil {
			dialogs.ShowError(mps.app.GetMainWindow(), loc.Get("peers.apply_failed"), fmt.Sprintf(loc.Get("peers.apply_failed_msg"), err))
			log.Printf("Failed to hot reload peers: %v", err)
			return
		}

		mps.hasUnsavedChanges = false
		mps.updateApplyButtonVisibility()
		// NOTE: Don't call updatePeersList() here to avoid scroll reset
		// Peer connection statuses will update automatically via connectionCallback

		dialogs.ShowSuccess(mps.app.GetMainWindow(), loc.Get("peers.changes_applied"), loc.Get("peers.config_updated"))
		log.Printf("Hot reloaded peers successfully")
	} else {
		// Service stopped, just clear the flag
		mps.hasUnsavedChanges = false
		mps.updateApplyButtonVisibility()
		dialogs.ShowSuccess(mps.app.GetMainWindow(), loc.Get("peers.changes_applied"), loc.Get("peers.config_saved"))
	}
}
