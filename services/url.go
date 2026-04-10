package services

import (
	"database/sql"
	"errors"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/lib/pq"
)

const maxShortCodeGenRetries = 5

var (
	ErrInvalidURL    = errors.New("invalid URL")
	ErrExpiryInvalid = errors.New("expiry_seconds must be >= 0")
	ErrExpired       = errors.New("expired")
	ErrNotFound      = errors.New("not found")
)

type URLRepo interface {
	Save(string, string, *time.Time) error
	GetByOriginalURL(string) (string, *time.Time, error)
	UpdateExpiry(string, *time.Time) error
	GetByShortCode(string) (string, *time.Time, error)
}

type URLService struct {
	repo URLRepo
}

func NewURLService(repo URLRepo) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) CreateShortURL(originalURL string, expirySeconds int64) (string, error) {
	parsed, err := url.ParseRequestURI(originalURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", ErrInvalidURL
	}

	normalizedURL, err := normalizeURL(originalURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	if expirySeconds < 0 {
		return "", ErrExpiryInvalid
	}

	now := time.Now().UTC()

	var newExpiry *time.Time
	if expirySeconds > 0 {
		t := now.Add(time.Duration(expirySeconds) * time.Second)
		newExpiry = &t
	}

	existingCode, existingExpiry, err := s.repo.GetByOriginalURL(normalizedURL)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if err != sql.ErrNoRows && existingCode != "" {
		if existingExpiry != nil && existingExpiry.Before(now) {
			existingCode = ""
		} else {
			if err := s.repo.UpdateExpiry(existingCode, newExpiry); err != nil {
				return "", err
			}
			return existingCode, nil
		}
	}

	for i := 0; i < maxShortCodeGenRetries; i++ {
		code := generateCode()
		err := s.repo.Save(code, normalizedURL, newExpiry)
		if err == nil {
			return code, nil
		}
		if isUniqueViolation(err) {
			continue
		}
		return "", err
	}

	return "", errors.New("could not generate unique short code")
}

func (s *URLService) GetOriginalURL(code string) (string, error) {
	originalURL, expiresAt, err := s.repo.GetByShortCode(code)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}

	now := time.Now().UTC()

	if expiresAt != nil && expiresAt.Before(now) {
		return "", ErrExpired
	}

	return originalURL, nil
}

func normalizeURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	u.Host = strings.ToLower(u.Host)
	if u.Path == "" {
		u.Path = "/"
	}
	return u.Scheme + "://" + u.Host + u.Path + "?" + u.RawQuery, nil
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func isUniqueViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23505"
}