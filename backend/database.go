package main

import (
	"database/sql"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

const sqliteBusyTimeoutMillis = 1000

func InitDatabase(path string) (*sql.DB, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	filePath := filepath.ToSlash(absPath)
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath // A Windows volume path is represented as /C:/... in file URIs.
	}
	// SQLite locks can briefly overlap between a page read and a device write.
	// The driver applies this bounded wait to each connection in the pool.
	params := url.Values{"_pragma": {"busy_timeout(" + strconv.Itoa(sqliteBusyTimeoutMillis) + ")"}}
	dsn := url.URL{Scheme: "file", Path: filePath, RawQuery: params.Encode()}
	db, err := sql.Open("sqlite", dsn.String())
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
	if err != nil { return nil, err }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS devices(
 id TEXT PRIMARY KEY,
 display_name TEXT NOT NULL,
 role TEXT,
 platform TEXT,
 browser TEXT,
 device_secret TEXT,
 last_seen INTEGER
 );`)
	if err != nil { return nil, err }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages(
 id TEXT PRIMARY KEY,
 type TEXT NOT NULL,
 content TEXT,
 file_id TEXT,
 sender_device_id TEXT,
 created_at DATETIME DEFAULT CURRENT_TIMESTAMP
 );`)
	if err != nil { return nil, err }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS message_targets(
 message_id TEXT NOT NULL,
 device_id TEXT NOT NULL,
 status TEXT NOT NULL,
 updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY(message_id, device_id)
 );`)
	if err != nil { return nil, err }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS files(
 id TEXT PRIMARY KEY,
 filename TEXT,
 path TEXT,
 mime TEXT,
 size INTEGER,
 created_at DATETIME DEFAULT CURRENT_TIMESTAMP
 );`)
	if err != nil { return nil, err }

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_history_created_at ON history(created_at DESC);`)
	return db, err
}
