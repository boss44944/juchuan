package main

import (
    "encoding/json"
    "net/http"
    "os"
    "path/filepath"
    "github.com/google/uuid"
)

type FileUploadResponse struct {
    FileID string `json:"file_id"`
}

func (s *Server) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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
    dir := "uploads"
    _ = os.MkdirAll(dir, 0755)
    path := filepath.Join(dir, id+filepath.Ext(header.Filename))

    out, err := os.Create(path)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, "FILE_SAVE_FAILED", nil)
        return
    }
    defer out.Close()

    buf := make([]byte, 32*1024)
    for {
        n, readErr := file.Read(buf)
        if n > 0 {
            _, _ = out.Write(buf[:n])
        }
        if readErr != nil {
            break
        }
    }

    _, err = s.db.Exec(`INSERT INTO files(id,filename,path,mime,size) VALUES(?,?,?,?,?)`, id, header.Filename, path, header.Header.Get("Content-Type"), header.Size)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, "FILE_RECORD_FAILED", nil)
        return
    }

    WriteJSON(w, http.StatusOK, APIResponse{Success:true, Data:FileUploadResponse{FileID:id}})
}

func (s *Server) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
    id := filepath.Base(r.URL.Path)
    var path string
    err := s.db.QueryRow(`SELECT path FROM files WHERE id=?`, id).Scan(&path)
    if err != nil {
        WriteError(w, http.StatusNotFound, "FILE_NOT_FOUND", nil)
        return
    }
    http.ServeFile(w, r, path)
}

var _ = json.NewEncoder
