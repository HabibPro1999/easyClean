package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetUserCacheDir(t *testing.T) {
	cacheDir, err := GetUserCacheDir()
	if err != nil {
		t.Fatalf("GetUserCacheDir() failed: %v", err)
	}

	if cacheDir == "" {
		t.Error("GetUserCacheDir() returned empty string")
	}

	if !strings.Contains(cacheDir, "easyClean") {
		t.Errorf("GetUserCacheDir() = %q, should contain 'easyClean'", cacheDir)
	}
}

func TestGetProjectHash(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		wantLen     int
		wantErr     bool
	}{
		{
			name:        "absolute path",
			projectRoot: "/home/user/project",
			wantLen:     12,
			wantErr:     false,
		},
		{
			name:        "relative path",
			projectRoot: ".",
			wantLen:     12,
			wantErr:     false,
		},
		{
			name:        "path with spaces",
			projectRoot: "/home/user/my project",
			wantLen:     12,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GetProjectHash(tt.projectRoot)

			if tt.wantErr {
				if err == nil {
					t.Error("GetProjectHash() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetProjectHash() unexpected error: %v", err)
			}

			if len(hash) != tt.wantLen {
				t.Errorf("GetProjectHash() hash length = %d, want %d", len(hash), tt.wantLen)
			}

			// Hash should be hex characters only
			for _, c := range hash {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
					t.Errorf("GetProjectHash() hash contains invalid character: %c", c)
				}
			}
		})
	}
}

func TestGetProjectHash_Consistency(t *testing.T) {
	// Same path should always produce same hash
	projectRoot := "/home/user/test-project"

	hash1, err := GetProjectHash(projectRoot)
	if err != nil {
		t.Fatalf("GetProjectHash() failed: %v", err)
	}

	hash2, err := GetProjectHash(projectRoot)
	if err != nil {
		t.Fatalf("GetProjectHash() failed: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("GetProjectHash() not consistent: %q != %q", hash1, hash2)
	}
}

func TestGetProjectHash_Uniqueness(t *testing.T) {
	// Different paths should produce different hashes
	path1 := "/home/user/project1"
	path2 := "/home/user/project2"

	hash1, err := GetProjectHash(path1)
	if err != nil {
		t.Fatalf("GetProjectHash() failed: %v", err)
	}

	hash2, err := GetProjectHash(path2)
	if err != nil {
		t.Fatalf("GetProjectHash() failed: %v", err)
	}

	if hash1 == hash2 {
		t.Errorf("GetProjectHash() produced same hash for different paths: %q", hash1)
	}
}

func TestGetProjectCacheDir(t *testing.T) {
	projectRoot := "/home/user/test-project"

	cacheDir, err := GetProjectCacheDir(projectRoot)
	if err != nil {
		t.Fatalf("GetProjectCacheDir() failed: %v", err)
	}

	if cacheDir == "" {
		t.Error("GetProjectCacheDir() returned empty string")
	}

	if !strings.Contains(cacheDir, "easyClean") {
		t.Errorf("GetProjectCacheDir() = %q, should contain 'easyClean'", cacheDir)
	}

	if !strings.Contains(cacheDir, "projects") {
		t.Errorf("GetProjectCacheDir() = %q, should contain 'projects'", cacheDir)
	}

	// Should end with a hash directory
	parts := strings.Split(cacheDir, string(filepath.Separator))
	lastPart := parts[len(parts)-1]
	if len(lastPart) != 12 {
		t.Errorf("GetProjectCacheDir() last directory should be 12-char hash, got: %q", lastPart)
	}
}

func TestGetScanResultsPath(t *testing.T) {
	projectRoot := "/home/user/test-project"

	scanPath, err := GetScanResultsPath(projectRoot)
	if err != nil {
		t.Fatalf("GetScanResultsPath() failed: %v", err)
	}

	if scanPath == "" {
		t.Error("GetScanResultsPath() returned empty string")
	}

	if !strings.HasSuffix(scanPath, "scan-results.json") {
		t.Errorf("GetScanResultsPath() = %q, should end with 'scan-results.json'", scanPath)
	}

	if !strings.Contains(scanPath, "easyClean") {
		t.Errorf("GetScanResultsPath() = %q, should contain 'easyClean'", scanPath)
	}
}

func TestEnsureCacheDirExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := filepath.Join(os.TempDir(), "asset-cleaner-test")
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "projects", "abc123", "subdir")

	// Directory should not exist initially
	if _, err := os.Stat(testPath); err == nil {
		t.Fatal("Test directory already exists")
	}

	// Create the directory
	if err := EnsureCacheDirExists(testPath); err != nil {
		t.Fatalf("EnsureCacheDirExists() failed: %v", err)
	}

	// Directory should now exist
	if _, err := os.Stat(testPath); err != nil {
		t.Errorf("EnsureCacheDirExists() did not create directory: %v", err)
	}

	// Calling again should not error (idempotent)
	if err := EnsureCacheDirExists(testPath); err != nil {
		t.Errorf("EnsureCacheDirExists() failed on second call: %v", err)
	}
}

func TestGetScanResultsPathOrDefault(t *testing.T) {
	projectRoot := "/home/user/test-project"
	customPath := "/custom/path/results.json"

	tests := []struct {
		name        string
		projectRoot string
		defaultPath string
		wantCustom  bool
	}{
		{
			name:        "use custom path",
			projectRoot: projectRoot,
			defaultPath: customPath,
			wantCustom:  true,
		},
		{
			name:        "use generated path",
			projectRoot: projectRoot,
			defaultPath: "",
			wantCustom:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := GetScanResultsPathOrDefault(tt.projectRoot, tt.defaultPath)
			if err != nil {
				t.Fatalf("GetScanResultsPathOrDefault() failed: %v", err)
			}

			if tt.wantCustom {
				if path != customPath {
					t.Errorf("GetScanResultsPathOrDefault() = %q, want %q", path, customPath)
				}
			} else {
				if !strings.Contains(path, "easyClean") {
					t.Errorf("GetScanResultsPathOrDefault() = %q, should contain 'easyClean'", path)
				}
			}
		})
	}
}
