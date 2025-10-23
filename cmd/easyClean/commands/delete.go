package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HabibPro1999/easyClean/internal/models"
	"github.com/HabibPro1999/easyClean/internal/ui"
	"github.com/HabibPro1999/easyClean/internal/utils"
	"github.com/spf13/cobra"
)

var (
	dryRun      bool
	interactive bool
	force       bool
	scanFile    string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [paths...]",
	Short: "Delete unused assets from filesystem",
	Long: `Delete removes unused assets from the filesystem.

By default, it deletes all unused assets from the last scan. You can also
specify individual paths to delete.

Safety features:
- Dry-run mode to preview deletions
- Confirmation prompts before deleting
- Git repository detection
- Recovery instructions`,
	RunE: runDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be deleted without deleting")
	deleteCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "prompt for confirmation before each file")
	deleteCmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompts")
	deleteCmd.Flags().StringVar(&scanFile, "scan-file", "", "load scan results from JSON file (default: scan-results.json)")
}

func runDelete(cmd *cobra.Command, args []string) error {
	if !quiet {
		ui.PrintHeader("Delete Unused Assets", "")
	}

	result, err := loadScanResultsOrFail()
	if err != nil {
		return err
	}

	filesToDelete := selectFilesToDelete(result, args)
	if len(filesToDelete) == 0 {
		if !quiet {
			fmt.Println("\n✓ No files to delete")
		}
		return nil
	}

	if dryRun {
		return showDryRun(filesToDelete, calculateTotalSize(filesToDelete))
	}

	isGitRepo := isGitRepository(result.ProjectRoot)

	if !force && !confirmDeletion(filesToDelete, isGitRepo) {
		if !quiet {
			fmt.Println("\n⊘ Deletion cancelled")
		}
		return nil
	}

	if interactive {
		return deleteInteractive(filesToDelete, isGitRepo)
	}

	return deleteBatch(filesToDelete, isGitRepo)
}

// loadScanResultsOrFail loads scan results or returns error with helpful message
func loadScanResultsOrFail() (*models.ScanResult, error) {
	if scanFile == "" {
		// Get current working directory
		projectRoot, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}

		// Get cache path for this project
		cachePath, err := utils.GetScanResultsPath(projectRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to get cache path: %w", err)
		}

		scanFile = cachePath

		// Check if scan results exist in cache
		if _, err := os.Stat(scanFile); err != nil {
			return nil, fmt.Errorf("no scan results found in cache for this project.\n" +
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
		return nil, fmt.Errorf("failed to load scan results: %w", err)
	}

	return result, nil
}

// selectFilesToDelete determines which files should be deleted based on args
func selectFilesToDelete(result *models.ScanResult, args []string) []models.AssetFile {
	if len(args) > 0 {
		return filterAssetsByPaths(result.UnusedAssets, args)
	}
	return result.UnusedAssets
}

// calculateTotalSize computes total size of asset files
func calculateTotalSize(files []models.AssetFile) int64 {
	totalSize := int64(0)
	for _, asset := range files {
		totalSize += asset.Size
	}
	return totalSize
}

// confirmDeletion shows warnings and prompts for confirmation
func confirmDeletion(files []models.AssetFile, isGitRepo bool) bool {
	if !quiet {
		totalSize := calculateTotalSize(files)
		fmt.Printf("\nFound %d unused assets (%s)\n\n", len(files), ui.FormatBytes(totalSize))

		if isGitRepo {
			fmt.Println("⚠️  You are about to delete files. Files will remain in git history.")
		} else {
			fmt.Println("⚠️  WARNING: Not in a git repository. Deletions are PERMANENT!")
			fmt.Println("   Consider backing up files before deletion.")
		}

		fmt.Println()
	}

	if !interactive {
		confirmed, err := promptConfirmation(fmt.Sprintf("Delete %d files?", len(files)))
		if err != nil || !confirmed {
			return false
		}
	}

	return true
}

func loadScanResults(path string) (*models.ScanResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result models.ScanResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func filterAssetsByPaths(assets []models.AssetFile, paths []string) []models.AssetFile {
	var filtered []models.AssetFile
	for _, asset := range assets {
		for _, path := range paths {
			if asset.Path == path || asset.RelativePath == path {
				filtered = append(filtered, asset)
				break
			}
		}
	}
	return filtered
}

func isGitRepository(root string) bool {
	gitDir := filepath.Join(root, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

func showDryRun(files []models.AssetFile, totalSize int64) error {
	if !quiet {
		fmt.Println("\n🧪 Dry Run Mode - No files will be deleted")
		fmt.Println("Would delete:")

		for _, asset := range files {
			fmt.Printf("  • %s (%s)\n", asset.RelativePath, ui.FormatBytes(asset.Size))
		}

		fmt.Printf("\nTotal: %d files (%s)\n", len(files), ui.FormatBytes(totalSize))
		fmt.Println("\nRun without --dry-run to actually delete files.")
	}
	return nil
}

func promptConfirmation(message string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", message)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

func deleteBatch(files []models.AssetFile, isGitRepo bool) error {
	if !quiet {
		fmt.Println("\nDeleting files...")
	}

	deletedCount, totalFreed, errors := performDeletion(files)

	printDeletionSummary(deletedCount, totalFreed, errors, isGitRepo)

	if len(errors) > 0 {
		return fmt.Errorf("%d files failed to delete", len(errors))
	}

	return nil
}

// performDeletion deletes files and tracks results
func performDeletion(files []models.AssetFile) (int, int64, []string) {
	deletedCount := 0
	totalFreed := int64(0)
	var errors []string

	for _, asset := range files {
		if err := os.Remove(asset.Path); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", asset.RelativePath, err))
			if !quiet {
				fmt.Printf("  ✗ %s (error)\n", asset.RelativePath)
			}
		} else {
			deletedCount++
			totalFreed += asset.Size
			if !quiet && verbose {
				fmt.Printf("  ✓ %s (%s)\n", asset.RelativePath, ui.FormatBytes(asset.Size))
			}
		}
	}

	return deletedCount, totalFreed, errors
}

// printDeletionSummary shows results and next steps
func printDeletionSummary(deletedCount int, totalFreed int64, errors []string, isGitRepo bool) {
	if !quiet {
		fmt.Println("\n" + strings.Repeat("━", 45))
		if deletedCount > 0 {
			fmt.Printf("\n✅ Deleted %d files (%s freed)\n", deletedCount, ui.FormatBytes(totalFreed))
		}

		if len(errors) > 0 {
			fmt.Printf("\n⚠️  %d errors occurred:\n", len(errors))
			for _, err := range errors {
				fmt.Printf("  • %s\n", err)
			}
		}

		if isGitRepo && deletedCount > 0 {
			printGitNextSteps()
		}
	}
}

// printGitNextSteps shows git recovery instructions
func printGitNextSteps() {
	fmt.Println("\nNext steps:")
	fmt.Println("  git add -u")
	fmt.Println("  git commit -m \"Remove unused assets\"")
	fmt.Println("\nTo recover deleted files:")
	fmt.Println("  git checkout HEAD -- <file-path>")
}

func deleteInteractive(files []models.AssetFile, isGitRepo bool) error {
	if !quiet {
		fmt.Println("\nInteractive deletion mode (y=yes, n=no, q=quit):")
	}

	deletedCount, skippedCount, totalFreed := promptAndDeleteFiles(files)

	printInteractiveSummary(deletedCount, skippedCount, totalFreed, isGitRepo)

	return nil
}

// promptAndDeleteFiles prompts user for each file and performs deletion
func promptAndDeleteFiles(files []models.AssetFile) (int, int, int64) {
	deletedCount := 0
	skippedCount := 0
	totalFreed := int64(0)

	reader := bufio.NewReader(os.Stdin)

	for _, asset := range files {
		action := promptFileAction(reader, asset)

		switch action {
		case "quit":
			if !quiet {
				fmt.Printf("\n⊘ Cancelled (%d files deleted, %d skipped)\n", deletedCount, skippedCount)
			}
			return deletedCount, skippedCount, totalFreed
		case "delete":
			if err := os.Remove(asset.Path); err != nil {
				if !quiet {
					fmt.Printf("  ✗ Error: %v\n", err)
				}
			} else {
				deletedCount++
				totalFreed += asset.Size
				if !quiet {
					fmt.Println("  ✓ Deleted")
				}
			}
		default:
			skippedCount++
			if !quiet {
				fmt.Println("  ⊘ Skipped")
			}
		}

		if !quiet {
			fmt.Print("\n")
		}
	}

	return deletedCount, skippedCount, totalFreed
}

// promptFileAction asks user what to do with a file
func promptFileAction(reader *bufio.Reader, asset models.AssetFile) string {
	if !quiet {
		fmt.Printf("Delete %s (%s)? [y/N/q]: ", asset.RelativePath, ui.FormatBytes(asset.Size))
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		return "skip"
	}

	response = strings.TrimSpace(strings.ToLower(response))

	switch response {
	case "q", "quit":
		return "quit"
	case "y", "yes":
		return "delete"
	default:
		return "skip"
	}
}

// printInteractiveSummary shows interactive deletion results
func printInteractiveSummary(deletedCount, skippedCount int, totalFreed int64, isGitRepo bool) {
	if !quiet {
		fmt.Println(strings.Repeat("━", 45))
		fmt.Printf("\n✅ Deleted %d files (%s freed)\n", deletedCount, ui.FormatBytes(totalFreed))
		fmt.Printf("   Skipped %d files\n", skippedCount)

		if isGitRepo && deletedCount > 0 {
			printGitNextSteps()
		}
	}
}
