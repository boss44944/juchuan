package main

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/system/info", SystemInfoHandler)

	mux.HandleFunc("/api/auth/status", s.AuthStatusHandler)
	mux.HandleFunc("/api/auth/login", s.LoginHandler)
	mux.HandleFunc("/api/auth/logout", s.LogoutHandler)
	mux.HandleFunc("/api/config", s.ConfigHandler)
	mux.HandleFunc("/api/qr", s.QRCodeHandler)

	mux.HandleFunc("/api/history", s.requireAuth(s.HistoryHandler))
	mux.HandleFunc("/upload", s.requireAuth(s.UploadHandler))
	mux.HandleFunc("/download/", s.requireAuth(s.DownloadHandler))

	mux.HandleFunc("/api/device/register", s.requireAuth(s.DeviceRegisterHandler))
	mux.HandleFunc("/api/device/heartbeat", s.requireAuth(s.DeviceHeartbeatHandler))
	mux.HandleFunc("/api/devices", s.requireAuth(s.DevicesHandler))
	mux.HandleFunc("/api/send/file", s.requireAuth(s.SendFileHandler))
}
