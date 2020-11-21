package load

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type (
	ImageLoader interface {
		Image([]byte) (image.Image, error)
		Time([]byte) (time.Time, error)
	}

	Loader struct {
		data []byte
		il   ImageLoader
	}
)

var (
	ErrNotSupported = errors.New("not supported")
)

func Load(file string) (*Loader, error) {
	var il ImageLoader

	switch strings.ToLower(filepath.Ext(file)) {
	case ".jpg":
		il = &JpegLoader{}
	case ".jpeg":
		il = &JpegLoader{}
	case ".heic":
		il = &HeicLoader{}
	default:
		return nil, ErrNotSupported
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %w", err)
	}

	return &Loader{data: data, il: il}, nil
}

func (l *Loader) Bytes() []byte {
	return l.data
}
func (l *Loader) Image() (image.Image, error) {
	return l.il.Image(l.data)
}
func (l *Loader) Time() (time.Time, error) {
	return l.il.Time(l.data)
}
