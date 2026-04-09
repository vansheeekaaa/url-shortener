package services

import (
	"math/rand"
	"time"
	"github.com/lib/pq"
	"errors"

	"urlshortener/repository"
)

const maxShortCodeGenRetries = 5

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) CreateShortURL(originalURL string) (string, error) {
	var shortCode string

	for i:= 0; i<maxShortCodeGenRetries; i++ {
		shortCode = generateCode()
		err := s.repo.Save(shortCode, originalURL)
		if err == nil {
			return shortCode, nil
		}
		if isUniqueViolation(err) {
			continue 
		}
		return "", err
	}

	return "",errors.New("could not generate unique short code after retries")
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

func isUniqueViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23505"
}

func (s *URLService) GetOriginalURL(code string) (string, error) {
	return s.repo.GetOriginalURL(code)
}

