package main

import (
 "embed"
 "net/http"
)

type Server struct {
 addr string
 static embed.FS
 hub *Hub
}

func NewServer() (*Server, error) {
 return &Server{
  addr: ":8000",
  hub: NewHub(),
 }, nil
}

func (s *Server) Address() string {
 return s.addr
}

func (s *Server) Start() error {
 mux := http.NewServeMux()
 s.registerRoutes(mux)

 mux.HandleFunc("/ws", s.hub.Handler)

 return http.ListenAndServe(s.addr, mux)
}
