package main

import (
 "fmt"
 "net/http"
)

type Server struct { addr string }

func NewServer() (*Server,error) { return &Server{addr: ":8000"}, nil }
func (s *Server) Address() string { return s.addr }
func (s *Server) Start() error { return http.ListenAndServe(s.addr, http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){fmt.Fprintln(w,"Juchuan") })) }
