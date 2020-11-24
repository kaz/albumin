package api

import (
	"github.com/labstack/echo/v4"
)

func ContentTypeJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Content-Type", "application/json")
		return next(c)
	}
}
