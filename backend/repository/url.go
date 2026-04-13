package repository

import (
	"database/sql"
	"time"

	"urlshortener/models"
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

// RecordClick atomically increments the click counter and sets last_accessed_at.
func (r *URLRepository) RecordClick(code string) error {
	query := `UPDATE urls SET click_count = click_count + 1, last_accessed_at = NOW() WHERE short_code = $1`
	_, err := r.DB.Exec(query, code)
	return err
}

// GetStats returns the full analytics record for a short code.
func (r *URLRepository) GetStats(code string) (*models.URLStats, error) {
	var stats models.URLStats
	var expiresAt sql.NullTime
	var lastAccessedAt sql.NullTime

	query := `
		SELECT short_code, original_url, click_count, created_at, last_accessed_at, expires_at
		FROM urls
		WHERE short_code = $1
	`
	err := r.DB.QueryRow(query, code).Scan(
		&stats.ShortCode,
		&stats.OriginalURL,
		&stats.ClickCount,
		&stats.CreatedAt,
		&lastAccessedAt,
		&expiresAt,
	)
	if err != nil {
		return nil, err
	}

	if lastAccessedAt.Valid {
		stats.LastAccessedAt = &lastAccessedAt.Time
	}
	if expiresAt.Valid {
		stats.ExpiresAt = &expiresAt.Time
	}

	return &stats, nil
}