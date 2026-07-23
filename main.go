package main

import (
 "log"
 "strings"
)

func main() {
 s, err := NewServer()
 if err != nil {
  log.Fatal(err)
 }

 port := strings.TrimPrefix(s.Address(), ":")
 url := LocalURL(port)

 log.Println("==============================")
 log.Println("Juchuan 菊传")
 log.Println("")
 log.Println("访问地址:")
 log.Println(url)
 log.Println("==============================")

 OpenBrowser(url)

 quit := make(chan struct{}, 1)
 go StartTray(url, quit)

 serverErr := make(chan error, 1)
 go func() {
  serverErr <- s.Start()
 }()

 select {
 case <-quit:
  if err := s.Shutdown(); err != nil {
   log.Println("shutdown error:", err)
  }
 case err := <-serverErr:
  if err != nil {
   log.Fatal(err)
  }
 }
}
