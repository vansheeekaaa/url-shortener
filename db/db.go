package db

import (
	"database/sql"
	"os"
	"errors"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, errors.New("DATABASE_URL not set")
	}

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}