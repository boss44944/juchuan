package main

import (
 "os"
 "path/filepath"
 "strconv"
)

// SafeFileName keeps the original filename.
// If the same name already exists, append (1), (2), etc.
func SafeFileName(dir, name string) string {
 base := filepath.Base(name)
 ext := filepath.Ext(base)
 filename := base[:len(base)-len(ext)]

 candidate := base
 index := 1

 for {
  if _, err := os.Stat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
   return candidate
  }
  candidate = filename + "(" + strconv.Itoa(index) + ")" + ext
  index++
 }
}
