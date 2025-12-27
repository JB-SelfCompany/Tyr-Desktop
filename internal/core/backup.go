package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// BackupData represents the structure of an encrypted backup file
// Serialized to JSON before encryption
type BackupData struct {
	// Version is the backup format version for compatibility checking
	Version string `json:"version"`

	// Timestamp is the ISO 8601 timestamp when the backup was created
	Timestamp string `json:"timestamp"`

	// Config contains the complete application configuration
	Config ConfigBackup `json:"config"`

	// Database contains the optional yggmail.db contents (base64 encoded)
	Database string `json:"database,omitempty"`

	// IncludesDatabase indicates if the database is included in this backup
	IncludesDatabase bool `json:"includes_database"`
}

// ConfigBackup represents the configuration data to be backed up
// This is a flattened version of Config for easier serialization
type ConfigBackup struct {
	// OnboardingComplete indicates if onboarding has been completed
	OnboardingComplete bool `json:"onboarding_complete"`

	// SMTPAddress is the SMTP server address
	SMTPAddress string `json:"smtp_address"`

	// IMAPAddress is the IMAP server address
	IMAPAddress string `json:"imap_address"`

	// DatabasePath is the database file path
	DatabasePath string `json:"database_path"`

	// Peers contains the list of network peers
	Peers []PeerBackup `json:"peers"`

	// Theme is the UI theme preference
	Theme string `json:"theme"`

	// Language is the UI language preference
	Language string `json:"language"`

	// AutoStart indicates if auto-start is enabled
	AutoStart bool `json:"auto_start"`

	// Password is the encrypted password (if available)
	// Note: This is encrypted separately by the keyring, not by backup encryption
	Password string `json:"password,omitempty"`
}

// PeerBackup represents a peer configuration for backup
type PeerBackup struct {
	Address string `json:"address"`
	Enabled bool   `json:"enabled"`
}

// Backup format version constants
const (
	// CurrentBackupVersion is the current backup format version
	CurrentBackupVersion = "1.0"

	// BackupFileExtension is the standard file extension for backups
	BackupFileExtension = ".tb"

	// MinBackupPasswordLength is the minimum password length for backups
	MinBackupPasswordLength = 8
)

// CreateBackup creates an encrypted backup of the configuration and optionally the database
// Returns encrypted backup data in format: [32-byte salt] + [12-byte nonce] + [encrypted JSON] + [16-byte tag]
// Uses AES-256-GCM with PBKDF2 key derivation (100,000 iterations)
// Thread-safe and validates all inputs
func CreateBackup(config *Config, includeDB bool, password string) ([]byte, error) {
	// Validate inputs
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if len(password) < MinBackupPasswordLength {
		return nil, fmt.Errorf("backup password must be at least %d characters", MinBackupPasswordLength)
	}

	// Read database if requested
	var databaseB64 string
	if includeDB {
		log.Printf("[Backup] Reading database from: %s", config.ServiceSettings.DatabasePath)
		dbData, err := readDatabase(config.ServiceSettings.DatabasePath)
		if err != nil {
			// Database read failed - include error in backup metadata but continue
			// This allows backing up configuration even if database is inaccessible
			log.Printf("[Backup] WARNING: Failed to read database: %v - backup will not include database", err)
			includeDB = false
			databaseB64 = ""
		} else {
			// Encode database to base64
			databaseB64 = base64.StdEncoding.EncodeToString(dbData)
			log.Printf("[Backup] Database encoded successfully (%d bytes raw, %d bytes base64)", len(dbData), len(databaseB64))
		}
	} else {
		log.Println("[Backup] Database not included in backup (includeDB=false)")
	}

	// Get password from keyring if available
	configPassword, err := config.GetPassword()
	if err != nil {
		// If password retrieval fails, continue without it
		// User will need to reconfigure password after restore
		configPassword = ""
	}

	// Create backup data structure
	backupData := BackupData{
		Version:          CurrentBackupVersion,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		IncludesDatabase: includeDB,
		Database:         databaseB64,
		Config: ConfigBackup{
			OnboardingComplete: config.OnboardingComplete,
			SMTPAddress:        config.ServiceSettings.SMTPAddress,
			IMAPAddress:        config.ServiceSettings.IMAPAddress,
			DatabasePath:       config.ServiceSettings.DatabasePath,
			Theme:              config.UIPreferences.Theme,
			Language:           config.UIPreferences.Language,
			AutoStart:          config.UIPreferences.AutoStart,
			Password:           configPassword,
			Peers:              make([]PeerBackup, 0, len(config.NetworkPeers)),
		},
	}

	// Copy peer configurations
	config.mu.RLock()
	for _, peer := range config.NetworkPeers {
		backupData.Config.Peers = append(backupData.Config.Peers, PeerBackup{
			Address: peer.Address,
			Enabled: peer.Enabled,
		})
	}
	config.mu.RUnlock()

	// Serialize to JSON
	jsonData, err := json.Marshal(backupData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize backup data: %w", err)
	}

	// Encrypt using AES-256-GCM
	encrypted, err := EncryptAESGCM(string(jsonData), password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt backup: %w", err)
	}

	return encrypted, nil
}

// RestoreBackup decrypts and restores a backup file
// Returns the restored configuration and optional database bytes
// Validates backup version compatibility and data integrity
// Thread-safe and provides detailed error messages
func RestoreBackup(data []byte, password string) (*Config, []byte, error) {
	// Validate inputs
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("backup data is empty")
	}
	if len(password) < MinBackupPasswordLength {
		return nil, nil, fmt.Errorf("backup password must be at least %d characters", MinBackupPasswordLength)
	}

	// Decrypt backup data
	decrypted, err := DecryptAESGCM(data, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt backup (invalid password or corrupted file): %w", err)
	}

	// Parse JSON
	var backupData BackupData
	if err := json.Unmarshal([]byte(decrypted), &backupData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse backup data (corrupted backup): %w", err)
	}

	// Validate backup version
	if !isCompatibleVersion(backupData.Version) {
		return nil, nil, fmt.Errorf("incompatible backup version %s (current version: %s)", backupData.Version, CurrentBackupVersion)
	}

	// Create config from backup data
	config := &Config{
		OnboardingComplete: backupData.Config.OnboardingComplete,
		ServiceSettings: ServiceSettings{
			SMTPAddress:  backupData.Config.SMTPAddress,
			IMAPAddress:  backupData.Config.IMAPAddress,
			DatabasePath: backupData.Config.DatabasePath,
		},
		UIPreferences: UIPreferences{
			Theme:     backupData.Config.Theme,
			Language:  backupData.Config.Language,
			AutoStart: backupData.Config.AutoStart,
		},
		NetworkPeers: make([]PeerConfig, 0, len(backupData.Config.Peers)),
	}

	// Restore peer configurations
	for _, peer := range backupData.Config.Peers {
		config.NetworkPeers = append(config.NetworkPeers, PeerConfig{
			Address: peer.Address,
			Enabled: peer.Enabled,
		})
	}

	// Apply defaults for any missing values
	config.applyDefaults()

	// Decode database if included
	var databaseBytes []byte
	log.Printf("[Restore] Backup includesDatabase=%v, database field length=%d", backupData.IncludesDatabase, len(backupData.Database))
	if backupData.IncludesDatabase && backupData.Database != "" {
		log.Println("[Restore] Decoding database from base64...")
		databaseBytes, err = base64.StdEncoding.DecodeString(backupData.Database)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode database from backup: %w", err)
		}
		log.Printf("[Restore] Database decoded successfully (%d bytes)", len(databaseBytes))
	} else {
		log.Println("[Restore] No database to restore from backup")
	}

	// Restore password to keyring if available
	if backupData.Config.Password != "" {
		if err := config.SetPassword(backupData.Config.Password); err != nil {
			// Password restore failed, but continue with config restore
			// User will need to set password manually
		}
	}

	return config, databaseBytes, nil
}

// WriteBackupFile writes encrypted backup data to a file
// Sets appropriate file permissions (0600 for user-only access)
// Thread-safe and validates path
func WriteBackupFile(path string, data []byte) error {
	// Validate inputs
	if path == "" {
		return fmt.Errorf("backup file path cannot be empty")
	}
	if len(data) == 0 {
		return fmt.Errorf("backup data is empty")
	}

	// Write file with restricted permissions (user read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// ReadBackupFile reads encrypted backup data from a file
// Validates file exists and is readable
// Thread-safe
func ReadBackupFile(path string) ([]byte, error) {
	// Validate input
	if path == "" {
		return nil, fmt.Errorf("backup file path cannot be empty")
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("backup file does not exist: %s", path)
		}
		return nil, fmt.Errorf("failed to access backup file: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("backup path is not a regular file: %s", path)
	}

	// Read file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("backup file is empty")
	}

	return data, nil
}

// GenerateBackupFilename generates a default backup filename with timestamp
// Format: tbackup-dd-mm-yy.tb
func GenerateBackupFilename() string {
	now := time.Now()
	day := fmt.Sprintf("%02d", now.Day())
	month := fmt.Sprintf("%02d", now.Month())
	year := fmt.Sprintf("%02d", now.Year()%100)
	return fmt.Sprintf("tbackup-%s-%s-%s%s", day, month, year, BackupFileExtension)
}

// VerifyBackupPassword verifies a backup password without full restoration
// Useful for validating password before proceeding with restore
// Thread-safe
func VerifyBackupPassword(data []byte, password string) bool {
	// Try to decrypt - if successful, password is correct
	_, err := DecryptAESGCM(data, password)
	return err == nil
}

// RestoreDatabase writes database bytes to the configured database path
// Overwrites existing database file
// Thread-safe and validates inputs
func RestoreDatabase(config *Config, databaseBytes []byte) error {
	// Validate inputs
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if len(databaseBytes) == 0 {
		return fmt.Errorf("database bytes are empty")
	}
	if config.ServiceSettings.DatabasePath == "" {
		return fmt.Errorf("database path not configured")
	}

	log.Printf("[RestoreDatabase] Writing %d bytes to %s", len(databaseBytes), config.ServiceSettings.DatabasePath)

	// Write database file with user-only read/write permissions
	if err := os.WriteFile(config.ServiceSettings.DatabasePath, databaseBytes, 0600); err != nil {
		log.Printf("[RestoreDatabase] ERROR: Failed to write database: %v", err)
		return fmt.Errorf("failed to write database file: %w", err)
	}

	log.Printf("[RestoreDatabase] Database file written successfully")
	return nil
}

// readDatabase reads the yggmail.db file and returns its contents
// Returns error if file doesn't exist or can't be read
func readDatabase(dbPath string) ([]byte, error) {
	// Validate input
	if dbPath == "" {
		return nil, fmt.Errorf("database path is empty")
	}

	// Check if database file exists
	info, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("database file does not exist: %s", dbPath)
		}
		return nil, fmt.Errorf("failed to access database file: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("database path is not a regular file: %s", dbPath)
	}

	// Read database file
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %w", err)
	}

	return data, nil
}

// isCompatibleVersion checks if a backup version is compatible with current version
// Currently supports version 1.0 only
// Future versions may support backward compatibility
func isCompatibleVersion(version string) bool {
	// For now, only exact version match is supported
	// Future: implement semver comparison for backward compatibility
	return version == CurrentBackupVersion
}

// GetBackupInfo returns metadata about a backup file without full decryption
// Returns timestamp and version information
// Requires correct password to decrypt metadata
func GetBackupInfo(data []byte, password string) (version, timestamp string, includesDB bool, err error) {
	// Decrypt backup
	decrypted, err := DecryptAESGCM(data, password)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to decrypt backup: %w", err)
	}

	// Parse JSON to extract metadata only
	var backupData BackupData
	if err := json.Unmarshal([]byte(decrypted), &backupData); err != nil {
		return "", "", false, fmt.Errorf("failed to parse backup metadata: %w", err)
	}

	return backupData.Version, backupData.Timestamp, backupData.IncludesDatabase, nil
}

// DeleteDatabase deletes the yggmail database file
// Used during key regeneration to force creation of new keys
// Thread-safe and validates path
func DeleteDatabase(dbPath string) error {
	// Validate input
	if dbPath == "" {
		return fmt.Errorf("database path is empty")
	}

	// Check if file exists
	info, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, nothing to delete
			return nil
		}
		return fmt.Errorf("failed to access database file: %w", err)
	}

	// Verify it's a regular file (safety check)
	if !info.Mode().IsRegular() {
		return fmt.Errorf("database path is not a regular file: %s", dbPath)
	}

	// Delete the file
	if err := os.Remove(dbPath); err != nil {
		return fmt.Errorf("failed to delete database file: %w", err)
	}

	return nil
}
