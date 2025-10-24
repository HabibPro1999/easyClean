# Changelog

All notable changes to easyClean are documented in this file.

## [1.0.1] - 2025-10-24

### Added

**Multi-Project Review Support**
- Run multiple concurrent review servers on different ports with intelligent port allocation
- `--list` flag to display all active servers with project name, port, PID, and uptime
- `--kill <port>` flag to gracefully stop specific servers
- Server registry tracking in `~/.cache/easyClean/servers.json`
- Graceful shutdown with Ctrl+C signal handling (SIGINT/SIGTERM)
- Automatic cleanup of stale/dead servers from registry
- See [MULTI_PROJECT_REVIEW.md](MULTI_PROJECT_REVIEW.md) for full documentation

**Framework-Specific Asset Detection**
- **Pattern Provider System**: Framework-specific detection patterns with fallback to generic
- **React/Next.js Patterns**:
  - `React.lazy()` dynamic imports
  - Next.js public folder conventions
  - Webpack magic comments in dynamic imports
- **Angular Patterns**:
  - `@Component` decorator `templateUrl` and `styleUrls`
  - Lazy route loading with `loadChildren`
  - Template bindings `[src]="path"`
- **Vue/Nuxt Patterns**:
  - `defineAsyncComponent` dynamic imports
  - Template `:src` bindings and require() calls
  - Nuxt static/public folder conventions
- **Flutter Patterns**:
  - Enhanced `Image.asset()` and `AssetImage()` detection
  - `pubspec.yaml` asset declarations
  - Font family references
- **AST-Based JavaScript/TypeScript Parsing**:
  - Static and dynamic import analysis
  - JSX image reference extraction
  - Object property asset detection
  - Export statement analysis
  - Confidence scoring (AST: 1.0, patterns: 0.95, generic: 0.7)
- See [FRAMEWORK_DETECTION.md](FRAMEWORK_DETECTION.md) for technical details

### Fixed
- Global state issue preventing multiple concurrent review servers
- Removed global `currentScanResult` variable causing server conflicts
- Server process cleanup on shutdown to prevent orphaned entries

### Improved
- Asset detection accuracy: 85% → 95%+ (framework-aware patterns + AST parsing)
- Server architecture: instance-based design instead of global state
- Reference deduplication to avoid counting the same reference twice
- HTTP server configuration with proper timeouts (Read: 15s, Write: 15s, Idle: 60s)
- Framework-specific file extension support per project type

## [1.0.0] - 2025-10-23

### Initial Release

**Core Features**
- Automated detection of unused asset files in codebases
- Support for 10+ project types (React, Vue, Angular, Svelte, Flutter, iOS, Android, Go, Rust)
- Multi-pattern reference detection with confidence scoring
- Three-tier asset classification system:
  - **Used**: Has active code references (safe to keep)
  - **Unused**: No references found (safe to delete)
  - **Potentially Unused**: Only found in comments/dead code (review first)

**Scanning Capabilities**
- Automatic project type detection
- Asset discovery across configurable directories
- Multi-file reference detection:
  - Import/require statements
  - CSS url() references
  - HTML src/href attributes
  - String literals with asset paths
  - Template literals with variables
  - Flutter-specific patterns (Image.asset, AssetImage, rootBundle.load)
- Comment detection to separate code from documentation
- Dynamic reference detection (string concatenation, template interpolation)
- Configurable asset extensions and exclude patterns

**Safety Features**
- Three-tier classification for risk-aware deletion
- Dry-run mode by default (preview before deletion)
- Confirmation prompts before any file deletion
- Git history preservation (deletes from filesystem only)
- Automatic result caching in OS cache directory (~/.cache/easyClean/)

**Commands**
- `scan`: Detect unused assets with progress indicators
- `review`: Interactive web UI for browsing and previewing results
- `delete`: Safe removal with dry-run and confirmation options
- `init`: Generate configuration files with templates
- `info`: Display project and configuration information
- `ignore`: Manage ignore patterns

**Configuration**
- YAML-based `.unusedassets.yaml` configuration file
- Per-project configuration support
- Framework-specific default settings
- Global flags for verbosity, quiet mode, color control

**Web UI**
- Interactive browser-based review interface
- Asset filtering and sorting
- File preview and metadata display
- Batch deletion with confirmation
- Export capabilities

**Performance**
- Scans 1,000 files in <10 seconds
- Handles projects up to 100,000 files
- <100MB memory usage for typical projects
- Concurrent file system scanning for large projects
- Efficient caching of results

**Output Formats**
- Text: Human-readable terminal output
- JSON: Structured data for integration
- CSV: Spreadsheet-compatible format
- Web UI: Interactive browser interface

---

## Version Format

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR** version for breaking changes
- **MINOR** version for new features (backward compatible)
- **PATCH** version for bug fixes (backward compatible)
