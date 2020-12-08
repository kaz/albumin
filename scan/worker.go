package scan

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/scan/load"
	"golang.org/x/sync/errgroup"
)

func thread(reqCh chan string, resCh chan result) {
	for target := range reqCh {
		photo, err := process(target)
		resCh <- result{
			target: target,
			photo:  photo,
			err:    err,
		}
	}
}

func process(target string) (*model.Photo, error) {
	stat, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("os.Stat: %w", err)
	}

	loader, err := load.Load(target)
	if errors.Is(err, load.ErrNotSupported) {
		fmt.Println("ErrNotSupported:", target)
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("load.Load: %w", err)
	}

	photo := &model.Photo{
		Path:    target,
		FsTime:  stat.ModTime(),
		Deleted: false,
	}

	eg := &errgroup.Group{}

	eg.Go(func() error {
		photo.Hash = calcHash(loader)
		return nil
	})
	eg.Go(func() error {
		var err error
		photo.PHash, err = calcPHash(loader)
		if err != nil {
			return fmt.Errorf("calcPHash: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		photo.ExifTime, err = getTimestamp(loader)
		if errors.Is(err, load.ErrNoEXIF) {
			fmt.Println("ErrNoEXIF:", target)
		} else if err != nil {
			return fmt.Errorf("getTimestamp: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("target=%v: %w", target, err)
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
