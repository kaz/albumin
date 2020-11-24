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
		Strategy string
		Layout   string
		Execute  bool
	}
	PostMoveResponse struct {
		Moves    []*move.Move
		Executed bool
	}
)

func PostMove(c echo.Context) error {
	req := &PostMoveRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Bind: %w", err))
	}

	var strategy move.Strategy
	switch req.Strategy {
	case "exif":
		strategy = move.StrategyExif(req.Layout)
	case "seq":
		strategy = move.StrategySequential(req.Layout)
	default:
		return fmt.Errorf("no such strategy: %v", req.Strategy)
	}

	photos := c.Get("photos").([]*model.Photo)
	moves, err := move.Plan(photos, strategy)
	if err != nil {
		return fmt.Errorf("move.Plan: %w", err)
	}

	if req.Execute {
		if err := move.Execute(moves); err != nil {
			return fmt.Errorf("move.Execute: %w", err)
		}
	}

	return c.JSON(http.StatusOK, &PostMoveResponse{Moves: moves, Executed: req.Execute})
}
