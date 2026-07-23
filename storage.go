package main

import (
 "os"
 "path/filepath"
)

type Storage struct {
 Root string
 UploadDir string
 DB string
}

func NewStorage() (*Storage,error){
 root,_:=os.UserConfigDir()
 root=filepath.Join(root,"Juchuan")
 if err:=os.MkdirAll(filepath.Join(root,"uploads"),0755);err!=nil{return nil,err}
 return &Storage{Root:root,UploadDir:filepath.Join(root,"uploads"),DB:filepath.Join(root,"juchuan.db")},nil
}
