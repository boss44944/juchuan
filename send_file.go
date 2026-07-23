package main

import (
	"encoding/json"
	"net/http"
)

type SendFileRequest struct {
	Device    string `json:"device"`
	HistoryID int64  `json:"history_id"`
}

func (s *Server) SendFileHandler(w http.ResponseWriter, r *http.Request) {
	var req SendFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	var filename string
	err := s.db.QueryRow("SELECT filename FROM history WHERE id=?", req.HistoryID).Scan(&filename)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	err = s.hub.SendTo(req.Device, WSMessage{
		Type:     "file",
		URL:      "/download/" + json.Number(req.HistoryID).String(),
		Filename: filename,
	})
	if err != nil {
		http.Error(w, "send failed", 500)
		return
	}

	w.Write([]byte("ok"))
}
