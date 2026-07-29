package main

import (
	"fmt"
	"net"
)

// FindAvailablePort returns the configured port if available,
// otherwise finds another free port automatically.
func FindAvailablePort(preferred int) int {
	if isPortAvailable(preferred) {
		return preferred
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return preferred
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
