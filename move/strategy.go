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
	pwd string
)

func init() {
	pwd, _ = os.Getwd()
}

func resolve(photo *model.Photo, to string) string {
	if !filepath.IsAbs(to) {
		to = filepath.Join(pwd, to)
	}
	return strings.Replace(to, ".png", filepath.Ext(photo.Path), 1)
}

func StrategyExif(layout string) Strategy {
	return func(photo *model.Photo) string {
		return resolve(photo, photo.Timestamp.Format(layout))
	}
}

func StrategySequential(layout string) Strategy {
	seq := 0
	return func(photo *model.Photo) string {
		seq++
		return resolve(photo, fmt.Sprintf(layout, seq))
	}
}
