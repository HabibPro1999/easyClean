// Package utils - Server registry for tracking active review servers
//
// Maintains a registry of all running review servers to support multi-project workflows.
// Registry is stored in ~/.cache/easyClean/servers.json
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	registryFileName = "servers.json"
)

// ServerInfo contains information about a running review server
type ServerInfo struct {
	ProjectPath string    `json:"project_path"` // Absolute path to project
	ProjectName string    `json:"project_name"` // Base directory name
	Port        int       `json:"port"`         // Server port
	PID         int       `json:"pid"`          // Process ID
	StartTime   time.Time `json:"start_time"`   // When server started
}

// serverRegistry holds all active servers
type serverRegistry struct {
	Servers []ServerInfo `json:"servers"`
	mu      sync.Mutex   // Protects file operations
}

// getRegistryPath returns the path to the server registry file
func getRegistryPath() (string, error) {
	cacheDir, err := GetUserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, registryFileName), nil
}

// loadRegistry loads the registry from disk, creating it if it doesn't exist
func loadRegistry() (*serverRegistry, error) {
	registryPath, err := getRegistryPath()
	if err != nil {
		return nil, err
	}

	// Ensure cache directory exists
	cacheDir := filepath.Dir(registryPath)
	if err := EnsureCacheDirExists(cacheDir); err != nil {
		return nil, err
	}

	// If registry doesn't exist, return empty registry
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		return &serverRegistry{Servers: []ServerInfo{}}, nil
	}

	// Read existing registry
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}

	var registry serverRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	return &registry, nil
}

// saveRegistry saves the registry to disk
func (r *serverRegistry) save() error {
	registryPath, err := getRegistryPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// RegisterServer adds a server to the registry
func RegisterServer(info ServerInfo) error {
	registry, err := loadRegistry()
	if err != nil {
		return err
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	// Remove any existing entry for this PID (in case of duplicate)
	for i := len(registry.Servers) - 1; i >= 0; i-- {
		if registry.Servers[i].PID == info.PID {
			registry.Servers = append(registry.Servers[:i], registry.Servers[i+1:]...)
		}
	}

	// Add new entry
	registry.Servers = append(registry.Servers, info)

	return registry.save()
}

// UnregisterServer removes a server from the registry by PID
func UnregisterServer(pid int) error {
	registry, err := loadRegistry()
	if err != nil {
		return err
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	// Remove entry with matching PID
	found := false
	for i := len(registry.Servers) - 1; i >= 0; i-- {
		if registry.Servers[i].PID == pid {
			registry.Servers = append(registry.Servers[:i], registry.Servers[i+1:]...)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("server with PID %d not found in registry", pid)
	}

	return registry.save()
}

// UnregisterServerByPort removes a server from the registry by port
func UnregisterServerByPort(port int) error {
	registry, err := loadRegistry()
	if err != nil {
		return err
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	// Remove entry with matching port
	found := false
	for i := len(registry.Servers) - 1; i >= 0; i-- {
		if registry.Servers[i].Port == port {
			registry.Servers = append(registry.Servers[:i], registry.Servers[i+1:]...)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("server on port %d not found in registry", port)
	}

	return registry.save()
}

// isProcessAlive checks if a process with given PID is still running
func isProcessAlive(pid int) bool {
	// Try to send signal 0 (no-op signal) to check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix, FindProcess always succeeds, so we need to actually signal
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// GetActiveServers returns all currently active servers (with PID alive check)
func GetActiveServers() ([]ServerInfo, error) {
	registry, err := loadRegistry()
	if err != nil {
		return nil, err
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	var activeServers []ServerInfo

	// Filter out dead servers
	for _, server := range registry.Servers {
		if isProcessAlive(server.PID) {
			activeServers = append(activeServers, server)
		}
	}

	return activeServers, nil
}

// CleanupDeadServers removes servers with dead PIDs from the registry
func CleanupDeadServers() error {
	registry, err := loadRegistry()
	if err != nil {
		return err
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	// Keep only servers with alive PIDs
	var aliveServers []ServerInfo
	for _, server := range registry.Servers {
		if isProcessAlive(server.PID) {
			aliveServers = append(aliveServers, server)
		}
	}

	registry.Servers = aliveServers
	return registry.save()
}

// GetServerByPort finds a server in the registry by port
func GetServerByPort(port int) (*ServerInfo, error) {
	servers, err := GetActiveServers()
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.Port == port {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("no active server found on port %d", port)
}
