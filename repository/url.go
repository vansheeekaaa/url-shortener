package repository

import (
	"database/sql"
	"time"
)

type URLRepository struct {
	DB *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (r *URLRepository) Save(shortCode, originalURL string, expiresAt *time.Time) error {
	query := `INSERT INTO urls (short_code, original_url, expires_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, shortCode, originalURL, expiresAt)
	return err
}

func (r *URLRepository) GetByOriginalURL(originalURL string) (string, *time.Time, error) {
	var code string
	var expiresAt sql.NullTime

	query := `
		SELECT short_code, expires_at
		FROM urls
		WHERE original_url = $1
		ORDER BY created_at DESC
		LIMIT 1
		`

	err := r.DB.QueryRow(query, originalURL).Scan(&code, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, sql.ErrNoRows
		}
		return "", nil, err
	}

	var expiryPtr *time.Time
	if expiresAt.Valid {
		expiryPtr = &expiresAt.Time
	}

	return code, expiryPtr, nil
}

func (r *URLRepository) GetByShortCode(code string) (string, *time.Time, error) {
	var originalURL string
	var expiresAt sql.NullTime

	query := `SELECT original_url, expires_at FROM urls WHERE short_code = $1`
	err := r.DB.QueryRow(query, code).Scan(&originalURL, &expiresAt)
	if err != nil {
		return "", nil, err
	}

	var expiryPtr *time.Time
	if expiresAt.Valid {
		expiryPtr = &expiresAt.Time
	}

	return originalURL, expiryPtr, nil
}