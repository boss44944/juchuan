package main

import (
 "crypto/rand"
 "encoding/hex"
 "path/filepath"
)

func SafeFileName(name string) string {
 ext:=filepath.Ext(name)
 b:=make([]byte,8)
 _,_=rand.Read(b)
 return hex.EncodeToString(b)+ext
}
