package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MessageListItem struct {
	MessageID      string `json:"message_id"`
	Type           string `json:"type"`
	Content        string `json:"content,omitempty"`
	FileID         string `json:"file_id,omitempty"`
	SenderDeviceID string `json:"sender_device_id"`
	TargetDeviceID string `json:"target_device_id,omitempty"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}

type MessageListResponse struct {
	Items []MessageListItem `json:"items"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

func (s *Server) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}

	typeFilter := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("type")))
	statusFilter := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("status")))
	deviceID := strings.TrimSpace(r.URL.Query().Get("device_id"))
	senderDeviceID := strings.TrimSpace(r.URL.Query().Get("sender_device_id"))
	targetDeviceID := strings.TrimSpace(r.URL.Query().Get("target_device_id"))

	where := []string{"1=1"}
	args := make([]any, 0, 8)

	if typeFilter != "" {
		where = append(where, "m.type=?")
		args = append(args, typeFilter)
	}
	if statusFilter != "" {
		where = append(where, "COALESCE(mt.status, 'CREATED')=?")
		args = append(args, statusFilter)
	}
	if deviceID != "" {
		where = append(where, "(m.sender_device_id=? OR mt.device_id=?)")
		args = append(args, deviceID, deviceID)
	}
	if senderDeviceID != "" {
		where = append(where, "m.sender_device_id=?")
		args = append(args, senderDeviceID)
	}
	if targetDeviceID != "" {
		where = append(where, "mt.device_id=?")
		args = append(args, targetDeviceID)
	}

	whereSQL := strings.Join(where, " AND ")

	countQuery := `
SELECT COUNT(*)
FROM messages m
LEFT JOIN message_targets mt ON mt.message_id=m.id
WHERE ` + whereSQL

	var total int
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_LIST_FAILED", nil)
		return
	}

	listQuery := `
SELECT
  m.id,
  m.type,
	COALESCE(m.content, ''),
	COALESCE(m.file_id, ''),
	COALESCE(m.sender_device_id, ''),
  COALESCE(mt.device_id, ''),
  COALESCE(mt.status, 'CREATED'),
	COALESCE(m.created_at, '')
FROM messages m
LEFT JOIN message_targets mt ON mt.message_id=m.id
WHERE ` + whereSQL + `
ORDER BY m.created_at DESC
LIMIT ? OFFSET ?`

	listArgs := append([]any{}, args...)
	listArgs = append(listArgs, size, (page-1)*size)

	rows, err := s.db.Query(listQuery, listArgs...)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_LIST_FAILED", nil)
		return
	}
	defer rows.Close()

	items := make([]MessageListItem, 0, size)
	for rows.Next() {
		var item MessageListItem
		if err := rows.Scan(
			&item.MessageID,
			&item.Type,
			&item.Content,
			&item.FileID,
			&item.SenderDeviceID,
			&item.TargetDeviceID,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			WriteError(w, http.StatusInternalServerError, "MESSAGE_LIST_FAILED", nil)
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		WriteError(w, http.StatusInternalServerError, "MESSAGE_LIST_FAILED", nil)
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: MessageListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	}})
}
