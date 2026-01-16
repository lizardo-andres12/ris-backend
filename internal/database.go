package internal

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// ConnectDB attempts to connect to the postgres database with the images
func ConnectDB() (*sql.DB, error) {
	connStr := "user=postgres password=password dbname=images sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

