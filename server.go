package main

import (
 "database/sql"
 "embed"
 "net/http"
)

type Server struct {
 addr string
 static embed.FS
 hub *Hub
 storage *Storage
 db *sql.DB
 clipboard *Clipboard
}

func NewServer() (*Server, error) {
 storage, err := NewStorage()
 if err != nil { return nil, err }

 db, err := InitDatabase(storage.DB)
 if err != nil { return nil, err }

 return &Server{
  addr: ":8000",
  hub: NewHub(),
  storage: storage,
  db: db,
  clipboard: &Clipboard{},
 }, nil
}

func (s *Server) Address() string {
 return s.addr
}

func (s *Server) Start() error {
 mux := http.NewServeMux()
 s.registerRoutes(mux)

 mux.HandleFunc("/ws", s.hub.Handler)
 mux.HandleFunc("/api/text", s.TextHandler)

 return http.ListenAndServe(s.addr, mux)
}
