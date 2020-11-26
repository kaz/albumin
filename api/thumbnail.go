package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/kaz/albumin/scan/load"
	"github.com/labstack/echo/v4"
)

func GetThumbnail(c echo.Context) error {
	width, err := strconv.Atoi(c.QueryParams().Get("width"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parameter `width` is invalid or null")
	}

	loader, err := load.Load(c.QueryParams().Get("path"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("load.Load: %w", err))
	}

	img, err := loader.Image()
	if err != nil {
		return fmt.Errorf("loader.Image: %w", err)
	}

	r, w := io.Pipe()
	go func() {
		imaging.Encode(w, imaging.Resize(img, width, 0, imaging.NearestNeighbor), imaging.JPEG, imaging.JPEGQuality(80))
		w.Close()
		r.Close()
	}()

	return c.Stream(http.StatusOK, "image/jpeg", r)
}
