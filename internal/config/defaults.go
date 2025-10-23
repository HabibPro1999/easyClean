// Package config handles configuration loading and default values.
//
// Configuration can be loaded from .unusedassets.yaml or defaults are provided.
// Supports project-type-specific defaults for asset paths and extensions.
package config

import "github.com/HabibPro1999/easyClean/internal/models"

// DefaultConfig returns the default project configuration
func DefaultConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		AssetPaths: []string{
			"assets/",
			"public/",
			"static/",
			"src/assets/",
		},
		Extensions: []string{
			// Images
			".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico", ".bmp",
			// Fonts
			".ttf", ".woff", ".woff2", ".eot", ".otf",
			// Videos
			".mp4", ".webm", ".mov", ".avi", ".mkv",
			// Audio
			".mp3", ".wav", ".ogg", ".m4a", ".flac",
		},
		ExcludePaths: []string{
			"node_modules/",
			"dist/",
			"build/",
			".next/",
			"target/",
			"vendor/",
			".git/",
			"__pycache__/",
			// Platform-specific assets (always exclude)
			"android/app/src/main/res/",
			"android/app/src/main/",
			"ios/Runner/Assets.xcassets/",
			"ios/Runner/Assets/",
			"ios/",
			"android/",
		},
		ConstantFiles:         []string{},
		BasePathVars:          []string{},
		CustomPatterns:        []string{},
		FollowSymlinks:        false,
		AutoDetectProjectType: true,
		ProjectType:           models.ProjectTypeUnknown,
		MaxWorkers:            0, // Auto-detect
		MemoryLimit:           0, // No limit
		Verbose:               false,
		ShowProgress:          true,
		ColorOutput:           true,
	}
}

// projectAssetPaths maps project types to their default asset paths
var projectAssetPaths = map[models.ProjectType][]string{
	models.ProjectTypeWebReact:    {"public/", "src/assets/", "static/"},
	models.ProjectTypeWebVue:      {"public/", "src/assets/"},
	models.ProjectTypeWebAngular:  {"src/assets/"},
	models.ProjectTypeWebSvelte:   {"static/", "src/assets/"},
	models.ProjectTypeReactNative: {"assets/", "src/assets/"},
	models.ProjectTypeFlutter:     {"assets/", "lib/assets/"},
	models.ProjectTypeIOS:         {"Assets.xcassets/", "Resources/"},
	models.ProjectTypeAndroid:     {"res/drawable/", "res/raw/", "assets/"},
	models.ProjectTypeGo:          {"assets/", "static/", "web/"},
	models.ProjectTypeRust:        {"assets/", "static/", "resources/"},
}

// defaultAssetPaths is the fallback for unknown project types
var defaultAssetPaths = []string{"assets/", "public/", "static/"}

// DefaultAssetPathsForProjectType returns default asset paths based on project type
func DefaultAssetPathsForProjectType(pt models.ProjectType) []string {
	if paths, ok := projectAssetPaths[pt]; ok {
		return paths
	}
	return defaultAssetPaths
}

// baseExtensions are the common asset extensions used by most project types
var baseExtensions = []string{
	".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp",
	".ttf", ".woff", ".woff2",
	".mp4", ".webm",
	".mp3", ".wav", ".ogg",
}

// projectSpecificExtensions maps project types to additional extensions
var projectSpecificExtensions = map[models.ProjectType][]string{
	models.ProjectTypeIOS:     {".heic", ".caf", ".aiff"},
	models.ProjectTypeAndroid: {".9.png", ".xml"},
}

// DefaultExtensionsForProjectType returns relevant extensions based on project type
func DefaultExtensionsForProjectType(pt models.ProjectType) []string {
	// Start with base extensions
	extensions := make([]string, len(baseExtensions))
	copy(extensions, baseExtensions)

	// Add platform-specific extensions if available
	if specific, ok := projectSpecificExtensions[pt]; ok {
		extensions = append(extensions, specific...)
	}

	return extensions
}
