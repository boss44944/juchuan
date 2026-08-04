package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SendTextRequest struct {
	Content        string   `json:"content"`
	SenderDeviceID string   `json:"sender_device_id"`
	Targets        []string `json:"targets"`
}

// resolveMessageTargets 把请求中的目标设备解析为最终的设备 ID 列表。
// 特殊值 "server" 会展开为所有 role=server 的设备（即运行菊传的电脑），
// 供手机（客户端）直接发送给电脑服务端使用。
func (s *Server) resolveMessageTargets(reqTargets []string) []string {
	out := make([]string, 0, len(reqTargets))
	seen := make(map[string]bool)
	add := func(id string) {
		if id != "" && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	for _, id := range reqTargets {
		id = strings.TrimSpace(id)
		if id == "server" {
			for _, d := range s.devices.List() {
				if d.Role == "server" {
					add(d.ID)
				}
			}
			continue
		}
		add(id)
	}
	return out
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
	if req.Content == "" {
		WriteError(w, http.StatusBadRequest, "MESSAGE_REQUIRED", nil)
		return
	}

	targets := s.resolveMessageTargets(req.Targets)
	if len(targets) == 0 {
		WriteError(w, http.StatusBadRequest, "MESSAGE_TARGET_REQUIRED", nil)
		return
	}

	message := NewTextMessage(uuid.NewString(), req.Content, req.SenderDeviceID)
	if err := SaveMessage(s.db, message); err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_SAVE_FAILED", nil)
		return
	}

	for _, id := range targets {
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
	if req.FileID == "" {
		WriteError(w, http.StatusBadRequest, "MESSAGE_REQUIRED", nil)
		return
	}

	targets := s.resolveMessageTargets(req.Targets)
	if len(targets) == 0 {
		WriteError(w, http.StatusBadRequest, "MESSAGE_TARGET_REQUIRED", nil)
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

	for _, id := range targets {
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

type DeleteMessageRequest struct {
	MessageID string `json:"message_id"`
	DeviceID  string `json:"device_id"`
}

func (s *Server) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	req.MessageID = strings.TrimSpace(req.MessageID)
	req.DeviceID = strings.TrimSpace(req.DeviceID)
	if req.MessageID == "" || req.DeviceID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	_, err := s.db.Exec(`DELETE FROM message_targets WHERE message_id=? AND device_id=?`, req.MessageID, req.DeviceID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_DELETE_FAILED", nil)
		return
	}

	var count int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM message_targets WHERE message_id=?`, req.MessageID).Scan(&count)
	if count == 0 {
		_, _ = s.db.Exec(`DELETE FROM messages WHERE id=?`, req.MessageID)
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true})
}

type ClearMessagesRequest struct {
	DeviceID string `json:"device_id"`
}

func (s *Server) ClearMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ClearMessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	req.DeviceID = strings.TrimSpace(req.DeviceID)
	if req.DeviceID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}

	_, err := s.db.Exec(`DELETE FROM message_targets WHERE device_id=?`, req.DeviceID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_CLEAR_FAILED", nil)
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true})
}
