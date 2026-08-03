package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileUploadResponse struct {
	FileID string `json:"file_id"`
}

func (s *Server) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_FILE", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "FILE_REQUIRED", nil)
		return
	}
	defer file.Close()

	id := uuid.NewString()
	dir := s.storage.FileDir
	_ = os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, id+filepath.Ext(header.Filename))

	out, err := os.Create(path)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "FILE_SAVE_FAILED", nil)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		WriteError(w, http.StatusInternalServerError, "FILE_SAVE_FAILED", nil)
		return
	}

	_, err = s.db.Exec(`INSERT INTO files(id,filename,path,mime,size) VALUES(?,?,?,?,?)`, id, header.Filename, path, header.Header.Get("Content-Type"), header.Size)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "FILE_RECORD_FAILED", nil)
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: FileUploadResponse{FileID: id}})
}

func (s *Server) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := filepath.Base(r.URL.Path)
	var path string
	var filename string
	err := s.db.QueryRow(`SELECT path, filename FROM files WHERE id=?`, id).Scan(&path, &filename)
	if err != nil {
		WriteError(w, http.StatusNotFound, "FILE_NOT_FOUND", nil)
		return
	}
	if filename != "" {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	}
	http.ServeFile(w, r, path)
}

var _ = json.NewEncoder
