package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/HabibPro1999/easyClean/internal/models"
	"github.com/HabibPro1999/easyClean/internal/parser"
	"github.com/HabibPro1999/easyClean/internal/utils"
)

// ReferenceFinder scans source files for asset references
type ReferenceFinder struct {
	config          *models.ProjectConfig
	root            string
	patterns        []parser.ReferencePattern
	projectType     models.ProjectType
	patternProvider parser.PatternProvider
}

// NewReferenceFinder creates a new ReferenceFinder instance
func NewReferenceFinder(root string, config *models.ProjectConfig) *ReferenceFinder {
	// Use project type from config, or Unknown if not set
	projectType := config.ProjectType
	if projectType == models.ProjectTypeUnknown && config.AutoDetectProjectType {
		// Auto-detection should be done before creating ReferenceFinder
		// For now, fallback to Unknown
		projectType = models.ProjectTypeUnknown
	}

	provider := parser.GetPatternProvider(projectType)

	return &ReferenceFinder{
		config:          config,
		root:            root,
		patterns:        provider.GetPatterns(),
		projectType:     projectType,
		patternProvider: provider,
	}
}

// FindReferences scans source files and finds references to assets
func (rf *ReferenceFinder) FindReferences() (map[string][]*models.Reference, error) {
	references := make(map[string][]*models.Reference)

	err := filepath.WalkDir(rf.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories and excluded paths
		if d.IsDir() {
			if shouldExcludeDir(path, rf.root, rf.config.ExcludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		// Only scan source files
		if rf.isSourceFile(path) {
			refs, err := rf.scanFile(path)
			if err == nil {
				// Group references by the asset path they reference
				for _, ref := range refs {
					assetPath := rf.resolveAssetPath(ref.MatchedText)
					if assetPath != "" {
						references[assetPath] = append(references[assetPath], ref)
					}
				}
			}
		}

		return nil
	})

	return references, err
}

// sourceExtensions maps file extensions to source code files
// Declared at package level to avoid repeated map creation
var sourceExtensions = map[string]bool{
	".js": true, ".jsx": true, ".ts": true, ".tsx": true,
	".vue": true, ".svelte": true,
	".css": true, ".scss": true, ".sass": true, ".less": true,
	".html": true, ".htm": true,
	".dart": true, ".yaml": true, // Flutter/Dart files
	".swift": true,
	".kt": true, ".java": true,
	".go": true,
	".rs": true,
}

// isSourceFile checks if a file is a source code file
func (rf *ReferenceFinder) isSourceFile(path string) bool {
	ext := filepath.Ext(path)

	// Check framework-specific extensions first
	supportedExts := rf.patternProvider.SupportedFileExtensions()
	for _, supported := range supportedExts {
		if ext == supported {
			return true
		}
	}

	// Fallback to generic source extensions
	return sourceExtensions[ext]
}

// scanFile scans a single file for asset references
func (rf *ReferenceFinder) scanFile(path string) ([]*models.Reference, error) {
	var references []*models.Reference

	// Check if we should use AST parsing for this file
	ext := filepath.Ext(path)
	useAST := rf.patternProvider.UseASTParsing() &&
		(ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx")

	if useAST {
		// Use AST parser for deep analysis
		astParser := parser.NewASTParser(path)
		astRefs, err := astParser.ParseFile()
		if err == nil && len(astRefs) > 0 {
			references = append(references, astRefs...)
		}
		// Continue with regex patterns as fallback/supplement
	}

	// Regex-based scanning (works for all files)
	file, err := os.Open(path)
	if err != nil {
		return references, err // Return AST results if available
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check if line is a comment
		isComment := rf.isCommentLine(line)

		// Try each pattern
		for _, patternDef := range rf.patterns {
			matches := patternDef.Pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 1 {
					ref := &models.Reference{
						SourceFile:  path,
						LineNumber:  lineNumber,
						MatchedText: match[1],
						Context:     strings.TrimSpace(line),
						Type:        rf.stringToRefType(patternDef.Type),
						Confidence:  patternDef.Confidence,
						IsComment:   isComment,
						IsDynamic:   rf.isDynamicReference(line),
					}
					references = append(references, ref)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return references, err
	}

	// De-duplicate references (AST + regex may find same references)
	references = rf.deduplicateReferences(references)

	return references, nil
}

// isCommentLine checks if a line is primarily a comment
func (rf *ReferenceFinder) isCommentLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "//") ||
		strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "/*") ||
		strings.HasPrefix(trimmed, "*") ||
		strings.HasPrefix(trimmed, "<!--")
}

// isDynamicReference checks if a reference appears to be dynamically constructed
func (rf *ReferenceFinder) isDynamicReference(line string) bool {
	// Simple heuristic: check for string concatenation or variable interpolation
	return strings.Contains(line, "+") ||
		strings.Contains(line, "${") ||
		strings.Contains(line, "concat") ||
		strings.Contains(line, "join")
}

// resolveAssetPath attempts to resolve a matched reference to an actual asset path
func (rf *ReferenceFinder) resolveAssetPath(matched string) string {
	cleaned := rf.cleanPath(matched)

	// Try strategies in order
	if path := rf.tryExactMatch(cleaned); path != "" {
		return path
	}
	if path := rf.tryAssetPathMatch(cleaned); path != "" {
		return path
	}
	if path := rf.tryBasenameMatch(cleaned); path != "" {
		return path
	}

	return cleaned
}

// cleanPath removes leading ./ or / from path
func (rf *ReferenceFinder) cleanPath(path string) string {
	cleaned := strings.TrimPrefix(path, "./")
	cleaned = strings.TrimPrefix(cleaned, "/")
	return cleaned
}

// tryExactMatch tries to match the path exactly from project root
func (rf *ReferenceFinder) tryExactMatch(cleaned string) string {
	fullPath := filepath.Join(rf.root, cleaned)
	if utils.Exists(fullPath) {
		return fullPath
	}

	// If it's a directory reference (ends with /), return it as-is
	if strings.HasSuffix(cleaned, "/") {
		return filepath.Join(rf.root, cleaned)
	}

	return ""
}

// tryAssetPathMatch tries to find the asset in configured asset paths
func (rf *ReferenceFinder) tryAssetPathMatch(cleaned string) string {
	for _, assetPath := range rf.config.AssetPaths {
		// Try with asset path prefix
		fullPath := filepath.Join(rf.root, assetPath, cleaned)
		if utils.Exists(fullPath) {
			return fullPath
		}

		// Try removing asset path prefix from cleaned (in case it's already there)
		if strings.HasPrefix(cleaned, assetPath) {
			withoutPrefix := strings.TrimPrefix(cleaned, assetPath)
			withoutPrefix = strings.TrimPrefix(withoutPrefix, "/")
			fullPath := filepath.Join(rf.root, assetPath, withoutPrefix)
			if utils.Exists(fullPath) {
				return fullPath
			}
		}
	}
	return ""
}

// tryBasenameMatch tries to find asset by basename in configured asset paths
func (rf *ReferenceFinder) tryBasenameMatch(cleaned string) string {
	basename := filepath.Base(cleaned)
	for _, assetPath := range rf.config.AssetPaths {
		assetDir := filepath.Join(rf.root, assetPath)
		if !utils.Exists(assetDir) {
			continue
		}

		// Walk the asset directory to find matching basenames
		var foundPath string
		filepath.WalkDir(assetDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && filepath.Base(path) == basename {
				foundPath = path
				return filepath.SkipAll
			}
			return nil
		})

		if foundPath != "" {
			return foundPath
		}
	}
	return ""
}

// stringToRefType converts a string type to ReferenceType
func (rf *ReferenceFinder) stringToRefType(typeStr string) models.ReferenceType {
	switch typeStr {
	case "Import", "DynamicImport":
		return models.RefTypeImport
	case "CSSUrl":
		return models.RefTypeCSSUrl
	case "HTMLAttribute":
		return models.RefTypeHTMLAttribute
	case "TemplateLiteral":
		return models.RefTypeTemplateLiteral
	case "FlutterImageAsset", "FlutterAssetImage", "FlutterAssetLoad":
		return models.RefTypeImport
	case "YAMLAsset":
		return models.RefTypeConfig
	default:
		return models.RefTypeStringLiteral
	}
}

// deduplicateReferences removes duplicate references (same file + line + matched text)
func (rf *ReferenceFinder) deduplicateReferences(refs []*models.Reference) []*models.Reference {
	seen := make(map[string]bool)
	var unique []*models.Reference

	for _, ref := range refs {
		// Create a key from source file, line number, and matched text
		key := filepath.Join(ref.SourceFile, string(rune(ref.LineNumber)), ref.MatchedText)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, ref)
		}
	}

	return unique
}
