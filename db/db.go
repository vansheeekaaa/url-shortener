package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connStr := "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")
	return database
}