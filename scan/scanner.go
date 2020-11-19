package scan

import (
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
			photos = append(photos, photo)
		}
	}()

	errs := []error{}
	go func() {
		for err := range s.errCh {
			errs = append(errs, err)
		}
	}()

	for i := 0; i < preference.ScanThread; i++ {
		go s.thread()
	}

}

func (s *Scanner) scan(target string) error {

}
