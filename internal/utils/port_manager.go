// Package utils - Port management for multi-project review servers
//
// Provides intelligent port allocation to support multiple concurrent review sessions.
// Automatically finds available ports in the range 3000-3009.
package utils

import (
	"fmt"
	"net"
)

const (
	// DefaultPort is the preferred starting port for review servers
	DefaultPort = 3000

	// MaxPort is the maximum port to try (3009 = 10 concurrent servers)
	MaxPort = 3009

	// DefaultHost is the default host for review servers
	DefaultHost = "localhost"
)

// IsPortAvailable checks if a specific port is available for use
func IsPortAvailable(port int) bool {
	addr := fmt.Sprintf("%s:%d", DefaultHost, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// FindAvailablePort finds the first available port starting from preferredPort
// Returns the available port or an error if all ports in range are taken
func FindAvailablePort(preferredPort int) (int, error) {
	// Validate preferred port is in range
	if preferredPort < DefaultPort || preferredPort > MaxPort {
		preferredPort = DefaultPort
	}

	// Try preferred port first
	if IsPortAvailable(preferredPort) {
		return preferredPort, nil
	}

	// Try remaining ports in range
	for port := DefaultPort; port <= MaxPort; port++ {
		if port == preferredPort {
			continue // Already tried
		}
		if IsPortAvailable(port) {
			return port, nil
		}
	}

	// All ports exhausted
	return 0, fmt.Errorf("no available ports in range %d-%d (all %d ports are in use)",
		DefaultPort, MaxPort, MaxPort-DefaultPort+1)
}

// GetPortRange returns the min and max ports for review servers
func GetPortRange() (int, int) {
	return DefaultPort, MaxPort
}
