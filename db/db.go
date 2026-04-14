package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "files.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS files (
            id            TEXT PRIMARY KEY,
            original_name TEXT,
            size          INTEGER,
            mime_type     TEXT,
            created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}
