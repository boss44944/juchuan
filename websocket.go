package main

import (
 "encoding/json"
 "net/http"
 "sync"
)

type WSMessage struct {
 Type string `json:"type"`
 Content string `json:"content,omitempty"`
 URL string `json:"url,omitempty"`
}

type Hub struct {
 clients map[chan []byte]bool
 mu sync.Mutex
}

func NewHub()*Hub{
 return &Hub{clients:make(map[chan []byte]bool)}
}

func (h *Hub) Broadcast(v WSMessage){
 b,_:=json.Marshal(v)
 h.mu.Lock()
 defer h.mu.Unlock()
 for c:=range h.clients { c<-b }
}

func (h *Hub) Handler(w http.ResponseWriter,r *http.Request){
 ch:=make(chan []byte,8)
 h.mu.Lock(); h.clients[ch]=true; h.mu.Unlock()
 defer func(){h.mu.Lock();delete(h.clients,ch);h.mu.Unlock()}()
 for {
  select {
  case b:=<-ch:
   w.Write(b)
   if f,ok:=w.(http.Flusher);ok {f.Flush()}
  case <-r.Context().Done(): return
  }
 }
}
