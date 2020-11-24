package main

import (
	"github.com/kaz/albumin/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Debug = true
	e.Use(middleware.Logger())

	apiGroup := e.Group("/api", api.ContentTypeJSON)
	apiGroup.POST("/scan", api.PostScan)

	e.Logger.Fatal(e.Start(":20000"))
}
