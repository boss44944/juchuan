package main

import "time"

type MessageType string

const (
	MessageText  MessageType = "TEXT"
	MessageFile  MessageType = "FILE"
	MessageImage MessageType = "IMAGE"
)

type Message struct {
	ID              string      `json:"id"`
	Type            MessageType `json:"type"`
	Content         string      `json:"content,omitempty"`
	FileID          string      `json:"file_id,omitempty"`
	SenderDeviceID  string      `json:"sender_device_id"`
	CreatedAt       time.Time   `json:"created_at"`
}

type MessageTarget struct {
	MessageID string `json:"message_id"`
	DeviceID  string `json:"device_id"`
	Status    string `json:"status"`
}
