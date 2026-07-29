package main

import (
	"database/sql"
)

func ReceiveText(db *sql.DB, clipboard *Clipboard, hub *Hub, content string) error {
	if err := AddHistory(db, History{
		Type:      "text",
		Direction: "phone_to_pc",
		Content:   content,
	}); err != nil {
		return err
	}

	if clipboard != nil {
		_ = clipboard.Copy(content)
	}

	if hub != nil {
		hub.Broadcast(WSMessage{
			Type:    "text",
			Content: content,
		})
	}
	return nil
}
