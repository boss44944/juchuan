package main

import (
	"net/http"
)

func SystemInfoHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"role": "server",
			"version": "1.0.0",
			"capabilities": []string{
				"device_manager",
				"config",
				"message_router",
			},
		},
	})
}
