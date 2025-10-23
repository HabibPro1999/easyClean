package classifier

import (
	"testing"

	"github.com/HabibPro1999/easyClean/internal/models"
)

func TestClassifyAsset(t *testing.T) {
	tests := []struct {
		name           string
		references     []*models.Reference
		expectedStatus models.AssetStatus
	}{
		{
			name:           "No references",
			references:     []*models.Reference{},
			expectedStatus: models.StatusUnused,
		},
		{
			name: "Active reference",
			references: []*models.Reference{
				{
					SourceFile:  "src/app.js",
					LineNumber:  10,
					MatchedText: "./assets/logo.png",
					Type:        models.RefTypeImport,
					IsComment:   false,
					IsDynamic:   false,
				},
			},
			expectedStatus: models.StatusUsed,
		},
		{
			name: "Only comment reference",
			references: []*models.Reference{
				{
					SourceFile:  "src/app.js",
					LineNumber:  10,
					MatchedText: "./assets/logo.png",
					Type:        models.RefTypeStringLiteral,
					IsComment:   true,
					IsDynamic:   false,
				},
			},
			expectedStatus: models.StatusPotentiallyUnused,
		},
		{
			name: "Dynamic reference",
			references: []*models.Reference{
				{
					SourceFile:  "src/app.js",
					LineNumber:  10,
					MatchedText: "./assets/logo.png",
					Type:        models.RefTypeTemplateLiteral,
					IsComment:   false,
					IsDynamic:   true,
				},
			},
			expectedStatus: models.StatusNeedsManualReview,
		},
		{
			name: "Mixed references with dynamic",
			references: []*models.Reference{
				{
					SourceFile:  "src/app.js",
					LineNumber:  10,
					MatchedText: "./assets/logo.png",
					Type:        models.RefTypeImport,
					IsComment:   false,
					IsDynamic:   false,
				},
				{
					SourceFile:  "src/util.js",
					LineNumber:  20,
					MatchedText: "./assets/logo.png",
					Type:        models.RefTypeTemplateLiteral,
					IsComment:   false,
					IsDynamic:   true,
				},
			},
			expectedStatus: models.StatusNeedsManualReview,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := &models.AssetFile{
				Path:       "/project/assets/logo.png",
				Name:       "logo.png",
				Extension:  ".png",
				References: tt.references,
				RefCount:   len(tt.references),
			}

			result := ClassifyAsset(asset)
			if result != tt.expectedStatus {
				t.Errorf("ClassifyAsset() = %v, want %v", result, tt.expectedStatus)
			}
		})
	}
}

func TestMatchesAssetPath(t *testing.T) {
	tests := []struct {
		name     string
		asset    models.AssetFile
		refPath  string
		expected bool
	}{
		{
			name: "Exact path match",
			asset: models.AssetFile{
				Path:         "/project/assets/logo.png",
				RelativePath: "assets/logo.png",
				Name:         "logo.png",
			},
			refPath:  "/project/assets/logo.png",
			expected: true,
		},
		{
			name: "Relative path match",
			asset: models.AssetFile{
				Path:         "/project/assets/logo.png",
				RelativePath: "assets/logo.png",
				Name:         "logo.png",
			},
			refPath:  "assets/logo.png",
			expected: true,
		},
		{
			name: "Name only match",
			asset: models.AssetFile{
				Path:         "/project/assets/logo.png",
				RelativePath: "assets/logo.png",
				Name:         "logo.png",
			},
			refPath:  "logo.png",
			expected: true,
		},
		{
			name: "Suffix match",
			asset: models.AssetFile{
				Path:         "/project/src/assets/images/logo.png",
				RelativePath: "src/assets/images/logo.png",
				Name:         "logo.png",
			},
			refPath:  "images/logo.png",
			expected: true,
		},
		{
			name: "No match",
			asset: models.AssetFile{
				Path:         "/project/assets/logo.png",
				RelativePath: "assets/logo.png",
				Name:         "logo.png",
			},
			refPath:  "other.png",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesAssetPath(&tt.asset, tt.refPath)
			if result != tt.expected {
				t.Errorf("matchesAssetPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
