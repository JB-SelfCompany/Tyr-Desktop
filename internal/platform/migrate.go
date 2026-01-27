package platform

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// MigrationResult contains information about the migration process
type MigrationResult struct {
	Migrated       bool     // Whether migration was performed
	SourceDir      string   // Source directory that was migrated from
	DestDir        string   // Destination directory (portable data dir)
	MigratedFiles  []string // List of files that were migrated
	CleanedUp      bool     // Whether the old directory was cleaned up
	Errors         []error  // Any errors encountered during migration
}

// MigrateFromLegacy checks for data in legacy configuration directories
// and migrates it to the portable data directory next to the executable.
// After successful migration, the old directories are removed.
func MigrateFromLegacy() (*MigrationResult, error) {
	result := &MigrationResult{
		DestDir: GetDataDir(),
	}

	// Get list of legacy directories to check
	legacyDirs := GetLegacyConfigDirs()
	if len(legacyDirs) == 0 {
		log.Println("No legacy directories to check for migration")
		return result, nil
	}

	// Check if destination already has data (skip migration if so)
	destConfigPath := GetConfigPath()
	if _, err := os.Stat(destConfigPath); err == nil {
		log.Printf("Portable data directory already contains config, skipping migration")
		return result, nil
	}

	// Find a legacy directory with data
	var sourceDir string
	for _, dir := range legacyDirs {
		if hasData(dir) {
			sourceDir = dir
			break
		}
	}

	if sourceDir == "" {
		log.Println("No legacy data found to migrate")
		return result, nil
	}

	result.SourceDir = sourceDir
	log.Printf("Found legacy data in: %s", sourceDir)
	log.Printf("Migrating to portable directory: %s", result.DestDir)

	// Ensure destination directory exists
	if err := EnsureDirectories(); err != nil {
		return result, fmt.Errorf("failed to create portable data directory: %w", err)
	}

	// Migrate files
	migratedFiles, errors := migrateDirectory(sourceDir, result.DestDir)
	result.MigratedFiles = migratedFiles
	result.Errors = errors

	if len(migratedFiles) > 0 {
		result.Migrated = true
		log.Printf("Successfully migrated %d files", len(migratedFiles))

		// Clean up old directory
		if len(errors) == 0 {
			if err := os.RemoveAll(sourceDir); err != nil {
				log.Printf("Warning: Failed to remove old directory %s: %v", sourceDir, err)
				result.Errors = append(result.Errors, err)
			} else {
				result.CleanedUp = true
				log.Printf("Cleaned up old directory: %s", sourceDir)
			}
		} else {
			log.Printf("Skipping cleanup due to migration errors")
		}
	}

	if len(result.Errors) > 0 {
		return result, fmt.Errorf("migration completed with %d errors", len(result.Errors))
	}

	return result, nil
}

// hasData checks if the directory exists and contains any files
func hasData(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check for config.toml or yggmail.db
	configPath := filepath.Join(dir, "config.toml")
	dbPath := filepath.Join(dir, "yggmail.db")

	if _, err := os.Stat(configPath); err == nil {
		return true
	}
	if _, err := os.Stat(dbPath); err == nil {
		return true
	}

	return false
}

// migrateDirectory copies all files from source to destination
func migrateDirectory(srcDir, destDir string) ([]string, []error) {
	var migratedFiles []string
	var errors []error

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, fmt.Errorf("error accessing %s: %w", path, err))
			return nil // Continue walking
		}

		// Calculate relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			errors = append(errors, fmt.Errorf("error calculating relative path for %s: %w", path, err))
			return nil
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				errors = append(errors, fmt.Errorf("error creating directory %s: %w", destPath, err))
			}
			return nil
		}

		// Copy file
		if err := copyFile(path, destPath, info.Mode()); err != nil {
			errors = append(errors, fmt.Errorf("error copying %s to %s: %w", path, destPath, err))
		} else {
			migratedFiles = append(migratedFiles, relPath)
			log.Printf("Migrated: %s", relPath)
		}

		return nil
	})

	if err != nil {
		errors = append(errors, fmt.Errorf("error walking source directory: %w", err))
	}

	return migratedFiles, errors
}

// copyFile copies a single file from src to dest
func copyFile(src, dest string, mode os.FileMode) error {
	// Ensure destination directory exists
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return destFile.Sync()
}
