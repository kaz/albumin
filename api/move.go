package api

import (
	"fmt"
	"net/http"

	"github.com/kaz/albumin/model"
	"github.com/kaz/albumin/move"
	"github.com/labstack/echo/v4"
)

type (
	PostMoveRequest struct {
		Layout  string
		Execute bool
	}
	PostMoveResponse struct {
		Moves []*move.Move
	}
)

func PostMove(c echo.Context) error {
	req := &PostMoveRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Bind: %w", err))
	}

	photos := c.Get("photos").([]*model.Photo)
	moves, err := move.Plan(photos, req.Layout)
	if err != nil {
		return fmt.Errorf("move.Plan: %w", err)
	}

	return c.JSON(http.StatusOK, &PostMoveResponse{Moves: moves})
}
