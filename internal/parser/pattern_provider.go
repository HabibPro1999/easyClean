// Package parser provides pattern matching for asset references.
//
// This file implements the PatternProvider system which routes to framework-specific
// patterns based on detected project type. It supports both regex-based pattern
// matching and AST parsing for JavaScript/TypeScript projects.
package parser

import (
	"github.com/HabibPro1999/easyClean/internal/models"
)

// PatternProvider defines the interface for framework-specific pattern providers
type PatternProvider interface {
	// GetPatterns returns the regex patterns for this framework
	GetPatterns() []ReferencePattern

	// UseASTParsing indicates whether AST parsing should be used for this project type
	UseASTParsing() bool

	// SupportedFileExtensions returns file extensions that should be parsed
	SupportedFileExtensions() []string
}

// GetPatternProvider returns the appropriate pattern provider for a project type
func GetPatternProvider(projectType models.ProjectType) PatternProvider {
	switch projectType {
	case models.ProjectTypeWebReact, models.ProjectTypeReactNative:
		return &ReactPatternProvider{}
	case models.ProjectTypeWebAngular:
		return &AngularPatternProvider{}
	case models.ProjectTypeWebVue:
		return &VuePatternProvider{}
	case models.ProjectTypeFlutter:
		return &FlutterPatternProvider{}
	case models.ProjectTypeWebSvelte:
		return &SveltePatternProvider{}
	default:
		// Fallback to generic patterns for unknown types
		return &GenericPatternProvider{}
	}
}

// GenericPatternProvider provides basic patterns for unknown project types
type GenericPatternProvider struct{}

func (g *GenericPatternProvider) GetPatterns() []ReferencePattern {
	return []ReferencePattern{
		{Pattern: ImportPattern, Type: "Import", Confidence: 1.0},
		{Pattern: RequirePattern, Type: "Import", Confidence: 1.0},
		{Pattern: CSSUrlPattern, Type: "CSSUrl", Confidence: 0.95},
		{Pattern: HTMLSrcPattern, Type: "HTMLAttribute", Confidence: 0.95},
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.7},
		{Pattern: TemplateLiteralPattern, Type: "TemplateLiteral", Confidence: 0.7},
	}
}

func (g *GenericPatternProvider) UseASTParsing() bool {
	return false
}

func (g *GenericPatternProvider) SupportedFileExtensions() []string {
	return []string{".js", ".ts", ".jsx", ".tsx", ".css", ".html"}
}

// SveltePatternProvider provides patterns for Svelte projects
// Uses mostly generic patterns with minor enhancements
type SveltePatternProvider struct{}

func (s *SveltePatternProvider) GetPatterns() []ReferencePattern {
	return []ReferencePattern{
		{Pattern: ImportPattern, Type: "Import", Confidence: 1.0},
		{Pattern: RequirePattern, Type: "Import", Confidence: 1.0},
		{Pattern: CSSUrlPattern, Type: "CSSUrl", Confidence: 0.95},
		{Pattern: HTMLSrcPattern, Type: "HTMLAttribute", Confidence: 0.95},
		{Pattern: StringLiteralPattern, Type: "StringLiteral", Confidence: 0.75},
		{Pattern: TemplateLiteralPattern, Type: "TemplateLiteral", Confidence: 0.75},
	}
}

func (s *SveltePatternProvider) UseASTParsing() bool {
	return true // Svelte uses JS/TS, benefit from AST parsing
}

func (s *SveltePatternProvider) SupportedFileExtensions() []string {
	return []string{".js", ".ts", ".svelte", ".css"}
}
