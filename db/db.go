package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

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