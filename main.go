package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := loadConfig()
	if err := saveConfig(cfg); err != nil {
		log.Printf("save config: %v", err)
	}

	server := NewServer(cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)

	log.Printf("Juchuan 菊传 started on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
