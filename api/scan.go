package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/scan"
	"github.com/labstack/echo/v4"
)

type (
	PostScanRequest struct {
		Directory string
	}
	PostScanResponse struct {
		Scanned []string
	}
)

func GetScanPwd(c echo.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}
	return c.String(http.StatusOK, pwd)
}

func PostScan(c echo.Context) error {
	req := &PostScanRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Bind: %w", err))
	}
	if req.Directory == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`directory` is not specified")
	}

	m, err := model.Default()
	if err != nil {
		return fmt.Errorf("model.Default: %w", err)
	}
	if err := m.InitPhoto(); err != nil {
		return fmt.Errorf("InitPhoto: %w", err)
	}

	ents, err := scan.Scan(req.Directory)
	if err != nil {
		return fmt.Errorf("Scan: %w", err)
	}

	res := &PostScanResponse{Scanned: make([]string, len(ents))}
	for i, ent := range ents {
		if err := m.UpdatePhoto(ent); err != nil {
			return fmt.Errorf("InsertPhoto: %w", err)
		}
		res.Scanned[i] = ent.Path
	}

	return c.JSON(http.StatusOK, res)
}
