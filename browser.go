package main

import (
 "fmt"
 "os/exec"
 "runtime"
)

func OpenBrowser(url string) {
 var cmd *exec.Cmd
 switch runtime.GOOS {
 case "darwin":
  cmd=exec.Command("open",url)
 case "windows":
  cmd=exec.Command("cmd","/c","start",url)
 default:
  cmd=exec.Command("xdg-open",url)
 }
 if cmd!=nil { _=cmd.Start() }
}

func LocalURL(port string) string {
 return fmt.Sprintf("http://%s:%s",GetLocalIP(),port)
}
