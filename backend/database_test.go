package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseAppliesBusyTimeoutToEveryConnection(t *testing.T) {
	// Spaces and URI metacharacters in a legitimate storage filename need not be escaped by callers.
	path := filepath.Join(t.TempDir(), "Juchuan fixture #%.db")
	db, err := InitDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("database did not use the exact filename: %v", err)
	}

	db.SetMaxOpenConns(2)
	first, err := db.Conn(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := db.Conn(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()

	for i, connection := range []*sql.Conn{first, second} {
		var timeout int
		if err := connection.QueryRowContext(t.Context(), "PRAGMA busy_timeout").Scan(&timeout); err != nil {
			t.Fatal(err)
		}
		if timeout != sqliteBusyTimeoutMillis {
			t.Errorf("connection %d busy_timeout = %d ms; want %d ms", i+1, timeout, sqliteBusyTimeoutMillis)
		}
	}
}
