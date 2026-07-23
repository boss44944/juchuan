package main

import (
 "context"
 "database/sql"
 "embed"
 "fmt"
 "io/fs"
 "net/http"
)

type Server struct {
 addr string
 static embed.FS
 hub *Hub
 storage *Storage
 db *sql.DB
 clipboard *Clipboard
 devices *DeviceManager
 httpServer *http.Server
}

func NewServer() (*Server, error) {
 storage, err := NewStorage()
 if err != nil { return nil, err }

 db, err := InitDatabase(storage.DB)
 if err != nil { return nil, err }

 port:=FindAvailablePort(8000)

 return &Server{
  addr: fmt.Sprintf(":%d",port),
  static: StaticFiles,
  hub: NewHub(),
  storage: storage,
  db: db,
  clipboard: &Clipboard{},
  devices: NewDeviceManager(),
 }, nil
}

func (s *Server) Address() string {
 return s.addr
}

func (s *Server) Start() error {
 mux := http.NewServeMux()
 s.registerRoutes(mux)

 staticFS,_:=fs.Sub(s.static,"static")
 mux.Handle("/static/",http.StripPrefix("/static/",http.FileServer(http.FS(staticFS))))
 mux.HandleFunc("/ws", s.hub.Handler)
 mux.HandleFunc("/api/text", s.TextHandler)

 s.httpServer=&http.Server{Addr:s.addr,Handler:mux}
 return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
 if s.httpServer==nil { return nil }
 if err:=s.httpServer.Shutdown(context.Background());err!=nil{return err}
 if s.db!=nil{return s.db.Close()}
 return nil
}
