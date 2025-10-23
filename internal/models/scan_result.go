package models

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ScanStatistics holds computed statistics from a scan
type ScanStatistics struct {
	TotalAssets            int     `json:"total_assets"`
	TotalSize              int64   `json:"total_size_bytes"`
	UnusedCount            int     `json:"unused_count"`
	UnusedSize             int64   `json:"unused_size_bytes"`
	PotentiallyUnusedCount int     `json:"potentially_unused_count"`
	NeedsReviewCount       int     `json:"needs_review_count"`
	FilesScanned           int     `json:"files_scanned"`
	ReferencesFound        int     `json:"references_found"`
	AvgScanSpeed           float64 `json:"avg_scan_speed_files_per_sec"`
}

// ScanResult represents the complete output of scanning a project
type ScanResult struct {
	// Metadata
	Timestamp   time.Time   `json:"timestamp"`
	ProjectRoot string      `json:"project_root"`
	ProjectType ProjectType `json:"project_type"`
	Duration    int64       `json:"duration_ms"`

	// Assets
	Assets                  []AssetFile `json:"assets"`
	UsedAssets              []AssetFile `json:"used_assets,omitempty"`
	UnusedAssets            []AssetFile `json:"unused_assets,omitempty"`
	PotentiallyUnusedAssets []AssetFile `json:"potentially_unused_assets,omitempty"`
	NeedsReviewAssets       []AssetFile `json:"needs_review_assets,omitempty"`

	// Statistics
	Stats ScanStatistics `json:"statistics"`

	// Configuration
	Config *ProjectConfig `json:"config,omitempty"`
}

// ComputeStatistics calculates all statistics from the Assets slice
func (sr *ScanResult) ComputeStatistics() {
	sr.Stats = ScanStatistics{
		TotalAssets: len(sr.Assets),
	}

	for _, asset := range sr.Assets {
		sr.Stats.TotalSize += asset.Size
		sr.Stats.ReferencesFound += len(asset.References)

		switch asset.Status {
		case StatusUnused:
			sr.Stats.UnusedCount++
			sr.Stats.UnusedSize += asset.Size
		case StatusPotentiallyUnused:
			sr.Stats.PotentiallyUnusedCount++
		case StatusNeedsManualReview:
			sr.Stats.NeedsReviewCount++
		}
	}

	// Calculate average scan speed
	if sr.Duration > 0 {
		durationSeconds := float64(sr.Duration) / 1000.0
		sr.Stats.AvgScanSpeed = float64(sr.Stats.FilesScanned) / durationSeconds
	}
}

// FilterByStatus returns assets matching the given status
func (sr *ScanResult) FilterByStatus(status AssetStatus) []AssetFile {
	var filtered []AssetFile
	for _, asset := range sr.Assets {
		if asset.Status == status {
			filtered = append(filtered, asset)
		}
	}
	return filtered
}

// PopulateFilteredLists populates the filtered asset lists based on status
func (sr *ScanResult) PopulateFilteredLists() {
	sr.UsedAssets = sr.FilterByStatus(StatusUsed)
	sr.UnusedAssets = sr.FilterByStatus(StatusUnused)
	sr.PotentiallyUnusedAssets = sr.FilterByStatus(StatusPotentiallyUnused)
	sr.NeedsReviewAssets = sr.FilterByStatus(StatusNeedsManualReview)
}

// ToJSON exports the scan result as JSON
func (sr *ScanResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(sr, "", "  ")
}

// ToCSV exports the scan result as CSV
func (sr *ScanResult) ToCSV() (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write header
	header := []string{"Status", "Path", "Size", "Category", "References", "ModTime"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, asset := range sr.Assets {
		row := []string{
			asset.Status.String(),
			asset.RelativePath,
			strconv.FormatInt(asset.Size, 10),
			asset.Category.String(),
			strconv.Itoa(asset.RefCount),
			asset.ModTime.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return builder.String(), nil
}
