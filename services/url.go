package services 

import (
	"math/rand"
	"time"
	storage "urlshortener/storage"
)

//generate shorturl from url
func CreateShortURL(URL string) string {
	shortCode := generateCode()
	storage.URL[shortCode] = URL
	return shortCode
}

//generate random shortcode
func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	shortCode := make([]byte, 6)

	for i:= range shortCode {
		shortCode[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortCode)
}

//fetch original url
func GetOriginalURL(code string) (string, bool) {
	url, exists:=storage.URL[code]
	return url, exists
}
