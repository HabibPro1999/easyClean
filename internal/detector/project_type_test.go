package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HabibPro1999/easyClean/internal/models"
)

func TestDetectProjectType_React(t *testing.T) {
	tmpDir := t.TempDir()

	// Create package.json with React dependency
	packageJSON := `{
		"dependencies": {
			"react": "^18.0.0",
			"react-dom": "^18.0.0"
		}
	}`
	writeFile(t, filepath.Join(tmpDir, "package.json"), packageJSON)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeWebReact {
		t.Errorf("Expected React project type, got %v", projectType)
	}
}

func TestDetectProjectType_ReactNative(t *testing.T) {
	tmpDir := t.TempDir()

	packageJSON := `{
		"dependencies": {
			"react": "^18.0.0",
			"react-native": "^0.70.0"
		}
	}`
	writeFile(t, filepath.Join(tmpDir, "package.json"), packageJSON)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeReactNative {
		t.Errorf("Expected React Native project type, got %v", projectType)
	}
}

func TestDetectProjectType_Vue(t *testing.T) {
	tmpDir := t.TempDir()

	packageJSON := `{
		"dependencies": {
			"vue": "^3.0.0"
		}
	}`
	writeFile(t, filepath.Join(tmpDir, "package.json"), packageJSON)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeWebVue {
		t.Errorf("Expected Vue project type, got %v", projectType)
	}
}

func TestDetectProjectType_Angular(t *testing.T) {
	tmpDir := t.TempDir()

	packageJSON := `{
		"dependencies": {
			"@angular/core": "^15.0.0"
		}
	}`
	writeFile(t, filepath.Join(tmpDir, "package.json"), packageJSON)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeWebAngular {
		t.Errorf("Expected Angular project type, got %v", projectType)
	}
}

func TestDetectProjectType_Svelte(t *testing.T) {
	tmpDir := t.TempDir()

	packageJSON := `{
		"devDependencies": {
			"svelte": "^3.0.0"
		}
	}`
	writeFile(t, filepath.Join(tmpDir, "package.json"), packageJSON)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeWebSvelte {
		t.Errorf("Expected Svelte project type, got %v", projectType)
	}
}

func TestDetectProjectType_Flutter(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pubspec.yaml (Flutter marker file)
	pubspec := `name: my_flutter_app
description: A Flutter project
`
	writeFile(t, filepath.Join(tmpDir, "pubspec.yaml"), pubspec)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeFlutter {
		t.Errorf("Expected Flutter project type, got %v", projectType)
	}
}

func TestDetectProjectType_iOS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .xcodeproj directory
	xcodeprojPath := filepath.Join(tmpDir, "MyApp.xcodeproj")
	if err := os.MkdirAll(xcodeprojPath, 0755); err != nil {
		t.Fatalf("Failed to create xcodeproj: %v", err)
	}

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeIOS {
		t.Errorf("Expected iOS project type, got %v", projectType)
	}
}

func TestDetectProjectType_Android(t *testing.T) {
	tmpDir := t.TempDir()

	// Create build.gradle
	buildGradle := `plugins {
		id 'com.android.application'
	}`
	writeFile(t, filepath.Join(tmpDir, "build.gradle"), buildGradle)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeAndroid {
		t.Errorf("Expected Android project type, got %v", projectType)
	}
}

func TestDetectProjectType_Go(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	goMod := `module github.com/example/project

go 1.21
`
	writeFile(t, filepath.Join(tmpDir, "go.mod"), goMod)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeGo {
		t.Errorf("Expected Go project type, got %v", projectType)
	}
}

func TestDetectProjectType_Rust(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Cargo.toml
	cargoToml := `[package]
name = "my_project"
version = "0.1.0"
`
	writeFile(t, filepath.Join(tmpDir, "Cargo.toml"), cargoToml)

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeRust {
		t.Errorf("Expected Rust project type, got %v", projectType)
	}
}

func TestDetectProjectType_Unknown(t *testing.T) {
	tmpDir := t.TempDir()

	// Empty directory - no project markers
	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeUnknown {
		t.Errorf("Expected Unknown project type, got %v", projectType)
	}
}

func TestDetectProjectType_InvalidPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Invalid JSON should fall back to Unknown
	writeFile(t, filepath.Join(tmpDir, "package.json"), "invalid json {{{")

	projectType := DetectProjectType(tmpDir)

	if projectType != models.ProjectTypeUnknown {
		t.Errorf("Expected Unknown project type for invalid JSON, got %v", projectType)
	}
}

func TestDetectProjectType_PriorityOrder(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both package.json and pubspec.yaml
	// Should prioritize package.json check first
	packageJSON := `{
		"dependencies": {
			"react": "^18.0.0"
		}
	}`
	writeFile(t, filepath.Join(tmpDir, "package.json"), packageJSON)
	writeFile(t, filepath.Join(tmpDir, "pubspec.yaml"), "name: test")

	projectType := DetectProjectType(tmpDir)

	// Should detect as React since package.json is checked first
	if projectType != models.ProjectTypeWebReact {
		t.Errorf("Expected React (package.json priority), got %v", projectType)
	}
}

// Helper function to write test files
func writeFile(t *testing.T, path, content string) {
	t.Helper()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
}
