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
	apiGroup.GET("/scan/pwd", api.GetScanPwd)
	apiGroup.POST("/scan", api.PostScan)
	apiGroup.GET("/dedup/hash", api.GetDedupHash, api.DedupMiddleware)
	apiGroup.GET("/dedup/phash", api.GetDedupPHash, api.DedupMiddleware)

	e.Logger.Fatal(e.Start(":20000"))
}
