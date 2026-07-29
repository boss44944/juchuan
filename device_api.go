package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type DeviceRequest struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	Platform    string `json:"platform"`
	Browser     string `json:"browser"`
}

func (s *Server) DeviceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	if req.ID == "" || req.DisplayName == "" {
		WriteError(w, http.StatusBadRequest, "DEVICE_INFO_REQUIRED", nil)
		return
	}

	if s.devices.NameExists(req.DisplayName) {
		WriteError(w, http.StatusConflict, "DEVICE_NAME_EXISTS", map[string]interface{}{
			"name": req.DisplayName,
		})
		return
	}

	if req.Role == "" {
		req.Role = "client"
	}

	device := &Device{
		ID:          req.ID,
		DisplayName: req.DisplayName,
		Role:        req.Role,
		Platform:    req.Platform,
		Browser:     req.Browser,
		LastSeen:    time.Now().Unix(),
	}

	s.devices.Add(device)
	if s.events != nil {
		s.events.Publish(DeviceEvent{
			Type: DeviceOnlineEvent,
			Data: *device,
		})
	}

	WriteJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]string{
			"device_id": req.ID,
		},
	})
}

func (s *Server) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    s.devices.List(),
	})
}

func (s *Server) DeviceHeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	device, ok := s.devices.Get(req.ID)
	if !ok {
		WriteError(w, http.StatusNotFound, "DEVICE_NOT_FOUND", nil)
		return
	}

	device.LastSeen = time.Now().Unix()

	if s.events != nil {
		s.events.Publish(DeviceEvent{
			Type: DeviceOnlineEvent,
			Data: *device,
		})
	}

	WriteJSON(w, http.StatusOK, APIResponse{
		Success: true,
	})
}
