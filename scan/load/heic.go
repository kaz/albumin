package load

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/strukturag/libheif/go/heif"
	mediaheif "go4.org/media/heif"
)

type (
	HeicLoader struct{}
)

func (l *HeicLoader) Image(data []byte) (image.Image, error) {
	ctx, err := heif.NewContext()
	if err != nil {
		return nil, fmt.Errorf("heif.NewContext: %w", err)
	}

	if err := ctx.ReadFromMemory(data); err != nil {
		return nil, fmt.Errorf("ReadFromMemory: %w", err)
	}

	handle, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return nil, fmt.Errorf("GetPrimaryImageHandle: %w", err)
	}

	hImg, err := handle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	if err != nil {
		return nil, fmt.Errorf("DecodeImage: %w", err)
	}

	img, err := hImg.GetImage()
	if err != nil {
		return nil, fmt.Errorf("GetImage: %w", err)
	}

	return img, nil
}

func (l *HeicLoader) Time(data []byte) (time.Time, error) {
	rawExif, err := mediaheif.Open(bytes.NewReader(data)).EXIF()
	if err != nil {
		return time.Time{}, fmt.Errorf("mediaheif.Open: %w", err)
	}

	meta, err := exif.Decode(bytes.NewReader(rawExif))
	if err != nil {
		return time.Time{}, fmt.Errorf("exif.Decode: %w", err)
	}

	dt, err := meta.DateTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("DateTime: %w", err)
	}

	return dt, nil
}
