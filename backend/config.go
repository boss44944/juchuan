package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port     int    `json:"port"`
	AutoOpen bool   `json:"auto_open"`
	Password string `json:"password"`
}

func configFilePath() (string, error) {
	root, err := AppDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	c := &Config{Port: 8000, AutoOpen: true}
	path, err := configFilePath()
	if err != nil {
		return c, err
	}

	b, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(b, c)
	}
	return c, nil
}

func SaveConfig(c *Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	b, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(path, b, 0644)
}
