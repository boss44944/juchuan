package main

import (
 "fmt"
 "net/http"
 "os"
 "path/filepath"
 "time"
)

func UploadHandler(storage *Storage) http.HandlerFunc {
 return func(w http.ResponseWriter,r *http.Request){
  if err:=r.ParseMultipartForm(64<<20);err!=nil {http.Error(w,err.Error(),400);return}
  file,header,err:=r.FormFile("file")
  if err!=nil {http.Error(w,err.Error(),400);return}
  defer file.Close()

  name:=fmt.Sprintf("%d_%s",time.Now().Unix(),header.Filename)
  path:=filepath.Join(storage.UploadDir,name)
  dst,err:=os.Create(path)
  if err!=nil {http.Error(w,err.Error(),500);return}
  defer dst.Close()
  _,_=dst.ReadFrom(file)
  w.Write([]byte(name))
 }
}
