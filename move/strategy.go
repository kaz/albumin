package move

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaz/albumin/model"
)

type (
	Strategy func(*model.Photo) string
)

var (
	Pwd string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Pwd = pwd
}

func resolve(photo *model.Photo, to string) string {
	if !filepath.IsAbs(to) {
		to = filepath.Join(Pwd, to)
	}

	realExt := strings.ToLower(filepath.Ext(photo.Path))
	return strings.Replace(to, ".png", realExt, 1)
}

func StrategyExif(layout string) Strategy {
	return func(photo *model.Photo) string {
		return resolve(photo, photo.ExifTime.Format(layout))
	}
}

func StrategySequential(layout string) Strategy {
	seq := 0
	return func(photo *model.Photo) string {
		seq++
		return resolve(photo, fmt.Sprintf(layout, seq))
	}
}
