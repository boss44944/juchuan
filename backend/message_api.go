package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type SendTextRequest struct {
	Content        string   `json:"content"`
	SenderDeviceID string   `json:"sender_device_id"`
	Targets        []string `json:"targets"`
}

func (s *Server) SendTextMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" || len(req.Targets) == 0 {
		WriteError(w, http.StatusBadRequest, "MESSAGE_REQUIRED", nil)
		return
	}

	message := NewTextMessage(uuid.NewString(), req.Content, req.SenderDeviceID)
	if err := SaveMessage(s.db, message); err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_SAVE_FAILED", nil)
		return
	}

	for _, id := range req.Targets {
		_ = SaveMessageTarget(s.db, &MessageTarget{
			MessageID: message.ID,
			DeviceID:  id,
			Status:    "CREATED",
		})
		err := s.hub.SendTo(id, WSMessage{
			Type: "MESSAGE_RECEIVED",
			Data: message,
		})
		if err == nil {
			_ = s.updateMessageTargetStatus(message.ID, id, "DELIVERED")
			statusPayload := map[string]string{
				"message_id": message.ID,
				"device_id":  id,
				"status":     "DELIVERED",
			}
			s.broadcastMessageStatus(message.SenderDeviceID, id, statusPayload)
		}
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: message})
}

type FileMessageRequest struct {
	SenderDeviceID string   `json:"sender_device_id"`
	FileID         string   `json:"file_id"`
	Targets        []string `json:"targets"`
}

func (s *Server) SendFileMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FileMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}
	if req.FileID == "" || len(req.Targets) == 0 {
		WriteError(w, http.StatusBadRequest, "MESSAGE_REQUIRED", nil)
		return
	}

	message := &Message{
		ID:             uuid.NewString(),
		Type:           MessageFile,
		FileID:         req.FileID,
		SenderDeviceID: req.SenderDeviceID,
		CreatedAt:      time.Now(),
	}

	if err := SaveMessage(s.db, message); err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_SAVE_FAILED", nil)
		return
	}

	for _, id := range req.Targets {
		_ = SaveMessageTarget(s.db, &MessageTarget{MessageID: message.ID, DeviceID: id, Status: "CREATED"})
		err := s.hub.SendTo(id, WSMessage{Type: "MESSAGE_RECEIVED", Data: message})
		if err == nil {
			_ = s.updateMessageTargetStatus(message.ID, id, "DELIVERED")
			statusPayload := map[string]string{
				"message_id": message.ID,
				"device_id":  id,
				"status":     "DELIVERED",
			}
			s.broadcastMessageStatus(message.SenderDeviceID, id, statusPayload)
		}
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: message})
}

func (s *Server) updateMessageTargetStatus(messageID string, deviceID string, status string) error {
	_, err := s.db.Exec(`UPDATE message_targets SET status=?, updated_at=CURRENT_TIMESTAMP WHERE message_id=? AND device_id=?`, status, messageID, deviceID)
	return err
}

func (s *Server) broadcastMessageStatus(senderID string, targetID string, payload map[string]string) {
	msg := WSMessage{Type: "MESSAGE_STATUS_UPDATED", Data: payload}
	if strings.TrimSpace(senderID) != "" {
		_ = s.hub.SendTo(senderID, msg)
	}
	if targetID != senderID {
		_ = s.hub.SendTo(targetID, msg)
	}
}

func isValidMessageStatus(status string) bool {
	switch status {
	case "CREATED", "DELIVERED", "READ":
		return true
	default:
		return false
	}
}

func (s *Server) messageSenderID(messageID string) (string, error) {
	var senderID string
	err := s.db.QueryRow(`SELECT sender_device_id FROM messages WHERE id=?`, messageID).Scan(&senderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return senderID, nil
}
