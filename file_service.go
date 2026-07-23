package main

import (
 "database/sql"
 "path/filepath"
)

func SaveFileHistory(db *sql.DB, filePath, filename, mime string, size int64, direction string) error {
 return AddHistory(db, History{
  Type:"file",
  Direction:direction,
  Filename:filename,
  Filepath:filePath,
  Mime:mime,
  Size:size,
 })
}

func IsImageFile(filename string) bool {
 ext:=filepath.Ext(filename)
 switch ext {
 case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic":
  return true
 }
 return false
}
