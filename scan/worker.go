package scan

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/scan/load"
)

func thread(reqCh chan string, resCh chan result) {
	for target := range reqCh {
		photo, errs := process(target)
		resCh <- result{
			target: target,
			photo:  photo,
			errs:   errs,
		}
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

func calcHash(loader *load.Loader) []byte {
	hash := md5.Sum(loader.Bytes())
	return hash[:]
}

func calcPHash(loader *load.Loader) ([]byte, error) {
	img, err := loader.Image()
	if err != nil {
		return nil, fmt.Errorf("loader.Image: %w", err)
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, fmt.Errorf("goimagehash.PerceptionHash: %w", err)
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, hash.GetHash())
	return buf, nil
}

func getTimestamp(loader *load.Loader) (time.Time, error) {
	t, err := loader.Time()
	if err != nil {
		return time.Time{}, fmt.Errorf("loader.Time: %w", err)
	}
	return t, nil
}
