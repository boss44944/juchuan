package main

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code   string                 `json:"code"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, status int, code string, params map[string]interface{}) {
	WriteJSON(w, status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:   code,
			Params: params,
		},
	})
}
