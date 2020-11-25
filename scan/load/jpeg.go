package load

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type (
	JpegLoader struct{}
)

func (l *JpegLoader) Image(data []byte) (image.Image, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("jpeg.Decode: %w", err)
	}

	return img, nil
}

func (l *JpegLoader) Time(data []byte) (time.Time, error) {
	meta, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrNoEXIF, err)
	}

	dt, err := meta.DateTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("DateTime: %w", err)
	}

	return dt, nil
}
