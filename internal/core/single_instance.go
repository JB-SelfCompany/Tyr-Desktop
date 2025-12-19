package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// SingleInstance manages ensuring only one instance of the application runs
type SingleInstance struct {
	lockFilePath string
	lockFile     *os.File
	locked       bool
}

// NewSingleInstance creates a new single instance manager
// Uses a lock file in the user's config directory
func NewSingleInstance() (*SingleInstance, error) {
	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Create lock file path
	lockFilePath := filepath.Join(configDir, "tyr.lock")

	si := &SingleInstance{
		lockFilePath: lockFilePath,
		locked:       false,
	}

	return si, nil
}

// Lock attempts to acquire the single instance lock
// Returns true if lock was acquired, false if another instance is running
func (si *SingleInstance) Lock() (bool, error) {
	// Platform-specific lock implementation
	locked, err := si.acquireLock()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	si.locked = locked
	return locked, nil
}

// Unlock releases the single instance lock
func (si *SingleInstance) Unlock() error {
	if !si.locked {
		return nil
	}

	// Platform-specific unlock implementation
	if err := si.releaseLock(); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	si.locked = false
	log.Println("Single instance lock released")
	return nil
}

// IsLocked returns whether the lock is currently held
func (si *SingleInstance) IsLocked() bool {
	return si.locked
}

// GetLockFilePath returns the path to the lock file
func (si *SingleInstance) GetLockFilePath() string {
	return si.lockFilePath
}
