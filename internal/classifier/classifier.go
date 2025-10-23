// Package classifier determines asset usage status based on code references.
//
// It provides a conservative classification system:
// - Used: Has active code references
// - Unused: No references found
// - PotentiallyUnused: Only referenced in comments
// - NeedsManualReview: Dynamic path construction detected
package classifier

import "github.com/HabibPro1999/easyClean/internal/models"

// ClassifyAsset determines the status of an asset based on its references
func ClassifyAsset(asset *models.AssetFile) models.AssetStatus {
	if len(asset.References) == 0 {
		return models.StatusUnused
	}

	hasActiveRef := false
	allInComments := true
	hasDynamicRef := false

	for _, ref := range asset.References {
		// If any reference is dynamic, mark for manual review
		if ref.IsDynamic {
			hasDynamicRef = true
		}

		// Track if we have non-comment references
		if !ref.IsComment {
			allInComments = false
			hasActiveRef = true
		}
	}

	// Conservative approach: dynamic references need manual review
	if hasDynamicRef {
		return models.StatusNeedsManualReview
	}

	// All references are in comments
	if allInComments {
		return models.StatusPotentiallyUnused
	}

	// Has active (non-comment) references
	if hasActiveRef {
		return models.StatusUsed
	}

	return models.StatusUnused
}

// ClassifyAssets classifies multiple assets at once
func ClassifyAssets(assets []models.AssetFile) []models.AssetFile {
	for i := range assets {
		assets[i].Status = ClassifyAsset(&assets[i])
	}
	return assets
}

// MatchReferencesToAssets matches found references to asset files
func MatchReferencesToAssets(assets []models.AssetFile, references map[string][]*models.Reference) []models.AssetFile {
	// Match references to assets using path matching
	for i := range assets {
		for refPath, refs := range references {
			if matchesAssetPath(&assets[i], refPath) {
				assets[i].References = append(assets[i].References, refs...)
				assets[i].RefCount = len(assets[i].References)
				break
			}
		}
	}

	return assets
}

// matchesAssetPath checks if a reference path matches an asset
func matchesAssetPath(asset *models.AssetFile, refPath string) bool {
	// Try exact matches
	if asset.Path == refPath || asset.RelativePath == refPath || asset.Name == refPath {
		return true
	}

	// Try suffix matching (e.g., "images/logo.png" matches "src/assets/images/logo.png")
	if len(refPath) > 0 && len(asset.Path) >= len(refPath) {
		suffix := asset.Path[len(asset.Path)-len(refPath):]
		if suffix == refPath {
			return true
		}
	}

	// Try relative path suffix matching
	if len(refPath) > 0 && len(asset.RelativePath) >= len(refPath) {
		suffix := asset.RelativePath[len(asset.RelativePath)-len(refPath):]
		if suffix == refPath {
			return true
		}
	}

	return false
}
