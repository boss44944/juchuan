package main

import "database/sql"

// SaveFileHistory stores only file metadata.
// The actual file is always kept on disk.
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
