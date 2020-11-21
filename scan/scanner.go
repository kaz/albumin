package scan

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/preference"
)

type (
	Scanner struct {
		inCh  chan string
		outCh chan *model.Photo
		errCh chan error
	}
)

func New() *Scanner {
	return &Scanner{
		inCh:  make(chan string),
		outCh: make(chan *model.Photo),
		errCh: make(chan error),
	}
}

func (s *Scanner) Scan(target string) ([]*model.Photo, error) {
	photos := []*model.Photo{}
	go func() {
		for photo := range s.outCh {
			fmt.Println("OK:", photo.Path)
			photos = append(photos, photo)
		}
	}()

	errs := []error{}
	go func() {
		for err := range s.errCh {
			fmt.Println("NG:", err)
			errs = append(errs, err)
		}
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < preference.ScanThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.thread()
		}()
	}

	if err := s.scan(target); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	close(s.inCh)

	wg.Wait()
	close(s.outCh)
	close(s.errCh)

	if len(errs) > 0 {
		return nil, fmt.Errorf("errors occurred in scanning: %v", errs)
	}
	return photos, nil
}

func (s *Scanner) scan(target string) error {
	ents, err := ioutil.ReadDir(target)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir: %w", err)
	}

	for _, ent := range ents {
		entPath, err := filepath.Abs(filepath.Join(target, ent.Name()))
		if err != nil {
			return fmt.Errorf("filepath.Abs: %w", err)
		}

		if ent.IsDir() {
			if err := s.scan(entPath); err != nil {
				return fmt.Errorf("child=%v: %w", ent.Name(), err)
			}
			continue
		}

		s.inCh <- entPath
	}

	return nil
}
