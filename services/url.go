package services

import (
	"math/rand"
	"time"

	"urlshortener/repository"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) CreateShortURL(originalURL string) (string, error) {
	shortCode := generateCode()

	err := s.repo.Save(shortCode, originalURL)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	shortCode := make([]byte, 6)
	for i := range shortCode {
		shortCode[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortCode)
}

func (s *URLService) GetOriginalURL(code string) (string, error) {
	return s.repo.GetOriginalURL(code)
}

