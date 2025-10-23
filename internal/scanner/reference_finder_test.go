package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HabibPro1999/easyClean/internal/config"
	"github.com/HabibPro1999/easyClean/internal/models"
)

func TestReferenceFinder_FindReferences(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a source file with asset references
	sourceCode := `
import logo from './assets/logo.png';
const icon = require('./assets/icon.svg');
const bgImage = "assets/background.jpg";
`
	createTestFile(t, filepath.Join(tmpDir, "app.js"))
	writeContent(t, filepath.Join(tmpDir, "app.js"), sourceCode)

	cfg := config.DefaultConfig()
	cfg.AssetPaths = []string{"assets/"}

	finder := NewReferenceFinder(tmpDir, cfg)
	references, err := finder.FindReferences()

	if err != nil {
		t.Fatalf("FindReferences() failed: %v", err)
	}

	if len(references) == 0 {
		t.Error("Expected to find references, got none")
	}

	// Should find at least the three references
	totalRefs := 0
	for _, refs := range references {
		totalRefs += len(refs)
	}

	if totalRefs < 3 {
		t.Errorf("Expected at least 3 references, got %d", totalRefs)
	}
}

func TestReferenceFinder_ImportPattern(t *testing.T) {
	tmpDir := t.TempDir()

	sourceCode := `import logo from './logo.png';`
	createTestFile(t, filepath.Join(tmpDir, "app.js"))
	writeContent(t, filepath.Join(tmpDir, "app.js"), sourceCode)

	cfg := config.DefaultConfig()
	finder := NewReferenceFinder(tmpDir, cfg)
	references, err := finder.FindReferences()

	if err != nil {
		t.Fatalf("FindReferences() failed: %v", err)
	}

	found := false
	for _, refs := range references {
		for _, ref := range refs {
			if ref.Type == models.RefTypeImport {
				found = true
			}
		}
	}

	if !found {
		t.Error("Expected to find import reference")
	}
}

func TestReferenceFinder_CSSUrlPattern(t *testing.T) {
	tmpDir := t.TempDir()

	cssCode := `
.logo {
	background: url('./assets/logo.png');
}
`
	createTestFile(t, filepath.Join(tmpDir, "style.css"))
	writeContent(t, filepath.Join(tmpDir, "style.css"), cssCode)

	cfg := config.DefaultConfig()
	cfg.AssetPaths = []string{"assets/"}

	finder := NewReferenceFinder(tmpDir, cfg)
	references, err := finder.FindReferences()

	if err != nil {
		t.Fatalf("FindReferences() failed: %v", err)
	}

	found := false
	for _, refs := range references {
		for _, ref := range refs {
			if ref.Type == models.RefTypeCSSUrl {
				found = true
			}
		}
	}

	if !found {
		t.Error("Expected to find CSS url() reference")
	}
}

func TestReferenceFinder_CommentDetection(t *testing.T) {
	tmpDir := t.TempDir()

	sourceCode := `
// This is a comment with logo.png
const actualRef = "assets/icon.svg";
`
	createTestFile(t, filepath.Join(tmpDir, "app.js"))
	writeContent(t, filepath.Join(tmpDir, "app.js"), sourceCode)

	cfg := config.DefaultConfig()
	finder := NewReferenceFinder(tmpDir, cfg)
	references, err := finder.FindReferences()

	if err != nil {
		t.Fatalf("FindReferences() failed: %v", err)
	}

	commentRefs := 0
	nonCommentRefs := 0

	for _, refs := range references {
		for _, ref := range refs {
			if ref.IsComment {
				commentRefs++
			} else {
				nonCommentRefs++
			}
		}
	}

	if nonCommentRefs == 0 {
		t.Error("Expected to find non-comment references")
	}
}

func TestReferenceFinder_DynamicReferenceDetection(t *testing.T) {
	tmpDir := t.TempDir()

	sourceCode := "const path = basePath + \"/logo.png\";\nconst img = `${assetPath}/icon.svg`;"
	createTestFile(t, filepath.Join(tmpDir, "app.js"))
	writeContent(t, filepath.Join(tmpDir, "app.js"), sourceCode)

	cfg := config.DefaultConfig()
	finder := NewReferenceFinder(tmpDir, cfg)
	references, err := finder.FindReferences()

	if err != nil {
		t.Fatalf("FindReferences() failed: %v", err)
	}

	dynamicRefs := 0
	for _, refs := range references {
		for _, ref := range refs {
			if ref.IsDynamic {
				dynamicRefs++
			}
		}
	}

	if dynamicRefs == 0 {
		t.Error("Expected to detect dynamic references")
	}
}

func TestReferenceFinder_ExcludeDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source files in different directories
	createTestFile(t, filepath.Join(tmpDir, "src", "app.js"))
	writeContent(t, filepath.Join(tmpDir, "src", "app.js"), `import logo from './logo.png';`)

	createTestFile(t, filepath.Join(tmpDir, "node_modules", "lib.js"))
	writeContent(t, filepath.Join(tmpDir, "node_modules", "lib.js"), `import icon from './icon.png';`)

	cfg := config.DefaultConfig()
	cfg.ExcludePaths = []string{"node_modules/"}

	finder := NewReferenceFinder(tmpDir, cfg)
	references, err := finder.FindReferences()

	if err != nil {
		t.Fatalf("FindReferences() failed: %v", err)
	}

	// Count references - should only find from src/, not node_modules/
	for refPath, refs := range references {
		for _, ref := range refs {
			if filepath.Dir(ref.SourceFile) == filepath.Join(tmpDir, "node_modules") {
				t.Errorf("Found reference in excluded directory: %s", refPath)
			}
		}
	}
}

func TestReferenceFinder_isSourceFile(t *testing.T) {
	cfg := config.DefaultConfig()
	finder := NewReferenceFinder(".", cfg)

	tests := []struct {
		path     string
		expected bool
	}{
		{"app.js", true},
		{"component.tsx", true},
		{"style.css", true},
		{"page.vue", true},
		{"main.dart", true},
		{"index.html", true},
		{"config.json", false},
		{"data.txt", false},
		{"image.png", false},
	}

	for _, tt := range tests {
		result := finder.isSourceFile(tt.path)
		if result != tt.expected {
			t.Errorf("isSourceFile(%s) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}

func TestReferenceFinder_resolveAssetPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create actual asset files
	createTestFile(t, filepath.Join(tmpDir, "assets", "logo.png"))
	createTestFile(t, filepath.Join(tmpDir, "public", "icon.svg"))

	cfg := config.DefaultConfig()
	cfg.AssetPaths = []string{"assets/", "public/"}

	finder := NewReferenceFinder(tmpDir, cfg)

	tests := []struct {
		input    string
		expected string
	}{
		{"./assets/logo.png", filepath.Join(tmpDir, "assets", "logo.png")},
		{"/assets/logo.png", filepath.Join(tmpDir, "assets", "logo.png")},
		{"assets/logo.png", filepath.Join(tmpDir, "assets", "logo.png")},
		{"logo.png", filepath.Join(tmpDir, "assets", "logo.png")}, // basename match
	}

	for _, tt := range tests {
		result := finder.resolveAssetPath(tt.input)
		if result != tt.expected {
			t.Errorf("resolveAssetPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// Helper to write content to file
func writeContent(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write to file %s: %v", path, err)
	}
}
