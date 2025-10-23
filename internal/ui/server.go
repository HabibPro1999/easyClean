package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/HabibPro1999/easyClean/internal/models"
)

//go:embed web/*
var webFiles embed.FS

var currentScanResult *models.ScanResult

// StartWebServer starts the web UI server
func StartWebServer(result *models.ScanResult, host string, port int) error {
	currentScanResult = result

	// Serve embedded static files from web subdirectory
	webFS, err := fs.Sub(webFiles, "web")
	if err != nil {
		return fmt.Errorf("failed to load web files: %w", err)
	}
	http.Handle("/", http.FileServer(http.FS(webFS)))

	// API endpoints
	http.HandleFunc("/api/results", handleGetResults)
	http.HandleFunc("/api/delete", handleDelete)
	http.HandleFunc("/api/asset", handleServeAsset)

	addr := fmt.Sprintf("%s:%d", host, port)
	return http.ListenAndServe(addr, nil)
}

func handleGetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentScanResult)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var request struct {
		Paths []string `json:"paths"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Delete files
	deletedCount := 0
	totalFreed := int64(0)
	var errors []string

	for _, path := range request.Paths {
		// Find asset in scan results
		var assetToDelete *models.AssetFile
		for _, asset := range currentScanResult.UnusedAssets {
			if asset.Path == path || asset.RelativePath == path {
				assetToDelete = &asset
				break
			}
		}

		if assetToDelete == nil {
			errors = append(errors, fmt.Sprintf("%s: not found in unused assets", path))
			continue
		}

		// Delete file
		if err := os.Remove(assetToDelete.Path); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", path, err))
		} else {
			deletedCount++
			totalFreed += assetToDelete.Size
		}
	}

	// Send response
	response := struct {
		Success      bool     `json:"success"`
		DeletedCount int      `json:"deleted_count"`
		TotalFreed   int64    `json:"total_freed"`
		Errors       []string `json:"errors,omitempty"`
	}{
		Success:      len(errors) == 0,
		DeletedCount: deletedCount,
		TotalFreed:   totalFreed,
		Errors:       errors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleServeAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get path from query parameter
	assetPath := r.URL.Query().Get("path")
	if assetPath == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	// Security: Validate the path is in our scan results (whitelist approach)
	if currentScanResult == nil {
		http.Error(w, "No scan results available", http.StatusNotFound)
		return
	}

	// Check if the requested path is in our asset list
	isValidAsset := false
	for _, asset := range currentScanResult.Assets {
		if asset.Path == assetPath {
			isValidAsset = true
			break
		}
	}

	if !isValidAsset {
		http.Error(w, "Asset not found in scan results", http.StatusForbidden)
		return
	}

	// Read the file
	data, err := os.ReadFile(assetPath)
	if err != nil {
		http.Error(w, "Failed to read asset file", http.StatusInternalServerError)
		return
	}

	// Set appropriate Content-Type based on file extension
	contentType := getContentType(assetPath)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=3600")

	w.Write(data)
}

func getContentType(path string) string {
	ext := ""
	for i := len(path) - 1; i >= 0 && i > len(path)-10; i-- {
		if path[i] == '.' {
			ext = path[i:]
			break
		}
	}

	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	case ".bmp":
		return "image/bmp"
	default:
		return "application/octet-stream"
	}
}
