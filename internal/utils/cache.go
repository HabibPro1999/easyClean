package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const (
	appName         = "easyClean"
	projectsSubdir  = "projects"
	scanResultsFile = "scan-results.json"
)

// GetUserCacheDir returns the OS-specific cache directory for the application
// Linux/macOS: ~/.cache/easyClean/
// Windows: %LOCALAPPDATA%\easyClean\cache\
func GetUserCacheDir() (string, error) {
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user cache directory: %w", err)
	}

	appCacheDir := filepath.Join(cacheRoot, appName)
	return appCacheDir, nil
}

// GetProjectHash creates a 12-character hash of the absolute project path
// This is used to create unique cache directories for different projects
func GetProjectHash(projectRoot string) (string, error) {
	// Ensure we're using absolute path for consistency
	absPath, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Clean the path to normalize separators
	absPath = filepath.Clean(absPath)

	// Create SHA256 hash
	hash := sha256.Sum256([]byte(absPath))

	// Take first 12 characters of hex encoding
	// Collision probability: ~1 in 281 trillion for 12 hex chars
	hashStr := hex.EncodeToString(hash[:])[:12]

	return hashStr, nil
}

// GetProjectCacheDir returns the cache directory for a specific project
// Format: ~/.cache/asset-cleaner/projects/<hash>/
func GetProjectCacheDir(projectRoot string) (string, error) {
	appCacheDir, err := GetUserCacheDir()
	if err != nil {
		return "", err
	}

	projectHash, err := GetProjectHash(projectRoot)
	if err != nil {
		return "", err
	}

	projectCacheDir := filepath.Join(appCacheDir, projectsSubdir, projectHash)
	return projectCacheDir, nil
}

// GetScanResultsPath returns the full path to the scan results file for a project
func GetScanResultsPath(projectRoot string) (string, error) {
	projectCacheDir, err := GetProjectCacheDir(projectRoot)
	if err != nil {
		return "", err
	}

	return filepath.Join(projectCacheDir, scanResultsFile), nil
}

// EnsureCacheDirExists creates the cache directory if it doesn't exist
func EnsureCacheDirExists(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	return nil
}

// GetScanResultsPathOrDefault returns the scan results path for a project,
// or uses the provided default path if not empty
func GetScanResultsPathOrDefault(projectRoot, defaultPath string) (string, error) {
	if defaultPath != "" {
		return defaultPath, nil
	}
	return GetScanResultsPath(projectRoot)
}
