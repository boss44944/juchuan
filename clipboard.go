package main

import (
 "fmt"
 "os/exec"
 "runtime"
)

func SetClipboard(text string) error {
 var cmd *exec.Cmd
 switch runtime.GOOS {
 case "darwin":
  cmd=exec.Command("pbcopy")
 case "windows":
  cmd=exec.Command("clip")
 default:
  cmd=exec.Command("sh","-c","xclip -selection clipboard")
 }
 in,err:=cmd.StdinPipe()
 if err!=nil{return err}
 if err:=cmd.Start();err!=nil{return err}
 _,_=in.Write([]byte(text))
 _=in.Close()
 if err:=cmd.Wait();err!=nil{return fmt.Errorf("clipboard: %w",err)}
 return nil
}
