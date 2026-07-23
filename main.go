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

 if err := s.Start(); err != nil {
  log.Fatal(err)
 }
}
