package main

import (
 "os"
 "path/filepath"
 "runtime"
)

func AppDataDir() (string,error) {
 switch runtime.GOOS {
 case "darwin":
  home,err:=os.UserHomeDir()
  if err!=nil{return "",err}
  return filepath.Join(home,"Library","Application Support","Juchuan"),nil
 case "windows":
  base:=os.Getenv("APPDATA")
  if base=="" { base=os.TempDir() }
  return filepath.Join(base,"Juchuan"),nil
 default:
  home,err:=os.UserHomeDir()
  if err!=nil{return "",err}
  return filepath.Join(home,".juchuan"),nil
 }
}
