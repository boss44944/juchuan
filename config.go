package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     int    `json:"port"`
	AutoOpen bool   `json:"auto_open"`
	Password string `json:"password"`
}

func LoadConfig() (*Config, error) {
	c := &Config{Port: 8000, AutoOpen: true}
	b, err := os.ReadFile("config.json")
	if err == nil {
		_ = json.Unmarshal(b, c)
	}
	return c, nil
}

func SaveConfig(c *Config) error {
	b, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile("config.json", b, 0644)
}
