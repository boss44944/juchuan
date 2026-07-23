package main

import (
 "net"
 "strings"
)

func LocalIP() string {
 addrs,_:=net.InterfaceAddrs()
 for _,a:=range addrs {
  ip:=a.String()
  if strings.HasPrefix(ip,"192.168.")||strings.HasPrefix(ip,"10.") { return strings.Split(ip,"/")[0] }
 }
 return "127.0.0.1"
}
