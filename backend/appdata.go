package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func AppDataDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "Juchuan"), nil
	case "windows":
		// Windows portable mode: keep data beside the executable.
		// This makes backup and migration easier for a standalone app.
		dir, err := os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Dir(dir), nil
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".juchuan"), nil
	}
}
