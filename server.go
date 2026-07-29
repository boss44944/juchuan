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
	events     *EventBus
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

	hub := NewHub()
	events := NewEventBus()
	events.Subscribe(func(data []byte) {
		var event DeviceEvent
		if err := json.Unmarshal(data, &event); err == nil {
			hub.BroadcastEvent(event)
		}
	})

	server := &Server{
		addr:      fmt.Sprintf(":%d", FindAvailablePort(cfg.Port)),
		static:    StaticFiles,
		hub:       hub,
		events:    events,
		storage:   storage,
		db:        db,
		clipboard: &Clipboard{},
		devices:   devices,
		config:    cfg,
		sessions:  make(map[string]time.Time),
	}

	StartDeviceMonitor(devices, func(event DeviceEvent) {
		events.Publish(event)
	})

	return server, nil
}

func (s *Server) Address() string { return s.addr }

func (s *Server) ShouldAutoOpen() bool {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	return s.config != nil && s.config.AutoOpen
}

func (s *Server) HasPassword() bool {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	return s.config != nil && strings.TrimSpace(s.config.Password) != ""
}

func (s *Server) CurrentPassword() string {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	if s.config == nil { return "" }
	return s.config.Password
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	staticFS, _ := fs.Sub(s.static, "static")
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "login.html")
	})
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "app.html")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if s.HasPassword() && !s.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/app", http.StatusFound)
	})
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/ws", s.requireAuth(s.hub.Handler))
	mux.HandleFunc("/api/text", s.requireAuth(s.TextHandler))

	s.httpServer = &http.Server{Addr: s.addr, Handler: mux}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	if s.httpServer == nil { return nil }
	if err := s.httpServer.Shutdown(context.Background()); err != nil { return err }
	if s.db != nil { return s.db.Close() }
	return nil
}
