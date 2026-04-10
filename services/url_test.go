package services

import (
	"database/sql"
	"testing"
	"time"
)

type mockRepo struct {
	store map[string]struct {
		url       string
		expiresAt *time.Time
	}
}

func newMockRepo() *mockRepo {
	return &mockRepo{store: make(map[string]struct {
		url       string
		expiresAt *time.Time
	})}
}

func (m *mockRepo) Save(code, url string, expiry *time.Time) error {
	m.store[code] = struct {
		url       string
		expiresAt *time.Time
	}{url, expiry}
	return nil
}

func (m *mockRepo) GetByOriginalURL(url string) (string, *time.Time, error) {
	for code, v := range m.store {
		if v.url == url {
			return code, v.expiresAt, nil
		}
	}
	return "", nil, sql.ErrNoRows
}

func (m *mockRepo) UpdateExpiry(code string, expiry *time.Time) error {
	v := m.store[code]
	v.expiresAt = expiry
	m.store[code] = v
	return nil
}

func (m *mockRepo) GetByShortCode(code string) (string, *time.Time, error) {
	v, ok := m.store[code]
	if !ok {
		return "", nil, sql.ErrNoRows
	}
	return v.url, v.expiresAt, nil
}

func TestNormalization_Idempotency(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	code1, err := service.CreateShortURL("https://google.com", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	code2, err := service.CreateShortURL("https://google.com/", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if code1 != code2 {
		t.Fatalf("expected same code, got %s and %s", code1, code2)
	}
}

func TestNormalization_DifferentURLs(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	code1, _ := service.CreateShortURL("https://google.com?a=1", 0)
	code2, _ := service.CreateShortURL("https://google.com?a=2", 0)

	if code1 == code2 {
		t.Fatalf("expected different codes for different query params")
	}
}

func TestExpiry_ExpiredURLCreatesNew(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	code1, _ := service.CreateShortURL("https://google.com", 1)

	time.Sleep(2 * time.Second)

	code2, _ := service.CreateShortURL("https://google.com", 1)

	if code1 == code2 {
		t.Fatalf("expected new code after expiry")
	}
}

func TestExpiry_ActiveURLReusesCode(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	code1, _ := service.CreateShortURL("https://google.com", 100)
	code2, _ := service.CreateShortURL("https://google.com", 100)

	if code1 != code2 {
		t.Fatalf("expected same code for active URL")
	}
}

func TestInvalidURL(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	_, err := service.CreateShortURL("invalid-url", 0)
	if err == nil {
		t.Fatalf("expected error for invalid URL")
	}
}

func TestNegativeExpiry(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	_, err := service.CreateShortURL("https://google.com", -1)
	if err == nil {
		t.Fatalf("expected error for negative expiry")
	}
}

func TestGetOriginalURL_Expired(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	code, _ := service.CreateShortURL("https://google.com", 1)

	time.Sleep(2 * time.Second)

	_, err := service.GetOriginalURL(code)
	if err != ErrExpired {
		t.Fatalf("expected expired error")
	}
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	_, err := service.GetOriginalURL("abcd")
	if err != ErrNotFound {
		t.Fatalf("expected not found error")
	}
}