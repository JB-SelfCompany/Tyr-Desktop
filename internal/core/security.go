package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/pbkdf2"
)

// Security constants for password storage and encryption
const (
	// KeyringService is the service name used for OS keyring storage
	KeyringService = "Tyr"

	// KeyringUsername is the username used for keyring entries
	KeyringUsername = "default"

	// PBKDF2Iterations is the number of iterations for key derivation
	// Using 100,000 iterations as per OWASP recommendations for 2025
	PBKDF2Iterations = 100000

	// AESKeySize is the key size for AES-256-GCM encryption (32 bytes)
	AESKeySize = 32

	// GCMNonceSize is the nonce size for GCM mode (12 bytes is standard)
	GCMNonceSize = 12

	// SaltSize is the size of the random salt (32 bytes for strong security)
	SaltSize = 32

	// fallbackPasswordFile is the encrypted password file for keyring fallback
	fallbackPasswordFile = ".password.enc"
)

// SavePassword stores a password in the OS keyring
// Falls back to encrypted file storage if keyring is unavailable
// Thread-safe and validates all inputs
func SavePassword(service, username, password string) error {
	// Input validation
	if service == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Try OS keyring first
	err := keyring.Set(service, username, password)
	if err == nil {
		return nil
	}

	// Keyring failed - use encrypted file fallback on Linux only
	// Windows keyring should always work, so we return the error
	if runtime.GOOS == "windows" {
		return fmt.Errorf("failed to save password to Windows Credential Manager: %w", err)
	}

	// Linux fallback: encrypted file with machine-ID based key
	return saveFallbackPassword(service, username, password)
}

// GetPassword retrieves a password from the OS keyring
// Falls back to encrypted file storage if keyring is unavailable
// Thread-safe and validates all inputs
func GetPassword(service, username string) (string, error) {
	// Input validation
	if service == "" {
		return "", fmt.Errorf("service name cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	// Try OS keyring first
	password, err := keyring.Get(service, username)
	if err == nil {
		return password, nil
	}

	// Keyring failed - try encrypted file fallback on Linux only
	if runtime.GOOS == "windows" {
		// On Windows, if keyring fails, we don't have a fallback
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("password not found in Windows Credential Manager")
		}
		return "", fmt.Errorf("failed to retrieve password from Windows Credential Manager: %w", err)
	}

	// Linux fallback: try encrypted file
	return getFallbackPassword(service, username)
}

// DeletePassword removes a password from the OS keyring
// Also removes fallback encrypted file if it exists
// Thread-safe and validates all inputs
func DeletePassword(service, username string) error {
	// Input validation
	if service == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Try to delete from keyring
	err := keyring.Delete(service, username)

	// Also try to delete fallback file on Linux
	if runtime.GOOS == "linux" {
		fallbackErr := deleteFallbackPassword(service, username)
		// If keyring delete succeeded but fallback delete failed, return fallback error
		if err == nil && fallbackErr != nil {
			return fallbackErr
		}
	}

	// Return keyring error if any
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete password: %w", err)
	}

	return nil
}

// DeriveKey derives an encryption key from a password using PBKDF2
// Uses SHA-256 as the hash function and returns a 32-byte key for AES-256
// salt: random salt (must be at least 16 bytes, recommended 32 bytes)
// iterations: number of PBKDF2 iterations (recommended: 100,000)
func DeriveKey(password string, salt []byte, iterations int) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, AESKeySize, sha256.New)
}

// EncryptAESGCM encrypts plaintext using AES-256-GCM with PBKDF2 key derivation
// Returns ciphertext in format: [32-byte salt] + [12-byte nonce] + [encrypted data + 16-byte tag]
// Uses 100,000 PBKDF2 iterations for strong key derivation
// Thread-safe and cryptographically secure
func EncryptAESGCM(plaintext, password string) ([]byte, error) {
	// Input validation
	if plaintext == "" {
		return nil, fmt.Errorf("plaintext cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Generate random salt (32 bytes)
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive encryption key from password using PBKDF2
	key := DeriveKey(password, salt, PBKDF2Iterations)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce (12 bytes for GCM)
	nonce := make([]byte, GCMNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt plaintext
	// GCM automatically appends 16-byte authentication tag
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Construct final output: salt + nonce + ciphertext
	result := make([]byte, 0, SaltSize+GCMNonceSize+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptAESGCM decrypts ciphertext using AES-256-GCM with PBKDF2 key derivation
// Expects ciphertext in format: [32-byte salt] + [12-byte nonce] + [encrypted data + 16-byte tag]
// Uses 100,000 PBKDF2 iterations to derive key from password
// Thread-safe and validates authentication tag
func DecryptAESGCM(ciphertext []byte, password string) (string, error) {
	// Input validation
	if len(ciphertext) < SaltSize+GCMNonceSize+16 {
		return "", fmt.Errorf("ciphertext too small (corrupted or invalid)")
	}
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// Extract salt, nonce, and encrypted data
	salt := ciphertext[:SaltSize]
	nonce := ciphertext[SaltSize : SaltSize+GCMNonceSize]
	encrypted := ciphertext[SaltSize+GCMNonceSize:]

	// Derive decryption key from password using PBKDF2
	key := DeriveKey(password, salt, PBKDF2Iterations)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt and verify authentication tag
	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		// GCM authentication failed - wrong password or corrupted data
		return "", fmt.Errorf("decryption failed (invalid password or corrupted data): %w", err)
	}

	return string(plaintext), nil
}

// getMachineID retrieves a unique machine identifier for Linux fallback encryption
// Tries /etc/machine-id first, then /var/lib/dbus/machine-id
// Returns error if neither file exists (non-Linux systems or misconfigured Linux)
func getMachineID() (string, error) {
	// Try /etc/machine-id first (systemd standard location)
	data, err := os.ReadFile("/etc/machine-id")
	if err == nil {
		return string(data), nil
	}

	// Fallback to /var/lib/dbus/machine-id (older systems)
	data, err = os.ReadFile("/var/lib/dbus/machine-id")
	if err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("failed to read machine ID (not found in /etc/machine-id or /var/lib/dbus/machine-id)")
}

// saveFallbackPassword stores password in encrypted file (Linux keyring fallback)
// Uses machine ID as encryption key to prevent plaintext storage
// File is stored in config directory with 0600 permissions
func saveFallbackPassword(service, username, password string) error {
	// Get machine ID for encryption key
	machineID, err := getMachineID()
	if err != nil {
		return fmt.Errorf("cannot use fallback password storage: %w", err)
	}

	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Encrypt password using machine ID as key
	encrypted, err := EncryptAESGCM(password, machineID)
	if err != nil {
		return fmt.Errorf("failed to encrypt password for fallback storage: %w", err)
	}

	// Write to file with restricted permissions
	fallbackPath := filepath.Join(configDir, fallbackPasswordFile)
	if err := os.WriteFile(fallbackPath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write fallback password file: %w", err)
	}

	return nil
}

// getFallbackPassword retrieves password from encrypted file (Linux keyring fallback)
// Uses machine ID as decryption key
func getFallbackPassword(service, username string) (string, error) {
	// Get machine ID for decryption key
	machineID, err := getMachineID()
	if err != nil {
		return "", fmt.Errorf("cannot use fallback password storage: %w", err)
	}

	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	// Read encrypted file
	fallbackPath := filepath.Join(configDir, fallbackPasswordFile)
	encrypted, err := os.ReadFile(fallbackPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("password not found in fallback storage")
		}
		return "", fmt.Errorf("failed to read fallback password file: %w", err)
	}

	// Decrypt password using machine ID as key
	password, err := DecryptAESGCM(encrypted, machineID)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt fallback password: %w", err)
	}

	return password, nil
}

// deleteFallbackPassword removes encrypted password file (Linux keyring fallback)
func deleteFallbackPassword(service, username string) error {
	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Delete fallback file
	fallbackPath := filepath.Join(configDir, fallbackPasswordFile)
	if err := os.Remove(fallbackPath); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, consider it deleted
			return nil
		}
		return fmt.Errorf("failed to delete fallback password file: %w", err)
	}

	return nil
}
