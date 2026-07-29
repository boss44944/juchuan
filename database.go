package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func InitDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS history(
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 type TEXT NOT NULL,
 direction TEXT,
 content TEXT,
 filepath TEXT,
 filename TEXT,
 mime TEXT,
 size INTEGER,
 created_at DATETIME DEFAULT CURRENT_TIMESTAMP
 );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS devices(
 id TEXT PRIMARY KEY,
 display_name TEXT NOT NULL,
 role TEXT,
 platform TEXT,
 browser TEXT,
 device_secret TEXT,
 last_seen INTEGER
 );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages(
 id TEXT PRIMARY KEY,
 type TEXT NOT NULL,
 content TEXT,
 file_id TEXT,
 sender_device_id TEXT,
 created_at DATETIME DEFAULT CURRENT_TIMESTAMP
 );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS message_targets(
 message_id TEXT NOT NULL,
 device_id TEXT NOT NULL,
 status TEXT NOT NULL,
 PRIMARY KEY(message_id, device_id)
 );`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_history_created_at ON history(created_at DESC);`)
	return db, err
}
