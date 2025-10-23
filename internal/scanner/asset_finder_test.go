package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HabibPro1999/easyClean/internal/config"
	"github.com/HabibPro1999/easyClean/internal/models"
)

func TestAssetFinder_FindAssets(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create test assets
	createTestFile(t, filepath.Join(tmpDir, "logo.png"))
	createTestFile(t, filepath.Join(tmpDir, "icon.svg"))
	createTestFile(t, filepath.Join(tmpDir, "assets", "photo.jpg"))
	createTestFile(t, filepath.Join(tmpDir, "script.js")) // Not an asset

	cfg := config.DefaultConfig()
	cfg.Extensions = []string{".png", ".svg", ".jpg"}
	cfg.ExcludePaths = []string{}

	finder := NewAssetFinder(tmpDir, cfg)
	assets, err := finder.FindAssets()

	if err != nil {
		t.Fatalf("FindAssets() failed: %v", err)
	}

	if len(assets) != 3 {
		t.Errorf("Expected 3 assets, got %d", len(assets))
	}

	// Verify asset properties
	for _, asset := range assets {
		if asset.Path == "" {
			t.Error("Asset path is empty")
		}
		if asset.Size == 0 {
			t.Error("Asset size is 0")
		}
		if asset.Category == models.CategoryOther {
			t.Errorf("Asset %s has invalid category", asset.Name)
		}
	}
}

func TestAssetFinder_ExcludeDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create assets in different directories
	createTestFile(t, filepath.Join(tmpDir, "public", "logo.png"))
	createTestFile(t, filepath.Join(tmpDir, "node_modules", "lib.png"))
	createTestFile(t, filepath.Join(tmpDir, "dist", "bundle.png"))

	cfg := config.DefaultConfig()
	cfg.Extensions = []string{".png"}
	cfg.ExcludePaths = []string{"node_modules/", "dist/"}

	finder := NewAssetFinder(tmpDir, cfg)
	assets, err := finder.FindAssets()

	if err != nil {
		t.Fatalf("FindAssets() failed: %v", err)
	}

	// Should only find logo.png, not the excluded ones
	if len(assets) != 1 {
		t.Errorf("Expected 1 asset (excluded 2), got %d", len(assets))
	}

	if len(assets) > 0 && assets[0].Name != "logo.png" {
		t.Errorf("Expected logo.png, got %s", assets[0].Name)
	}
}

func TestAssetFinder_CountAssets(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, filepath.Join(tmpDir, "a.png"))
	createTestFile(t, filepath.Join(tmpDir, "b.png"))
	createTestFile(t, filepath.Join(tmpDir, "c.jpg"))

	cfg := config.DefaultConfig()
	cfg.Extensions = []string{".png", ".jpg"}

	finder := NewAssetFinder(tmpDir, cfg)
	count, err := finder.CountAssets()

	if err != nil {
		t.Fatalf("CountAssets() failed: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestAssetFinder_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := config.DefaultConfig()
	finder := NewAssetFinder(tmpDir, cfg)
	assets, err := finder.FindAssets()

	if err != nil {
		t.Fatalf("FindAssets() failed: %v", err)
	}

	if len(assets) != 0 {
		t.Errorf("Expected 0 assets in empty directory, got %d", len(assets))
	}
}

func TestAssetFinder_CategoryDetection(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		filename string
		category models.AssetCategory
	}{
		{"logo.png", models.CategoryImage},
		{"icon.svg", models.CategoryImage},
		{"font.ttf", models.CategoryFont},
		{"video.mp4", models.CategoryVideo},
		{"audio.mp3", models.CategoryAudio},
	}

	for _, tt := range tests {
		createTestFile(t, filepath.Join(tmpDir, tt.filename))
	}

	cfg := config.DefaultConfig()
	finder := NewAssetFinder(tmpDir, cfg)
	assets, err := finder.FindAssets()

	if err != nil {
		t.Fatalf("FindAssets() failed: %v", err)
	}

	if len(assets) != len(tests) {
		t.Fatalf("Expected %d assets, got %d", len(tests), len(assets))
	}

	// Create map of filename to asset
	assetMap := make(map[string]models.AssetFile)
	for _, asset := range assets {
		assetMap[asset.Name] = asset
	}

	for _, tt := range tests {
		asset, found := assetMap[tt.filename]
		if !found {
			t.Errorf("Asset %s not found", tt.filename)
			continue
		}
		if asset.Category != tt.category {
			t.Errorf("Asset %s: expected category %v, got %v",
				tt.filename, tt.category, asset.Category)
		}
	}
}

// Helper function to create test files
func createTestFile(t *testing.T, path string) {
	t.Helper()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
}
