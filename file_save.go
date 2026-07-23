package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveUploadedFile(header *multipart.FileHeader, src multipart.File, storage *Storage, isImage bool) (string, error) {
	dir := storage.FileDir
	if isImage {
		dir = storage.ImageDir
	}

	name := randomName(filepath.Ext(header.Filename))
	dstPath := filepath.Join(dir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return dstPath, nil
}

func randomName(ext string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x%s", b, ext)
}
