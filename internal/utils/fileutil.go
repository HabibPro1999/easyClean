// Package utils provides common file system utility functions.
//
// Includes helpers for:
// - File existence and type checking
// - Extension and pattern matching
// - Symlink detection
// - File size queries
package utils

import (
	"os"
	"path/filepath"
)

// Exists checks if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir checks if the path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile checks if the path is a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// HasExtension checks if the file has one of the given extensions
func HasExtension(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// ShouldExclude checks if a path matches any of the exclude patterns
func ShouldExclude(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		if matchesExcludePattern(path, pattern) {
			return true
		}
	}
	return false
}

// matchesExcludePattern checks if a path matches an exclusion pattern
// Uses three strategies: glob matching, basename matching, and prefix matching
func matchesExcludePattern(path, pattern string) bool {
	// Try glob pattern match
	// Error only occurs on malformed pattern (compile-time issue, not runtime)
	if matched, _ := filepath.Match(pattern, path); matched {
		return true
	}

	// Try basename match (e.g., "node_modules" matches "/path/to/node_modules")
	if filepath.Base(path) == filepath.Base(pattern) {
		return true
	}

	// Try directory match (e.g., "vendor/" matches "/path/vendor/")
	if filepath.Dir(path) == filepath.Dir(pattern) {
		return true
	}

	// Try prefix match (e.g., "dist/" matches "dist/bundle.js")
	if len(path) >= len(pattern) && path[:len(pattern)] == pattern {
		return true
	}

	return false
}

// IsSymlink checks if the path is a symbolic link
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// GetFileSize returns the size of the file in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
