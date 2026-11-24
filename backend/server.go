package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

type Server struct {
	configMu sync.RWMutex
	config   *Config

	sessionsMu sync.RWMutex
	sessions   map[string]time.Time

	storage   *Storage
	db        *sql.DB
	clipboard *Clipboard
	devices   *DeviceManager
	events    *EventBus
	hub       *Hub

	staticFS embed.FS

	httpServer *http.Server
	listenAddr string
}

func NewServer() (*Server, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	cfg.Port = FindAvailablePort(cfg.Port)

	storage, err := NewStorage()
	if err != nil {
		return nil, err
	}

	db, err := InitDatabase(storage.DB)
	if err != nil {
		return nil, err
	}

	devices := NewDeviceManager()
	if err := LoadDevices(db, devices); err != nil {
		return nil, err
	}

	s := &Server{
		config:     cfg,
		sessions:   make(map[string]time.Time),
		storage:    storage,
		db:         db,
		clipboard:  &Clipboard{},
		devices:    devices,
		events:     NewEventBus(),
		hub:        NewHub(),
		staticFS:   StaticFiles,
		listenAddr: fmt.Sprintf(":%d", cfg.Port),
	}

	// Bridge internal device events to WebSocket clients.
	s.events.Subscribe(func(payload []byte) {
		var evt DeviceEvent
		if err := json.Unmarshal(payload, &evt); err != nil {
			return
		}
		s.hub.BroadcastEvent(evt)
	})

	StartDeviceMonitor(s.devices, func(evt DeviceEvent) {
		if s.db != nil {
			d := evt.Data
			_ = SaveDevice(s.db, &d)
		}
		if s.events != nil {
			s.events.Publish(evt)
		}
	})

	return s, nil
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	s.registerRoutes(mux)
	mux.HandleFunc("/ws", s.requireAuth(s.hub.Handler))
	s.registerWebRoutes(mux)

	s.httpServer = &http.Server{
		Addr:    s.listenAddr,
		Handler: mux,
	}

	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}
	}

	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

func (s *Server) Address() string {
	return s.listenAddr
}

func (s *Server) CurrentPassword() string {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	return strings.TrimSpace(s.config.Password)
}

func (s *Server) HasPassword() bool {
	return s.CurrentPassword() != ""
}

func (s *Server) ShouldAutoOpen() bool {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	return s.config.AutoOpen
}

// The standalone build sets this at link time to serve exactly the compiled bundle,
// even when it is launched next to a different development checkout.
var frontendEmbeddedOnly = "false"

func (s *Server) registerWebRoutes(mux *http.ServeMux) {
	staticFS, err := fs.Sub(s.staticFS, "static")
	if err != nil {
		log.Printf("static fs unavailable: %v", err)
		return
	}

	staticHandler := http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)
		http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
	}))
	mux.Handle("/static/", staticHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws") {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/app" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		frontend := frontendFiles(staticFS)
		name := strings.TrimPrefix(r.URL.Path, "/")
		switch strings.TrimSuffix(r.URL.Path, "/") {
		case "", "/login", "/devices", "/messages", "/send", "/config":
			name = "index.html"
		}
		if info, err := fs.Stat(frontend, name); err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}
		setNoCacheHeaders(w)
		http.ServeFileFS(w, r, frontend, name)
	})
}

func frontendFiles(embedded fs.FS) fs.FS {
	if frontendEmbeddedOnly != "true" {
		dist := os.DirFS(filepath.Join("..", "frontend", "dist"))
		if info, err := fs.Stat(dist, "index.html"); err == nil && !info.IsDir() {
			return dist
		}
	}
	return embedded
}
