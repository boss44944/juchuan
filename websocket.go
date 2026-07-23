package main

import (
 "encoding/json"
 "net/http"
 "sync"
 "time"

 "github.com/gorilla/websocket"
)

type WSMessage struct {
 Type string `json:"type"`
 Content string `json:"content,omitempty"`
 URL string `json:"url,omitempty"`
 Filename string `json:"filename,omitempty"`
}

type Hub struct {
 clients map[*websocket.Conn]bool
 mu sync.Mutex
}

var upgrader = websocket.Upgrader{
 CheckOrigin: func(r *http.Request) bool { return true },
}

func NewHub()*Hub{
 return &Hub{clients:make(map[*websocket.Conn]bool)}
}

func (h *Hub) Broadcast(v WSMessage){
 b,_:=json.Marshal(v)
 h.mu.Lock()
 defer h.mu.Unlock()
 for c:=range h.clients {
  _=c.WriteMessage(websocket.TextMessage,b)
 }
}

func (h *Hub) Handler(w http.ResponseWriter,r *http.Request){
 conn,err:=upgrader.Upgrade(w,r)
 if err!=nil{return}

 h.mu.Lock()
 h.clients[conn]=true
 h.mu.Unlock()

 defer func(){
  h.mu.Lock()
  delete(h.clients,conn)
  h.mu.Unlock()
  conn.Close()
 }()

 for {
  conn.SetReadDeadline(time.Now().Add(60*time.Second))
  if _,_,err:=conn.ReadMessage();err!=nil{
   return
  }
 }
}
