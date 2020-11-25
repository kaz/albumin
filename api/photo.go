package api

import (
	"fmt"
	"net/http"

	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/scan"
	"github.com/labstack/echo/v4"
)

type (
	DeletePhotoRequest struct {
		Path string
	}
	DeletePhotoResponse struct {
		Photo *model.Photo
	}
	PostPhotoScanRequest struct {
		Directory string
	}
	PostPhotoScanResponse struct {
		Scanned []string
	}
)

func DeletePhoto(c echo.Context) error {
	req := &DeletePhotoRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Bind: %w", err))
	}

	m := c.Get("model").(*model.Model)

	photo, err := m.GetPhotoByPath(req.Path)
	if err != nil {
		return fmt.Errorf("GetPhotoByPath: %w", err)
	}
	if photo == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("no such photo: %v", req.Path))
	}

	photo.Deleted = true
	if err := m.UpdatePhoto(photo); err != nil {
		return fmt.Errorf("UpdatePhoto: %w", err)
	}

	return c.JSON(http.StatusOK, &DeletePhotoResponse{Photo: photo})
}

func PostPhotoScan(c echo.Context) error {
	req := &PostPhotoScanRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Bind: %w", err))
	}
	if req.Directory == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`directory` is not specified")
	}

	ents, err := scan.Scan(req.Directory)
	if err != nil {
		return fmt.Errorf("Scan: %w", err)
	}

	m := c.Get("model").(*model.Model)

	res := &PostPhotoScanResponse{Scanned: make([]string, len(ents))}
	for i, ent := range ents {
		if err := m.UpdatePhoto(ent); err != nil {
			return fmt.Errorf("InsertPhoto: %w", err)
		}
		res.Scanned[i] = ent.Path
	}

	return c.JSON(http.StatusOK, res)
}
