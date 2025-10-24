// Package parser - React-specific patterns
//
// Detects asset references in React and React Native projects including:
// - Public folder conventions
// - React.lazy and dynamic imports
// - Next.js specific patterns
// - JSX image references
package parser

import "regexp"

var (
	// React.lazy dynamic imports
	ReactLazyPattern = regexp.MustCompile(`React\.lazy\s*\(\s*\(\s*\)\s*=>\s*import\s*\(\s*['"]([^'"]+)['"]`)

	// Next.js public folder references (implicit /public prefix)
	NextPublicPattern = regexp.MustCompile(`['"]/(images?|assets?|static|fonts?|videos?|media)/([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ico|ttf|woff|woff2|mp4|mp3))['"]`)

	// Dynamic imports with webpack magic comments
	WebpackDynamicImport = regexp.MustCompile(`import\s*\(\s*/\*.*?\*/\s*['"]([^'"]+\.(jpg|jpeg|png|svg))['"]`)
)

// ReactPatternProvider provides patterns for React projects
type ReactPatternProvider struct{}

func (r *ReactPatternProvider) GetPatterns() []ReferencePattern {
	return []ReferencePattern{
		// React-specific patterns
		{Pattern: ReactLazyPattern, Type: "DynamicImport", Confidence: 1.0},
		{Pattern: NextPublicPattern, Type: "PublicFolder", Confidence: 0.95},
		{Pattern: WebpackDynamicImport, Type: "DynamicImport", Confidence: 0.9},

		// Standard patterns
		{Pattern: ImportPattern, Type: "Import", Confidence: 1.0},
		{Pattern: RequirePattern, Type: "Import", Confidence: 1.0},
		{Pattern: CSSUrlPattern, Type: "CSSUrl", Confidence: 0.95},
		{Pattern: HTMLSrcPattern, Type: "HTMLAttribute", Confidence: 0.95},
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.75},
		{Pattern: TemplateLiteralPattern, Type: "TemplateLiteral", Confidence: 0.75},
	}
}

func (r *ReactPatternProvider) UseASTParsing() bool {
	return true // React benefits from AST parsing for JSX
}

func (r *ReactPatternProvider) SupportedFileExtensions() []string {
	return []string{".js", ".jsx", ".ts", ".tsx", ".css", ".scss", ".less"}
}
