# CLI Interface Contract: Unused Assets Detector

**Tool Name**: `asset-cleaner`
**Version**: 1.0.0

## Overview

This document defines the command-line interface contract for the Unused Assets Detector. All commands, flags, and output formats are specified here to ensure consistent behavior and enable contract testing.

---

## Global Flags

Available for all commands:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | `-c` | string | `.unusedassets.yaml` | Path to configuration file |
| `--verbose` | `-v` | bool | false | Enable verbose logging |
| `--quiet` | `-q` | bool | false | Suppress all output except errors |
| `--no-color` | | bool | false | Disable colored output |
| `--version` | | bool | false | Print version and exit |
| `--help` | `-h` | bool | false | Show help message |

---

## Commands

### 1. `asset-cleaner scan [directory]`

Scan a project directory for unused assets.

**Arguments**:
- `directory` (optional): Project root to scan. Default: current directory `.`

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--extensions` | []string | (auto) | Asset extensions to scan (e.g., `.png,.jpg`) |
| `--exclude` | []string | (auto) | Paths to exclude (glob patterns) |
| `--output` | string | - | Export results to file (JSON/CSV based on extension) |
| `--format` | string | `text` | Output format: `text`, `json`, `csv` |
| `--no-progress` | bool | false | Disable progress bar |

**Exit Codes**:
- `0`: Success (unused assets found or not found)
- `1`: Error (invalid config, directory not found, permission denied)
- `2`: Invalid arguments

**Output (text format)**:
```
╭─────────────────────────────────────────────╮
│  🔍 Asset Cleaner v1.0                      │
╰─────────────────────────────────────────────╯

⠋ Detecting project type...
✓ Found: React (Web)

📁 Scanning asset directories...

✓ public/images/ (123 files)
✓ public/fonts/ (8 files)
✓ src/assets/ (116 files)

⠙ Analyzing code references...

✓ Direct imports (234 found)
✓ String literals (567 found)
✓ Dynamic paths (89 found)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 Scan Complete

  Total Assets:           247

  ✓ Used Assets:          189
  ⚠️  Unused Assets:       58

  💾 Potential Savings:   12.4 MB

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✨ Run 'asset-cleaner review' to inspect unused assets
```

**Output (JSON format)**:
```json
{
  "version": "1.0.0",
  "timestamp": "2025-10-23T14:30:00Z",
  "project_root": "/path/to/project",
  "project_type": "React (Web)",
  "duration_ms": 2340,
  "statistics": {
    "total_assets": 247,
    "used_count": 189,
    "unused_count": 58,
    "potentially_unused_count": 0,
    "needs_review_count": 0,
    "total_size_bytes": 45000000,
    "unused_size_bytes": 13002752,
    "files_scanned": 1432
  },
  "used_assets": [...],
  "unused_assets": [
    {
      "path": "/path/to/project/public/images/old-banner.jpg",
      "relative_path": "public/images/old-banner.jpg",
      "name": "old-banner.jpg",
      "extension": ".jpg",
      "size_bytes": 890240,
      "mod_time": "2025-07-15T10:20:30Z",
      "category": "Image",
      "status": "Unused",
      "reference_count": 0
    }
  ],
  "potentially_unused_assets": [],
  "needs_review_assets": []
}
```

**Output (CSV format)**:
```csv
Status,Path,Size,Category,References,ModTime
Unused,public/images/old-banner.jpg,890240,Image,0,2025-07-15T10:20:30Z
Unused,public/fonts/unused.woff,45000,Font,0,2025-06-10T14:30:00Z
PotentiallyUnused,src/assets/logo-v1.png,120000,Image,1,2025-08-01T09:15:00Z
```

**Example Usage**:
```bash
# Scan current directory
asset-cleaner scan

# Scan specific directory
asset-cleaner scan ./my-project

# Export results to JSON
asset-cleaner scan --output results.json

# Custom extensions only
asset-cleaner scan --extensions .png,.jpg,.svg

# Scan with custom exclusions
asset-cleaner scan --exclude "node_modules/,dist/,**/*.test.js"
```

---

### 2. `asset-cleaner review`

Launch web UI to review unused assets interactively.

**Arguments**: None

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--port` | int | `3000` | HTTP server port |
| `--host` | string | `localhost` | HTTP server host |
| `--no-browser` | bool | false | Don't auto-open browser |
| `--scan-file` | string | - | Load scan results from JSON file (skip rescan) |

**Exit Codes**:
- `0`: Success (server started)
- `1`: Error (port in use, no scan results found)
- `2`: Invalid arguments

**Output**:
```
╭─────────────────────────────────────────────╮
│  🌐 Asset Cleaner Review UI                 │
╰─────────────────────────────────────────────╯

🔍 Using scan results from previous run

🌐 Server started at http://localhost:3000
🚀 Opening browser...

Press Ctrl+C to stop server
```

**Web UI Endpoints**:
- `GET /` - Main review interface (HTML)
- `GET /api/results` - Get scan results (JSON)
- `POST /api/delete` - Delete selected assets (JSON body: `{"paths": [...]}`)
- `POST /api/ignore` - Add paths to ignore list (JSON body: `{"patterns": [...]}`)

**Example Usage**:
```bash
# Launch review UI (uses last scan)
asset-cleaner review

# Load specific scan results
asset-cleaner review --scan-file results.json

# Use custom port
asset-cleaner review --port 8080

# Don't auto-open browser
asset-cleaner review --no-browser
```

---

### 3. `asset-cleaner delete [paths...]`

Delete unused assets from filesystem.

**Arguments**:
- `paths...` (optional): Specific paths to delete. If omitted, deletes all unused assets from last scan.

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dry-run` | bool | false | Show what would be deleted without deleting |
| `--interactive` | `-i` | bool | false | Prompt for confirmation before each file |
| `--force` | bool | false | Skip confirmation prompts |
| `--scan-file` | string | - | Load scan results from JSON file |

**Exit Codes**:
- `0`: Success (files deleted or dry-run completed)
- `1`: Error (file not found, permission denied, not in git repo)
- `2`: Invalid arguments
- `3`: User cancelled operation

**Output (default)**:
```
╭─────────────────────────────────────────────╮
│  🗑️  Delete Unused Assets                   │
╰─────────────────────────────────────────────╯

Found 58 unused assets (12.4 MB)

⚠️  You are about to delete 58 files. This action removes files from your filesystem.
Files will remain in git history and can be recovered via git commands.

Continue? [y/N]: y

Deleting files...
  ✓ public/images/old-banner.jpg (890 KB)
  ✓ public/fonts/unused.woff (45 KB)
  ...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Deleted 58 files (12.4 MB freed)

Next steps:
  git add -u
  git commit -m "Remove unused assets"
```

**Output (dry-run)**:
```
🧪 Dry Run Mode - No files will be deleted

Would delete:
  • public/images/old-banner.jpg (890 KB)
  • public/fonts/unused.woff (45 KB)
  ...

Total: 58 files (12.4 MB)

Run without --dry-run to actually delete files.
```

**Output (interactive)**:
```
Delete public/images/old-banner.jpg (890 KB)? [y/N/q]: y
  ✓ Deleted

Delete public/fonts/unused.woff (45 KB)? [y/N/q]: n
  ⊘ Skipped

Delete src/assets/old-logo.png (120 KB)? [y/N/q]: q
  ⊘ Cancelled (3 files deleted, 55 skipped)
```

**Example Usage**:
```bash
# Delete all unused assets (with confirmation)
asset-cleaner delete

# Dry run to preview
asset-cleaner delete --dry-run

# Delete specific files
asset-cleaner delete public/images/old.jpg src/assets/unused.png

# Delete with confirmation for each file
asset-cleaner delete --interactive

# Force delete without confirmation (dangerous!)
asset-cleaner delete --force
```

---

### 4. `asset-cleaner info`

Display project information and detection results.

**Arguments**: None

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--show-config` | bool | false | Display current configuration |
| `--show-paths` | bool | false | List detected asset paths |

**Exit Codes**:
- `0`: Success
- `1`: Error (directory not found)

**Output**:
```
╭─────────────────────────────────────────────╮
│  ℹ️  Project Information                     │
╰─────────────────────────────────────────────╯

📁 Project Root: /Users/user/my-project
🏷️  Project Type: React (Web)
🔧 Config File:  .unusedassets.yaml (found)

📂 Asset Directories:
  • public/images/ (123 files)
  • public/fonts/ (8 files)
  • src/assets/ (116 files)

🚫 Excluded Paths:
  • node_modules/
  • dist/
  • build/

📄 File Extensions:
  Images: .jpg, .jpeg, .png, .gif, .svg, .webp
  Fonts:  .ttf, .woff, .woff2, .eot
  Videos: .mp4, .webm, .mov
  Audio:  .mp3, .wav, .ogg

✅ Configuration valid
```

**Output (with --show-config)**:
```yaml
# Current Configuration

asset_paths:
  - public/
  - src/assets/

extensions:
  - .jpg
  - .png
  - .svg
  - .woff
  - .ttf

exclude_paths:
  - node_modules/
  - dist/
  - "**/__tests__/**"

max_workers: 8
show_progress: true
```

**Example Usage**:
```bash
# Show project info
asset-cleaner info

# Show full configuration
asset-cleaner info --show-config

# List all detected asset paths
asset-cleaner info --show-paths
```

---

### 5. `asset-cleaner init`

Initialize configuration file for current project.

**Arguments**: None

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | false | Overwrite existing config file |
| `--template` | string | `default` | Config template: `default`, `minimal`, `comprehensive` |

**Exit Codes**:
- `0`: Success (config created)
- `1`: Error (config already exists and --force not used)
- `2`: Invalid arguments

**Output**:
```
╭─────────────────────────────────────────────╮
│  📝 Initialize Configuration                 │
╰─────────────────────────────────────────────╯

✓ Detected project type: React (Web)
✓ Created .unusedassets.yaml

Configuration includes:
  • 3 asset directories
  • 15 file extensions
  • 6 exclusion patterns

Edit .unusedassets.yaml to customize settings.
Run 'asset-cleaner scan' to start scanning.
```

**Example Usage**:
```bash
# Create default config
asset-cleaner init

# Overwrite existing config
asset-cleaner init --force

# Use minimal template
asset-cleaner init --template minimal
```

---

### 6. `asset-cleaner ignore <patterns...>`

Add patterns to ignore list.

**Arguments**:
- `patterns...` (required): Glob patterns to ignore

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--global` | bool | false | Add to global ignore list (~/.asset-cleaner-ignore) |
| `--remove` | bool | false | Remove patterns from ignore list |

**Exit Codes**:
- `0`: Success (patterns added/removed)
- `1`: Error (config file not found)
- `2`: Invalid patterns

**Output**:
```
╭─────────────────────────────────────────────╮
│  🚫 Update Ignore Patterns                   │
╰─────────────────────────────────────────────╯

Adding to .unusedassets.yaml:
  • legacy/**/*.png
  • backup/**/*

✓ Updated configuration

Run 'asset-cleaner scan' to re-scan with new exclusions.
```

**Example Usage**:
```bash
# Add patterns to ignore list
asset-cleaner ignore "legacy/**/*.png" "backup/**/*"

# Remove patterns
asset-cleaner ignore --remove "backup/**/*"

# Add to global ignore list
asset-cleaner ignore --global "*.backup.*"
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ASSET_CLEANER_CONFIG` | `.unusedassets.yaml` | Default config file path |
| `ASSET_CLEANER_NO_COLOR` | `false` | Disable colored output (`true`/`false`) |
| `ASSET_CLEANER_LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `NO_COLOR` | - | Standard env var to disable colors (any value) |

---

## Configuration File Format

**Location**: `.unusedassets.yaml` (project root) or `~/.asset-cleaner/config.yaml` (global)

**Schema**:
```yaml
# Project-specific asset paths
asset_paths:
  - public/
  - src/assets/
  - static/media/

# File extensions to scan
extensions:
  - .jpg
  - .jpeg
  - .png
  - .gif
  - .svg
  - .webp
  - .ttf
  - .woff
  - .woff2
  - .mp4
  - .mp3

# Paths/patterns to exclude
exclude_paths:
  - node_modules/
  - dist/
  - build/
  - .next/
  - "**/__tests__/**"
  - "*.test.js"

# Asset constant files to analyze
constant_files:
  - src/constants/assets.ts
  - lib/assets.dart
  - app/config/AssetPaths.swift

# Base path variable names to track
base_path_vars:
  - ASSETS_BASE
  - PUBLIC_URL
  - ASSET_PREFIX

# Advanced settings
max_workers: 8              # Concurrent workers (0 = auto)
follow_symlinks: false      # Follow symbolic links
show_progress: true         # Show progress bar
color_output: true          # Enable colors
```

---

## Error Messages

### Standard Error Format

```
Error: [error message]

  Context: [additional context if available]

  Suggestion: [suggested action to fix]
```

### Common Errors

**Config file not found**:
```
Error: Configuration file not found

  Looked for: .unusedassets.yaml

  Suggestion: Run 'asset-cleaner init' to create a configuration file.
```

**Directory not found**:
```
Error: Directory does not exist: /path/to/project

  Suggestion: Check the path and try again.
```

**Permission denied**:
```
Error: Permission denied reading: /path/to/file

  Context: Scanning requires read access to all project files.

  Suggestion: Check file permissions or run with appropriate privileges.
```

**Invalid pattern**:
```
Error: Invalid glob pattern: "**[invalid"

  Context: In configuration: exclude_paths

  Suggestion: Fix the pattern syntax. See glob documentation.
```

---

## Contract Testing

### Test Cases

All commands should be tested with the following scenarios:

1. **Success paths**: Normal operation with valid inputs
2. **Error paths**: Invalid inputs, missing files, permission errors
3. **Edge cases**: Empty projects, very large projects, no unused assets
4. **Flag combinations**: Compatible and incompatible flag combinations
5. **Exit codes**: Verify correct exit codes for all scenarios

### Example Test Suite

```bash
# Test scan command
assert_exit_code 0 "asset-cleaner scan"
assert_exit_code 1 "asset-cleaner scan /nonexistent"
assert_output_contains "asset-cleaner scan" "Scan Complete"

# Test delete dry-run
assert_exit_code 0 "asset-cleaner delete --dry-run"
assert_output_contains "asset-cleaner delete --dry-run" "Dry Run Mode"
assert_files_unchanged "asset-cleaner delete --dry-run"

# Test config init
asset-cleaner init
assert_file_exists ".unusedassets.yaml"
assert_exit_code 1 "asset-cleaner init"  # Fail without --force
assert_exit_code 0 "asset-cleaner init --force"  # Succeed with --force
```

---

## Versioning & Compatibility

**Semantic Versioning**: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking CLI changes (removed commands, incompatible flags)
- **MINOR**: New commands, new flags (backward compatible)
- **PATCH**: Bug fixes, output format improvements (no behavior change)

**Version Display**:
```bash
$ asset-cleaner --version
asset-cleaner version 1.0.0
  Go version: go1.23.1
  Build date: 2025-10-23
  Commit: abc123
```

---

## Summary

This CLI interface contract ensures:

✅ **Consistent behavior** across all commands
✅ **Predictable output formats** (text, JSON, CSV)
✅ **Clear error messages** with actionable suggestions
✅ **Standard exit codes** for automation
✅ **Testable contracts** for quality assurance
✅ **Extensible design** for future commands

All implementations must conform to this contract to ensure a reliable and intuitive user experience.
