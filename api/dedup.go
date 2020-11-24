package api

import (
	"net/http"
	"strconv"

	"github.com/kaz/albumin/dedup"
	"github.com/kaz/albumin/model"
	"github.com/labstack/echo/v4"
)

type (
	GetDedupHashResponse struct {
		Groups [][]*model.Photo
	}
	GetDedupPHashResponse struct {
		Groups [][]*model.Photo
	}
)

func GetDedupHash(c echo.Context) error {
	photos := c.Get("photos").([]*model.Photo)
	groups := dedup.GroupByHash(photos)
	return c.JSON(http.StatusOK, &GetDedupHashResponse{Groups: groups})
}

func GetDedupPHash(c echo.Context) error {
	tolerance, err := strconv.Atoi(c.QueryParams().Get("tolerance"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "parameter `tolerance` is invalid or null")
	}

	photos := c.Get("photos").([]*model.Photo)
	groups := dedup.GroupByPHash(photos, tolerance)
	return c.JSON(http.StatusOK, &GetDedupPHashResponse{Groups: groups})
}
