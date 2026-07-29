package main

import (
	"database/sql"
	"time"
)

func SaveMessage(db *sql.DB, m *Message) error {
	_, err := db.Exec(`INSERT INTO messages(id,type,content,file_id,sender_device_id,created_at) VALUES(?,?,?,?,?,?)`, m.ID, m.Type, m.Content, m.FileID, m.SenderDeviceID, m.CreatedAt)
	return err
}

func SaveMessageTarget(db *sql.DB, target *MessageTarget) error {
	_, err := db.Exec(`INSERT INTO message_targets(message_id,device_id,status) VALUES(?,?,?)`, target.MessageID, target.DeviceID, target.Status)
	return err
}

func NewTextMessage(id string, content string, sender string) *Message {
	return &Message{
		ID: id,
		Type: MessageText,
		Content: content,
		SenderDeviceID: sender,
		CreatedAt: time.Now(),
	}
}
