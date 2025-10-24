# 🧹 easyClean

**Version:** 1.0.1 | **Status:** Stable

Automatically detect and safely remove unused asset files from your codebase.

**easyClean** scans your project for unused images, fonts, videos, and other assets. It's fast, safe, and supports 10+ project types including React, Vue, Flutter, iOS, and Android.

---

## 🚀 Installation

### macOS / Linux (From Source)
```bash
git clone https://github.com/yourusername/easyClean.git
cd easyClean
go build -o easyClean ./cmd/easyClean
sudo mv easyClean /usr/local/bin/
```

### Windows (From Source)
```bash
git clone https://github.com/yourusername/easyClean.git
cd easyClean
go build -o easyClean.exe ./cmd/easyClean
# Move easyClean.exe to your PATH
```

---

## ⚡ Quick Start

### 1️⃣ Scan Your Project
```bash
easyClean scan
```
Automatically detects your project type and scans for unused assets.

### 2️⃣ Review Results
```bash
easyClean review
```
Opens an interactive web UI to browse and preview unused assets.

### 3️⃣ Delete Safely
```bash
easyClean delete --dry-run    # Preview what will be deleted
easyClean delete              # Delete with confirmation
```

That's it! Results are cached automatically between commands.

---

## 🔍 How It Works

easyClean uses a sophisticated multi-layer detection system to identify unused assets with high accuracy:

### Detection Pipeline

1. **Project Type Detection**
   - Automatically identifies your project type (React, Vue, Angular, Flutter, etc.)
   - Analyzes `package.json`, `pubspec.yaml`, `go.mod`, and other framework markers
   - Applies framework-specific detection patterns

2. **Asset Discovery**
   - Scans configured asset directories (public/, assets/, static/, etc.)
   - Identifies all asset files by extension (.png, .jpg, .svg, .woff, .mp4, etc.)
   - Catalogs file metadata (size, path, modification time)

3. **Reference Detection** (Hybrid Approach)

   **For JavaScript/TypeScript Projects:**
   - **AST Parsing**: Deep code analysis using esbuild parser
     - Static imports: `import logo from './logo.png'`
     - Dynamic imports: `import('./image.png')`
     - JSX references: `<img src={logo} />`
     - Object properties: `{ background: './bg.jpg' }`

   **Framework-Specific Patterns:**
   - **React**: `React.lazy()`, Next.js public folder conventions
   - **Angular**: `templateUrl`, `styleUrls`, lazy route loading
   - **Vue**: `defineAsyncComponent`, template bindings
   - **Flutter**: `Image.asset()`, `AssetImage()`, pubspec declarations

   **Generic Patterns** (all projects):
   - Import/require statements
   - CSS url() references
   - HTML src/href attributes
   - String literals with asset paths

4. **Smart Classification**
   - **Used**: Active code references found → Keep
   - **Unused**: No references anywhere → Safe to delete
   - **Potentially Unused**: Only in comments/dead code → Review
   - **Needs Manual Review**: Dynamic paths detected → Exclude from auto-delete

5. **Confidence Scoring**
   - AST-detected references: **100% confidence**
   - Framework patterns: **95-100% confidence**
   - Generic string matches: **70-80% confidence**

### Why It's Accurate

✅ **Framework-aware**: Uses React.lazy, Angular lazy routes, Vue async components
✅ **Deep analysis**: AST parsing finds complex patterns regex would miss
✅ **Multi-pattern**: Combines 15+ detection patterns per framework
✅ **Conservative**: Prefers false negatives over false positives (safety first)

---

## 📋 Available Commands

| Command | Purpose | Example |
|---------|---------|---------|
| **scan** | Detect unused assets | `easyClean scan ./my-project` |
| **review** | Web UI to browse results | `easyClean review --port 3000` |
| **delete** | Remove unused files | `easyClean delete --dry-run` |
| **init** | Create config file | `easyClean init --template default` |
| **info** | Show project details | `easyClean info --show-config` |

---

## 🔧 Scan Options

```bash
easyClean scan [directory] [flags]

Flags:
  --extensions string    Assets to scan (.png, .jpg, .svg, etc.)
  --exclude string       Paths to exclude (glob patterns)
  -f, --format string    Output format: text, json, csv (default: text)
  -o, --output string    Save results to file
  --no-progress          Disable progress bar
```

### Example
```bash
# Scan with custom extensions
easyClean scan . --extensions .png,.jpg,.svg

# Export to JSON
easyClean scan . --format json --output results.json

# Exclude specific paths
easyClean scan . --exclude "node_modules/*" --exclude "dist/*"
```

---

## 🗑️ Delete Options

```bash
easyClean delete [paths...] [flags]

Flags:
  --dry-run              Preview deletions without removing files
  -i, --interactive      Prompt before deleting each file
  --force                Skip confirmation (use with caution!)
  --scan-file string     Use specific scan results file
```

### Examples
```bash
# Preview what would be deleted
easyClean delete --dry-run

# Delete with confirmation
easyClean delete

# Delete specific files
easyClean delete path/to/unused1.png path/to/unused2.jpg

# Interactive mode (confirm each file)
easyClean delete --interactive
```

---

## 🌐 Multi-Project Review

Run review servers for multiple projects simultaneously:

```bash
# Terminal 1 - Project A
cd ~/my-react-app
easyClean review
# 🌐 Starting server at http://localhost:3000

# Terminal 2 - Project B
cd ~/other-project
easyClean review
# ⚠️  Port 3000 already in use, using port 3001
# 🌐 Starting server at http://localhost:3001
```

**List all active servers:**
```bash
easyClean review --list

Active Review Servers:

┌─────────────────────────────────┬──────┬─────────┬──────────┐
│ Project                         │ Port │ PID     │ Uptime   │
├─────────────────────────────────┼──────┼─────────┼──────────┤
│ my-react-app                    │ 3000 │ 12345   │ 5m 23s   │
│ other-project                   │ 3001 │ 12346   │ 1m 10s   │
└─────────────────────────────────┴──────┴─────────┴──────────┘
```

**Stop a specific server:**
```bash
easyClean review --kill 3001
```

See [MULTI_PROJECT_REVIEW.md](MULTI_PROJECT_REVIEW.md) for full documentation.

---

## 🎯 Features

✅ **Smart Detection**
- Framework-aware pattern recognition (React, Angular, Vue, Flutter, etc.)
- Hybrid AST parsing + regex-based detection for accuracy
- Multi-pattern reference detection (15+ patterns per framework)
- Supports 10+ project types
- Handles dynamic paths with confidence scoring

✅ **Three-Tier Classification**
- Used (safe to keep)
- Unused (safe to delete)
- Potentially unused (review first)
- Needs manual review (dynamic references)

✅ **Safety First**
- Dry-run mode by default
- Confirmation prompts
- Git-aware warnings
- Conservative accuracy (prefers false negatives)

✅ **Multi-Project Support**
- Run concurrent review servers on different ports
- Intelligent port allocation (3000→3001→3002...)
- List and manage all active servers
- Graceful shutdown with signal handling

✅ **No Project Pollution**
- Results cached in OS temp directory
- Never creates files in your project
- Automatic cleanup

---

## 🌍 Supported Projects

- React / Next.js / React Native
- Vue / Nuxt
- Angular / Svelte
- Flutter
- iOS (Swift)
- Android (Kotlin/Java)
- Go / Rust
- And more...

---

## 🛡️ Safety

All deletions require explicit confirmation. Results are cached automatically to your OS cache directory:

- **macOS/Linux:** `~/.cache/easyClean/`
- **Windows:** `%LOCALAPPDATA%\easyClean\cache\`

---

## 📖 Global Flags

Available on all commands:

```
  -c, --config string    Config file path (default: .unusedassets.yaml)
  -v, --verbose          Enable verbose logging
  -q, --quiet            Suppress all output except errors
      --no-color         Disable colored output
      --help             Show command help
```

---

## ⚙️ Configuration

Create a `.unusedassets.yaml` in your project root to customize behavior:

```yaml
asset_paths:
  - public/
  - src/assets/
  - static/

extensions:
  - .jpg
  - .png
  - .svg
  - .woff

exclude_paths:
  - node_modules/
  - dist/
  - build/

max_workers: 8
show_progress: true
```

Generate default config:
```bash
easyClean init
easyClean init --template minimal      # Minimal config
easyClean init --template comprehensive # All options
```

---

## 📊 Performance

- Scans 1,000 files in < 10 seconds
- Handles projects up to 100,000 files
- < 100MB memory usage for typical projects

---

## 🤝 Contributing

Found a bug? Have a suggestion? Open an issue or submit a pull request!

---

## 📄 License

MIT License - See LICENSE file for details

---

## ❓ Need Help?

```bash
easyClean --help
easyClean scan --help
easyClean review --help
easyClean delete --help
```

**Happy cleaning! 🧹**
