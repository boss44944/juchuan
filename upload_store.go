package main

import (
 "database/sql"
 "io"
 "mime/multipart"
 "os"
 "path/filepath"
)

func SaveUploadedFile(db *sql.DB, storage *Storage, header *multipart.FileHeader) (string,error){
 src,err:=header.Open()
 if err!=nil{return "",err}
 defer src.Close()

 name:=SafeFileName(header.Filename)
 path:=filepath.Join(storage.FileDir,name)

 dst,err:=os.Create(path)
 if err!=nil{return "",err}
 defer dst.Close()

 if _,err=io.Copy(dst,src);err!=nil{return "",err}

 if err:=SaveFileHistory(db,path,header.Filename,"",header.Size,"phone_to_pc");err!=nil{return "",err}

 return path,nil
}
