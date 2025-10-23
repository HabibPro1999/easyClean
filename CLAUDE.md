# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**easyClean** - A Go-based CLI tool that automatically detects and safely removes unused asset files (images, fonts, videos, etc.) from codebases. The tool uses smart scanning with multi-pattern reference detection and supports multiple project types (React, Vue, Flutter, iOS, Android).

**Current Status**: Planning phase complete. Ready for implementation (see tasks.md).

## Technology Stack

- **Language**: Go 1.23+
- **CLI Framework**: Cobra v1.9.1+ (command structure)
- **TUI**: Bubble Tea v1.3.10 (progress indicators, interactive UI)
- **Styling**: Lip Gloss v2.x (terminal layouts and colors)
- **Config**: Viper v1.21.0 (YAML configuration management)
- **Parsing**: esbuild parser (JavaScript/TypeScript AST parsing)
- **Testing**: Go standard library with table-driven tests

## Project Structure

```
cmd/asset-cleaner/          # CLI entry point and commands
  commands/                 # Command implementations (scan, delete, review, etc.)
  web/                      # Embedded web UI files

internal/                   # Application logic (not importable)
  scanner/                  # File system scanning and asset discovery
  parser/                   # Code parsing (JS/TS AST, regex fallback)
  detector/                 # Project type detection
  classifier/               # Asset status classification
  config/                   # Configuration loading and defaults
  models/                   # Core data structures
  ui/                       # TUI components and web server
  utils/                    # File system and path utilities

specs/001-unused-assets-detector/  # Complete feature specification
  spec.md                   # Requirements and user stories
  plan.md                   # Technical architecture and decisions
  tasks.md                  # Actionable implementation tasks
  data-model.md             # Core data structures
  contracts/cli-interface.md # CLI command contracts

tests/
  fixtures/                 # Test project fixtures (React, Flutter, Vue)
  integration/              # End-to-end tests
  unit/                     # Unit tests
```

## Core Architecture

### Data Model (Data Structures Over Algorithms)

The system is built around these core entities:

1. **AssetFile** - Represents a discovered asset with metadata (path, size, status)
2. **Reference** - A code location referencing an asset (source file, line, type)
3. **ScanResult** - Complete scan output with statistics
4. **ProjectConfig** - User configuration (.unusedassets.yaml)
5. **ProjectType** - Detected project type (React, Flutter, iOS, etc.)

### Asset Classification

Assets are classified into four states:
- **StatusUsed** - Has active code references (keep)
- **StatusUnused** - No references found (safe to delete)
- **StatusPotentiallyUnused** - Only in comments/dead code (warn)
- **StatusNeedsManualReview** - Dynamic paths detected (exclude from auto-delete)

### Scan Workflow

```
Load Config → Detect Project Type → Walk Filesystem (collect assets) →
Parse Source Files (extract references) → Match References to Assets →
Classify Asset Status → Generate Report
```

## Commands

The tool provides these commands:

```bash
# Initialize configuration file
easyClean init [--force] [--template default|minimal|comprehensive]

# Scan project for unused assets (results auto-saved to cache)
easyClean scan [directory] [--output FILE] [--format text|json|csv]

# Review results in web UI (auto-loads from cache)
easyClean review [--port 3000] [--no-browser] [--scan-file FILE]

# Delete unused assets (auto-loads from cache)
easyClean delete [paths...] [--dry-run] [--interactive] [--force] [--scan-file FILE]

# Show project information
easyClean info [--show-config] [--show-paths]

# Manage ignore patterns
easyClean ignore <patterns...> [--remove] [--global]
```

## Configuration

Projects use `.unusedassets.yaml` in the project root:

```yaml
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
  - build/

constant_files:
  - src/constants/assets.ts

base_path_vars:
  - ASSETS_BASE
  - PUBLIC_URL

max_workers: 8
show_progress: true
```

## Development Guidelines

### Core Principles (from ultimate_principles.md)

1. **KISS** - Keep It Simple, Stupid. The simplest solution that works is the right solution.
2. **DRY (Rule of Three)** - Don't abstract until you've seen a pattern 3 times.
3. **YAGNI** - You Aren't Gonna Need It. Don't build for imaginary futures.
4. **Process** - Make It Work → Make It Right → Make It Fast

### Code Quality Standards

- **Functions**: Each does one thing. Max 50 lines.
- **Indirection**: Maximum 3 levels (main → command → service → scanner)
- **Dependencies**: Pass explicitly, no hidden state
- **Data First**: Design data structures before algorithms
- **Testability**: Functions accept interfaces (mockable for testing)

### Testing Strategy

- **Unit Tests**: Per component (scanner, parser, classifier)
  - Test happy path, edge cases, errors
  - Use table-driven tests

- **Integration Tests**: Per user story
  - Scan fixture projects with known unused files
  - Verify 95%+ accuracy, <5% false positives

- **Benchmarks**: Performance validation
  - Target: 1K files in <10s
  - Handle 100K files without crashing

## Implementation Phases

**MVP (User Story 1 - P1)**:
1. Setup (T001-T006): Go module, dependencies, directory structure
2. Foundational (T007-T014): Core data models and utilities
3. Scan Command (T015-T048): Project type detection, asset discovery, reference detection, classification, reporting

**Subsequent Features**:
- User Story 2 (P2): Preview and deletion with safety mechanisms
- User Story 3 (P3): Ignore patterns and config management
- User Story 4 (P4): Report generation and export formats

See `specs/001-unused-assets-detector/tasks.md` for complete task breakdown with dependencies.

## Performance Goals

- Scan 1,000 files in <10 seconds
- Handle projects up to 100,000 files
- Progress updates every 2 seconds
- <100MB memory for typical projects (1K-10K files)

## Safety Mechanisms

- **Three-tier classification** - Separate unused, potentially unused, and needs review
- **Dry-run by default** - Preview before deletion
- **Git history preservation** - Deletion removes from filesystem only
- **Confirmation prompts** - Multiple safety checks before deletion
- **Recovery instructions** - Clear guidance on git restore

## Key Technical Decisions

1. **Go Language** - Single binary, fast file I/O, excellent concurrency
2. **Bubble Tea TUI** - Beautiful progress indicators, Elm Architecture for state
3. **Concurrent Scanning** - Single-threaded for <10K files, parallel for larger projects
4. **Embedded Web UI** - go:embed for zero external dependencies
5. **Conservative Detection** - Prefer false negatives over false positives

## Common Development Tasks

```bash
# Initialize Go module (first time)
go mod init github.com/yourusername/asset-cleaner

# Install dependencies
go get github.com/spf13/cobra@v1.9.1
go get github.com/charmbracelet/bubbletea@v1.3.10
go get github.com/charmbracelet/lipgloss@latest
go get github.com/spf13/viper@v1.21.0

# Run tests
go test ./...

# Run specific test
go test ./internal/scanner -v -run TestAssetFinder

# Build binary
go build -o easyClean ./cmd/easyClean

# Run locally
go run ./cmd/easyClean scan

# Format code
go fmt ./...

# Run linter
golangci-lint run

# Benchmarks
go test -bench=. ./internal/scanner
```

## References

- **Full Spec**: `specs/001-unused-assets-detector/spec.md` - Requirements, user stories, success criteria
- **Implementation Plan**: `specs/001-unused-assets-detector/plan.md` - Architecture, technical decisions
- **Task List**: `specs/001-unused-assets-detector/tasks.md` - 87 actionable tasks with dependencies
- **Data Model**: `specs/001-unused-assets-detector/data-model.md` - Core structs and relationships
- **CLI Contract**: `specs/001-unused-assets-detector/contracts/cli-interface.md` - Command specifications
- **Engineering Principles**: `ultimate_principles.md` - Development philosophy and code standards
