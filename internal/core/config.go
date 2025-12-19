package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the application configuration
// Stored in TOML format with passwords in OS keyring
// Thread-safe for concurrent access
type Config struct {
	// OnboardingComplete indicates if the initial setup wizard has been completed
	OnboardingComplete bool `toml:"onboarding_complete"`

	// ServiceSettings contains SMTP/IMAP server configuration
	ServiceSettings ServiceSettings `toml:"service_settings"`

	// NetworkPeers contains the list of Yggdrasil network peers
	NetworkPeers []PeerConfig `toml:"network_peers"`

	// UIPreferences contains user interface preferences
	UIPreferences UIPreferences `toml:"ui_preferences"`

	// CachedDiscoveredPeers contains cached discovered peers (TTL: 24 hours)
	CachedDiscoveredPeers []DiscoveredPeer `toml:"cached_discovered_peers,omitempty"`

	// CacheTimestamp is the Unix timestamp when peers were cached
	CacheTimestamp int64 `toml:"cache_timestamp,omitempty"`

	// Mutex for thread-safe access to configuration
	mu sync.RWMutex `toml:"-"`
}

// ServiceSettings contains SMTP/IMAP server and database configuration
type ServiceSettings struct {
	// SMTPAddress is the local SMTP server address (default: 127.0.0.1:1025)
	SMTPAddress string `toml:"smtp_address"`

	// IMAPAddress is the local IMAP server address (default: 127.0.0.1:1143)
	IMAPAddress string `toml:"imap_address"`

	// DatabasePath is the path to the yggmail.db file
	DatabasePath string `toml:"database_path"`

	// PasswordInitialized indicates if the password has been set in yggmail database
	// Used to avoid calling SetPassword on every restart
	PasswordInitialized bool `toml:"password_initialized"`
}

// PeerConfig represents a Yggdrasil network peer configuration
type PeerConfig struct {
	// Address is the peer URI in one of the supported formats:
	//   tcp://host:port, tls://host:port, quic://host:port
	//   socks://proxy:port/host:port, sockstls://proxy:port/host:port
	//   unix:///path/to/sock.sock
	//   ws://host:port[/path], wss://host:port[/path]
	Address string `toml:"address"`

	// Enabled indicates if this peer should be used for connections
	Enabled bool `toml:"enabled"`
}

// UIPreferences contains user interface configuration
type UIPreferences struct {
	// Theme is the UI theme ("light", "dark", or "system")
	Theme string `toml:"theme"`

	// Language is the UI language code ("en", "ru", etc.)
	Language string `toml:"language"`

	// AutoStart indicates if the service should start on system boot
	AutoStart bool `toml:"auto_start"`

	// WindowState contains the window size and position
	WindowState WindowState `toml:"window_state"`
}

// WindowState contains window size and position information
type WindowState struct {
	// Width of the window in pixels
	Width int `toml:"width"`

	// Height of the window in pixels
	Height int `toml:"height"`

	// X position of the window in pixels
	X int `toml:"x"`

	// Y position of the window in pixels
	Y int `toml:"y"`
}

// Default configuration values
const (
	// DefaultSMTPAddress is the default SMTP server listen address
	DefaultSMTPAddress = "127.0.0.1:1025"

	// DefaultIMAPAddress is the default IMAP server listen address
	DefaultIMAPAddress = "127.0.0.1:1143"

	// DefaultTheme is the default UI theme
	DefaultTheme = "system"

	// DefaultLanguage is the default UI language
	DefaultLanguage = "en"

	// ConfigFileName is the name of the configuration file
	ConfigFileName = "config.toml"

	// Default window dimensions (compact for efficient use)
	DefaultWindowWidth  = 1080
	DefaultWindowHeight = 800
	DefaultWindowX      = -1 // -1 means center on screen
	DefaultWindowY      = -1 // -1 means center on screen

	// Window constraints for security validation
	MinWindowWidth  = 1080
	MinWindowHeight = 800
	MaxWindowWidth  = 4096
	MaxWindowHeight = 2160
)

// DefaultPeers is the list of default Yggdrasil network peers
var DefaultPeers = []string{
	"tcp://bra.zbin.eu:7743",
}

// peerAddressRegex validates peer address format according to Yggdrasil documentation
// Supported formats: tcp://, tls://, quic://, socks://, sockstls://, unix://, ws://, wss://
var peerAddressRegex = regexp.MustCompile(
	`^(tcp|tls|quic)://[a-zA-Z0-9.-]+:[0-9]+$|` + // Standard protocols: tcp, tls, quic
		`^(socks|sockstls)://[a-zA-Z0-9.-]+:[0-9]+/[a-zA-Z0-9.-]+:[0-9]+$|` + // SOCKS proxies
		`^unix:///[^\s]+$|` + // Unix sockets
		`^(ws|wss)://[a-zA-Z0-9.-]+:[0-9]+(/[^\s]*)?$`, // WebSockets with optional path
)

// GetConfigDir returns the platform-specific configuration directory path
// Windows: %APPDATA%\Tyr
// Linux: ~/.config/tyr
func GetConfigDir() (string, error) {
	// Get user config directory (cross-platform)
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	// Platform-specific subdirectory structure
	var configDir string
	if filepath.Separator == '\\' {
		// Windows: %APPDATA%\Tyr
		configDir = filepath.Join(userConfigDir, "Tyr")
	} else {
		// Linux: ~/.config/tyr
		configDir = filepath.Join(userConfigDir, "tyr")
	}

	return configDir, nil
}

// EnsureConfigDir creates the configuration directory if it doesn't exist
// Sets appropriate permissions (0700 for user-only access)
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create directory with user-only permissions (rwx------)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// Load reads the configuration from the TOML file
// Creates a default configuration if the file doesn't exist
// Thread-safe and validates configuration values
func Load() (*Config, error) {
	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return nil, err
	}

	// Get config file path
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(configDir, ConfigFileName)

	// Check if config file exists
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		// Create default configuration
		config := newDefaultConfig()
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse TOML
	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and apply defaults for missing values
	config.applyDefaults()

	return &config, nil
}

// Save writes the configuration to the TOML file
// Creates the config directory if it doesn't exist
// Thread-safe with write lock to prevent concurrent modifications
func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	// Get config file path
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, ConfigFileName)

	// Marshal config to TOML
	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config to TOML: %w", err)
	}

	// Write to file with user-only read/write permissions (rw-------)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Configuration saved successfully to: %s\n", configPath)
	fmt.Printf("  - Theme: %s\n", c.UIPreferences.Theme)
	fmt.Printf("  - Language: %s\n", c.UIPreferences.Language)

	return nil
}

// GetPassword retrieves the password from the OS keyring
// Returns empty string if no password is set
// Thread-safe and uses OS-specific secure storage
func (c *Config) GetPassword() (string, error) {
	password, err := GetPassword(KeyringService, KeyringUsername)
	if err != nil {
		// Check if it's just not found (valid case for new installations)
		if strings.Contains(err.Error(), "not found") {
			return "", nil
		}
		return "", fmt.Errorf("failed to retrieve password from keyring: %w", err)
	}
	return password, nil
}

// SetPassword stores the password in the OS keyring
// Validates password meets minimum requirements (6 characters)
// Thread-safe and uses OS-specific secure storage
func (c *Config) SetPassword(password string) error {
	// Validate password length
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	// Save to keyring
	if err := SavePassword(KeyringService, KeyringUsername, password); err != nil {
		return fmt.Errorf("failed to save password to keyring: %w", err)
	}

	return nil
}

// AddPeer adds a new peer to the configuration
// Validates peer address format and prevents duplicates
// If peer already exists, enables it instead of returning error
// Thread-safe with write lock
func (c *Config) AddPeer(address string) error {
	// Validate peer address format
	if !peerAddressRegex.MatchString(address) {
		return fmt.Errorf("invalid peer address format. Supported: tcp://, tls://, quic://, socks://, sockstls://, unix://, ws://, wss://")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if peer already exists
	for i, peer := range c.NetworkPeers {
		if peer.Address == address {
			// Peer already exists - enable it if not already enabled
			if !peer.Enabled {
				c.NetworkPeers[i].Enabled = true
				log.Printf("Peer already exists, enabled it: %s", address)
			} else {
				log.Printf("Peer already exists and is enabled: %s", address)
			}
			return nil
		}
	}

	// Add new peer (enabled by default when added through this method)
	c.NetworkPeers = append(c.NetworkPeers, PeerConfig{
		Address: address,
		Enabled: true,
	})

	return nil
}

// RemovePeer removes a peer from the configuration
// Thread-safe with write lock
func (c *Config) RemovePeer(address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find and remove peer
	found := false
	newPeers := make([]PeerConfig, 0, len(c.NetworkPeers))
	for _, peer := range c.NetworkPeers {
		if peer.Address != address {
			newPeers = append(newPeers, peer)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("peer not found: %s", address)
	}

	c.NetworkPeers = newPeers
	return nil
}

// EnablePeer enables a peer in the configuration
// Thread-safe with write lock
func (c *Config) EnablePeer(address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find and enable peer
	for i := range c.NetworkPeers {
		if c.NetworkPeers[i].Address == address {
			c.NetworkPeers[i].Enabled = true
			return nil
		}
	}

	return fmt.Errorf("peer not found: %s", address)
}

// DisablePeer disables a peer in the configuration
// Thread-safe with write lock
func (c *Config) DisablePeer(address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find and disable peer
	for i := range c.NetworkPeers {
		if c.NetworkPeers[i].Address == address {
			c.NetworkPeers[i].Enabled = false
			return nil
		}
	}

	return fmt.Errorf("peer not found: %s", address)
}

// GetEnabledPeers returns a list of enabled peer addresses
// Thread-safe with read lock
func (c *Config) GetEnabledPeers() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	enabled := make([]string, 0, len(c.NetworkPeers))
	for _, peer := range c.NetworkPeers {
		if peer.Enabled {
			enabled = append(enabled, peer.Address)
		}
	}

	return enabled
}

// newDefaultConfig creates a new configuration with default values
func newDefaultConfig() *Config {
	// Get config directory for database path
	configDir, err := GetConfigDir()
	if err != nil {
		// Fallback to current directory if config dir unavailable
		configDir = "."
	}

	// Create default peers list
	defaultPeers := make([]PeerConfig, len(DefaultPeers))
	for i, address := range DefaultPeers {
		defaultPeers[i] = PeerConfig{
			Address: address,
			Enabled: true,
		}
	}

	return &Config{
		OnboardingComplete: false,
		ServiceSettings: ServiceSettings{
			SMTPAddress:  DefaultSMTPAddress,
			IMAPAddress:  DefaultIMAPAddress,
			DatabasePath: filepath.Join(configDir, "yggmail.db"),
		},
		NetworkPeers: defaultPeers,
		UIPreferences: UIPreferences{
			Theme:     DefaultTheme,
			Language:  DefaultLanguage,
			AutoStart: false,
		},
	}
}

// applyDefaults fills in missing configuration values with defaults
// Used when loading older config files that may be missing new fields
func (c *Config) applyDefaults() {
	// Apply service settings defaults
	if c.ServiceSettings.SMTPAddress == "" {
		c.ServiceSettings.SMTPAddress = DefaultSMTPAddress
	}
	if c.ServiceSettings.IMAPAddress == "" {
		c.ServiceSettings.IMAPAddress = DefaultIMAPAddress
	}
	if c.ServiceSettings.DatabasePath == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			configDir = "."
		}
		c.ServiceSettings.DatabasePath = filepath.Join(configDir, "yggmail.db")
	}

	// Apply UI preferences defaults
	if c.UIPreferences.Theme == "" {
		c.UIPreferences.Theme = DefaultTheme
	}
	if c.UIPreferences.Language == "" {
		c.UIPreferences.Language = DefaultLanguage
	}

	// Apply window state defaults if not set
	if c.UIPreferences.WindowState.Width == 0 {
		c.UIPreferences.WindowState.Width = DefaultWindowWidth
	}
	if c.UIPreferences.WindowState.Height == 0 {
		c.UIPreferences.WindowState.Height = DefaultWindowHeight
	}
	if c.UIPreferences.WindowState.X == 0 {
		c.UIPreferences.WindowState.X = DefaultWindowX
	}
	if c.UIPreferences.WindowState.Y == 0 {
		c.UIPreferences.WindowState.Y = DefaultWindowY
	}

	// Validate window dimensions (security: prevent out-of-bounds values)
	c.ValidateWindowState()

	// Ensure at least one peer exists
	if len(c.NetworkPeers) == 0 {
		for _, address := range DefaultPeers {
			c.NetworkPeers = append(c.NetworkPeers, PeerConfig{
				Address: address,
				Enabled: true,
			})
		}
	}
}

// SaveWindowState updates the window state in the configuration
// Validates dimensions before saving
// Thread-safe with write lock
func (c *Config) SaveWindowState(width, height, x, y int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update window state
	c.UIPreferences.WindowState.Width = width
	c.UIPreferences.WindowState.Height = height
	c.UIPreferences.WindowState.X = x
	c.UIPreferences.WindowState.Y = y

	// Validate before saving
	c.validateWindowStateUnsafe()

	return nil
}

// ValidateWindowState validates and corrects window dimensions
// Ensures window is within acceptable bounds and on visible screen area
// Thread-safe with write lock
func (c *Config) ValidateWindowState() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.validateWindowStateUnsafe()
}

// validateWindowStateUnsafe validates window state without locking
// Must be called with write lock held
// Security: Prevents invalid window dimensions and positions
func (c *Config) validateWindowStateUnsafe() {
	ws := &c.UIPreferences.WindowState

	// Validate width (minimum and maximum bounds)
	if ws.Width < MinWindowWidth {
		ws.Width = DefaultWindowWidth
	}
	if ws.Width > MaxWindowWidth {
		ws.Width = MaxWindowWidth
	}

	// Validate height (minimum and maximum bounds)
	if ws.Height < MinWindowHeight {
		ws.Height = DefaultWindowHeight
	}
	if ws.Height > MaxWindowHeight {
		ws.Height = MaxWindowHeight
	}

	// Validate position (prevent negative coordinates that would hide window)
	// Note: -1 is special value meaning "center on screen"
	if ws.X < -1 {
		ws.X = DefaultWindowX
	}
	if ws.Y < -1 {
		ws.Y = DefaultWindowY
	}

	// Prevent extremely large coordinates that would place window off-screen
	// This is a conservative check; actual screen size is checked at runtime
	if ws.X > MaxWindowWidth {
		ws.X = DefaultWindowX
	}
	if ws.Y > MaxWindowHeight {
		ws.Y = DefaultWindowY
	}
}

// ==============================================================================
// Peer Discovery Cache Management
// ==============================================================================

const (
	// CacheTTLHours is the time-to-live for discovered peers cache (24 hours)
	CacheTTLHours = 24
)

// GetCachedDiscoveredPeers returns cached discovered peers if within TTL
// Returns nil if cache is expired or empty
// Thread-safe with read lock
func (c *Config) GetCachedDiscoveredPeers() []DiscoveredPeer {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if cache is empty
	if len(c.CachedDiscoveredPeers) == 0 || c.CacheTimestamp == 0 {
		return nil
	}

	// Check TTL
	currentTime := time.Now().Unix()
	cacheAge := currentTime - c.CacheTimestamp
	cacheTTL := int64(CacheTTLHours * 60 * 60) // Convert hours to seconds

	if cacheAge > cacheTTL {
		log.Printf("Discovered peers cache expired (age: %d hours)", cacheAge/3600)
		return nil
	}

	log.Printf("Retrieved %d cached discovered peers (age: %d minutes)",
		len(c.CachedDiscoveredPeers), cacheAge/60)

	// Return a copy to prevent external modification
	result := make([]DiscoveredPeer, len(c.CachedDiscoveredPeers))
	copy(result, c.CachedDiscoveredPeers)
	return result
}

// CacheDiscoveredPeers stores discovered peers with current timestamp
// Thread-safe with write lock
func (c *Config) CacheDiscoveredPeers(peers []DiscoveredPeer) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Store peers and timestamp
	c.CachedDiscoveredPeers = make([]DiscoveredPeer, len(peers))
	copy(c.CachedDiscoveredPeers, peers)
	c.CacheTimestamp = time.Now().Unix()

	log.Printf("Cached %d discovered peers", len(peers))

	return nil
}

// ClearCachedDiscoveredPeers removes all cached discovered peers
// Thread-safe with write lock
func (c *Config) ClearCachedDiscoveredPeers() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CachedDiscoveredPeers = nil
	c.CacheTimestamp = 0

	log.Println("Cleared discovered peers cache")
}
