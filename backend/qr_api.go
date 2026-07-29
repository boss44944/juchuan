package main

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

func (s *Server) QRCodeHandler(w http.ResponseWriter, r *http.Request) {
	urlText := r.URL.Query().Get("url")
	if urlText == "" {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		urlText = scheme + "://" + r.Host + "/"
	}

	png, err := qrcode.Encode(urlText, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "qrcode encode failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(png)
}
