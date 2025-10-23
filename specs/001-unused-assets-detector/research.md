# Research: Unused Assets Detector

**Date**: 2025-10-23
**Feature**: Unused Assets Detector CLI Tool with Smart Detection

## Technical Stack Research

### Language & Version

**Decision**: Go 1.23+

**Rationale**:
- Native compilation produces single binary executables (ideal for CLI distribution)
- Excellent concurrency support for parallel file scanning
- Fast execution speeds critical for scanning large codebases
- Strong standard library for file system operations
- Growing ecosystem of beautiful CLI/TUI libraries (Charm.sh)

**Alternatives Considered**:
- **Rust**: Excellent performance but steeper learning curve, smaller CLI ecosystem
- **Node.js/TypeScript**: Easy JavaScript parsing but requires Node runtime, slower file I/O
- **Python**: Great for prototyping but slower performance, distribution complexity

---

### Core Dependencies

#### 1. CLI Framework: Cobra v1.9.1+

**Decision**: Use Cobra v1.9.1 (latest stable, September 2025)

**Rationale**:
- Industry standard used by Kubernetes, Docker, Hugo, GitHub CLI
- Excellent command/subcommand structure (`scan`, `review`, `delete`, `info`)
- Built-in help generation and command completion
- Works seamlessly with Viper for configuration
- Imported by 181,299+ Go packages (proven stability)

**Implementation Notes**:
- Use cobra-cli generator for initial structure
- Define persistent flags for common options (--config, --verbose)
- Leverage command PreRun hooks for validation

---

#### 2. TUI Framework: Bubble Tea v1.3.10 (stable) or v2.0.0-beta.3

**Decision**: Start with Bubble Tea v1.3.10 (stable), prepare for v2 migration

**Rationale**:
- **v1.3.10** is production-ready (September 2025 release)
- **v2** adds progressive keyboard enhancements (shift+enter, key releases)
- Battle-tested in multiple large projects
- Elm Architecture makes complex state management predictable
- Imported by 9,060+ packages

**Implementation Notes**:
- Use for progress indicators during scan operations
- Interactive file selection in `review` command
- Real-time updates during analysis

**Alternatives Considered**:
- **tview**: More widget-based but less compositional than Bubble Tea
- **termui**: Dashboard-oriented, overkill for our use case

---

#### 3. Styling: Lip Gloss v2.x

**Decision**: Use Lip Gloss v2.x (latest, March 2025)

**Rationale**:
- Declarative CSS-like API familiar to web developers
- True color (24-bit) and 256 color support
- Layout system with padding, margin, borders
- Table component perfect for scan results display
- Designed as Bubble Tea companion

**Implementation Notes**:
- Define style constants for asset categories (used/unused/needs-review)
- Use table component for result summaries
- Apply consistent color scheme (green=used, yellow=warning, red=unused)

---

#### 4. Configuration: Viper v1.21.0

**Decision**: Use Viper v1.21.0 (September 2025)

**Rationale**:
- Supports YAML, JSON, TOML, and other formats
- Automatic `.unusedassets` config file loading
- Environment variable binding
- Integrates with Cobra flags
- Used by 99,496+ packages

**Configuration File Format**:
```yaml
# .unusedassets.yaml
asset_paths:
  - public/
  - src/assets/
  - static/

ignore:
  - "*.test.js"
  - "**/__tests__/**"
  - node_modules/
  - dist/

extensions:
  - .jpg
  - .png
  - .svg
  - .woff
  - .woff2
  - .ttf
  - .mp4
  - .mp3
```

---

### File System Scanning

**Decision**: Use `filepath.WalkDir` from Go 1.16+ standard library, with concurrent wrapper for large projects

**Rationale**:
- `filepath.WalkDir` is faster than `filepath.Walk` (avoids os.Lstat per file)
- For projects >10K files, wrap in concurrent scanner using `fastwalk` pattern
- Respects filesystem boundaries, doesn't follow symlinks by default

**Performance Strategy**:
1. **Small projects (<10K files)**: Standard `filepath.WalkDir`
2. **Large projects (>10K files)**: Concurrent walker with worker pool
3. **Excluded directories**: Skip early via `fs.SkipDir` (node_modules, dist, etc.)

**Implementation Pattern**:
```go
// Pseudo-code
func ScanAssets(root string, config *Config) (*ScanResult, error) {
    assets := []AssetFile{}

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if shouldSkipDir(path, config.Ignore) {
            return filepath.SkipDir
        }

        if isAssetFile(path, config.Extensions) {
            assets = append(assets, NewAssetFile(path))
        }

        return nil
    })

    return &ScanResult{Assets: assets}, err
}
```

---

### AST Parsing for Multi-Language Support

#### JavaScript/TypeScript

**Decision**: Use `github.com/evanw/esbuild` parser (extracted from esbuild bundler)

**Rationale**:
- Mature, used in production by esbuild (one of fastest JS bundlers)
- Handles both JavaScript and TypeScript
- Fast parsing (written in Go)
- Supports ESM, CommonJS, JSX/TSX

**Alternative Considered**:
- **go-fAST**: Newer, less battle-tested
- **typescript-ast-go**: TypeScript-only, less comprehensive
- **otto/parser**: JavaScript interpreter parser, slower for static analysis

**Implementation Notes**:
- Parse import statements: `import logo from './assets/logo.png'`
- Extract string literals containing asset paths
- Detect template literals: `` `${ASSETS_BASE}/icon.svg` ``
- Track import chains for asset constant files

#### Other Languages

**Strategy**: Start with string/regex scanning, add proper parsers as patterns emerge (Rule of Three)

**Phase 1 (MVP)**: Regex-based scanning for common patterns
```go
// Detect common asset reference patterns
patterns := []string{
    `"[./]*assets/[^"]+\.(png|jpg|svg|woff|mp4)"`,  // String literals
    `'[./]*public/[^']+\.(png|jpg|svg|woff|mp4)'`,  // Single quotes
    `src="[^"]+"`,                                  // HTML/JSX src attributes
    `url\(['"]?([^'")]+)['"]?\)`,                  // CSS url()
}
```

**Phase 2 (If needed)**: Add specific language parsers
- **Go**: Use `go/parser` and `go/ast` (standard library)
- **Dart**: String scanning + pubspec.yaml asset declarations
- **Swift**: String scanning + Assets.xcassets discovery
- **Kotlin/Java**: String scanning + R.drawable references

---

### Project Type Detection

**Decision**: Heuristic-based detection using project files

**Detection Logic**:
```go
type ProjectType int

const (
    Unknown ProjectType = iota
    WebReact
    WebVue
    WebAngular
    ReactNative
    Flutter
    iOS
    Android
)

func DetectProjectType(root string) ProjectType {
    // Check package.json
    if pkg := readPackageJSON(root); pkg != nil {
        if pkg.HasDependency("react") {
            if pkg.HasDependency("react-native") {
                return ReactNative
            }
            return WebReact
        }
        if pkg.HasDependency("vue") {
            return WebVue
        }
        // ... etc
    }

    // Check pubspec.yaml
    if fileExists(root, "pubspec.yaml") {
        return Flutter
    }

    // Check .xcodeproj
    if hasXcodeProject(root) {
        return iOS
    }

    // Check gradle files
    if fileExists(root, "build.gradle") {
        return Android
    }

    return Unknown
}
```

**Asset Path Defaults by Project Type**:
```go
var assetPaths = map[ProjectType][]string{
    WebReact:     {"public/", "src/assets/", "static/"},
    WebVue:       {"public/", "src/assets/"},
    ReactNative:  {"assets/", "src/assets/"},
    Flutter:      {"assets/", "lib/assets/"},
    iOS:          {"Assets.xcassets/", "Resources/"},
    Android:      {"res/drawable/", "res/raw/", "assets/"},
}
```

---

### Web UI Server

**Decision**: Use `net/http` standard library + simple static file server

**Rationale**:
- No need for full web framework (overkill for review UI)
- Standard library is sufficient for serving static HTML/JSON
- Embed HTML/CSS/JS using `//go:embed` directive
- Lightweight, zero external dependencies for HTTP

**Architecture**:
```
/cmd/asset-cleaner/
  └── web/
      ├── index.html       (embedded)
      ├── styles.css       (embedded)
      └── app.js           (embedded)

// Go code
//go:embed web/*
var webFiles embed.FS

func StartWebServer(result *ScanResult) error {
    http.Handle("/", http.FileServer(http.FS(webFiles)))
    http.HandleFunc("/api/results", serveResults(result))
    http.HandleFunc("/api/delete", handleDelete(result))

    fmt.Println("Review at http://localhost:3000")
    return http.ListenAndServe(":3000", nil)
}
```

**UI Technology**: Vanilla JavaScript (no framework)
- **Rationale**: Small UI scope doesn't justify React/Vue build pipeline
- Use native Web Components for reusability
- Minimal JavaScript for filtering/sorting
- TailwindCSS CDN for styling (optional, or inline CSS)

---

### Dynamic Path Analysis

**Challenge**: Detect constructs like `const path = BASE + '/icon.svg'`

**Strategy**: Multi-pass analysis with symbol tracking

**Phase 1**: Collect base path constants
```go
// Scan for declarations like:
// const ASSETS_BASE = '/public/images'
// const ICON_PATH = ASSETS_BASE + '/icons'

type ConstantDeclaration struct {
    Name  string
    Value string    // Literal value
    Deps  []string  // Dependencies (other constants used)
}

func ExtractConstants(file *ast.File) []ConstantDeclaration {
    // Parse const declarations
    // Track assignments
    // Resolve dependencies
}
```

**Phase 2**: Find constant usage + string concatenation
```go
// Look for patterns:
// - Variable + string: ASSETS_BASE + '/icon.svg'
// - Template literal: `${ASSETS_BASE}/icon.svg`
// - Function call: getAssetPath('icon.svg')

func FindDynamicReferences(file *ast.File, constants []ConstantDeclaration) []PossibleReference {
    // Track variable usage in string contexts
    // Resolve template literals
    // Flag unresolved patterns as "needs manual review"
}
```

**Fallback**: If resolution is ambiguous, flag asset as "needs manual review" (exclude from auto-deletion)

---

### Git Integration

**Decision**: Use `os.Remove` for filesystem deletion, rely on git CLI for version control awareness

**Rationale**:
- Don't need full git library (no go-git dependency)
- Users manage git operations themselves (commit deletions)
- Tool provides helpful instructions: `git add -u && git commit -m "Remove unused assets"`

**Safety Checks**:
```go
func IsGitRepository(root string) bool {
    _, err := os.Stat(filepath.Join(root, ".git"))
    return err == nil
}

func WarnIfNotVersionControlled(assets []string) {
    if !IsGitRepository(".") {
        fmt.Println("⚠️  Warning: Not in a git repository. Deletions are permanent!")
        fmt.Println("Consider backing up files before deletion.")
    }
}
```

---

## Performance Goals & Constraints

### Target Performance (from Success Criteria)

1. **Scan 1,000 files in <10 seconds** (SC-001)
2. **Handle up to 100,000 files** without crashing (SC-006)
3. **Progress updates every 2 seconds** during long scans (SC-010)

### Optimization Strategy

**For MVP**:
- Single-threaded `filepath.WalkDir` (sufficient for 1-10K files)
- In-memory storage of scan results (no database needed)

**For large projects (10K+ files)**:
- Worker pool pattern (8-16 workers)
- Stream results instead of buffering (reduce memory)
- Incremental progress updates via channel

**Memory Constraints**:
- Assume 100 bytes per asset file entry
- 100K files × 100 bytes = 10MB (well within limits)
- Use streaming for report generation if >1M files

---

## Testing Strategy

**Unit Tests**:
- File pattern matching (extensions, ignore patterns)
- Project type detection
- AST parsing for common reference patterns
- Constant resolution logic

**Integration Tests**:
- Scan fixture projects with known unused assets
- Verify accuracy: 95%+ detection rate, <5% false positives (SC-002, SC-003)
- Test edge cases: comments, dynamic paths, template literals

**Benchmark Tests**:
- Measure scan performance on projects of varying sizes (1K, 10K, 100K files)
- Validate <10s target for 1K files
- Profile memory usage and optimize hot paths

**Manual Testing**:
- Real-world projects: React, Vue, React Native, Flutter
- Cross-platform testing: macOS, Linux, Windows
- UX testing: Can users complete workflow in <5 minutes? (SC-005, SC-007)

---

## Architecture Decision Summary

| Decision Point | Choice | Key Rationale |
|----------------|--------|---------------|
| **Language** | Go 1.23+ | Single binary, fast file I/O, strong concurrency |
| **CLI Framework** | Cobra v1.9.1 | Industry standard, 181K+ dependents |
| **TUI Framework** | Bubble Tea v1.3.10 | Stable Elm architecture, beautiful TUIs |
| **Styling** | Lip Gloss v2.x | Declarative CSS-like API, table support |
| **Configuration** | Viper v1.21.0 | Multi-format support, 99K+ dependents |
| **File Scanning** | filepath.WalkDir + concurrent wrapper | Fast, standard library, scalable pattern |
| **JS/TS Parsing** | esbuild parser | Battle-tested, fast, handles both JS/TS |
| **Web UI** | net/http + embedded files | Zero dependencies, simple static server |
| **Git Integration** | Filesystem ops + git CLI | Simple, user-controlled, no heavy deps |

---

## Risk Mitigation

### Risk: False Positives (marking used assets as unused)

**Mitigation**:
- Conservative detection (prefer false negatives over false positives)
- "Needs manual review" category for ambiguous cases
- Dry-run mode by default
- Preview before deletion with clear warnings

### Risk: Poor Performance on Large Codebases

**Mitigation**:
- Concurrent scanning for projects >10K files
- Early directory skipping (node_modules, dist)
- Streaming results to avoid memory bloat
- Progress indicators to maintain user confidence

### Risk: Incomplete Language Support

**Mitigation**:
- Start with JavaScript/TypeScript (widest use case)
- Regex fallback for other languages (covers 80% of patterns)
- Configuration option to add custom patterns
- Document limitations clearly

### Risk: Cross-Platform Compatibility Issues

**Mitigation**:
- Use `filepath` package (abstracts OS-specific paths)
- Test on macOS, Linux, Windows
- Avoid shell commands (rely on Go standard library)

---

## Next Steps (Phase 1: Design)

1. Define data model (AssetFile, Reference, ScanResult, etc.)
2. Design CLI command structure and flags
3. Create quickstart guide with example workflows
4. Design configuration file schema
5. Sketch out concurrent scanning architecture
