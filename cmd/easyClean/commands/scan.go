package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/HabibPro1999/easyClean/internal/classifier"
	"github.com/HabibPro1999/easyClean/internal/config"
	"github.com/HabibPro1999/easyClean/internal/detector"
	"github.com/HabibPro1999/easyClean/internal/models"
	"github.com/HabibPro1999/easyClean/internal/scanner"
	"github.com/HabibPro1999/easyClean/internal/ui"
	"github.com/HabibPro1999/easyClean/internal/utils"
	"github.com/spf13/cobra"
)

var (
	extensions  []string
	exclude     []string
	outputFile  string
	format      string
	noProgress  bool
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [directory]",
	Short: "Scan a project directory for unused assets",
	Long: `Scan scans a project directory recursively to find unused asset files.

It identifies assets (images, fonts, videos, etc.) and searches for references
in source code. Assets without references are marked as unused.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringSliceVar(&extensions, "extensions", nil, "asset extensions to scan (e.g., .png,.jpg)")
	scanCmd.Flags().StringSliceVar(&exclude, "exclude", nil, "paths to exclude (glob patterns)")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "export results to file (JSON/CSV based on extension)")
	scanCmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text, json, csv")
	scanCmd.Flags().BoolVar(&noProgress, "no-progress", false, "disable progress bar")
}

func runScan(cmd *cobra.Command, args []string) error {
	// Determine project root
	projectRoot := "."
	if len(args) > 0 {
		projectRoot = args[0]
	}

	// Make path absolute
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if directory exists
	if info, err := os.Stat(absRoot); err != nil || !info.IsDir() {
		return fmt.Errorf("directory does not exist: %s", absRoot)
	}

	// Load configuration from file or use defaults
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with command-line flags
	if len(extensions) > 0 {
		cfg.Extensions = extensions
	}
	if len(exclude) > 0 {
		cfg.ExcludePaths = exclude
	}
	cfg.ShowProgress = !noProgress && !quiet

	// Print header
	if !quiet {
		ui.PrintHeader("easyClean", "1.0.1")
	}

	// Detect project type
	if !quiet {
		fmt.Println("\n🔍 Detecting project type...")
	}
	projectType := detector.DetectProjectType(absRoot)
	if !quiet {
		fmt.Printf("✓ Found: %s\n", projectType)
	}

	// Adjust config based on project type
	if cfg.AutoDetectProjectType && projectType != models.ProjectTypeUnknown {
		cfg.AssetPaths = config.DefaultAssetPathsForProjectType(projectType)
	}

	// Start scan
	startTime := time.Now()

	if !quiet {
		fmt.Println("\n📁 Scanning asset directories...")
	}

	// Find assets
	assetFinder := scanner.NewAssetFinder(absRoot, cfg)
	assets, err := assetFinder.FindAssets()
	if err != nil {
		return fmt.Errorf("failed to scan assets: %w", err)
	}

	if !quiet {
		fmt.Printf("✓ Found %d asset files\n", len(assets))
	}

	// Find references
	if !quiet {
		fmt.Println("\n🔎 Analyzing code references...")
	}

	referenceFinder := scanner.NewReferenceFinder(absRoot, cfg)
	references, err := referenceFinder.FindReferences()
	if err != nil {
		return fmt.Errorf("failed to scan references: %w", err)
	}

	if !quiet {
		fmt.Printf("✓ Found %d references\n", len(references))
	}

	// Match references to assets
	assets = classifier.MatchReferencesToAssets(assets, references)

	// Classify assets
	assets = classifier.ClassifyAssets(assets)

	// Create scan result
	duration := time.Since(startTime)
	result := &models.ScanResult{
		Timestamp:   time.Now(),
		ProjectRoot: absRoot,
		ProjectType: projectType,
		Duration:    duration.Milliseconds(),
		Assets:      assets,
		Config:      cfg,
	}

	// Compute statistics
	result.ComputeStatistics()
	result.PopulateFilteredLists()

	// Display results based on format
	var displayErr error
	switch format {
	case "json":
		displayErr = outputJSON(result, outputFile)
	case "csv":
		displayErr = outputCSV(result, outputFile)
	default:
		displayErr = outputText(result, outputFile)
	}

	// Always auto-save JSON results to cache for review/delete commands
	cachePath, err := utils.GetScanResultsPath(absRoot)
	if err != nil {
		if !quiet {
			fmt.Printf("\n⚠️  Warning: Failed to get cache path: %v\n", err)
		}
	} else {
		// Ensure cache directory exists
		cacheDir := filepath.Dir(cachePath)
		if err := utils.EnsureCacheDirExists(cacheDir); err != nil {
			if !quiet {
				fmt.Printf("\n⚠️  Warning: Failed to create cache directory: %v\n", err)
			}
		} else {
			// Save to cache
			if err := autoSaveJSON(result, cachePath); err != nil {
				if !quiet {
					fmt.Printf("\n⚠️  Warning: Failed to save results to cache: %v\n", err)
				}
			} else if !quiet {
				fmt.Printf("\n💾 Scan results saved to cache:\n")
				fmt.Printf("   %s\n", cachePath)
				fmt.Printf("   Use 'asset-cleaner review' or 'asset-cleaner delete' to proceed\n")
			}
		}
	}

	return displayErr
}

func outputText(result *models.ScanResult, file string) error {
	output := ui.FormatScanResult(result)

	if file != "" {
		return os.WriteFile(file, []byte(output), 0644)
	}

	fmt.Println(output)
	return nil
}

func outputJSON(result *models.ScanResult, file string) error {
	data, err := result.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to generate JSON: %w", err)
	}

	if file != "" {
		return os.WriteFile(file, data, 0644)
	}

	fmt.Println(string(data))
	return nil
}

func outputCSV(result *models.ScanResult, file string) error {
	data, err := result.ToCSV()
	if err != nil {
		return fmt.Errorf("failed to generate CSV: %w", err)
	}

	if file != "" {
		return os.WriteFile(file, []byte(data), 0644)
	}

	fmt.Println(data)
	return nil
}

// autoSaveJSON saves scan results to JSON file silently (for review/delete commands)
func autoSaveJSON(result *models.ScanResult, filename string) error {
	data, err := result.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
