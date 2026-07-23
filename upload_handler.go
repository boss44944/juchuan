package main

import (
 "net/http"
)

func (s *Server) UploadHandler(w http.ResponseWriter,r *http.Request){
 if err:=r.ParseMultipartForm(1024*1024*1024);err!=nil{
  http.Error(w,"upload failed",400)
  return
 }

 files:=r.MultipartForm.File["file"]
 if len(files)==0{
  http.Error(w,"no file",400)
  return
 }

 w.Write([]byte("ok"))
}
