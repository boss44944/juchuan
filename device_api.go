package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type DeviceRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) DeviceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if req.ID == "" {
		http.Error(w, "missing id", 400)
		return
	}

	s.devices.Add(&Device{
		ID:       req.ID,
		Name:     req.Name,
		LastSeen: time.Now().Unix(),
	})

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.devices.List())
}
