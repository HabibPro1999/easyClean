package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/HabibPro1999/easyClean/internal/ui"
	"github.com/HabibPro1999/easyClean/internal/utils"
	"github.com/spf13/cobra"
)

var (
	port      int
	host      string
	noBrowser bool
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Launch web UI to review unused assets interactively",
	Long: `Review launches a web-based UI for interactively reviewing scan results.

The web interface allows you to:
- Browse unused assets with filtering and sorting
- Preview file details and references
- Select assets for deletion
- Export results`,
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().IntVar(&port, "port", 3000, "HTTP server port")
	reviewCmd.Flags().StringVar(&host, "host", "localhost", "HTTP server host")
	reviewCmd.Flags().BoolVar(&noBrowser, "no-browser", false, "don't auto-open browser")
	reviewCmd.Flags().StringVar(&scanFile, "scan-file", "", "load scan results from JSON file (default: scan-results.json)")
}

func runReview(cmd *cobra.Command, args []string) error {
	// Print header
	if !quiet {
		ui.PrintHeader("Asset Cleaner Review UI", "")
	}

	// Auto-discover scan file if not specified
	if scanFile == "" {
		// Get current working directory
		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Get cache path for this project
		cachePath, err := utils.GetScanResultsPath(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to get cache path: %w", err)
		}

		scanFile = cachePath

		// Check if scan results exist in cache
		if _, err := os.Stat(scanFile); err != nil {
			return fmt.Errorf("no scan results found in cache for this project.\n" +
				"Run 'asset-cleaner scan' first, or use --scan-file to specify a custom file.\n" +
				"Expected cache location: %s", cachePath)
		}

		if !quiet {
			fmt.Printf("\n📂 Loading scan results from cache:\n")
			fmt.Printf("   %s\n", scanFile)
		}
	} else {
		if !quiet {
			fmt.Printf("\n📂 Loading scan results from %s\n", scanFile)
		}
	}

	result, err := loadScanResults(scanFile)
	if err != nil {
		return fmt.Errorf("failed to load scan results: %w", err)
	}

	if !quiet {
		totalToReview := result.Stats.UnusedCount + result.Stats.PotentiallyUnusedCount + result.Stats.NeedsReviewCount
		fmt.Printf("\n🔍 Loaded scan results: %d total assets\n", result.Stats.TotalAssets)
		fmt.Printf("   • %d unused, %d potentially unused, %d needs review (%d total to review)\n",
			result.Stats.UnusedCount,
			result.Stats.PotentiallyUnusedCount,
			result.Stats.NeedsReviewCount,
			totalToReview)
	}

	// Start web server
	serverURL := fmt.Sprintf("http://%s:%d", host, port)

	if !quiet {
		fmt.Printf("\n🌐 Starting server at %s\n", serverURL)
	}

	// Open browser
	if !noBrowser && !quiet {
		fmt.Println("🚀 Opening browser...")
		if err := openBrowser(serverURL); err != nil {
			fmt.Printf("⚠️  Failed to open browser: %v\n", err)
			fmt.Printf("   Please open %s manually\n", serverURL)
		}
	}

	if !quiet {
		fmt.Println("\nPress Ctrl+C to stop server")
		fmt.Println()
	}

	// Start server (this blocks)
	return ui.StartWebServer(result, host, port)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
