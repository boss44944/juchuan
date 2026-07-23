package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// HistoryRecord stores clipboard history items in SQLite.
type HistoryRecord struct {
	ID int64
	Type string
	Content string
	Filename string
	Size int64
	URL string
	CreatedAt time.Time
}

// HistoryStore manages persistent history records.
type HistoryStore struct {
	db *sql.DB
}

func NewHistoryStore(path string) (*HistoryStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil { return nil, err }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		content TEXT,
		filename TEXT,
		size INTEGER DEFAULT 0,
		url TEXT,
		created_at DATETIME NOT NULL
	)`)
	if err != nil { return nil, err }

	return &HistoryStore{db: db}, nil
}

func (h *HistoryStore) Add(r HistoryRecord) error {
	_, err := h.db.Exec(`INSERT INTO history(type,content,filename,size,url,created_at) VALUES(?,?,?,?,?,?)`, r.Type,r.Content,r.Filename,r.Size,r.URL,r.CreatedAt)
	return err
}

func (h *HistoryStore) List(page,size int) ([]HistoryRecord,error) {
	if page < 1 { page=1 }
	if size < 1 { size=20 }
	offset := (page-1)*size
	rows,err:=h.db.Query(`SELECT id,type,content,filename,size,url,created_at FROM history ORDER BY id DESC LIMIT ? OFFSET ?`,size,offset)
	if err!=nil{return nil,err}
	defer rows.Close()
	var result []HistoryRecord
	for rows.Next(){
		var r HistoryRecord
		if err:=rows.Scan(&r.ID,&r.Type,&r.Content,&r.Filename,&r.Size,&r.URL,&r.CreatedAt);err!=nil{return nil,err}
		result=append(result,r)
	}
	return result,nil
}

func (h *HistoryStore) Close() error { if h.db==nil{return fmt.Errorf("database closed")}; return h.db.Close() }
