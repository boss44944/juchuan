package main

import (
	"net/http"
	"strconv"
)

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	if err := r.ParseMultipartForm(1024 * 1024 * 1024); err != nil {
		http.Error(w, "upload failed", 400)
		return
	}

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "no file", 400)
		return
	}

	path, err := SaveUploadedFile(s.db, s.storage, files[0])
	if err != nil {
		http.Error(w, "save failed", 500)
		return
	}

	var historyID int64
	if err := s.db.QueryRow("SELECT id FROM history WHERE filepath=? ORDER BY id DESC LIMIT 1", path).Scan(&historyID); err != nil {
		http.Error(w, "history lookup failed", 500)
		return
	}

	s.hub.Broadcast(WSMessage{
		Type:     "file",
		URL:      "/download/" + strconv.FormatInt(historyID, 10),
		Filename: files[0].Filename,
	})

	w.Write([]byte("ok"))
}
