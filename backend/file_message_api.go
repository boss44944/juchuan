package main

import (
	"net/http"
	"os"
	"path/filepath"
)

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
