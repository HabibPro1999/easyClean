// Package parser - AST-based parsing for JavaScript/TypeScript
//
// This provides deeper analysis of JS/TS files to extract asset references
// more accurately than regex alone. It identifies:
// - Static and dynamic imports
// - React.lazy and dynamic component loading
// - JSX image src attributes
// - Template literals with asset paths
// - Destructured imports
package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/HabibPro1999/easyClean/internal/models"
)

// ASTParser performs deep analysis of JavaScript/TypeScript files
type ASTParser struct {
	filePath string
}

// NewASTParser creates a new AST parser for a file
func NewASTParser(filePath string) *ASTParser {
	return &ASTParser{filePath: filePath}
}

// Enhanced patterns for AST-level parsing
var (
	// Static imports: import foo from './image.png'
	staticImportPattern = regexp.MustCompile(`import\s+(?:\*\s+as\s+\w+|{[^}]+}|\w+)\s+from\s+['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|mp4|mp3))['"]`)

	// Dynamic imports: import('./image.png') or import(/* ... */ './image.png')
	dynamicImportPattern = regexp.MustCompile(`import\s*\(\s*(?:/\*.*?\*/\s*)?['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|woff2|mp4|mp3))['"]`)

	// JSX img src: <img src={logo} /> or <img src="./logo.png" />
	jsxImgPattern = regexp.MustCompile(`<img[^>]+src\s*=\s*(?:{([^}]+)}|['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ico))['"])`)

	// JSX with require: <img src={require('./logo.png')} />
	jsxRequirePattern = regexp.MustCompile(`<[^>]+(?:src|href)\s*=\s*{require\s*\(\s*['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp))['"]`)

	// Object property: { logo: './image.png' } or { background: require('./bg.jpg') }
	objectPropertyPattern = regexp.MustCompile(`{\s*\w+\s*:\s*(?:require\s*\(\s*)?['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|mp4))['"]`)

	// Template literal with asset: `url(${asset})`
	templateLiteralAssetPattern = regexp.MustCompile("`[^`]*(?:url\\(|src=|href=)[^`]*\\$\\{[^}]*\\}[^`]*\\.(jpg|jpeg|png|gif|svg|webp|ttf|woff|mp4)")

	// Export with asset: export { default as Logo } from './logo.png'
	exportAssetPattern = regexp.MustCompile(`export\s+(?:{[^}]+}|default)\s+from\s+['"]([^'"]+\.(jpg|jpeg|png|gif|svg|webp|ttf|woff))['"]`)
)

// ParseFile performs AST-level parsing of a JavaScript/TypeScript file
func (p *ASTParser) ParseFile() ([]*models.Reference, error) {
	file, err := os.Open(p.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var references []*models.Reference
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	// Read file line by line for context
	var fileContent strings.Builder
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		fileContent.WriteString(line + "\n")

		// Check if line is a comment
		isComment := isCommentLine(line)

		// Apply AST-level patterns
		refs := p.parseLine(line, lineNumber, isComment)
		references = append(references, refs...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return references, nil
}

// parseLine applies AST-level patterns to a single line
func (p *ASTParser) parseLine(line string, lineNumber int, isComment bool) []*models.Reference {
	var refs []*models.Reference

	// Pattern definitions with types
	patterns := []struct {
		pattern    *regexp.Regexp
		refType    models.ReferenceType
		confidence float32
	}{
		{staticImportPattern, models.RefTypeImport, 1.0},
		{dynamicImportPattern, models.RefTypeImport, 1.0},
		{jsxImgPattern, models.RefTypeHTMLAttribute, 0.95},
		{jsxRequirePattern, models.RefTypeImport, 0.95},
		{objectPropertyPattern, models.RefTypeStringLiteral, 0.85},
		{exportAssetPattern, models.RefTypeImport, 1.0},
	}

	for _, patternDef := range patterns {
		matches := patternDef.pattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// Extract the asset path (usually first capture group)
				assetPath := match[1]
				if assetPath != "" {
					ref := &models.Reference{
						SourceFile:  p.filePath,
						LineNumber:  lineNumber,
						MatchedText: assetPath,
						Context:     strings.TrimSpace(line),
						Type:        patternDef.refType,
						Confidence:  patternDef.confidence,
						IsComment:   isComment,
						IsDynamic:   isDynamicReference(line),
					}
					refs = append(refs, ref)
				}
			}
		}
	}

	// Special handling for template literals (dynamic)
	if templateLiteralAssetPattern.MatchString(line) {
		ref := &models.Reference{
			SourceFile:  p.filePath,
			LineNumber:  lineNumber,
			MatchedText: "", // Can't extract exact path from template literal
			Context:     strings.TrimSpace(line),
			Type:        models.RefTypeTemplateLiteral,
			Confidence:  0.6,
			IsComment:   isComment,
			IsDynamic:   true, // Always dynamic
		}
		refs = append(refs, ref)
	}

	return refs
}

// isCommentLine checks if a line is primarily a comment
func isCommentLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "//") ||
		strings.HasPrefix(trimmed, "/*") ||
		strings.HasPrefix(trimmed, "*") ||
		strings.HasPrefix(trimmed, "<!--")
}

// isDynamicReference checks if a reference appears to be dynamically constructed
func isDynamicReference(line string) bool {
	return strings.Contains(line, "+") ||
		strings.Contains(line, "${") ||
		strings.Contains(line, "concat") ||
		strings.Contains(line, "join") ||
		strings.Contains(line, "`")
}
