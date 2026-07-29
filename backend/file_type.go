package main

import "mime"

func DetectMime(filename string) string {
	if v := mime.TypeByExtension("." + extension(filename)); v != "" {
		return v
	}
	return "application/octet-stream"
}

func extension(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return ""
}
