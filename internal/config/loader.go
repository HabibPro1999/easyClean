package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/HabibPro1999/easyClean/internal/models"
	"github.com/spf13/viper"
)

// LoadConfig loads configuration from file or returns defaults
func LoadConfig(configPath string) (*models.ProjectConfig, error) {
	// Check if config file exists
	if configPath == "" {
		configPath = ".unusedassets.yaml"
	}

	// If file doesn't exist, return defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Set up Viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Start with empty config and unmarshal from file
	cfg := &models.ProjectConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults for fields not specified in config
	if len(cfg.AssetPaths) == 0 {
		cfg.AssetPaths = DefaultConfig().AssetPaths
	}
	if len(cfg.Extensions) == 0 {
		cfg.Extensions = DefaultConfig().Extensions
	}
	if len(cfg.ExcludePaths) == 0 {
		cfg.ExcludePaths = DefaultConfig().ExcludePaths
	}

	return cfg, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(cfg *models.ProjectConfig, configPath string) error {
	if configPath == "" {
		configPath = ".unusedassets.yaml"
	}

	v := viper.New()

	// Set all config values
	v.Set("asset_paths", cfg.AssetPaths)
	v.Set("extensions", cfg.Extensions)
	v.Set("exclude_paths", cfg.ExcludePaths)
	v.Set("constant_files", cfg.ConstantFiles)
	v.Set("base_path_vars", cfg.BasePathVars)
	v.Set("custom_patterns", cfg.CustomPatterns)
	v.Set("follow_symlinks", cfg.FollowSymlinks)
	v.Set("auto_detect_project_type", cfg.AutoDetectProjectType)
	v.Set("max_workers", cfg.MaxWorkers)
	v.Set("memory_limit", cfg.MemoryLimit)
	v.Set("show_progress", cfg.ShowProgress)
	v.Set("color_output", cfg.ColorOutput)

	// Write to file
	return v.WriteConfigAs(configPath)
}

// ConfigExists checks if a config file exists
func ConfigExists(configPath string) bool {
	if configPath == "" {
		configPath = ".unusedassets.yaml"
	}
	_, err := os.Stat(configPath)
	return err == nil
}

// GetConfigPath returns the absolute path to the config file
func GetConfigPath(configPath string) (string, error) {
	if configPath == "" {
		configPath = ".unusedassets.yaml"
	}
	return filepath.Abs(configPath)
}
