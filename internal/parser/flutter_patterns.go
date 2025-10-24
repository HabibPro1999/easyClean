// Package parser - Flutter-specific patterns
//
// Detects asset references in Flutter/Dart projects including:
// - Image.asset() and AssetImage()
// - pubspec.yaml asset declarations
// - Network vs local image distinction
// - Font family references
package parser

import "regexp"

var (
	// Existing Flutter patterns from patterns.go are already defined
	// We'll add additional ones here

	// Font family references in Flutter
	FlutterFontFamilyPattern = regexp.MustCompile(`fontFamily:\s*['"]([^'"]+)['"]`)

	// Asset path in pubspec.yaml (assets section)
	FlutterPubspecAssetPattern = regexp.MustCompile(`^\s*-\s+([^#\n]+\.(png|jpg|jpeg|gif|svg|webp|ttf|otf))`)

	// NetworkImage pattern (to distinguish from local assets)
	FlutterNetworkImagePattern = regexp.MustCompile(`NetworkImage\s*\(\s*['"]https?://`)
)

// FlutterPatternProvider provides patterns for Flutter projects
type FlutterPatternProvider struct{}

func (f *FlutterPatternProvider) GetPatterns() []ReferencePattern {
	return []ReferencePattern{
		// Flutter-specific patterns (high confidence - these are explicit API calls)
		{Pattern: FlutterImageAssetPattern, Type: "FlutterImageAsset", Confidence: 1.0},
		{Pattern: FlutterAssetImagePattern, Type: "FlutterAssetImage", Confidence: 1.0},
		{Pattern: FlutterAssetLoadPattern, Type: "FlutterAssetLoad", Confidence: 1.0},
		{Pattern: FlutterFontFamilyPattern, Type: "FlutterFont", Confidence: 0.9},
		{Pattern: FlutterPubspecAssetPattern, Type: "PubspecAsset", Confidence: 1.0},

		// Standard patterns (lower confidence for generic matches in Dart)
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.7},
	}
}

func (f *FlutterPatternProvider) UseASTParsing() bool {
	return false // Dart AST parsing not implemented yet, regex sufficient
}

func (f *FlutterPatternProvider) SupportedFileExtensions() []string {
	return []string{".dart", ".yaml"}
}
