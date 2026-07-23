package main

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
 mux.HandleFunc("/api/health", func(w http.ResponseWriter,r *http.Request){
  w.Write([]byte("ok"))
 })

 mux.HandleFunc("/api/history", s.HistoryHandler)
 mux.HandleFunc("/upload", s.UploadHandler)
 mux.HandleFunc("/download/", s.DownloadHandler)
}
