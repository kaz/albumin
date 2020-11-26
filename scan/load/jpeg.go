package load

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

type (
	JpegLoader struct{}
)

func (l *JpegLoader) Image(data []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("imaging.Decode: %w", err)
	}

	return img, nil
}

func (l *JpegLoader) Time(data []byte) (time.Time, error) {
	meta, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrNoEXIF, err)
	}

	dt, err := meta.DateTime()
	if exif.IsTagNotPresentError(err) {
		return time.Time{}, fmt.Errorf("%w: %v", ErrNoEXIF, err)
	} else if err != nil {
		return time.Time{}, fmt.Errorf("DateTime: %w", err)
	}

	return dt, nil
}
