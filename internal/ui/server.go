package ui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/HabibPro1999/easyClean/internal/models"
)

//go:embed web/*
var webFiles embed.FS

// ReviewServer wraps an HTTP server for the review UI
type ReviewServer struct {
	server     *http.Server
	scanResult *models.ScanResult
}

// NewReviewServer creates a new review server instance
func NewReviewServer(result *models.ScanResult, host string, port int) (*ReviewServer, error) {
	rs := &ReviewServer{
		scanResult: result,
	}

	// Serve embedded static files from web subdirectory
	webFS, err := fs.Sub(webFiles, "web")
	if err != nil {
		return nil, fmt.Errorf("failed to load web files: %w", err)
	}

	// Create HTTP mux
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(webFS)))

	// API endpoints (use closures to access scanResult)
	mux.HandleFunc("/api/results", rs.handleGetResults)
	mux.HandleFunc("/api/delete", rs.handleDelete)
	mux.HandleFunc("/api/asset", rs.handleServeAsset)

	// Create HTTP server
	rs.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return rs, nil
}

// Start starts the web server (blocking)
func (rs *ReviewServer) Start() error {
	return rs.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (rs *ReviewServer) Shutdown(ctx context.Context) error {
	return rs.server.Shutdown(ctx)
}

func (rs *ReviewServer) handleGetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rs.scanResult)
}

func (rs *ReviewServer) handleDelete(w http.ResponseWriter, r *http.Request) {
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
		for _, asset := range rs.scanResult.UnusedAssets {
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

func (rs *ReviewServer) handleServeAsset(w http.ResponseWriter, r *http.Request) {
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
	if rs.scanResult == nil {
		http.Error(w, "No scan results available", http.StatusNotFound)
		return
	}

	// Check if the requested path is in our asset list
	isValidAsset := false
	for _, asset := range rs.scanResult.Assets {
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
