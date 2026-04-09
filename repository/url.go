package repository

import "database/sql"

type URLRepository struct {
	DB *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (r *URLRepository) Save(shortCode string, originalURL string) error {
	query := `INSERT INTO urls (short_code, original_url) VALUES ($1, $2)`

	println("INSERTING INTO DB:", shortCode, originalURL)

	_, err := r.DB.Exec(query, shortCode, originalURL)
	return err
}

func (r *URLRepository) GetOriginalURL(shortCode string) (string, error) {
	query := `
		SELECT original_url FROM urls
		WHERE short_code = $1
	`

	var originalURL string
	err := r.DB.QueryRow(query, shortCode).Scan(&originalURL)

	return originalURL, err
}