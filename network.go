package main

import (
 "net"
 "strings"
)

func LocalIP() string {
 addrs,err:=net.InterfaceAddrs()
 if err!=nil { return "127.0.0.1" }

 for _,a:=range addrs {
  ipnet,ok:=a.(*net.IPNet)
  if !ok || ipnet.IP.IsLoopback() { continue }

  ip:=ipnet.IP.String()
  if strings.Contains(ip,":") { continue }

  if strings.HasPrefix(ip,"192.168.") ||
   strings.HasPrefix(ip,"10.") ||
   strings.HasPrefix(ip,"172.") {
   return ip
  }
 }
 return "127.0.0.1"
}
