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
	result struct {
		target string
		photo  *model.Photo
		errs   []error
	}
	progress struct {
		total   int
		current int
	}
)

var (
	mu      = &sync.RWMutex{}
	progMap = map[string]*progress{}
)

func GetProgress(key string) string {
	prog, ok := progMap[key]
	if !ok {
		return "preparing"
	}
	return fmt.Sprintf("%.2f %%", 100*float64(prog.current)/float64(prog.total))
}

func Scan(target string, progKey string) ([]*model.Photo, error) {
	ents, err := walk(target)
	if err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}

	wg := &sync.WaitGroup{}
	reqCh := make(chan string)
	resCh := make(chan result)

	for i := 0; i < preference.ScanThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			thread(reqCh, resCh)
		}()
	}

	go func() {
		for _, ent := range ents {
			reqCh <- ent
		}
		close(reqCh)
		wg.Wait()
		close(resCh)
	}()

	progMap[progKey] = &progress{total: len(ents)}
	photos := []*model.Photo{}
	errs := []error{}

	for res := range resCh {
		progMap[progKey].current++

		if res.errs != nil {
			for _, err := range res.errs {
				errs = append(errs, fmt.Errorf("target=%v: %w", res.target, err))
			}
			continue
		}
		if res.photo == nil {
			continue
		}

		photos = append(photos, res.photo)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("in thread: %v", errs)
	}
	return photos, nil
}

func walk(target string) ([]string, error) {
	ents, err := ioutil.ReadDir(target)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir: %w", err)
	}

	result := []string{}
	for _, ent := range ents {
		entPath, err := filepath.Abs(filepath.Join(target, ent.Name()))
		if err != nil {
			return nil, fmt.Errorf("filepath.Abs: %w", err)
		}

		if !ent.IsDir() {
			result = append(result, entPath)
			continue
		}

		children, err := walk(entPath)
		if err != nil {
			return nil, fmt.Errorf("child=%v: %w", ent.Name(), err)
		}

		result = append(result, children...)
	}

	return result, nil
}
