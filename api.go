package main

import (
	"encoding/json"
	"net/http"
)

func TextHandler(clipboard *Clipboard) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := clipboard.Copy(req.Content); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("ok"))
	}
}
