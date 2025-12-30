package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// StorageStats contains information about storage usage
type StorageStats struct {
	// DatabaseSizeMB is the size of the yggmail.db file in megabytes
	DatabaseSizeMB float64
	// FilesSizeMB is the total size of all stored message files in megabytes
	FilesSizeMB float64
	// TotalSizeMB is the total storage usage in megabytes
	TotalSizeMB float64
}

// GetStorageStats returns storage usage statistics
// Thread-safe with read lock on config
func GetStorageStats(config *Config) (*StorageStats, error) {
	config.mu.RLock()
	dbPath := config.ServiceSettings.DatabasePath
	config.mu.RUnlock()

	stats := &StorageStats{}

	// Get database size
	dbInfo, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Database doesn't exist yet (fresh install)
			stats.DatabaseSizeMB = 0
		} else {
			return nil, fmt.Errorf("failed to stat database: %w", err)
		}
	} else {
		stats.DatabaseSizeMB = float64(dbInfo.Size()) / (1024 * 1024)
	}

	// Get filestore directory size (if it exists)
	// FileStore directory is in the same directory as database: db-path-dir/filestore/
	dbDir := filepath.Dir(dbPath)
	filestoreDir := filepath.Join(dbDir, "filestore")

	var totalFilesSize int64
	err = filepath.WalkDir(filestoreDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If filestore directory doesn't exist, that's OK
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalFilesSize += info.Size()
		}
		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to calculate filestore size: %w", err)
	}

	stats.FilesSizeMB = float64(totalFilesSize) / (1024 * 1024)
	stats.TotalSizeMB = stats.DatabaseSizeMB + stats.FilesSizeMB

	return stats, nil
}
