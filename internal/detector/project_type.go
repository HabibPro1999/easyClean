// Package detector identifies project types by inspecting filesystem markers.
//
// It supports detection of 10+ project types including React, Vue, Flutter,
// iOS, Android, Go, and Rust by looking for characteristic files like
// package.json, pubspec.yaml, go.mod, etc.
package detector

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/HabibPro1999/easyClean/internal/models"
)

// PackageJSON represents a minimal package.json structure for detection
type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// DetectProjectType attempts to detect the project type from filesystem markers
func DetectProjectType(root string) models.ProjectType {
	// Check for package.json (JavaScript/TypeScript projects)
	if pkg, err := readPackageJSON(filepath.Join(root, "package.json")); err == nil {
		return detectFromPackageJSON(pkg)
	}

	// Check for pubspec.yaml (Flutter)
	if fileExists(filepath.Join(root, "pubspec.yaml")) {
		return models.ProjectTypeFlutter
	}

	// Check for .xcodeproj (iOS)
	if hasXcodeProject(root) {
		return models.ProjectTypeIOS
	}

	// Check for build.gradle (Android)
	if fileExists(filepath.Join(root, "build.gradle")) ||
	   fileExists(filepath.Join(root, "build.gradle.kts")) ||
	   fileExists(filepath.Join(root, "app/build.gradle")) {
		return models.ProjectTypeAndroid
	}

	// Check for go.mod (Go)
	if fileExists(filepath.Join(root, "go.mod")) {
		return models.ProjectTypeGo
	}

	// Check for Cargo.toml (Rust)
	if fileExists(filepath.Join(root, "Cargo.toml")) {
		return models.ProjectTypeRust
	}

	return models.ProjectTypeUnknown
}

// detectFromPackageJSON determines project type from package.json dependencies
func detectFromPackageJSON(pkg *PackageJSON) models.ProjectType {
	allDeps := make(map[string]bool)
	for dep := range pkg.Dependencies {
		allDeps[dep] = true
	}
	for dep := range pkg.DevDependencies {
		allDeps[dep] = true
	}

	// Check for React Native first (has both react and react-native)
	if allDeps["react-native"] {
		return models.ProjectTypeReactNative
	}

	// Check for React
	if allDeps["react"] {
		return models.ProjectTypeWebReact
	}

	// Check for Vue
	if allDeps["vue"] {
		return models.ProjectTypeWebVue
	}

	// Check for Angular
	if allDeps["@angular/core"] {
		return models.ProjectTypeWebAngular
	}

	// Check for Svelte
	if allDeps["svelte"] {
		return models.ProjectTypeWebSvelte
	}

	return models.ProjectTypeUnknown
}

// readPackageJSON reads and parses a package.json file
func readPackageJSON(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// hasXcodeProject checks if there's an .xcodeproj directory in the root
func hasXcodeProject(root string) bool {
	entries, err := os.ReadDir(root)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() && filepath.Ext(entry.Name()) == ".xcodeproj" {
			return true
		}
	}

	return false
}
