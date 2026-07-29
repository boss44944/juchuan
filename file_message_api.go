package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"github.com/google/uuid"
)

type SendFileMessageRequest struct {
	FileID string `json:"file_id"`
	SenderDeviceID string `json:"sender_device_id"`
	Targets []string `json:"targets"`
}

func (s *Server) SendFileMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req SendFileMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", nil)
		return
	}
	if req.FileID == "" || len(req.Targets) == 0 {
		WriteError(w, http.StatusBadRequest, "FILE_MESSAGE_REQUIRED", nil)
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

	for _, target := range req.Targets {
		SaveMessageTarget(s.db, &MessageTarget{
			MessageID: message.ID,
			DeviceID: target,
			Status: "CREATED",
		})
		s.hub.SendTo(target, WSMessage{Type:"MESSAGE_RECEIVED", Data:message})
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success:true, Data:message})
}

func (s *Server) DownloadMessageFileHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	if id == "" {
		http.NotFound(w, r)
		return
	}
	_, err := os.Stat(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, id)
}
