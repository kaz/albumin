package load

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/disintegration/imaging"
)

type (
	PngLoader struct{}
)

func (l *PngLoader) Image(data []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("imaging.Decode: %w", err)
	}

	return img, nil
}

func (l *PngLoader) Time(data []byte) (time.Time, error) {
	return time.Time{}, ErrNoEXIF
}
