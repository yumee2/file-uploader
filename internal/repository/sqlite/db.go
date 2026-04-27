package sqlite

import (
	"database/sql"
	"file-uploader/models"
	"fmt"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func NewDBConnection() (*Repository, error) {
	db, err := sql.Open("sqlite", "data/files.db")
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

	return &Repository{db: db}, nil
}

func (d *Repository) AddFile(file *models.File) error {
	_, err := d.db.Exec("INSERT INTO files (id, original_name, size, mime_type) VALUES (?, ?, ?, ?)",
		file.ID, file.OriginalName, file.Size, file.MimeType)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	return nil
}

func (d *Repository) GetFile(id string) (*models.File, error) {
	rows, err := d.db.Query("SELECT id, original_name, size, mime_type FROM files WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query file: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("file not found")
	}

	var file models.File
	if err := rows.Scan(&file.ID, &file.OriginalName, &file.Size, &file.MimeType); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return &file, nil
}

func (d *Repository) GetFiles() ([]*models.File, error) {
	rows, err := d.db.Query("SELECT id, original_name, size, mime_type FROM files")
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.OriginalName, &file.Size, &file.MimeType); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, &file)
	}

	return files, nil
}

func (d *Repository) DeleteFile(id string) error {
	result, err := d.db.Exec("delete from files where id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("file not found")
	}

	return nil
}

func (d *Repository) Close() error {
	return d.db.Close()
}
