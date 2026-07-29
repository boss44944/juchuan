package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type MessageStatusRequest struct {
	MessageID string `json:"message_id"`
	DeviceID  string `json:"device_id"`
	Status    string `json:"status"`
}

func (s *Server) MessageStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MessageStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	req.MessageID = strings.TrimSpace(req.MessageID)
	req.DeviceID = strings.TrimSpace(req.DeviceID)
	req.Status = strings.TrimSpace(strings.ToUpper(req.Status))
	if req.MessageID == "" || req.DeviceID == "" || !isValidMessageStatus(req.Status) {
		WriteError(w, http.StatusBadRequest, "INVALID_STATUS", nil)
		return
	}

	_, err := s.db.Exec(`UPDATE message_targets SET status=?, updated_at=CURRENT_TIMESTAMP WHERE message_id=? AND device_id=?`, req.Status, req.MessageID, req.DeviceID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "STATUS_UPDATE_FAILED", nil)
		return
	}

	senderID, err := s.messageSenderID(req.MessageID)
	if err == nil {
		s.broadcastMessageStatus(senderID, req.DeviceID, map[string]string{
			"message_id": req.MessageID,
			"device_id":  req.DeviceID,
			"status":     req.Status,
		})
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true})
}
