package scan

import (
	"fmt"
	"sync"
	"time"

	"github.com/kaz/albumin/model"
)

type (
	worker struct {
	}
)

func (s *Scanner) thread() {
	for targer := range s.inCh {
		photo, errs := process(targer)
		if errs != nil {
			for _, err := range errs {
				s.errCh <- fmt.Errorf("processing %v: %w", targer, err)
			}
			continue
		}
		s.outCh <- photo
	}
}

func process(target string) (*model.Photo, []error) {
	photo := &model.Photo{Path: target}
	errs := []error{}

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		photo.Hash, err = calcHash(target)
		if err != nil {
			mu.Lock()
			defer mu.Unlock()

			errs = append(errs, fmt.Errorf("calcHash: %w", err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		photo.PHash, err = calcPHash(target)
		if err != nil {
			mu.Lock()
			defer mu.Unlock()

			errs = append(errs, fmt.Errorf("calcPHash: %w", err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		photo.Timestamp, err = getTimestamp(target)
		if err != nil {
			mu.Lock()
			defer mu.Unlock()

			errs = append(errs, fmt.Errorf("getTimestamp: %w", err))
		}
	}()

	wg.Wait()
	if len(errs) > 0 {
		return nil, errs
	}
	return photo, nil
}

func calcHash(target string) (string, error) {
	return "", nil
}

func calcPHash(target string) (uint64, error) {
	return 0, nil
}

func getTimestamp(target string) (time.Time, error) {
	return time.Unix(0, 0), nil
}
