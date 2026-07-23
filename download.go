package main

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/download/")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var path string
	var filename string
	err = s.db.QueryRow("SELECT filepath, filename FROM history WHERE id=?", id).Scan(&path, &filename)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	if filename == "" {
		filename = filepath.Base(path)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	http.ServeFile(w, r, path)
}
