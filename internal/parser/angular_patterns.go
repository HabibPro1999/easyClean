// Package parser - Angular-specific patterns
//
// Detects asset references in Angular projects including:
// - @Component templateUrl and styleUrls
// - Lazy loading with loadChildren
// - Assets in angular.json
// - Template asset references
package parser

import "regexp"

var (
	// @Component decorator templateUrl
	AngularTemplateUrlPattern = regexp.MustCompile(`templateUrl:\s*['"]([^'"]+\.html)['"]`)

	// @Component decorator styleUrls (can be array)
	AngularStyleUrlsPattern = regexp.MustCompile(`styleUrls:\s*\[([^\]]+)\]`)

	// Lazy loading routes with loadChildren
	AngularLazyLoadPattern = regexp.MustCompile(`loadChildren:\s*\(\)\s*=>\s*import\s*\(\s*['"]([^'"]+)['"]`)

	// Assets in templates (Angular-specific syntax)
	AngularTemplateAssetPattern = regexp.MustCompile(`\[src\]=['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ico))['"]`)
)

// AngularPatternProvider provides patterns for Angular projects
type AngularPatternProvider struct{}

func (a *AngularPatternProvider) GetPatterns() []ReferencePattern {
	return []ReferencePattern{
		// Angular-specific patterns
		{Pattern: AngularTemplateUrlPattern, Type: "TemplateUrl", Confidence: 1.0},
		{Pattern: AngularStyleUrlsPattern, Type: "StyleUrls", Confidence: 1.0},
		{Pattern: AngularLazyLoadPattern, Type: "LazyLoad", Confidence: 1.0},
		{Pattern: AngularTemplateAssetPattern, Type: "TemplateBinding", Confidence: 0.95},

		// Standard patterns
		{Pattern: ImportPattern, Type: "Import", Confidence: 1.0},
		{Pattern: RequirePattern, Type: "Import", Confidence: 1.0},
		{Pattern: CSSUrlPattern, Type: "CSSUrl", Confidence: 0.95},
		{Pattern: HTMLSrcPattern, Type: "HTMLAttribute", Confidence: 0.95},
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.75},
		{Pattern: TemplateLiteralPattern, Type: "TemplateLiteral", Confidence: 0.75},
	}
}

func (a *AngularPatternProvider) UseASTParsing() bool {
	return true // Angular uses TypeScript, benefits from AST parsing
}

func (a *AngularPatternProvider) SupportedFileExtensions() []string {
	return []string{".ts", ".html", ".css", ".scss", ".sass", ".less"}
}
