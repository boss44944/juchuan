package main

import (
 "fmt"
 "net/http"
)

type Server struct { cfg Config }

func NewServer(cfg Config)*Server{return &Server{cfg:cfg}}

func(s *Server)Start()error{
 http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","text/html;charset=utf-8")
  fmt.Fprint(w,"<h1>菊传 Juchuan</h1><p>LAN Transfer Server</p>")
 })
 return http.ListenAndServe(fmt.Sprintf(":%d",s.cfg.Port),nil)
}
