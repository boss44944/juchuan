//go:build !darwin && !windows

package main

// StartTray is disabled on Linux.
// Linux keeps the command line mode.
func StartTray(url string, quit chan struct{}) {
}
