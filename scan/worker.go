package scan

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/scan/load"
)

func (s *Scanner) thread() {
	for target := range s.inCh {
		photo, errs := process(target)
		if errs != nil {
			for _, err := range errs {
				s.errCh <- fmt.Errorf("target=%v: %w", target, err)
			}
			continue
		}
		if photo == nil {
			continue
		}
		s.outCh <- photo
	}
}

func process(target string) (*model.Photo, []error) {
	loader, err := load.Load(target)
	if errors.Is(err, load.ErrNotSupported) {
		return nil, nil
	} else if err != nil {
		return nil, []error{fmt.Errorf("load.Load: %w", err)}
	}

	photo := &model.Photo{Path: target}
	errs := []error{}

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		photo.Hash = calcHash(loader)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var err error
		photo.PHash, err = calcPHash(loader)
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
		photo.Timestamp, err = getTimestamp(loader)
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

func calcHash(loader *load.Loader) string {
	hash := md5.Sum(loader.Bytes())
	return hex.EncodeToString(hash[:])
}

func calcPHash(loader *load.Loader) (string, error) {
	img, err := loader.Image()
	if err != nil {
		return "", fmt.Errorf("loader.Image: %w", err)
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", fmt.Errorf("goimagehash.PerceptionHash: %w", err)
	}

	return fmt.Sprintf("%016x", hash.GetHash()), nil
}

func getTimestamp(loader *load.Loader) (time.Time, error) {
	t, err := loader.Time()
	if err != nil {
		return time.Time{}, fmt.Errorf("goimagehash.PerceptionHash: %w", err)
	}
	return t, nil
}
