# Implementation Plan: Unused Assets Detector

**Branch**: `001-unused-assets-detector` | **Date**: 2025-10-23 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-unused-assets-detector/spec.md`

## Summary

Build a CLI tool that automatically detects and safely removes unused asset files (images, fonts, videos, etc.) from codebases. The tool uses smart scanning logic with multi-pattern reference detection, supports multiple project types (React, Vue, Flutter, iOS, Android), and provides both a beautiful CLI interface and optional web UI for reviewing results before deletion.

**Technical Approach**: Go-based CLI with Bubble Tea TUI framework, concurrent file scanning, AST parsing for JavaScript/TypeScript, and embedded web server for interactive review. The tool prioritizes safety (git history preservation, three-tier classification, dry-run mode) and performance (handles 100K+ files, <10s for 1K files).

---

## Technical Context

**Language/Version**: Go 1.23+
**Primary Dependencies**:
- Cobra v1.9.1+ (CLI framework, command structure)
- Bubble Tea v1.3.10 (TUI framework for progress indicators)
- Lip Gloss v2.x (Terminal styling and layouts)
- Viper v1.21.0 (Configuration management)
- esbuild parser (JavaScript/TypeScript AST parsing)

**Storage**: In-memory scan results (no database needed, ~50MB for 100K assets)
**Testing**: Go standard testing package + table-driven tests
**Target Platform**: Cross-platform (macOS, Linux, Windows) - single binary distribution
**Project Type**: CLI application with optional embedded web UI

**Performance Goals**:
- Scan 1,000 files in under 10 seconds
- Handle projects up to 100,000 files without crashing
- Progress updates every 2 seconds during scans
- File I/O optimization via filepath.WalkDir + concurrent scanning

**Constraints**:
- Offline-capable (no network required)
- Single binary executable (no runtime dependencies)
- < 100MB memory usage for typical projects (1K-10K files)
- Concurrent scanning limited to available CPU cores

**Scale/Scope**:
- Support 5+ project types (React, Vue, Flutter, iOS, Android)
- Handle 15+ asset file types (images, fonts, videos, audio)
- Detect 7+ reference patterns (imports, literals, CSS urls, etc.)
- Process projects from 100 to 100,000 files

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Simplicity First (KISS) ✅

- **CLI structure**: Simple command hierarchy (`scan`, `review`, `delete`, `info`, `init`)
- **No unnecessary abstractions**: Use standard library for file I/O, avoid framework overload
- **Data model**: Straightforward structs (AssetFile, Reference, ScanResult) with clear relationships
- **Configuration**: Single YAML file, no complex nested structures

**Compliance**: PASS - Solution uses proven libraries without over-engineering

---

### Principle II: Avoid Repetition (DRY with Rule of Three) ✅

- **Pattern detection**: Start with regex scanning, add language-specific parsers only after 3 use cases
- **Project type defaults**: Map-based configuration reduces duplication across project types
- **Code generation**: No premature code generation - start with explicit implementations

**Compliance**: PASS - Abstraction justified by proven patterns (project type detection, reference matching)

---

### Principle III: Build Only What's Needed (YAGNI) ✅

- **MVP focus**: User Story 1 (scanning) first, deletion and web UI second
- **Language support**: JavaScript/TypeScript only initially, other languages via regex fallback
- **No speculative features**: No plugin system, no GUI, no cloud sync (unless explicitly requested)

**Compliance**: PASS - Feature scope limited to spec requirements, no gold-plating

---

### Principle IV: Progressive Implementation (Make It Work → Right → Fast) ✅

**Phase 1**: Single-threaded scan, simple reference detection (Make It Work)
**Phase 2**: Refactor for clarity, add tests, improve error messages (Make It Right)
**Phase 3**: Add concurrent scanning for large projects after profiling (Make It Fast)

**Compliance**: PASS - Clear progression from working prototype to optimized solution

---

### Principle V: Explicit Over Magic ✅

- **Dependencies**: Configuration, file system access, and scanners passed explicitly to functions
- **No hidden state**: Functions are pure where possible, side effects clearly documented
- **Framework magic**: Viper/Cobra magic justified by widespread adoption (181K+ and 99K+ users)

**Compliance**: PASS - Limited framework magic, dependencies explicit in function signatures

---

### Principle VI: Data Structures Over Algorithms ✅

- **Core model**: AssetFile, Reference, ScanResult structs with clear relationships
- **Classification**: AssetStatus enum drives logic (Used, Unused, PotentiallyUnused, NeedsManualReview)
- **Indexing**: Hash maps for O(1) asset lookup by path

**Compliance**: PASS - Data model designed first, algorithms follow naturally

---

### Principle VII: Testability and Observability ✅

- **Testable**: Functions accept file system interface (mockable for testing)
- **Observable**: Structured logging at key decision points (project detection, reference matching)
- **No shared state**: ScanResult is immutable after creation, state transitions explicit

**Compliance**: PASS - Clean interfaces enable unit testing, logging captures key events

---

### Code Quality: Cognitive Load Limits ✅

- **Functions**: Each function does one thing (scan assets, match references, classify status)
- **Indirection**: Maximum 3 levels (main → command → service → scanner)
- **File navigation**: Understanding scan workflow requires reading 1-2 files, not 10+

**Compliance**: PASS - Simple, readable code structure

---

### Code Quality: Code Smells ⚠️

**Potential Violation**: Multiple AST parsers (JavaScript, TypeScript, Dart, Swift)

**Justification**: Each language requires different parsing logic. Starting with JavaScript/TypeScript only (YAGNI), adding others based on actual demand.

**Simpler Alternative Rejected**: Single regex-based scanner would miss complex patterns like template literals and constant imports, leading to false positives.

---

## Project Structure

### Documentation (this feature)

```text
specs/001-unused-assets-detector/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Technology research and decisions
├── data-model.md        # Core data structures
├── quickstart.md        # User getting-started guide
└── contracts/
    └── cli-interface.md # CLI command contracts
```

### Source Code (repository root)

```text
cmd/
└── asset-cleaner/
    ├── main.go                 # Entry point
    ├── commands/
    │   ├── root.go             # Root command setup
    │   ├── scan.go             # Scan command
    │   ├── review.go           # Review (web UI) command
    │   ├── delete.go           # Delete command
    │   ├── info.go             # Info command
    │   ├── init.go             # Init command
    │   └── ignore.go           # Ignore command
    └── web/
        ├── index.html          # Embedded web UI
        ├── styles.css          # UI styles
        └── app.js              # Client-side logic

internal/
├── scanner/
│   ├── scanner.go              # Main scanner orchestration
│   ├── asset_finder.go         # File system asset discovery
│   ├── reference_finder.go     # Code reference detection
│   └── concurrent.go           # Concurrent scanning logic
│
├── parser/
│   ├── javascript.go           # JS/TS AST parsing (esbuild)
│   ├── regex.go                # Regex-based fallback parser
│   └── patterns.go             # Common reference patterns
│
├── detector/
│   ├── project_type.go         # Project type detection
│   └── asset_paths.go          # Default paths by project type
│
├── classifier/
│   └── classifier.go           # Asset status classification logic
│
├── config/
│   ├── config.go               # Configuration structs
│   ├── loader.go               # Load from .unusedassets.yaml
│   └── defaults.go             # Default configurations
│
├── models/
│   ├── asset.go                # AssetFile struct
│   ├── reference.go            # Reference struct
│   ├── scan_result.go          # ScanResult struct
│   └── project_config.go       # ProjectConfig struct
│
├── ui/
│   ├── progress.go             # Bubble Tea progress indicators
│   ├── styles.go               # Lip Gloss styling
│   └── server.go               # Web UI HTTP server
│
└── utils/
    ├── fileutil.go             # File system utilities
    ├── pathutil.go             # Path manipulation
    └── hashutil.go             # File hashing (deduplication)

tests/
├── fixtures/                   # Test project fixtures
│   ├── react-project/
│   ├── flutter-project/
│   └── vue-project/
├── integration/
│   ├── scan_test.go            # End-to-end scan tests
│   ├── delete_test.go          # Deletion tests
│   └── accuracy_test.go        # False positive/negative tests
└── unit/
    ├── scanner_test.go         # Scanner unit tests
    ├── parser_test.go          # Parser unit tests
    └── classifier_test.go      # Classification logic tests

go.mod
go.sum
.unusedassets.yaml              # Example config file
README.md
LICENSE
```

**Structure Decision**: Single Go project structure with clear separation of concerns:
- `cmd/` for executable entry points
- `internal/` for application logic (not importable by external packages)
- `tests/` for test code and fixtures
- Embedded web UI files in `cmd/asset-cleaner/web/`

This structure follows Go best practices and keeps related code together while maintaining clear boundaries between components.

---

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Multiple language parsers (future) | Different languages have incompatible syntax requiring specialized parsers | Regex-only approach misses 30%+ references in testing (template literals, imports, constant files) |

**Note**: Starting with JavaScript/TypeScript parser only. Additional parsers added based on user demand (Rule of Three).

---

## Implementation Phases

### Phase 0: Research ✅ COMPLETE

See [research.md](./research.md) for:
- Technology stack selection (Go 1.23+, Cobra, Bubble Tea, etc.)
- Performance strategies (filepath.WalkDir, concurrent scanning)
- AST parsing approach (esbuild parser for JS/TS)
- Risk mitigation strategies

### Phase 1: Design ✅ COMPLETE

See artifacts:
- [data-model.md](./data-model.md) - Core data structures
- [contracts/cli-interface.md](./contracts/cli-interface.md) - CLI commands and flags
- [quickstart.md](./quickstart.md) - User guide and workflows

### Phase 2: Task Generation (Next)

Run `/speckit.tasks` to generate implementation tasks based on:
- User stories from spec.md (P1-P4 prioritization)
- Technical design from plan.md and data-model.md
- CLI contracts from contracts/cli-interface.md

**Expected task breakdown**:
1. **Setup**: Project scaffolding, dependency installation, CI/CD
2. **Core Scanning (US1)**: Asset discovery, reference detection, classification
3. **Deletion (US2)**: Preview mode, confirmation, filesystem operations
4. **Configuration (US3)**: Ignore patterns, config file loading
5. **Reporting (US4)**: JSON/CSV export, web UI
6. **Polish**: Error handling, logging, documentation

---

## Key Technical Decisions

### 1. Go Language Choice

**Decision**: Use Go 1.23+ as the implementation language

**Rationale**:
- Single binary distribution (no runtime dependencies)
- Excellent concurrency support (goroutines for parallel scanning)
- Fast file I/O operations (critical for large codebases)
- Growing ecosystem of beautiful CLI libraries (Charm.sh)
- Cross-platform compilation built-in

**Trade-offs**:
- ❌ Less familiar than Node.js/Python for some developers
- ✅ Much faster performance for file operations
- ✅ Zero runtime dependencies (distribute single executable)

---

### 2. Bubble Tea for TUI

**Decision**: Use Bubble Tea v1.3.10 for terminal user interface

**Rationale**:
- Elm Architecture makes complex state predictable
- Beautiful progress indicators for long-running scans
- Battle-tested (9K+ dependents, used in production)
- Composable components (reuse progress bar, spinner, etc.)

**Trade-offs**:
- ❌ Learning curve for Elm Architecture pattern
- ✅ Much better UX than plain text output
- ✅ Easy to add interactive features later

---

### 3. Concurrent Scanning Strategy

**Decision**: Single-threaded for <10K files, concurrent for larger projects

**Rationale**:
- filepath.WalkDir is fast enough for small projects
- Concurrent scanning adds complexity (synchronization, memory)
- Auto-detect project size and switch strategies

**Implementation**:
```go
if estimatedFiles < 10000 {
    return singleThreadedScan(root, config)
} else {
    return concurrentScan(root, config, runtime.NumCPU())
}
```

---

### 4. Safety-First Deletion

**Decision**: Three-tier classification + git history preservation

**Rationale**:
- False positives (deleting used assets) are catastrophic
- Three tiers: Used, Unused, PotentiallyUnused, NeedsManualReview
- Only "Unused" included in auto-deletion
- Git preserves history for recovery

**Implementation**:
- StatusUnused: Zero references found → safe to delete
- StatusPotentiallyUnused: Only in comments → warn user
- StatusNeedsManualReview: Dynamic paths detected → exclude from auto-delete
- StatusUsed: Active references found → keep

---

### 5. Embedded Web UI

**Decision**: Embed HTML/CSS/JS using `//go:embed`, serve via net/http

**Rationale**:
- Zero external dependencies for web UI
- Single binary includes everything
- Standard library HTTP server is sufficient
- No need for React/Vue build pipeline (overkill for simple UI)

**Implementation**:
```go
//go:embed web/*
var webFiles embed.FS

func StartWebServer(result *ScanResult, port int) error {
    http.Handle("/", http.FileServer(http.FS(webFiles)))
    http.HandleFunc("/api/results", serveJSON(result))
    return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
```

---

## Risk Assessment

### High Risk: False Positives (SC-003: <5% target)

**Mitigation**:
1. Conservative detection (prefer false negatives)
2. Three-tier classification system
3. Dry-run mode by default
4. Comprehensive integration tests with fixture projects
5. Clear warnings for ambiguous cases

**Validation**: Test against 10+ real-world projects, measure false positive rate

---

### Medium Risk: Performance on Large Projects (SC-006: 100K files)

**Mitigation**:
1. Concurrent scanning with worker pools
2. Early directory skipping (node_modules, dist)
3. Streaming results to avoid memory bloat
4. Memory profiling and optimization

**Validation**: Benchmark on synthetic projects with 10K, 50K, 100K files

---

### Medium Risk: Incomplete Language Support

**Mitigation**:
1. Start with JavaScript/TypeScript (80% of use cases)
2. Regex fallback for other languages (covers common patterns)
3. Document limitations clearly
4. Make parser system extensible for future languages

**Validation**: Test on React, Vue, Flutter, iOS, Android projects

---

### Low Risk: Cross-Platform Compatibility

**Mitigation**:
1. Use `filepath` package (abstracts OS differences)
2. Avoid shell commands (use Go standard library)
3. Test on macOS, Linux, Windows in CI

**Validation**: GitHub Actions matrix testing on all platforms

---

## Success Metrics (from Spec)

| Metric | Target | How to Validate |
|--------|--------|-----------------|
| **SC-001** | Scan 1K files in <10s | Benchmark test with timer |
| **SC-002** | 95%+ accuracy | Integration tests with hand-verified fixtures |
| **SC-003** | <5% false positive rate | Measure against real projects |
| **SC-004** | Safe deletion | Manual testing + user feedback |
| **SC-005** | Complete workflow <5 min | User testing with timer |
| **SC-006** | Handle 100K files | Stress test with synthetic project |
| **SC-007** | 90%+ success without docs | Usability testing with 10 users |
| **SC-008** | Git recovery works | Manual test: delete, restore via git |
| **SC-009** | 10-30% size reduction | Measure on 5 real projects with asset bloat |
| **SC-010** | Progress updates every 2s | Timer validation in long scans |

---

## Next Steps

1. **Generate Tasks**: Run `/speckit.tasks` to create detailed implementation tasks
2. **Set up Project**: Initialize Go module, install dependencies
3. **Implement US1 (MVP)**: Scanning and reporting only
4. **Integration Testing**: Test against fixture projects (React, Vue, etc.)
5. **Implement US2-US4**: Deletion, configuration, reporting
6. **Polish**: Error handling, documentation, CI/CD
7. **Release**: Package binaries, publish to GitHub Releases

---

**Plan Complete** - Ready for task generation and implementation.
