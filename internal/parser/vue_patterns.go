// Package parser - Vue-specific patterns
//
// Detects asset references in Vue projects including:
// - Template <img> and <source> references
// - defineAsyncComponent dynamic imports
// - Static folder conventions (Nuxt)
// - Vue SFC <style> and <script> references
package parser

import "regexp"

var (
	// Vue defineAsyncComponent pattern
	VueAsyncComponentPattern = regexp.MustCompile(`defineAsyncComponent\s*\(\s*\(\s*\)\s*=>\s*import\s*\(\s*['"]([^'"]+)['"]`)

	// Vue template img bindings
	VueTemplateImgPattern = regexp.MustCompile(`<img[^>]+:src=['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ico))['"]`)

	// Vue template require() in bindings
	VueTemplateRequirePattern = regexp.MustCompile(`:src="require\s*\(\s*['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp))['"]`)

	// Nuxt public/static folder references
	NuxtStaticPattern = regexp.MustCompile(`['"]/(_nuxt|static|public)/([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2))['"]`)
)

// VuePatternProvider provides patterns for Vue projects
type VuePatternProvider struct{}

func (v *VuePatternProvider) GetPatterns() []ReferencePattern {
	return []ReferencePattern{
		// Vue-specific patterns
		{Pattern: VueAsyncComponentPattern, Type: "AsyncComponent", Confidence: 1.0},
		{Pattern: VueTemplateImgPattern, Type: "TemplateBinding", Confidence: 0.95},
		{Pattern: VueTemplateRequirePattern, Type: "TemplateRequire", Confidence: 0.95},
		{Pattern: NuxtStaticPattern, Type: "StaticFolder", Confidence: 0.9},

		// Standard patterns
		{Pattern: ImportPattern, Type: "Import", Confidence: 1.0},
		{Pattern: RequirePattern, Type: "Import", Confidence: 1.0},
		{Pattern: CSSUrlPattern, Type: "CSSUrl", Confidence: 0.95},
		{Pattern: HTMLSrcPattern, Type: "HTMLAttribute", Confidence: 0.95},
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.75},
		{Pattern: TemplateLiteralPattern, Type: "TemplateLiteral", Confidence: 0.75},
	}
}

func (v *VuePatternProvider) UseASTParsing() bool {
	return true // Vue uses JS/TS, benefits from AST parsing
}

func (v *VuePatternProvider) SupportedFileExtensions() []string {
	return []string{".js", ".ts", ".vue", ".css", ".scss", ".sass", ".less"}
}
