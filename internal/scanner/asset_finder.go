// Package scanner provides filesystem traversal and asset discovery functionality.
//
// It includes two main components:
// - AssetFinder: discovers asset files based on file extensions and exclusion rules
// - ReferenceFinder: scans source code files to find asset references
//
// These work together to identify which assets are actually used in a codebase.
package scanner

import (
	"os"
	"path/filepath"

	"github.com/HabibPro1999/easyClean/internal/models"
	"github.com/HabibPro1999/easyClean/internal/utils"
)

// AssetFinder scans the filesystem for asset files
type AssetFinder struct {
	config *models.ProjectConfig
	root   string
}

// NewAssetFinder creates a new AssetFinder instance
func NewAssetFinder(root string, config *models.ProjectConfig) *AssetFinder {
	return &AssetFinder{
		config: config,
		root:   root,
	}
}

// FindAssets walks the filesystem and collects all asset files
func (af *AssetFinder) FindAssets() ([]models.AssetFile, error) {
	assets := []models.AssetFile{}

	err := filepath.WalkDir(af.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Skip paths we can't access
			return nil
		}

		// Skip symlinks unless configured to follow them
		if !af.config.FollowSymlinks && utils.IsSymlink(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this directory should be excluded
		if d.IsDir() {
			if shouldExcludeDir(path, af.root, af.config.ExcludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this file is an asset
		if af.isAssetFile(path) {
			asset, err := af.createAssetFile(path)
			if err == nil {
				assets = append(assets, asset)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return assets, nil
}

// isAssetFile checks if a file is an asset based on extension
func (af *AssetFinder) isAssetFile(path string) bool {
	ext := filepath.Ext(path)
	for _, e := range af.config.Extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// createAssetFile creates an AssetFile struct from a file path
func (af *AssetFinder) createAssetFile(path string) (models.AssetFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return models.AssetFile{}, err
	}

	relPath, err := filepath.Rel(af.root, path)
	if err != nil {
		relPath = path
	}

	ext := filepath.Ext(path)
	name := filepath.Base(path)

	return models.AssetFile{
		Path:         path,
		RelativePath: relPath,
		Name:         name,
		Extension:    ext,
		Size:         info.Size(),
		ModTime:      info.ModTime(),
		Category:     models.DetermineCategoryFromExtension(ext),
		Status:       models.StatusUnused, // Default status, will be updated during classification
		References:   []*models.Reference{},
		RefCount:     0,
	}, nil
}

// CountAssets returns the estimated number of asset files without collecting them
func (af *AssetFinder) CountAssets() (int, error) {
	count := 0

	err := filepath.WalkDir(af.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if shouldExcludeDir(path, af.root, af.config.ExcludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		if af.isAssetFile(path) {
			count++
		}

		return nil
	})

	return count, err
}
