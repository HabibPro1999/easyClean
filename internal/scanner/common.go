package scanner

import "path/filepath"

// shouldExcludeDir checks if a directory should be excluded from scanning
// This function is shared between AssetFinder and ReferenceFinder to avoid duplication
func shouldExcludeDir(path, root string, excludePatterns []string) bool {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}

	for _, pattern := range excludePatterns {
		// Check if the directory name matches the pattern
		if filepath.Base(path) == filepath.Base(pattern) {
			return true
		}
		// Check if the relative path matches
		// Error only occurs on malformed pattern (compile-time issue, not runtime)
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return true
		}
	}

	return false
}
