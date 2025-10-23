# Quickstart Guide: Unused Assets Detector

**Version**: 1.0.0
**Last Updated**: 2025-10-23

## Overview

Get started with the Unused Assets Detector in under 5 minutes. This guide walks you through installation, first scan, and cleanup.

---

## Installation

### Option 1: Install via Homebrew (macOS/Linux)

```bash
brew install asset-cleaner
```

### Option 2: Download Binary

Visit [Releases](https://github.com/yourusername/asset-cleaner/releases) and download for your platform:

- **macOS**: `asset-cleaner-darwin-amd64` or `asset-cleaner-darwin-arm64`
- **Linux**: `asset-cleaner-linux-amd64`
- **Windows**: `asset-cleaner-windows-amd64.exe`

Make executable (macOS/Linux):
```bash
chmod +x asset-cleaner
sudo mv asset-cleaner /usr/local/bin/
```

### Option 3: Build from Source

```bash
git clone https://github.com/yourusername/asset-cleaner.git
cd asset-cleaner
go build -o asset-cleaner ./cmd/asset-cleaner
sudo mv asset-cleaner /usr/local/bin/
```

### Verify Installation

```bash
asset-cleaner --version
```

Expected output:
```
asset-cleaner version 1.0.0
```

---

## Quick Start (3 Steps)

### Step 1: Navigate to Your Project

```bash
cd /path/to/your/project
```

### Step 2: Scan for Unused Assets

```bash
asset-cleaner scan
```

**What happens:**
- Auto-detects project type (React, Vue, Flutter, etc.)
- Finds asset directories (`public/`, `assets/`, etc.)
- Analyzes code references
- Generates report

**Example Output:**
```
╭─────────────────────────────────────────────╮
│  🔍 Asset Cleaner v1.0                      │
╰─────────────────────────────────────────────╯

✓ Found: React (Web)

📁 Scanning asset directories...
✓ public/images/ (123 files)
✓ src/assets/ (45 files)

⠙ Analyzing code references...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 Scan Complete

  Total Assets:           168
  ✓ Used Assets:          142
  ⚠️  Unused Assets:       26

  💾 Potential Savings:   5.2 MB

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Step 3: Review and Delete

**Option A: Interactive Review (Recommended)**

```bash
asset-cleaner review
```

Opens browser at `http://localhost:3000` with visual interface showing:
- Thumbnail previews
- File sizes
- Last modified dates
- Checkboxes to select files for deletion

**Option B: Delete All Unused (Advanced)**

```bash
# Preview what would be deleted
asset-cleaner delete --dry-run

# Delete after reviewing
asset-cleaner delete
```

Confirms before deletion:
```
⚠️  You are about to delete 26 files (5.2 MB).
Continue? [y/N]: y
```

### Step 4: Commit Changes

```bash
git add -u
git commit -m "Remove unused assets (5.2 MB freed)"
```

---

## Common Workflows

### Workflow 1: First-Time Cleanup

**Scenario**: You inherited a project with years of asset bloat.

```bash
# 1. Scan the project
asset-cleaner scan

# 2. Review interactively (safest)
asset-cleaner review
#    - Uncheck anything you want to keep
#    - Click "Delete Selected"

# 3. Test your app thoroughly
npm start
#    OR
flutter run

# 4. If something broke, restore from git
git restore public/images/important-file.png

# 5. Commit when confident
git add -u
git commit -m "Clean up unused assets"
```

---

### Workflow 2: CI/CD Integration

**Scenario**: Automatically detect unused assets in pull requests.

```yaml
# .github/workflows/asset-check.yml
name: Check Unused Assets

on: [pull_request]

jobs:
  asset-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install asset-cleaner
        run: |
          curl -L https://github.com/yourusername/asset-cleaner/releases/latest/download/asset-cleaner-linux-amd64 -o asset-cleaner
          chmod +x asset-cleaner
          sudo mv asset-cleaner /usr/local/bin/

      - name: Scan for unused assets
        run: asset-cleaner scan --output results.json --format json

      - name: Check for unused assets
        run: |
          UNUSED_COUNT=$(jq '.statistics.unused_count' results.json)
          if [ "$UNUSED_COUNT" -gt 0 ]; then
            echo "⚠️ Found $UNUSED_COUNT unused assets"
            jq -r '.unused_assets[].relative_path' results.json
            exit 1
          fi
```

---

### Workflow 3: Pre-Commit Hook

**Scenario**: Prevent committing new unused assets.

Create `.git/hooks/pre-commit`:
```bash
#!/bin/bash

echo "Checking for unused assets..."

asset-cleaner scan --quiet

UNUSED=$(asset-cleaner scan --format json --output /tmp/scan.json 2>/dev/null && jq '.statistics.unused_count' /tmp/scan.json)

if [ "$UNUSED" -gt 0 ]; then
  echo "⚠️  Warning: Found $UNUSED unused assets"
  echo "Run 'asset-cleaner review' to clean up before committing."

  # Uncomment to block commit:
  # exit 1
fi
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

### Workflow 4: Scheduled Cleanup (Monthly)

**Scenario**: Automatically run cleanup and create PR.

```yaml
# .github/workflows/monthly-cleanup.yml
name: Monthly Asset Cleanup

on:
  schedule:
    - cron: '0 0 1 * *'  # First day of each month
  workflow_dispatch:

jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install asset-cleaner
        run: |
          curl -L https://github.com/yourusername/asset-cleaner/releases/latest/download/asset-cleaner-linux-amd64 -o asset-cleaner
          chmod +x asset-cleaner
          sudo mv asset-cleaner /usr/local/bin/

      - name: Scan and delete unused assets
        run: |
          asset-cleaner scan --output before.json --format json
          UNUSED=$(jq '.statistics.unused_count' before.json)

          if [ "$UNUSED" -gt 0 ]; then
            asset-cleaner delete --force

            SIZE=$(jq '.statistics.unused_size_bytes' before.json | awk '{print $1/1024/1024 "MB"}')

            git config user.name "Asset Cleaner Bot"
            git config user.email "bot@example.com"
            git checkout -b cleanup/unused-assets-$(date +%Y%m)
            git add -u
            git commit -m "Remove $UNUSED unused assets ($SIZE freed)"
            git push origin cleanup/unused-assets-$(date +%Y%m)

            gh pr create --title "Clean up unused assets" --body "Automated cleanup removed $UNUSED files ($SIZE)"
          fi
```

---

## Configuration

### Create Configuration File

```bash
asset-cleaner init
```

Creates `.unusedassets.yaml`:
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
  - "**/__tests__/**"

max_workers: 8
show_progress: true
```

### Customize for Your Project

**Example: Flutter Project**
```yaml
asset_paths:
  - assets/images/
  - assets/icons/
  - assets/fonts/

constant_files:
  - lib/constants/assets.dart

extensions:
  - .png
  - .jpg
  - .svg
  - .ttf
```

**Example: React Native Project**
```yaml
asset_paths:
  - assets/
  - src/assets/

base_path_vars:
  - ASSETS_BASE
  - IMAGE_PATH

extensions:
  - .png
  - .jpg
  - '@2x.png'
  - '@3x.png'
```

---

## Troubleshooting

### Problem: False Positives (Used Assets Marked as Unused)

**Cause**: Dynamic path construction not detected.

**Solution**: Add to ignore list
```bash
asset-cleaner ignore "path/to/dynamic/**/*.png"
```

Or add to `.unusedassets.yaml`:
```yaml
exclude_paths:
  - path/to/dynamic/**/*.png
```

---

### Problem: Scan Takes Too Long

**Cause**: Large project (>100K files) or slow disk I/O.

**Solution**: Increase workers or exclude unnecessary directories
```yaml
max_workers: 16  # Increase parallelism

exclude_paths:
  - node_modules/
  - dist/
  - build/
  - .next/
  - coverage/
```

---

### Problem: Permission Denied Errors

**Cause**: Insufficient permissions to read files.

**Solution**: Run with appropriate permissions
```bash
sudo asset-cleaner scan
```

Or fix permissions:
```bash
chmod -R u+r ./public
```

---

### Problem: Web UI Won't Start (Port in Use)

**Cause**: Port 3000 already in use.

**Solution**: Use a different port
```bash
asset-cleaner review --port 8080
```

---

### Problem: Tool Not Detecting Project Type

**Cause**: Non-standard project structure.

**Solution**: Manually specify asset paths
```bash
asset-cleaner init
# Edit .unusedassets.yaml with your custom paths
asset-cleaner scan
```

---

## Best Practices

### ✅ DO

1. **Run in version control** - Always commit before running delete commands
2. **Test thoroughly** - Run your app after deletion to catch false positives
3. **Use dry-run first** - Always preview with `--dry-run` before actual deletion
4. **Commit incrementally** - Delete and test in batches, not all at once
5. **Review manually** - Use `asset-cleaner review` for visual confirmation
6. **Configure excludes** - Add dynamic asset directories to exclude list

### ❌ DON'T

1. **Don't use --force** - Skip the `--force` flag unless you're 100% confident
2. **Don't delete without backup** - Always have committed changes or a backup
3. **Don't skip testing** - Never deploy without testing after asset deletion
4. **Don't ignore warnings** - Pay attention to "needs manual review" flags
5. **Don't run on production** - Only run on development machines or CI
6. **Don't delete without understanding** - If unsure why an asset is unused, investigate first

---

## Advanced Usage

### Custom Patterns

Define custom regex patterns for detecting references:

```yaml
custom_patterns:
  - 'getImage\("([^"]+)"\)'  # Matches: getImage("logo.png")
  - "require\('([^']+\.png)'\)"  # Matches: require('icon.png')
```

### Multiple Projects

Scan multiple projects in one command:

```bash
# Create a script: scan-all.sh
#!/bin/bash

for project in project1 project2 project3; do
  echo "Scanning $project..."
  cd "$project"
  asset-cleaner scan --output "../reports/$project.json"
  cd ..
done
```

### Export and Analyze

Export results for external analysis:

```bash
# Export to JSON
asset-cleaner scan --output results.json

# Analyze with jq
jq '.unused_assets | sort_by(.size_bytes) | reverse | .[0:10]' results.json

# Export to CSV for Excel
asset-cleaner scan --output results.csv --format csv
```

---

## Real-World Examples

### Example 1: React Project Cleanup

**Before**:
```
public/images/    - 234 files, 45 MB
src/assets/       - 89 files, 12 MB
Total:            - 323 files, 57 MB
```

**Command**:
```bash
asset-cleaner scan
asset-cleaner review
# Selected 78 unused files
```

**After**:
```
public/images/    - 189 files, 32 MB
src/assets/       - 56 files, 8 MB
Total:            - 245 files, 40 MB

Savings: 78 files, 17 MB (30% reduction)
```

---

### Example 2: Flutter App Cleanup

**Before**:
```
assets/images/    - 156 files, 28 MB
assets/icons/     - 67 files, 2.1 MB
Total:            - 223 files, 30.1 MB
```

**Command**:
```bash
asset-cleaner scan --extensions .png,.jpg
asset-cleaner delete --dry-run
asset-cleaner delete
```

**After**:
```
assets/images/    - 134 files, 24 MB
assets/icons/     - 58 files, 1.8 MB
Total:            - 192 files, 25.8 MB

Savings: 31 files, 4.3 MB (14% reduction)
```

---

## Next Steps

1. **Integrate into CI/CD** - Add automated checks to your pipeline
2. **Set up pre-commit hooks** - Prevent new unused assets from being committed
3. **Schedule regular cleanups** - Run monthly scans to keep projects lean
4. **Share configuration** - Commit `.unusedassets.yaml` for team consistency
5. **Monitor trends** - Track asset count and size over time

---

## Get Help

- **Documentation**: [Full documentation](https://github.com/yourusername/asset-cleaner/docs)
- **Issues**: [Report bugs or request features](https://github.com/yourusername/asset-cleaner/issues)
- **Community**: [Join discussions](https://github.com/yourusername/asset-cleaner/discussions)

---

## Quick Reference

```bash
# Essential commands
asset-cleaner scan                  # Scan current directory
asset-cleaner review                # Launch web UI
asset-cleaner delete --dry-run      # Preview deletions
asset-cleaner delete                # Delete unused assets
asset-cleaner info                  # Show project info
asset-cleaner init                  # Create config file

# Useful flags
--output results.json               # Export results
--format json|csv|text              # Output format
--dry-run                           # Preview only
--interactive                       # Confirm each file
--verbose                           # Detailed logging
--no-progress                       # Hide progress bar

# Configuration
.unusedassets.yaml                  # Project config file
~/.asset-cleaner/config.yaml        # Global config file
```

---

**Happy cleaning!** 🧹
