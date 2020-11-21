package command

import (
	"fmt"

	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/scan"
)

type (
	Scan struct {
		Directory string `short:"d" long:"directory" required:"true"`
	}
)

func (s *Scan) Execute(args []string) error {
	m, err := model.Default()
	if err != nil {
		return fmt.Errorf("model.Default: %w", err)
	}
	if err := m.InitPhoto(); err != nil {
		return fmt.Errorf("InitPhoto: %w", err)
	}

	ents, err := scan.Scan(s.Directory)
	if err != nil {
		return fmt.Errorf("Scan: %w", err)
	}

	for _, ent := range ents {
		if err := m.UpdatePhoto(ent); err != nil {
			return fmt.Errorf("InsertPhoto: %w", err)
		}
	}

	return nil
}
