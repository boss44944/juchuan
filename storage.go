package main

import (
 "os"
 "path/filepath"
)

type Storage struct {
 Root string
 UploadDir string
 ImageDir string
 FileDir string
 DB string
}

func NewStorage() (*Storage,error){
 root,err:=os.Getwd()
 if err!=nil{return nil,err}

 upload:=filepath.Join(root,"uploads")
 s:=&Storage{
  Root:root,
  UploadDir:upload,
  ImageDir:filepath.Join(upload,"images"),
  FileDir:filepath.Join(upload,"files"),
  DB:filepath.Join(root,"data","juchuan.db"),
 }

 for _,dir:=range []string{s.ImageDir,s.FileDir,filepath.Dir(s.DB)}{
  if err:=os.MkdirAll(dir,0755);err!=nil{return nil,err}
 }
 return s,nil
}
