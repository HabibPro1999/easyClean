package models

import (
	"testing"
)

func TestDetermineCategoryFromExtension(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected AssetCategory
	}{
		{"JPEG image", ".jpg", CategoryImage},
		{"PNG image", ".png", CategoryImage},
		{"SVG image", ".svg", CategoryImage},
		{"TrueType font", ".ttf", CategoryFont},
		{"WOFF font", ".woff", CategoryFont},
		{"WOFF2 font", ".woff2", CategoryFont},
		{"MP4 video", ".mp4", CategoryVideo},
		{"WebM video", ".webm", CategoryVideo},
		{"MP3 audio", ".mp3", CategoryAudio},
		{"WAV audio", ".wav", CategoryAudio},
		{"Unknown extension", ".xyz", CategoryOther},
		{"No extension", "", CategoryOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineCategoryFromExtension(tt.ext)
			if result != tt.expected {
				t.Errorf("DetermineCategoryFromExtension(%q) = %v, want %v", tt.ext, result, tt.expected)
			}
		})
	}
}

func TestAssetCategoryString(t *testing.T) {
	tests := []struct {
		category AssetCategory
		expected string
	}{
		{CategoryImage, "Image"},
		{CategoryFont, "Font"},
		{CategoryVideo, "Video"},
		{CategoryAudio, "Audio"},
		{CategoryOther, "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.category.String()
			if result != tt.expected {
				t.Errorf("AssetCategory.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestAssetStatusString(t *testing.T) {
	tests := []struct {
		status   AssetStatus
		expected string
	}{
		{StatusUsed, "Used"},
		{StatusUnused, "Unused"},
		{StatusPotentiallyUnused, "PotentiallyUnused"},
		{StatusNeedsManualReview, "NeedsManualReview"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("AssetStatus.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}
