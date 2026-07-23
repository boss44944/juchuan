package main

import "log"

func main() {
 s, err := NewServer()
 if err != nil { log.Fatal(err) }
 log.Printf("Juchuan started: %s", s.Address())
 if err := s.Start(); err != nil { log.Fatal(err) }
}
