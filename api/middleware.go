package api

import (
	"fmt"

	"github.com/kaz/albumin/model"
	"github.com/labstack/echo/v4"
)

func ContentTypeJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Content-Type", "application/json")
		return next(c)
	}
}

func QueryPhotosMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := model.Default()
		if err != nil {
			return fmt.Errorf("model.Default: %w", err)
		}
		if err := m.InitPhoto(); err != nil {
			return fmt.Errorf("InitPhoto: %w", err)
		}

		photos, err := m.GetPhotos()
		if err != nil {
			return fmt.Errorf("GetPhotos: %w", err)
		}

		c.Set("photos", photos)
		return next(c)
	}
}
