package main

import "time"

type FileRecord struct {
	ID string `json:"id"`
	Filename string `json:"filename"`
	Path string `json:"path"`
	Mime string `json:"mime"`
	Size int64 `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
