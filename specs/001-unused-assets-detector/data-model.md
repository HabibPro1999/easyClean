# Data Model: Unused Assets Detector

**Feature**: Unused Assets Detector
**Date**: 2025-10-23

## Overview

This document defines the core data structures for the Unused Assets Detector. The model follows the principle of **Data Structures Over Algorithms** - by getting these structures right, the implementation logic becomes straightforward.

---

## Core Entities

### 1. AssetFile

Represents a single asset file discovered in the project.

**Attributes**:
```go
type AssetFile struct {
    // Identity
    Path         string    // Absolute path to asset file
    RelativePath string    // Path relative to project root
    Name         string    // File basename (e.g., "logo.png")
    Extension    string    // File extension (e.g., ".png")

    // Metadata
    Size         int64     // File size in bytes
    ModTime      time.Time // Last modification time
    Hash         string    // SHA-256 hash for duplicate detection (optional)

    // Classification
    Category     AssetCategory  // Image, Font, Video, Audio, Other
    Status       AssetStatus    // Used, Unused, PotentiallyUnused, NeedsManualReview

    // Usage Information
    References   []Reference    // Where this asset is referenced
    RefCount     int            // Cached count of references
}

type AssetCategory int

const (
    CategoryImage AssetCategory = iota  // .jpg, .png, .svg, .gif, .webp
    CategoryFont                        // .ttf, .woff, .woff2, .eot
    CategoryVideo                       // .mp4, .webm, .mov
    CategoryAudio                       // .mp3, .wav, .ogg
    CategoryOther                       // Unknown or miscellaneous
)

type AssetStatus int

const (
    StatusUsed AssetStatus = iota             // Has active code references
    StatusUnused                              // No references found
    StatusPotentiallyUnused                   // Only in comments/dead code
    StatusNeedsManualReview                   // Dynamic paths, ambiguous
)
```

**Validation Rules**:
- `Path` must be absolute and exist on filesystem
- `Extension` must be in configured extensions list
- `Size` must be >= 0
- `Status` determines inclusion in deletion candidates

**State Transitions**:
```
Initial -> StatusUsed (when reference found)
Initial -> StatusUnused (after scan with no references)
StatusUsed -> StatusPotentiallyUnused (if only comment refs found)
StatusUsed -> StatusNeedsManualReview (if dynamic path detected)
```

**Relationships**:
- Has many `Reference` objects (one-to-many)
- Belongs to one `ScanResult` (many-to-one)

---

### 2. Reference

Represents a single location in code where an asset is referenced.

**Attributes**:
```go
type Reference struct {
    // Location
    SourceFile   string    // File containing the reference
    LineNumber   int       // Line number in source file
    Column       int       // Column position (optional, for IDE integration)

    // Content
    MatchedText  string    // The actual matched string (e.g., "./assets/logo.png")
    Context      string    // Surrounding code (±2 lines for preview)

    // Classification
    Type         ReferenceType  // Import, StringLiteral, CSSUrl, etc.
    Confidence   float32        // 0.0-1.0 confidence score

    // Flags
    IsComment    bool       // True if reference is in comment
    IsDynamic    bool       // True if constructed at runtime
    IsDeadCode   bool       // True if in unreachable code branch (optional)
}

type ReferenceType int

const (
    RefTypeImport ReferenceType = iota        // import/require statement
    RefTypeStringLiteral                      // Plain string literal
    RefTypeTemplateLiteral                    // Template string with interpolation
    RefTypeCSSUrl                             // CSS url() function
    RefTypeHTMLAttribute                      // src/href attribute
    RefTypeConstant                           // Asset constant file
    RefTypeFunctionCall                       // getAsset('path') pattern
)
```

**Validation Rules**:
- `SourceFile` must exist at time of scan
- `LineNumber` must be > 0
- `Confidence` range: 0.0-1.0
- `Type` determines parsing strategy

**Usage**:
- Used to determine `AssetFile.Status`
- Displayed in detailed reports to explain why asset is marked used
- `IsComment` flag moves asset to `StatusPotentiallyUnused`

---

### 3. ScanResult

The complete output of scanning a project.

**Attributes**:
```go
type ScanResult struct {
    // Metadata
    Timestamp    time.Time         // When scan was performed
    ProjectRoot  string            // Root directory scanned
    ProjectType  ProjectType       // Detected project type
    Duration     time.Duration     // Time taken to scan

    // Assets
    Assets       []AssetFile       // All discovered assets
    UsedAssets   []AssetFile       // Filtered: StatusUsed
    UnusedAssets []AssetFile       // Filtered: StatusUnused
    PotentiallyUnusedAssets []AssetFile  // Filtered: StatusPotentiallyUnused
    NeedsReviewAssets []AssetFile   // Filtered: StatusNeedsManualReview

    // Statistics
    Stats        ScanStatistics    // Computed statistics

    // Configuration
    Config       *ProjectConfig    // Config used for this scan
}

type ScanStatistics struct {
    TotalAssets         int       // Total asset files found
    TotalSize           int64     // Total bytes of all assets
    UnusedCount         int       // Count of unused assets
    UnusedSize          int64     // Bytes of unused assets (savings)
    PotentiallyUnusedCount int    // Count of potentially unused
    NeedsReviewCount    int       // Count needing manual review
    FilesScanned        int       // Source files analyzed
    ReferencesFound     int       // Total references detected

    // Performance
    AvgScanSpeed        float64   // Files per second
}
```

**Methods**:
```go
// ComputeStatistics calculates all stats from Assets slice
func (sr *ScanResult) ComputeStatistics()

// FilterByStatus returns assets matching given status
func (sr *ScanResult) FilterByStatus(status AssetStatus) []AssetFile

// ToJSON exports result as JSON for web UI / reporting
func (sr *ScanResult) ToJSON() ([]byte, error)

// ToCSV exports result as CSV
func (sr *ScanResult) ToCSV() ([]byte, error)
```

**Relationships**:
- Has many `AssetFile` objects
- Has one `ProjectConfig`

---

### 4. ProjectConfig

User-defined configuration for scanning behavior.

**Attributes**:
```go
type ProjectConfig struct {
    // Asset Discovery
    AssetPaths     []string           // Directories to scan for assets
    Extensions     []string           // File extensions to consider (.png, .jpg, etc.)
    ExcludePaths   []string           // Paths to ignore (glob patterns)

    // Reference Detection
    ConstantFiles  []string           // Asset constant files (assets.ts, etc.)
    BasePathVars   []string           // Variable names for base paths
    CustomPatterns []string           // User-defined regex patterns

    // Behavior
    FollowSymlinks bool               // Follow symbolic links
    AutoDetectProjectType bool        // Auto-detect vs manual type
    ProjectType    ProjectType        // Manual project type override

    // Performance
    MaxWorkers     int                // Concurrent workers (0 = auto)
    MemoryLimit    int64              // Max memory usage in bytes (0 = no limit)

    // Output
    Verbose        bool               // Verbose logging
    ShowProgress   bool               // Display progress bar
    ColorOutput    bool               // Enable colored output
}
```

**Default Values**:
```go
func DefaultConfig() *ProjectConfig {
    return &ProjectConfig{
        AssetPaths:    []string{"assets/", "public/", "static/", "src/assets/"},
        Extensions:    []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp",
                                ".ttf", ".woff", ".woff2", ".eot",
                                ".mp4", ".webm", ".mov",
                                ".mp3", ".wav", ".ogg"},
        ExcludePaths:  []string{"node_modules/", "dist/", "build/", ".next/",
                                "target/", "vendor/", ".git/"},
        FollowSymlinks: false,
        AutoDetectProjectType: true,
        MaxWorkers:    0,  // Auto-detect CPU cores
        ShowProgress:  true,
        ColorOutput:   true,
    }
}
```

**Validation Rules**:
- At least one `AssetPath` must be specified
- At least one `Extension` must be specified
- `ExcludePaths` patterns must be valid glob patterns
- `MaxWorkers` must be >= 0 (0 means auto-detect)

**File Format** (.unusedassets.yaml):
```yaml
asset_paths:
  - public/
  - src/assets/

extensions:
  - .png
  - .jpg
  - .svg
  - .woff
  - .ttf

exclude_paths:
  - node_modules/
  - dist/
  - "**/__tests__/**"

constant_files:
  - src/constants/assets.ts
  - lib/assets.dart

base_path_vars:
  - ASSETS_BASE
  - PUBLIC_URL

max_workers: 8
show_progress: true
```

---

### 5. ProjectType

Enum for detected project types.

**Definition**:
```go
type ProjectType int

const (
    ProjectTypeUnknown ProjectType = iota
    ProjectTypeWebReact
    ProjectTypeWebVue
    ProjectTypeWebAngular
    ProjectTypeWebSvelte
    ProjectTypeReactNative
    ProjectTypeFlutter
    ProjectTypeIOS
    ProjectTypeAndroid
    ProjectTypeGo
    ProjectTypeRust
)

func (pt ProjectType) String() string {
    return [...]string{
        "Unknown",
        "React (Web)",
        "Vue (Web)",
        "Angular (Web)",
        "Svelte (Web)",
        "React Native",
        "Flutter",
        "iOS (Swift)",
        "Android (Kotlin/Java)",
        "Go",
        "Rust",
    }[pt]
}

// DefaultAssetPaths returns common asset locations for this project type
func (pt ProjectType) DefaultAssetPaths() []string

// DefaultExtensions returns relevant asset types for this project type
func (pt ProjectType) DefaultExtensions() []string
```

**Detection Logic**:
- Check for `package.json` → analyze dependencies
- Check for `pubspec.yaml` → Flutter
- Check for `.xcodeproj` → iOS
- Check for `build.gradle` → Android
- Check for `go.mod` → Go
- Check for `Cargo.toml` → Rust

---

### 6. IgnoreRule

Pattern matching rule for excluding files/directories.

**Attributes**:
```go
type IgnoreRule struct {
    Pattern    string            // Glob or regex pattern
    Type       IgnoreRuleType    // Glob or Regex
    Reason     string            // Optional: why this is ignored
    Enabled    bool              // Can be disabled without removing
}

type IgnoreRuleType int

const (
    IgnoreTypeGlob IgnoreRuleType = iota  // Glob pattern (default)
    IgnoreTypeRegex                       // Regular expression
)

// Matches checks if given path matches this rule
func (ir *IgnoreRule) Matches(path string) bool
```

**Examples**:
```go
rules := []IgnoreRule{
    {Pattern: "node_modules/", Type: IgnoreTypeGlob, Enabled: true},
    {Pattern: "*.test.js", Type: IgnoreTypeGlob, Enabled: true},
    {Pattern: `.*\.backup\..*`, Type: IgnoreTypeRegex, Enabled: false},
}
```

---

## Data Flow

### Scan Workflow

```
1. Load ProjectConfig from file or defaults
   ↓
2. Detect ProjectType (auto or manual)
   ↓
3. Walk filesystem, collect AssetFile objects
   ↓ (filter by extensions, apply ignore rules)
4. For each source file:
   - Parse AST (if applicable)
   - Extract string literals
   - Create Reference objects
   ↓
5. Match References to AssetFiles
   ↓
6. Classify AssetFile.Status based on references
   ↓
7. Build ScanResult with statistics
   ↓
8. Export to CLI, JSON, CSV, or Web UI
```

### Status Classification Logic

```go
func ClassifyAsset(asset *AssetFile) AssetStatus {
    if len(asset.References) == 0 {
        return StatusUnused
    }

    hasActiveRef := false
    allInComments := true

    for _, ref := range asset.References {
        if ref.IsDynamic {
            return StatusNeedsManualReview  // Conservative approach
        }

        if !ref.IsComment {
            allInComments = false
            hasActiveRef = true
        }
    }

    if allInComments {
        return StatusPotentiallyUnused
    }

    if hasActiveRef {
        return StatusUsed
    }

    return StatusUnused
}
```

---

## Data Persistence

### In-Memory (Default)

All data structures are held in memory during scan and review. No database needed.

**Memory Estimation**:
- AssetFile: ~500 bytes (with references)
- 100K assets × 500 bytes = 50MB (acceptable)

### Export Formats

**JSON** (for programmatic access):
```json
{
  "timestamp": "2025-10-23T14:30:00Z",
  "project_root": "/path/to/project",
  "project_type": "React (Web)",
  "duration_ms": 2340,
  "stats": {
    "total_assets": 247,
    "unused_count": 58,
    "unused_size": 13002752
  },
  "unused_assets": [
    {
      "path": "/path/to/project/public/images/old-banner.jpg",
      "relative_path": "public/images/old-banner.jpg",
      "size": 890240,
      "category": "Image",
      "status": "Unused"
    }
  ]
}
```

**CSV** (for spreadsheets):
```csv
Path,Size,Category,Status,References
public/images/old-banner.jpg,890240,Image,Unused,0
public/fonts/unused.woff,45000,Font,Unused,0
src/assets/logo-v1.png,120000,Image,PotentiallyUnused,1
```

---

## Data Integrity

### Validation

All data structures implement validation:

```go
type Validator interface {
    Validate() error
}

func (af *AssetFile) Validate() error {
    if af.Path == "" {
        return errors.New("AssetFile.Path cannot be empty")
    }
    if af.Size < 0 {
        return errors.New("AssetFile.Size cannot be negative")
    }
    // ... more validations
    return nil
}
```

### Immutability

Once created, core structures are immutable (except status updates):

```go
// Good: Create new slice for filtering
func (sr *ScanResult) GetUnused() []AssetFile {
    return sr.FilterByStatus(StatusUnused)  // Returns copy
}

// Bad: Don't mutate original Assets slice
// sr.Assets = append(sr.Assets, newAsset)  // ❌
```

---

## Performance Considerations

### Indexing

For large projects, use maps for O(1) lookups:

```go
type AssetIndex struct {
    ByPath map[string]*AssetFile     // Path -> AssetFile
    ByName map[string][]*AssetFile   // Name -> []AssetFile (duplicates)
}

func BuildIndex(assets []AssetFile) *AssetIndex {
    idx := &AssetIndex{
        ByPath: make(map[string]*AssetFile, len(assets)),
        ByName: make(map[string][]*AssetFile),
    }

    for i := range assets {
        asset := &assets[i]
        idx.ByPath[asset.Path] = asset
        idx.ByName[asset.Name] = append(idx.ByName[asset.Name], asset)
    }

    return idx
}
```

### Memory Optimization

For very large projects (>100K assets):

1. **Stream processing**: Don't load all assets at once
2. **Reference pooling**: Reuse Reference objects
3. **String interning**: Deduplicate common strings (extensions, categories)

---

## Summary

This data model provides:

✅ **Clear entity boundaries** - Each struct has a single responsibility
✅ **Explicit relationships** - One-to-many between ScanResult → AssetFile → Reference
✅ **Type safety** - Enums for categories, statuses, reference types
✅ **Extensibility** - Easy to add new ProjectTypes, AssetCategories, ReferenceTypes
✅ **Testability** - Pure data structures, easy to construct fixtures
✅ **Performance** - Optimized for memory usage and lookup speed

The model aligns with the **Data Structures Over Algorithms** principle: by modeling the problem domain accurately, the implementation code becomes straightforward and self-documenting.
