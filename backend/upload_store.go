package main

import (
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveUploadedFile(db *sql.DB, storage *Storage, header *multipart.FileHeader) (string, error) {
	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dayDir := filepath.Join(storage.FileDir, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(dayDir, 0755); err != nil {
		return "", err
	}

	name := SafeFileName(dayDir, header.Filename)
	path := filepath.Join(dayDir, name)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	if err := SaveFileHistory(db, path, name, "", header.Size, "phone_to_pc"); err != nil {
		return "", err
	}

	return path, nil
}
