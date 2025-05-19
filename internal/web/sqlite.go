// Lab 7: Implement a SQLite video metadata service

package web

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteVideoMetadataService struct {
	Db string
}

func (s *SQLiteVideoMetadataService) ensureTable(dbInstance *sql.DB) error {
	_, err := dbInstance.Exec(`
        CREATE TABLE IF NOT EXISTS videos (
            id TEXT PRIMARY KEY,
            uploaded_at DATETIME
        )
    `)
	return err
}

// Create implements VideoMetadataService.
func (s *SQLiteVideoMetadataService) Create(videoId string, uploadedAt time.Time) error {
	dbInstance, err := sql.Open("sqlite3", s.Db)
	if err != nil {
		return err
	}
	defer dbInstance.Close()

	if err := s.ensureTable(dbInstance); err != nil {
		return err
	}

	_, err = dbInstance.Exec(`
        INSERT INTO videos (id, uploaded_at) VALUES (?, ?)
    `, videoId, uploadedAt)
	if err != nil {
		return err
	}

	return nil
}

// List implements VideoMetadataService.
func (s *SQLiteVideoMetadataService) List() ([]VideoMetadata, error) {
	dbInstance, err := sql.Open("sqlite3", s.Db)
	if err != nil {
		return nil, err
	}
	defer dbInstance.Close()

	if err := s.ensureTable(dbInstance); err != nil {
		return nil, err
	}

	rows, err := dbInstance.Query("SELECT id, uploaded_at FROM videos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []VideoMetadata
	for rows.Next() {
		var video VideoMetadata
		if err := rows.Scan(&video.Id, &video.UploadedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// Read implements VideoMetadataService.
func (s *SQLiteVideoMetadataService) Read(id string) (*VideoMetadata, error) {
	dbInstance, err := sql.Open("sqlite3", s.Db)
	if err != nil {
		return nil, err
	}
	defer dbInstance.Close()

	if err := s.ensureTable(dbInstance); err != nil {
		return nil, err
	}

	var video VideoMetadata
	err = dbInstance.QueryRow("SELECT id, uploaded_at FROM videos WHERE id = ?", id).Scan(&video.Id, &video.UploadedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &video, nil
}

// Uncomment the following line to ensure SQLiteVideoMetadataService implements VideoMetadataService
var _ VideoMetadataService = (*SQLiteVideoMetadataService)(nil)
