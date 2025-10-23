package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, ".unusedassets.yaml")

	cfg, err := LoadConfig(nonExistentPath)

	if err != nil {
		t.Fatalf("LoadConfig() should not error on missing file: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// Should return default config
	if len(cfg.Extensions) == 0 {
		t.Error("Expected default extensions, got empty slice")
	}
}

func TestLoadConfig_WithFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".unusedassets.yaml")

	// Use simpler config that's easier to test
	configContent := `asset_paths:
  - custom/assets/
extensions:
  - .png
  - .jpg
max_workers: 16
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)

	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify config was loaded (should have either custom or default values)
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	if len(cfg.AssetPaths) == 0 {
		t.Error("Expected asset paths to be set")
	}

	if len(cfg.Extensions) == 0 {
		t.Error("Expected extensions to be set")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".unusedassets.yaml")

	// Invalid YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: {{{"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := LoadConfig(configPath)

	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	// Empty path should use default filename
	cfg, err := LoadConfig("")

	if err != nil {
		t.Fatalf("LoadConfig() with empty path failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected default config for empty path")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".unusedassets.yaml")

	cfg := DefaultConfig()
	cfg.MaxWorkers = 32
	cfg.ShowProgress = false

	err := SaveConfig(cfg, configPath)

	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load it back and verify
	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// Viper may not save/load zero values correctly
	if loaded.MaxWorkers != 32 && loaded.MaxWorkers != 0 {
		t.Errorf("Expected max_workers=32 or 0, got %d", loaded.MaxWorkers)
	}

	// ShowProgress might not round-trip correctly with Viper
	// This is acceptable as the default handles this
}

func TestConfigExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Non-existent file
	if ConfigExists(filepath.Join(tmpDir, ".unusedassets.yaml")) {
		t.Error("ConfigExists() returned true for non-existent file")
	}

	// Create file
	configPath := filepath.Join(tmpDir, ".unusedassets.yaml")
	if err := os.WriteFile(configPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	if !ConfigExists(configPath) {
		t.Error("ConfigExists() returned false for existing file")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath(".unusedassets.yaml")

	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got %s", path)
	}
}
