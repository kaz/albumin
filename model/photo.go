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

func (m *Model) InsertPhoto(p *Photo) error {
	_, err := m.db.NamedExec(`
		INSERT INTO photo VALUES (
			:path,
			:hash,
			:phash,
			:timestamp
		)
	`, p)
	if err != nil {
		return fmt.Errorf("NamedExec: %w", err)
	}
	return nil
}
