package main

import (
 "encoding/json"
 "net/http"
)

type TextRequest struct {
 Content string `json:"content"`
}

func (s *Server) TextHandler(w http.ResponseWriter,r *http.Request){
 var req TextRequest
 if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{
  http.Error(w,"bad request",400)
  return
 }

 if req.Content==""{
  http.Error(w,"empty content",400)
  return
 }

 if err:=ReceiveText(s.db,s.clipboard,s.hub,req.Content);err!=nil{
  http.Error(w,"internal error",500)
  return
 }

 w.Header().Set("Content-Type","application/json")
 json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}
