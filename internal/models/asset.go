// Package models defines core data structures for the asset cleaner.
//
// Key types:
// - AssetFile: Represents a discovered asset with metadata and usage status
// - Reference: A code location that references an asset
// - ScanResult: Complete output of a project scan
// - ProjectConfig: User configuration settings
// - ProjectType: Detected framework/platform type
package models

import (
	"time"
)

// AssetCategory represents the type of asset
type AssetCategory int

const (
	CategoryImage AssetCategory = iota
	CategoryFont
	CategoryVideo
	CategoryAudio
	CategoryOther
)

// String returns the string representation of AssetCategory
func (ac AssetCategory) String() string {
	return [...]string{
		"Image",
		"Font",
		"Video",
		"Audio",
		"Other",
	}[ac]
}

// AssetStatus represents the usage status of an asset
type AssetStatus int

const (
	StatusUsed AssetStatus = iota
	StatusUnused
	StatusPotentiallyUnused
	StatusNeedsManualReview
)

// String returns the string representation of AssetStatus
func (as AssetStatus) String() string {
	return [...]string{
		"Used",
		"Unused",
		"PotentiallyUnused",
		"NeedsManualReview",
	}[as]
}

// AssetFile represents a single asset file discovered in the project
type AssetFile struct {
	// Identity
	Path         string `json:"path"`
	RelativePath string `json:"relative_path"`
	Name         string `json:"name"`
	Extension    string `json:"extension"`

	// Metadata
	Size    int64     `json:"size_bytes"`
	ModTime time.Time `json:"mod_time"`
	Hash    string    `json:"hash,omitempty"`

	// Classification
	Category AssetCategory `json:"category"`
	Status   AssetStatus   `json:"status"`

	// Usage Information
	References []*Reference `json:"references,omitempty"`
	RefCount   int          `json:"reference_count"`
}

// DetermineCategoryFromExtension returns the asset category based on file extension
func DetermineCategoryFromExtension(ext string) AssetCategory {
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".svg": true, ".webp": true, ".ico": true, ".bmp": true,
	}
	fontExts := map[string]bool{
		".ttf": true, ".woff": true, ".woff2": true, ".eot": true, ".otf": true,
	}
	videoExts := map[string]bool{
		".mp4": true, ".webm": true, ".mov": true, ".avi": true, ".mkv": true,
	}
	audioExts := map[string]bool{
		".mp3": true, ".wav": true, ".ogg": true, ".m4a": true, ".flac": true,
	}

	if imageExts[ext] {
		return CategoryImage
	}
	if fontExts[ext] {
		return CategoryFont
	}
	if videoExts[ext] {
		return CategoryVideo
	}
	if audioExts[ext] {
		return CategoryAudio
	}
	return CategoryOther
}
