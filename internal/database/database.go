package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

// NewDB attempts to connect to the postgres database with the images
func NewDB(ctx context.Context) (*sql.DB, error) {
	connStr := "user=postgres password=password dbname=images sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

