package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	addr       string
	static     embed.FS
	hub        *Hub
	storage    *Storage
	db         *sql.DB
	clipboard  *Clipboard
	devices    *DeviceManager
	config     *Config
	configMu   sync.RWMutex
	sessions   map[string]time.Time
	sessionsMu sync.RWMutex
	httpServer *http.Server
}

func NewServer() (*Server, error) {
	storage, err := NewStorage()
	if err != nil {
		return nil, err
	}

	db, err := InitDatabase(storage.DB)
	if err != nil {
		return nil, err
	}

	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	cfg.Password = "123456"
	if err := SaveConfig(cfg); err != nil {
		return nil, err
	}

	devices := NewDeviceManager()
	if err := LoadDevices(db, devices); err != nil {
		return nil, err
	}

	startPort := cfg.Port
	if startPort <= 0 {
		startPort = 8000
	}
	port := FindAvailablePort(startPort)

	return &Server{
		addr:      fmt.Sprintf(":%d", port),
		static:    StaticFiles,
		hub:       NewHub(),
		storage:   storage,
		db:        db,
		clipboard: &Clipboard{},
		devices:   devices,
		config:    cfg,
		sessions:  make(map[string]time.Time),
	}, nil
}

func (s *Server) Address() string {
	return s.addr
}

func (s *Server) ShouldAutoOpen() bool {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	return s.config != nil && s.config.AutoOpen
}

func (s *Server) HasPassword() bool {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	if s.config == nil {
		return false
	}
	return strings.TrimSpace(s.config.Password) != ""
}

func (s *Server) CurrentPassword() string {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	if s.config == nil {
		return ""
	}
	return s.config.Password
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	staticFS, _ := fs.Sub(s.static, "static")
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFileFS(w, r, staticFS, "login.html")
	})
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/app" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFileFS(w, r, staticFS, "app.html")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if s.HasPassword() && !s.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/app", http.StatusFound)
	})
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))
	mux.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		staticHandler.ServeHTTP(w, r)
	}))
	mux.HandleFunc("/ws", s.requireAuth(s.hub.Handler))
	mux.HandleFunc("/api/text", s.requireAuth(s.TextHandler))

	s.httpServer = &http.Server{Addr: s.addr, Handler: mux}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	if s.httpServer == nil {
		return nil
	}
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
