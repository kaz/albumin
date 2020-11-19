package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kaz/albumin/preference"
	_ "github.com/mattn/go-sqlite3"
)

type (
	Model struct {
		db *sqlx.DB
	}
)

func New(dsn string) (*Model, error) {
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}
	return &Model{db: db}, nil
}

func Default() (*Model, error) {
	return New(preference.DatabasePath)
}
