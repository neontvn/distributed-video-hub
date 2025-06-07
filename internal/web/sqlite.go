// Lab 7: Implement a SQLite video metadata service

package web

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteVideoMetadataService struct {
	Instance *sql.DB
}

func (s *SQLiteVideoMetadataService) ensureTable() error {
	_, err := s.Instance.Exec(`
        CREATE TABLE IF NOT EXISTS videos (
            id TEXT PRIMARY KEY,
            uploaded_at DATETIME
        )
    `)
	return err
}

// Create implements VideoMetadataService.
func (s *SQLiteVideoMetadataService) Create(videoId string, uploadedAt time.Time) error {
	if err := s.ensureTable(); err != nil {
		return err
	}

	_, err := s.Instance.Exec(`INSERT INTO videos (id, uploaded_at) VALUES (?, ?)`, videoId, uploadedAt)
	if err != nil {
		return err
	}

	return nil
}

// List implements VideoMetadataService.
func (s *SQLiteVideoMetadataService) List() ([]VideoMetadata, error) {
	if err := s.ensureTable(); err != nil {
		return nil, err
	}

	rows, err := s.Instance.Query("SELECT id, uploaded_at FROM videos")
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

	if err := s.ensureTable(); err != nil {
		return nil, err
	}

	var video VideoMetadata
	err := s.Instance.QueryRow("SELECT id, uploaded_at FROM videos WHERE id = ?", id).Scan(&video.Id, &video.UploadedAt)
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
