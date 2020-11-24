package move

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaz/albumin/model"
)

type (
	Move struct {
		From string
		To   string
	}
)

func Plan(photos []*model.Photo, strategy Strategy) ([]*Move, error) {
	result := make([]*Move, 0, len(photos))
	count := make(map[string]int, len(photos))

	for _, photo := range photos {
		move := &Move{
			From: photo.Path,
			To:   strategy(photo),
		}

		count[move.To]++
		if count[move.To] > 1 {
			ext := filepath.Ext(move.To)
			move.To = strings.Replace(move.To, ext, fmt.Sprintf("(%d)%s", count[move.To], ext), 1)
		}

		result = append(result, move)
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
