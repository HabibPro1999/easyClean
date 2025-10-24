package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/HabibPro1999/easyClean/internal/ui"
	"github.com/HabibPro1999/easyClean/internal/utils"
	"github.com/spf13/cobra"
)

var (
	port       int
	host       string
	noBrowser  bool
	listServers bool
	killPort   int
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
- Export results

Multiple projects can run review servers simultaneously on different ports.
Use --list to see all active servers, or --kill to stop a specific server.`,
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().IntVar(&port, "port", 3000, "preferred HTTP server port (auto-increments if taken)")
	reviewCmd.Flags().StringVar(&host, "host", "localhost", "HTTP server host")
	reviewCmd.Flags().BoolVar(&noBrowser, "no-browser", false, "don't auto-open browser")
	reviewCmd.Flags().StringVar(&scanFile, "scan-file", "", "load scan results from JSON file (default: scan-results.json)")
	reviewCmd.Flags().BoolVar(&listServers, "list", false, "list all active review servers")
	reviewCmd.Flags().IntVar(&killPort, "kill", 0, "stop server running on specified port")
}

func runReview(cmd *cobra.Command, args []string) error {
	// Handle --list flag
	if listServers {
		return listActiveServers()
	}

	// Handle --kill flag
	if killPort > 0 {
		return killServerOnPort(killPort)
	}

	// Print header
	if !quiet {
		ui.PrintHeader("Asset Cleaner Review UI", "")
	}

	// Get project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Auto-discover scan file if not specified
	if scanFile == "" {
		// Get cache path for this project
		cachePath, err := utils.GetScanResultsPath(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to get cache path: %w", err)
		}

		scanFile = cachePath

		// Check if scan results exist in cache
		if _, err := os.Stat(scanFile); err != nil {
			return fmt.Errorf("no scan results found in cache for this project.\n" +
				"Run 'easyClean scan' first, or use --scan-file to specify a custom file.\n" +
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

	// Find available port
	actualPort, err := utils.FindAvailablePort(port)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	if actualPort != port && !quiet {
		fmt.Printf("\n⚠️  Port %d is already in use, using port %d instead\n", port, actualPort)
	}

	// Create server
	server, err := ui.NewReviewServer(result, host, actualPort)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	serverURL := fmt.Sprintf("http://%s:%d", host, actualPort)

	if !quiet {
		fmt.Printf("\n🌐 Starting server at %s\n", serverURL)
	}

	// Register server
	serverInfo := utils.ServerInfo{
		ProjectPath: projectRoot,
		ProjectName: filepath.Base(projectRoot),
		Port:        actualPort,
		PID:         os.Getpid(),
		StartTime:   time.Now(),
	}

	if err := utils.RegisterServer(serverInfo); err != nil {
		if !quiet {
			fmt.Printf("⚠️  Warning: failed to register server: %v\n", err)
		}
	}

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start()
	}()

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

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		// Graceful shutdown
		if !quiet {
			fmt.Println("\n\n🛑 Shutting down gracefully...")
		}

		// Unregister server
		if err := utils.UnregisterServer(os.Getpid()); err != nil {
			if !quiet {
				fmt.Printf("⚠️  Warning: failed to unregister server: %v\n", err)
			}
		}

		// Give server 5 seconds to shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		if !quiet {
			fmt.Println("✓ Server stopped successfully")
		}

		return nil

	case err := <-serverErr:
		// Server error
		utils.UnregisterServer(os.Getpid()) // Best effort cleanup
		return fmt.Errorf("server error: %w", err)
	}
}

func listActiveServers() error {
	// Cleanup dead servers first
	if err := utils.CleanupDeadServers(); err != nil {
		return fmt.Errorf("failed to cleanup dead servers: %w", err)
	}

	servers, err := utils.GetActiveServers()
	if err != nil {
		return fmt.Errorf("failed to get active servers: %w", err)
	}

	if len(servers) == 0 {
		fmt.Println("No active review servers")
		fmt.Println("\nStart a server with: easyClean review")
		return nil
	}

	fmt.Println("Active Review Servers:")

	// Print table header
	fmt.Println("┌─────────────────────────────────┬──────┬─────────┬──────────┐")
	fmt.Println("│ Project                         │ Port │ PID     │ Uptime   │")
	fmt.Println("├─────────────────────────────────┼──────┼─────────┼──────────┤")

	// Print each server
	for _, server := range servers {
		uptime := time.Since(server.StartTime)
		uptimeStr := formatUptime(uptime)
		projectName := server.ProjectName
		if len(projectName) > 31 {
			projectName = projectName[:28] + "..."
		}
		fmt.Printf("│ %-31s │ %4d │ %-7d │ %-8s │\n",
			projectName, server.Port, server.PID, uptimeStr)
	}

	fmt.Println("└─────────────────────────────────┴──────┴─────────┴──────────┘")

	// Print access URLs
	fmt.Println("\nOpen any server:")
	for _, server := range servers {
		fmt.Printf("  http://localhost:%d  (%s)\n", server.Port, server.ProjectName)
	}

	fmt.Println("\nStop a server:")
	fmt.Println("  easyClean review --kill <port>")

	return nil
}

func killServerOnPort(port int) error {
	server, err := utils.GetServerByPort(port)
	if err != nil {
		return fmt.Errorf("server on port %d not found or not active", port)
	}

	// Try to kill the process
	process, err := os.FindProcess(server.PID)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	// Send SIGTERM (graceful shutdown)
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send shutdown signal: %w", err)
	}

	// Unregister from registry
	if err := utils.UnregisterServerByPort(port); err != nil {
		// Best effort - process might already be dead
		fmt.Printf("⚠️  Warning: failed to unregister server: %v\n", err)
	}

	fmt.Printf("✓ Server on port %d (%s) stopped successfully\n", port, server.ProjectName)
	return nil
}

func formatUptime(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
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
