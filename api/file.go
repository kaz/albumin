package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func GetFile(c echo.Context) error {
	file, err := os.Open(c.QueryParams().Get("path"))
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	return c.Stream(http.StatusOK, "application/octet-stream", file)
}
