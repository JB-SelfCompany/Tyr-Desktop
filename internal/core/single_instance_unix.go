// +build !windows

package core

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

// acquireLock uses file locking (flock) to ensure single instance on Unix systems
func (si *SingleInstance) acquireLock() (bool, error) {
	// Create or open lock file
	file, err := os.OpenFile(si.lockFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return false, fmt.Errorf("failed to open lock file: %w", err)
	}

	// Try to acquire exclusive lock (non-blocking)
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		// Lock is held by another process
		file.Close()
		log.Println("Another instance of Tyr is already running")
		return false, nil
	}

	// Lock acquired
	si.lockFile = file
	log.Println("Single instance lock acquired (Unix flock)")
	return true, nil
}

// releaseLock releases the file lock
func (si *SingleInstance) releaseLock() error {
	if si.lockFile == nil {
		return nil
	}

	// Unlock and close the file
	syscall.Flock(int(si.lockFile.Fd()), syscall.LOCK_UN)
	err := si.lockFile.Close()
	si.lockFile = nil

	// Remove lock file
	os.Remove(si.lockFilePath)

	return err
}
