package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
)

type SendTextRequest struct {
	Content string `json:"content"`
	SenderDeviceID string `json:"sender_device_id"`
	Targets []string `json:"targets"`
}

func (s *Server) SendTextMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req SendTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}
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
			DeviceID: id,
			Status: "CREATED",
		})
		_ = s.hub.SendTo(id, WSMessage{
			Type: "MESSAGE_RECEIVED",
			Data: message,
		})
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: message})
}

type FileMessageRequest struct {
	SenderDeviceID string `json:"sender_device_id"`
	FileID string `json:"file_id"`
	Targets []string `json:"targets"`
}

func (s *Server) SendFileMessageHandler(w http.ResponseWriter, r *http.Request) {
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
		ID: uuid.NewString(),
		Type: MessageFile,
		FileID: req.FileID,
		SenderDeviceID: req.SenderDeviceID,
	}

	if err := SaveMessage(s.db, message); err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_SAVE_FAILED", nil)
		return
	}

	for _, id := range req.Targets {
		_ = SaveMessageTarget(s.db, &MessageTarget{MessageID: message.ID, DeviceID: id, Status: "CREATED"})
		_ = s.hub.SendTo(id, WSMessage{Type: "MESSAGE_RECEIVED", Data: message})
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success:true, Data:message})
}
