package move

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaz/albumin/model"
)

type (
	Move struct {
		From string
		To   string
	}
)

func Plan(photos []*model.Photo, layout string) ([]*Move, error) {
	result := make([]*Move, 0, len(photos))
	for _, photo := range photos {
		to := photo.Timestamp.Format(layout)
		if !filepath.IsAbs(to) {
			to = filepath.Join(filepath.Dir(photo.Path), to)
		}
		result = append(result, &Move{
			From: photo.Path,
			To:   to,
		})
	}
	return result, nil
}

func Execute(moves []*Move) error {
	for _, move := range moves {
		if err := os.MkdirAll(filepath.Dir(move.To), 0755); err != nil {
			return fmt.Errorf("os.MkdirAll: %w", err)
		}
		if err := os.Rename(move.From, move.To); err != nil {
			return fmt.Errorf("os.Rename: %w", err)
		}
	}
	return nil
}
