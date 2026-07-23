package main

import (
 "database/sql"
 _ "modernc.org/sqlite"
)

func InitDatabase(path string)(*sql.DB,error){
 db,err:=sql.Open("sqlite",path)
 if err!=nil{return nil,err}
 _,err=db.Exec(`CREATE TABLE IF NOT EXISTS history(
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 type TEXT,
 direction TEXT,
 content TEXT,
 filename TEXT,
 filepath TEXT,
 mime TEXT,
 size INTEGER,
 created_at DATETIME DEFAULT CURRENT_TIMESTAMP
 );`)
 return db,err
}
