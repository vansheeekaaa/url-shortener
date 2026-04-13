package models

import "time"

type ShortenRequest struct {
	URL           string `json:"url"`
	ExpirySeconds int64  `json:"expiry_seconds"`
}

type ShortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}

type URLStats struct {
	ShortCode      string     `json:"short_code"`
	OriginalURL    string     `json:"original_url"`
	ClickCount     int        `json:"click_count"`
	CreatedAt      time.Time  `json:"created_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
}