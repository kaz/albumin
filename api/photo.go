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
		Path        string
		ProgressKey string
	}
	PostPhotoScanResponse struct {
		Photos []*model.Photo
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
	if err := m.ReplacePhotos([]*model.Photo{photo}); err != nil {
		return fmt.Errorf("ReplacePhotos: %w", err)
	}

	return c.JSON(http.StatusOK, &DeletePhotoResponse{Photo: photo})
}

func PostPhotoScan(c echo.Context) error {
	req := &PostPhotoScanRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Bind: %w", err))
	}
	if req.Path == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`Path` is not specified")
	}
	if req.ProgressKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`ProgressKey` is not specified")
	}

	photos, err := scan.Scan(req.Path, req.ProgressKey)
	if err != nil {
		return fmt.Errorf("Scan: %w", err)
	}

	m := c.Get("model").(*model.Model)
	if err := m.ReplacePhotos(photos); err != nil {
		return fmt.Errorf("ReplacePhotos: %w", err)
	}

	return c.JSON(http.StatusOK, &PostPhotoScanResponse{Photos: photos})
}

func GetPhotoScanProgress(c echo.Context) error {
	return c.String(http.StatusOK, scan.GetProgress(c.QueryParams().Get("key")))
}
