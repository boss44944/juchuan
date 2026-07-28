package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const sessionCookieName = "juchuan_session"
const sessionTTL = 24 * time.Hour

type authStatusResponse struct {
	RequiresPassword bool `json:"requires_password"`
	Authenticated    bool `json:"authenticated"`
}

type loginRequest struct {
	Password string `json:"password"`
}

type configResponse struct {
	Port        int  `json:"port"`
	AutoOpen    bool `json:"auto_open"`
	HasPassword bool `json:"has_password"`
}

type configUpdateRequest struct {
	Port     int     `json:"port"`
	AutoOpen *bool   `json:"auto_open"`
	Password *string `json:"password"`
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.isAuthenticated(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		next(w, r)
	}
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	if !s.HasPassword() {
		return true
	}

	c, err := r.Cookie(sessionCookieName)
	if err != nil || c.Value == "" {
		return false
	}

	s.sessionsMu.RLock()
	expireAt, ok := s.sessions[c.Value]
	s.sessionsMu.RUnlock()
	if !ok {
		return false
	}
	return time.Now().Before(expireAt)
}

func (s *Server) createSession() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := hex.EncodeToString(buf)

	s.sessionsMu.Lock()
	s.sessions[token] = time.Now().Add(sessionTTL)
	s.sessionsMu.Unlock()
	return token, nil
}

func (s *Server) clearSession(token string) {
	if token == "" {
		return
	}
	s.sessionsMu.Lock()
	delete(s.sessions, token)
	s.sessionsMu.Unlock()
}

func (s *Server) AuthStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authStatusResponse{
		RequiresPassword: s.HasPassword(),
		Authenticated:    s.isAuthenticated(r),
	})
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.HasPassword() {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	s.configMu.RLock()
	want := s.config.Password
	s.configMu.RUnlock()
	if req.Password != want {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	token, err := s.createSession()
	if err != nil {
		http.Error(w, "create session failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookieName); err == nil {
		s.clearSession(c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.configMu.RLock()
		cfg := *s.config
		s.configMu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(configResponse{
			Port:        cfg.Port,
			AutoOpen:    cfg.AutoOpen,
			HasPassword: strings.TrimSpace(cfg.Password) != "",
		})
		return

	case http.MethodPost:
		if s.HasPassword() && !s.isAuthenticated(r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req configUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		s.configMu.Lock()
		if req.Port > 0 {
			s.config.Port = req.Port
		}
		if req.AutoOpen != nil {
			s.config.AutoOpen = *req.AutoOpen
		}
		if req.Password != nil {
			s.config.Password = *req.Password
		}
		cfg := *s.config
		s.configMu.Unlock()

		if err := SaveConfig(&cfg); err != nil {
			http.Error(w, "save config failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(configResponse{
			Port:        cfg.Port,
			AutoOpen:    cfg.AutoOpen,
			HasPassword: strings.TrimSpace(cfg.Password) != "",
		})
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
