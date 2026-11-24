package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestTextMessageDeliveryTracksConnectedTargets(t *testing.T) {
	db, err := InitDatabase(filepath.Join(t.TempDir(), "juchuan.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	s := &Server{db: db, hub: NewHub()}
	ws := httptest.NewServer(http.HandlerFunc(s.hub.Handler))
	t.Cleanup(ws.Close)

	connect := func(t *testing.T, device string) *websocket.Conn {
		t.Helper()
		target, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ws.URL, "http")+"?device="+device, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { target.Close() })
		for deadline := time.Now().Add(time.Second); time.Now().Before(deadline); {
			s.hub.mu.Lock()
			connected := s.hub.clients[device] != nil
			s.hub.mu.Unlock()
			if connected {
				return target
			}
			time.Sleep(time.Millisecond)
		}
		t.Fatalf("WebSocket client %s never joined", device)
		return nil
	}
	send := func(t *testing.T, target, content string) string {
		t.Helper()
		b, err := json.Marshal(SendTextRequest{Content: content, SenderDeviceID: "sender", Targets: []string{target}})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		s.SendTextMessageHandler(w, httptest.NewRequest(http.MethodPost, "/api/message/text", strings.NewReader(string(b))))
		if w.Code != http.StatusOK {
			t.Fatalf("send status = %d, body %s", w.Code, w.Body.String())
		}
		var body struct{ Data Message }
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Data.ID == "" {
			t.Fatalf("invalid send response %s (%v)", w.Body.String(), err)
		}
		return body.Data.ID
	}
	assertStatus := func(t *testing.T, id, target, want string) {
		t.Helper()
		var have string
		if err := db.QueryRow("SELECT status FROM message_targets WHERE message_id=? AND device_id=?", id, target).Scan(&have); err != nil {
			t.Fatal(err)
		}
		if have != want {
			t.Fatalf("status for %s = %q; want %q", target, have, want)
		}
	}

	t.Run("offline stays pending", func(t *testing.T) {
		id := send(t, "offline", "pending example")
		assertStatus(t, id, "offline", "CREATED")
	})
	t.Run("connected receives text and delivery status", func(t *testing.T) {
		target := connect(t, "online")
		id := send(t, "online", "hello")
		target.SetReadDeadline(time.Now().Add(2 * time.Second))
		var event struct {
			Type string
			Data Message
		}
		if err := target.ReadJSON(&event); err != nil {
			t.Fatal(err)
		}
		if event.Type != "MESSAGE_RECEIVED" || event.Data.ID != id || event.Data.Content != "hello" {
			t.Fatalf("unexpected received event: %+v", event)
		}
		assertStatus(t, id, "online", "DELIVERED")
	})
}

func TestLateDeliveryDoesNotOverwriteRead(t *testing.T) {
	db, err := InitDatabase(filepath.Join(t.TempDir(), "juchuan.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	s := &Server{db: db}
	message := NewTextMessage("late-delivery", "text", "sender")
	if err := SaveMessage(db, message); err != nil {
		t.Fatal(err)
	}
	if err := SaveMessageTarget(db, &MessageTarget{MessageID: message.ID, DeviceID: "receiver", Status: "READ"}); err != nil {
		t.Fatal(err)
	}
	// The recipient can acknowledge READ while the original send handler is still inside its WebSocket write.
	s.markMessageDelivered(message, "receiver")
	var status sql.NullString
	if err := db.QueryRow("SELECT status FROM message_targets WHERE message_id=? AND device_id=?", message.ID, "receiver").Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status.String != "READ" {
		t.Fatalf("late delivery changed READ into %q", status.String)
	}
}
