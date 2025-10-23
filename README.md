# 🧹 easyClean

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

## 🎯 Features

✅ **Smart Detection**
- Multi-pattern reference detection
- Supports 10+ project types
- Handles dynamic paths

✅ **Three-Tier Classification**
- Used (safe to keep)
- Unused (safe to delete)
- Potentially unused (review first)

✅ **Safety First**
- Dry-run mode by default
- Confirmation prompts
- Git-aware warnings

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
