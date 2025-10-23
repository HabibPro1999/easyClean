package commands

import (
	"fmt"
	"os"

	"github.com/HabibPro1999/easyClean/internal/config"
	"github.com/HabibPro1999/easyClean/internal/detector"
	"github.com/HabibPro1999/easyClean/internal/models"
	"github.com/HabibPro1999/easyClean/internal/ui"
	"github.com/spf13/cobra"
)

var (
	forceInit bool
	template  string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file for current project",
	Long: `Init creates a .unusedassets.yaml configuration file in the current directory.

The configuration file controls which directories are scanned, which file types
are considered assets, and which paths should be excluded.

Templates:
  default       - Standard configuration for most projects
  minimal       - Minimal configuration (fewer options)
  comprehensive - All available options with comments`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&forceInit, "force", false, "overwrite existing config file")
	initCmd.Flags().StringVar(&template, "template", "default", "config template: default, minimal, comprehensive")
}

func runInit(cmd *cobra.Command, args []string) error {
	configPath := ".unusedassets.yaml"

	// Print header
	if !quiet {
		ui.PrintHeader("Initialize Configuration", "")
	}

	// Check if config already exists
	if config.ConfigExists(configPath) && !forceInit {
		return fmt.Errorf("configuration file already exists. Use --force to overwrite")
	}

	// Detect project type
	projectType := detector.DetectProjectType(".")
	if !quiet {
		fmt.Printf("\n✓ Detected project type: %s\n", projectType)
	}

	// Create configuration based on template
	var cfg *models.ProjectConfig
	switch template {
	case "minimal":
		cfg = createMinimalConfig(projectType)
	case "comprehensive":
		cfg = createComprehensiveConfig(projectType)
	default:
		cfg = createDefaultConfig(projectType)
	}

	// Save configuration
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	if !quiet {
		fmt.Printf("✓ Created %s\n", configPath)
		fmt.Println("\nConfiguration includes:")
		fmt.Printf("  • %d asset directories\n", len(cfg.AssetPaths))
		fmt.Printf("  • %d file extensions\n", len(cfg.Extensions))
		fmt.Printf("  • %d exclusion patterns\n", len(cfg.ExcludePaths))
		fmt.Println("\nEdit .unusedassets.yaml to customize settings.")
		fmt.Println("Run 'asset-cleaner scan' to start scanning.")
	}

	return nil
}

func createMinimalConfig(projectType models.ProjectType) *models.ProjectConfig {
	cfg := config.DefaultConfig()

	// Customize based on project type
	if projectType != models.ProjectTypeUnknown {
		cfg.AssetPaths = config.DefaultAssetPathsForProjectType(projectType)
	} else {
		cfg.AssetPaths = []string{"assets/", "public/"}
	}

	// Minimal extensions
	cfg.Extensions = []string{
		".jpg", ".jpeg", ".png", ".gif", ".svg",
		".ttf", ".woff", ".woff2",
		".mp4", ".mp3",
	}

	// Basic excludes
	cfg.ExcludePaths = []string{
		"node_modules/",
		"dist/",
		"build/",
	}

	return cfg
}

func createDefaultConfig(projectType models.ProjectType) *models.ProjectConfig {
	cfg := config.DefaultConfig()

	// Customize based on project type
	if projectType != models.ProjectTypeUnknown {
		cfg.AssetPaths = config.DefaultAssetPathsForProjectType(projectType)
		cfg.ProjectType = projectType
	}

	return cfg
}

func createComprehensiveConfig(projectType models.ProjectType) *models.ProjectConfig {
	cfg := createDefaultConfig(projectType)

	// Add comprehensive options
	cfg.ConstantFiles = []string{
		"src/constants/assets.ts",
		"src/constants/assets.js",
		"lib/assets.dart",
		"app/config/AssetPaths.swift",
	}

	cfg.BasePathVars = []string{
		"ASSETS_BASE",
		"PUBLIC_URL",
		"ASSET_PREFIX",
		"CDN_URL",
	}

	cfg.FollowSymlinks = false
	cfg.MaxWorkers = 8
	cfg.ShowProgress = true
	cfg.ColorOutput = true

	return cfg
}

// WriteConfigToFile writes a config file with comments (for comprehensive template)
func WriteConfigToFile(cfg *models.ProjectConfig, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write commented YAML
	content := fmt.Sprintf(`# Asset Cleaner Configuration
# Generated for project type: %s

# Directories to scan for asset files
asset_paths:
`, cfg.ProjectType)

	for _, path := range cfg.AssetPaths {
		content += fmt.Sprintf("  - %s\n", path)
	}

	content += "\n# File extensions to consider as assets\nextensions:\n"
	for _, ext := range cfg.Extensions {
		content += fmt.Sprintf("  - %s\n", ext)
	}

	content += "\n# Paths and patterns to exclude from scanning\nexclude_paths:\n"
	for _, path := range cfg.ExcludePaths {
		content += fmt.Sprintf("  - %s\n", path)
	}

	if len(cfg.ConstantFiles) > 0 {
		content += "\n# Asset constant files to analyze\nconstant_files:\n"
		for _, file := range cfg.ConstantFiles {
			content += fmt.Sprintf("  - %s\n", file)
		}
	}

	content += fmt.Sprintf(`
# Advanced settings
max_workers: %d           # Concurrent workers (0 = auto-detect)
follow_symlinks: %t      # Follow symbolic links
show_progress: %t        # Show progress bar
color_output: %t         # Enable colored output
`, cfg.MaxWorkers, cfg.FollowSymlinks, cfg.ShowProgress, cfg.ColorOutput)

	_, err = file.WriteString(content)
	return err
}
