// Package ui provides terminal output formatting and presentation utilities.
//
// Includes functions for:
// - Formatted headers and separators
// - Scan result summaries
// - Human-readable byte formatting
// - Asset list display
package ui

import (
	"fmt"
	"strings"

	"github.com/HabibPro1999/easyClean/internal/models"
)

const (
	headerWidth           = 45
	maxDisplayedAssets    = 10
	bytesPerKilobyte      = 1024
	separatorWidth        = 45
)

// PrintHeader prints the application header
func PrintHeader(name, version string) {
	width := headerWidth
	topLine := "╭" + strings.Repeat("─", width-2) + "╮"
	bottomLine := "╰" + strings.Repeat("─", width-2) + "╯"

	title := fmt.Sprintf("🔍 %s v%s", name, version)
	padding := (width - len(title) - 2) / 2
	titleLine := "│" + strings.Repeat(" ", padding) + title + strings.Repeat(" ", width-len(title)-padding-2) + "│"

	fmt.Println(topLine)
	fmt.Println(titleLine)
	fmt.Println(bottomLine)
}

// FormatScanResult formats a scan result as text output
func FormatScanResult(result *models.ScanResult) string {
	var sb strings.Builder

	// Print separator
	separator := strings.Repeat("━", separatorWidth)
	sb.WriteString("\n" + separator + "\n\n")

	// Print summary
	sb.WriteString("📊 Scan Complete\n\n")
	sb.WriteString(fmt.Sprintf("  Total Assets:           %d\n", result.Stats.TotalAssets))
	sb.WriteString(fmt.Sprintf("  ✓ Used Assets:          %d\n", result.Stats.TotalAssets-result.Stats.UnusedCount-result.Stats.PotentiallyUnusedCount-result.Stats.NeedsReviewCount))
	sb.WriteString(fmt.Sprintf("  ⚠️  Unused Assets:       %d\n", result.Stats.UnusedCount))

	if result.Stats.PotentiallyUnusedCount > 0 {
		sb.WriteString(fmt.Sprintf("  🤔 Potentially Unused:  %d\n", result.Stats.PotentiallyUnusedCount))
	}

	if result.Stats.NeedsReviewCount > 0 {
		sb.WriteString(fmt.Sprintf("  👀 Needs Review:        %d\n", result.Stats.NeedsReviewCount))
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  💾 Potential Savings:   %s\n", FormatBytes(result.Stats.UnusedSize)))
	sb.WriteString(fmt.Sprintf("  ⏱️  Scan Duration:        %.2fs\n", float64(result.Duration)/1000.0))

	sb.WriteString("\n" + separator + "\n")

	// Show unused assets if any
	if result.Stats.UnusedCount > 0 {
		sb.WriteString("\n📝 Unused Assets:\n\n")
		count := 0
		for _, asset := range result.UnusedAssets {
			if count >= maxDisplayedAssets {
				remaining := result.Stats.UnusedCount - count
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", remaining))
				break
			}
			sb.WriteString(fmt.Sprintf("  • %s (%s)\n", asset.RelativePath, FormatBytes(asset.Size)))
			count++
		}
	}

	sb.WriteString("\n✨ Run 'asset-cleaner review' to inspect unused assets\n")
	sb.WriteString("✨ Run 'asset-cleaner delete --dry-run' to preview deletion\n")

	return sb.String()
}

// FormatBytes formats bytes as a human-readable string
func FormatBytes(bytes int64) string {
	const unit = bytesPerKilobyte
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatAssetList formats a list of assets for display
func FormatAssetList(assets []models.AssetFile) string {
	var sb strings.Builder

	for _, asset := range assets {
		sb.WriteString(fmt.Sprintf("%s (%s, %d refs)\n",
			asset.RelativePath,
			FormatBytes(asset.Size),
			asset.RefCount))
	}

	return sb.String()
}
