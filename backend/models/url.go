package models

type ShortenRequest struct {
	URL           string `json:"url"`
	ExpirySeconds int64  `json:"expiry_seconds"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}