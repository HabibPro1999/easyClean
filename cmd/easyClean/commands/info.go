package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HabibPro1999/easyClean/internal/config"
	"github.com/HabibPro1999/easyClean/internal/detector"
	"github.com/HabibPro1999/easyClean/internal/ui"
	"github.com/spf13/cobra"
)

var (
	showConfig bool
	showPaths  bool
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display project information and detection results",
	Long: `Info displays information about the current project including:
- Detected project type
- Configuration file location and status
- Asset directories and file counts
- Excluded paths
- Configured file extensions`,
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().BoolVar(&showConfig, "show-config", false, "display current configuration")
	infoCmd.Flags().BoolVar(&showPaths, "show-paths", false, "list detected asset paths")
}

func runInfo(cmd *cobra.Command, args []string) error {
	// Print header
	if !quiet {
		ui.PrintHeader("Project Information", "")
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Detect project type
	projectType := detector.DetectProjectType(currentDir)

	// Load configuration
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if config file exists
	configStatus := "not found (using defaults)"
	if config.ConfigExists(cfgFile) {
		absPath, _ := config.GetConfigPath(cfgFile)
		configStatus = fmt.Sprintf("found at %s", absPath)
	}

	// Display basic info
	if !quiet {
		fmt.Printf("\n📁 Project Root: %s\n", currentDir)
		fmt.Printf("🏷️  Project Type: %s\n", projectType)
		fmt.Printf("🔧 Config File:  %s\n", configStatus)
	}

	// Show asset directories
	if showPaths || !showConfig {
		fmt.Println("\n📂 Asset Directories:")
		for _, path := range cfg.AssetPaths {
			fullPath := filepath.Join(currentDir, path)
			count := countFilesInDir(fullPath, cfg.Extensions)
			if count > 0 {
				fmt.Printf("  • %s (%d files)\n", path, count)
			} else {
				fmt.Printf("  • %s (not found or empty)\n", path)
			}
		}

		fmt.Println("\n🚫 Excluded Paths:")
		for _, path := range cfg.ExcludePaths {
			fmt.Printf("  • %s\n", path)
		}

		fmt.Println("\n📄 File Extensions:")
		// Group by category
		images := []string{}
		fonts := []string{}
		videos := []string{}
		audio := []string{}
		other := []string{}

		for _, ext := range cfg.Extensions {
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico", ".bmp":
				images = append(images, ext)
			case ".ttf", ".woff", ".woff2", ".eot", ".otf":
				fonts = append(fonts, ext)
			case ".mp4", ".webm", ".mov", ".avi", ".mkv":
				videos = append(videos, ext)
			case ".mp3", ".wav", ".ogg", ".m4a", ".flac":
				audio = append(audio, ext)
			default:
				other = append(other, ext)
			}
		}

		if len(images) > 0 {
			fmt.Printf("  Images: %s\n", strings.Join(images, ", "))
		}
		if len(fonts) > 0 {
			fmt.Printf("  Fonts:  %s\n", strings.Join(fonts, ", "))
		}
		if len(videos) > 0 {
			fmt.Printf("  Videos: %s\n", strings.Join(videos, ", "))
		}
		if len(audio) > 0 {
			fmt.Printf("  Audio:  %s\n", strings.Join(audio, ", "))
		}
		if len(other) > 0 {
			fmt.Printf("  Other:  %s\n", strings.Join(other, ", "))
		}

		fmt.Println("\n✅ Configuration valid")
	}

	// Show full configuration if requested
	if showConfig {
		fmt.Print("\n# Current Configuration\n\n")
		fmt.Println("asset_paths:")
		for _, path := range cfg.AssetPaths {
			fmt.Printf("  - %s\n", path)
		}

		fmt.Println("\nextensions:")
		for _, ext := range cfg.Extensions {
			fmt.Printf("  - %s\n", ext)
		}

		fmt.Println("\nexclude_paths:")
		for _, path := range cfg.ExcludePaths {
			fmt.Printf("  - %s\n", path)
		}

		if len(cfg.ConstantFiles) > 0 {
			fmt.Println("\nconstant_files:")
			for _, file := range cfg.ConstantFiles {
				fmt.Printf("  - %s\n", file)
			}
		}

		if len(cfg.BasePathVars) > 0 {
			fmt.Println("\nbase_path_vars:")
			for _, varName := range cfg.BasePathVars {
				fmt.Printf("  - %s\n", varName)
			}
		}

		fmt.Printf("\nmax_workers: %d\n", cfg.MaxWorkers)
		fmt.Printf("show_progress: %t\n", cfg.ShowProgress)
		fmt.Printf("color_output: %t\n", cfg.ColorOutput)
	}

	return nil
}

func countFilesInDir(dir string, extensions []string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, e := range extensions {
				if ext == e {
					count++
					break
				}
			}
		}
		return nil
	})
	return count
}
