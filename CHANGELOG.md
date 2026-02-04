# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.2.0] - 2025-02-04

### Changed

- **Unified Visual Design** - Completely redesigned UI to design language: Slate + Emerald color palette, glassmorphism effects, consistent border-radius and spacing.
- **Dark-Only Theme** - Switched to dark-only mode. Removed light/system theme options from Settings. All `dark:` Tailwind class prefixes and light theme CSS variables removed.
- **Updated Base Font Size** - Increased base font size from 14px to 16px for better readability.
- **Sidebar Redesign** - Updated sidebar to 256px width with glassmorphism background, centered logo, and rounded active navigation state with emerald accent.
- **Component Styling Overhaul** - Updated core UI components:
  - Button: `rounded-xl`, removed borders from variants
  - Input: `rounded-xl`, darker background (`bg-slate-800`), refined focus ring
  - GlassCard: `rounded-2xl`, added `glass` class and `shadow-glass`
  - Modal: `rounded-2xl`, updated border styling
- **CSS Architecture** - Rewrote `style.css`: removed aggressive `!important` overrides, added proper CSS variables and `.glass`/`.shadow-glass` utility classes.
- **Legacy Theme Cleanup** - Replaced all old theme classes (`neon-cyan`, `space-blue`, `md-light-*`, `md-dark-*`) with Slate + Emerald equivalents across PeerDiscoveryModal, LogViewer, and ErrorBoundary components.

### Removed

- **Theme Selector** - Removed theme selection (Light/Dark/System) from General Settings. Application now uses dark theme exclusively.

## [2.1.1] - 2025-01-27

### Added

- **Portable Mode** - Application is now fully portable. All configuration files are stored in `data/` subdirectory next to the executable instead of system-specific locations (AppData on Windows, ~/.config on Linux).
- **Automatic Migration** - On first launch, data from legacy directories is automatically migrated to the new portable location and old directory is cleaned up.
- **Centralized Version Management** - Version is now defined in single source of truth (`internal/version/version.go`) and automatically used by build scripts.

### Fixed

- **System Tray Behavior** - Fixed tray menu and window show/hide behavior.
- **Backup Restore Error Handling** - Fixed incorrect password not showing error message. Now properly validates `ResultDTO.Success` and displays appropriate error to user.
- **UI Update After Restore** - Fixed onboarding screen not updating after successful backup restore. App now listens for `config:restored` event and refreshes onboarding state.
- **Database Path After Migration** - Fixed service not starting after migration due to old database path in config. Now always uses portable path regardless of saved value.

### Changed

- **Data Storage Location** - All application data now stored in `./data/` subdirectory next to executable.

## [2.1.0] - 2025-01-26

### Added

- **Message Size Limit Setting** - Added ability to adjust the maximum size of received messages, similar to the mobile version.
- **Stop Button for Peer Discovery** - Added Stop button during peer discovery process on Peers page for user-initiated cancellation.

### Changed

- **Backup Moved to Settings** - Backup & Restore functionality moved from sidebar to Settings page.
- **Updated Yggmail Library** - Updated to latest yggmail version with improvements.
- **DeltaChat Instructions** - Added correct instructions for configuring DeltaChat and other Email clients.
- **Toast Notifications** - Repositioned all toast notifications to bottom-center with reduced display time (2 seconds).
- **System Tray** - Adjusted system tray operation and translations.

### Removed

- **System Tray Logging** - Removed excessive logging from system tray operations.

## [2.0.2] - 2025-01-20

### Fixed

- **Key Regeneration Button** - Fixed Key Regeneration button not working correctly.
- **Password Change Button** - Fixed Password Change button functionality.
- **Copy Buttons** - Fixed Copy buttons in the dashboard.
- **Notifications** - Fixed notification display issues.

### Changed

- **Interface Adjustments** - Various UI/UX improvements.

## [2.0.1] - 2025-01-18

### Fixed

- **DeltaChat Quick Setup** - Adjusted the way to quickly configure DeltaChat.
- **Peer Apply Button** - Fixed bug where it was necessary to press "Apply" multiple times for peer changes.
- **Database Path** - Fixed work with database path, including display in the application.
- **Service Restart** - Adjusted service restart behavior.
- **System Tray** - Fixed tray when application was stuck and would not open.

### Changed

- **Build Scripts** - Adjusted naming of manually assembled files.

## [2.0.0] - 2025-01-15

### Added

- Initial release of Tyr Desktop
- P2P email over Yggdrasil network
- DeltaChat integration with one-click setup
- Local SMTP/IMAP server (ports 1025/1143)
- Encrypted backup and restore (AES-256-GCM)
- OS keyring integration for password storage (Windows Credential Manager / GNOME Keyring)
- System tray support with show/hide functionality
- Auto-start on system boot
- Real-time peer management with hot reload
- Dark theme with glassmorphism UI
- English and Russian localization
- Real-time service status monitoring
- Log viewer

[2.2.0]: https://github.com/JB-SelfCompany/Tyr-Desktop/compare/2.1.1...2.2.0
[2.1.1]: https://github.com/JB-SelfCompany/Tyr-Desktop/compare/2.1.0...2.1.1
[2.1.0]: https://github.com/JB-SelfCompany/Tyr-Desktop/compare/2.0.2...2.1.0
[2.0.2]: https://github.com/JB-SelfCompany/Tyr-Desktop/compare/2.0.1...2.0.2
[2.0.1]: https://github.com/JB-SelfCompany/Tyr-Desktop/compare/2.0.0...2.0.1
[2.0.0]: https://github.com/JB-SelfCompany/Tyr-Desktop/releases/tag/2.0.0
