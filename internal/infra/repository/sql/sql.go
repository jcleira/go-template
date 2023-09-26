// Package sql opens a connection to the database
package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/jcleira/go-template/config"
)

// New opens a connection to the database and returns a sqlx.DB.
func New(cfg config.DB) (*sqlx.DB, error) {
	db, err := sql.Open("postgres", cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	dbx := sqlx.NewDb(db, "postgres")

	dbx.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	dbx.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	dbx.SetMaxIdleConns(cfg.MaxIdleConns)
	dbx.SetMaxOpenConns(cfg.MaxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
	defer cancel()

	if err := dbx.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("dbx.PingContext: %w", err)
	}

	return dbx, nil
}
