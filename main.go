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

	e.Static("/", "./view")

	apiGroup := e.Group("/api", api.ContentTypeJSON)
	apiGroup.GET("/file", api.GetFile)
	apiGroup.DELETE("/photo", api.DeletePhoto, api.SetupModelMiddleware)
	apiGroup.POST("/photo/scan", api.PostPhotoScan, api.SetupModelMiddleware)
	apiGroup.GET("/dedup/hash", api.GetDedupHash, api.QueryPhotosMiddleware)
	apiGroup.GET("/dedup/phash", api.GetDedupPHash, api.QueryPhotosMiddleware)
	apiGroup.GET("/move/pwd", api.GetMovePwd)
	apiGroup.POST("/move", api.PostMove, api.QueryPhotosMiddleware)

	e.Logger.Fatal(e.Start(":20000"))
}
