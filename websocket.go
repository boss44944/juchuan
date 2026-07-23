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
 clients map[string]*websocket.Conn
 mu sync.Mutex
}

var upgrader = websocket.Upgrader{
 CheckOrigin: func(r *http.Request) bool { return true },
}

func NewHub()*Hub{
 return &Hub{clients:make(map[string]*websocket.Conn)}
}

func (h *Hub) Broadcast(v WSMessage){
 h.mu.Lock()
 defer h.mu.Unlock()
 b,_:=json.Marshal(v)
 for _,c:=range h.clients {
  _=c.WriteMessage(websocket.TextMessage,b)
 }
}

func (h *Hub) SendTo(deviceID string,v WSMessage) error {
 h.mu.Lock()
 defer h.mu.Unlock()
 c,ok:=h.clients[deviceID]
 if !ok { return nil }
 b,_:=json.Marshal(v)
 return c.WriteMessage(websocket.TextMessage,b)
}

func (h *Hub) Handler(w http.ResponseWriter,r *http.Request){
 deviceID:=r.URL.Query().Get("device")
 if deviceID=="" {
  http.Error(w,"missing device",400)
  return
 }

 conn,err:=upgrader.Upgrade(w,r)
 if err!=nil{return}

 h.mu.Lock()
 h.clients[deviceID]=conn
 h.mu.Unlock()

 defer func(){
  h.mu.Lock()
  delete(h.clients,deviceID)
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
