package model

import (
	"fmt"
	"time"
)

type (
	Photo struct {
		Path      string    `db:"path"`
		Hash      string    `db:"hash"`
		PHash     uint64    `db:"phash"`
		Timestamp time.Time `db:"timestamp"`
	}
)

func (m *Model) InitPhoto() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS photo (
			path TEXT PRIMARY KEY,
			hash TEXT,
			phash BIGINT UNSIGNED,
			timestamp DATETIME
		);
		CREATE INDEX IF NOT EXISTS photo_hash ON photo (hash);
		CREATE INDEX IF NOT EXISTS photo_phash ON photo (phash);
	`)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (m *Model) UpdatePhoto(p *Photo) error {
	_, err := m.db.NamedExec(`
		REPLACE INTO photo VALUES (
			:path,
			:hash,
			:phash,
			:timestamp
		)
	`, p)
	if err != nil {
		return fmt.Errorf("db.NamedExec: %w", err)
	}
	return nil
}
