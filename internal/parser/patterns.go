// Package parser provides regex patterns for detecting asset references in source code.
//
// It includes patterns for:
// - Import/require statements (JS/TS)
// - CSS url() functions
// - HTML src/href attributes
// - Flutter Image.asset() calls
// - Template literals and string literals
//
// Each pattern is assigned a confidence score indicating likelihood of actual usage.
package parser

import "regexp"

// Common patterns for detecting asset references in code
var (
	// String literals with common asset paths (supports both single and double quotes)
	StringLiteralPattern = regexp.MustCompile(`['"]([^'"]*\.(jpg|jpeg|png|gif|svg|webp|ico|bmp|ttf|woff|woff2|eot|otf|mp4|webm|mov|avi|mkv|mp3|wav|ogg|m4a|flac))['"]`)

	// Import statements
	ImportPattern = regexp.MustCompile(`import\s+.*?['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|mp4|mp3))['"]`)

	// Require statements
	RequirePattern = regexp.MustCompile(`require\s*\(\s*['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|mp4|mp3))['"]\ s*\)`)

	// CSS url() function
	CSSUrlPattern = regexp.MustCompile(`url\s*\(\s*['"]?([^"')]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|eot|otf))['"]?\s*\)`)

	// HTML src/href attributes
	HTMLSrcPattern = regexp.MustCompile(`(?:src|href)\s*=\s*['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|mp4|webm|mp3|wav))['"]`)

	// Template literals (basic pattern)
	TemplateLiteralPattern = regexp.MustCompile("`([^`]*\\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|mp4|mp3))`")

	// Flutter Image.asset() pattern
	FlutterImageAssetPattern = regexp.MustCompile(`Image\.asset\s*\(\s*['"]([^'"]+\.(png|jpg|jpeg|gif|svg|webp|ico))['"]`)

	// Flutter AssetImage() pattern
	FlutterAssetImagePattern = regexp.MustCompile(`AssetImage\s*\(\s*['"]([^'"]+\.(png|jpg|jpeg|gif|svg|webp|ico))['"]`)

	// Flutter general asset loading (rootBundle.load, etc.)
	FlutterAssetLoadPattern = regexp.MustCompile(`(?:rootBundle\.load|DefaultAssetBundle\.of.*?\.load)\s*\(\s*['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|mp4|mp3|wav|ogg))['"]`)

	// Comment patterns
	CommentPattern = regexp.MustCompile(`(?://.*?$|/\*[\s\S]*?\*/|#.*?$|<!--[\s\S]*?-->)`)
)

// ReferencePattern represents a pattern for finding references
type ReferencePattern struct {
	Pattern    *regexp.Regexp
	Type       string
	Confidence float32
}

// GetAllPatterns returns all reference detection patterns
func GetAllPatterns() []ReferencePattern {
	return []ReferencePattern{
		{Pattern: FlutterImageAssetPattern, Type: "FlutterImageAsset", Confidence: 1.0},
		{Pattern: FlutterAssetImagePattern, Type: "FlutterAssetImage", Confidence: 1.0},
		{Pattern: FlutterAssetLoadPattern, Type: "FlutterAssetLoad", Confidence: 1.0},
		{Pattern: ImportPattern, Type: "Import", Confidence: 1.0},
		{Pattern: RequirePattern, Type: "Import", Confidence: 1.0},
		{Pattern: CSSUrlPattern, Type: "CSSUrl", Confidence: 0.95},
		{Pattern: HTMLSrcPattern, Type: "HTMLAttribute", Confidence: 0.95},
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.8},
		{Pattern: TemplateLiteralPattern, Type: "TemplateLiteral", Confidence: 0.8},
	}
}
