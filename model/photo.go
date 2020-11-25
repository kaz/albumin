package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type (
	Photo struct {
		Path      string    `db:"path"`
		Hash      []byte    `db:"hash"`
		PHash     []byte    `db:"phash"`
		Deleted   bool      `db:"deleted"`
		Timestamp time.Time `db:"timestamp"`
	}
)

func (m *Model) InitPhoto() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS photo (
			path TEXT PRIMARY KEY,
			hash BLOB,
			phash BLOB,
			deleted BOOLEAN,
			timestamp DATETIME
		);
		CREATE INDEX IF NOT EXISTS photo_deleted ON photo (deleted);
	`)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (m *Model) GetPhotoByPath(path string) (*Photo, error) {
	photo := &Photo{}
	if err := m.db.Get(photo, "SELECT * FROM photo WHERE path = ?", path); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("db.Get: %w", err)
	}
	return photo, nil
}
func (m *Model) GetPhotos() ([]*Photo, error) {
	photos := []*Photo{}
	if err := m.db.Select(&photos, "SELECT * FROM photo WHERE deleted = ?", false); err != nil {
		return nil, fmt.Errorf("db.Select: %w", err)
	}
	return photos, nil
}

func (m *Model) UpdatePhoto(p *Photo) error {
	_, err := m.db.NamedExec(`
		REPLACE INTO photo VALUES (
			:path,
			:hash,
			:phash,
			:deleted,
			:timestamp
		)
	`, p)
	if err != nil {
		return fmt.Errorf("db.NamedExec: %w", err)
	}
	return nil
}
