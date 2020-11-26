package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type (
	Photo struct {
		Path     string    `db:"path"`
		Hash     []byte    `db:"hash"`
		PHash    []byte    `db:"phash"`
		Deleted  bool      `db:"deleted"`
		FsTime   time.Time `db:"fs_time"`
		ExifTime time.Time `db:"exif_time"`
	}
)

func (p *Photo) Timestamp() int64 {
	exifTime := p.ExifTime.UnixNano()
	if exifTime > 0 {
		return exifTime
	}
	return p.FsTime.UnixNano()
}

func (m *Model) InitPhoto() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS photo (
			path TEXT PRIMARY KEY,
			hash BLOB,
			phash BLOB,
			deleted BOOLEAN,
			fs_time DATETIME,
			exif_time DATETIME
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

func (m *Model) ReplacePhotos(photos []*Photo) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("db.Beginx: %w", err)
	}
	defer tx.Rollback()

	for _, photo := range photos {
		_, err := tx.NamedExec(`
			REPLACE INTO photo VALUES (
				:path,
				:hash,
				:phash,
				:deleted,
				:fs_time,
				:exif_time
			)
		`, photo)
		if err != nil {
			return fmt.Errorf("tx.NamedExec: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
