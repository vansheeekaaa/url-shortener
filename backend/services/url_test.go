package services

import (
	"database/sql"
	"testing"
	"time"

	"urlshortener/models"
)

// ---- mock repo ----

type mockRepo struct {
	store map[string]struct {
		url        string
		expiresAt  *time.Time
		clickCount int
	}
}

func newMockRepo() *mockRepo {
	return &mockRepo{store: make(map[string]struct {
		url        string
		expiresAt  *time.Time
		clickCount int
	})}
}

func (m *mockRepo) Save(code, url string, expiry *time.Time) error {
	m.store[code] = struct {
		url        string
		expiresAt  *time.Time
		clickCount int
	}{url, expiry, 0}
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

func (m *mockRepo) GetByShortCode(code string) (string, *time.Time, error) {
	v, ok := m.store[code]
	if !ok {
		return "", nil, sql.ErrNoRows
	}
	return v.url, v.expiresAt, nil
}

func (m *mockRepo) RecordClick(code string) error {
	if v, ok := m.store[code]; ok {
		v.clickCount++
		m.store[code] = v
	}
	return nil
}

func (m *mockRepo) GetStats(code string) (*models.URLStats, error) {
	v, ok := m.store[code]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return &models.URLStats{
		ShortCode:   code,
		OriginalURL: v.url,
		ClickCount:  v.clickCount,
		CreatedAt:   time.Now(),
		ExpiresAt:   v.expiresAt,
	}, nil
}

// ---- existing tests (unchanged) ----

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

// ---- new analytics tests ----

func TestGetStats_ClickCount(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	code, err := service.CreateShortURL("https://example.com", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Simulate 3 redirects. RecordClick runs in a goroutine inside
	// GetOriginalURL, so we yield briefly after each call.
	for i := 0; i < 3; i++ {
		if _, err := service.GetOriginalURL(code); err != nil {
			t.Fatalf("unexpected redirect error on call %d: %v", i+1, err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	stats, err := service.GetStats(code)
	if err != nil {
		t.Fatalf("unexpected stats error: %v", err)
	}

	if stats.ClickCount != 3 {
		t.Fatalf("expected 3 clicks, got %d", stats.ClickCount)
	}
}

func TestGetStats_NotFound(t *testing.T) {
	repo := newMockRepo()
	service := NewURLService(repo)

	_, err := service.GetStats("doesnotexist")
	if err != ErrNotFound {
		t.Fatalf("expected not found error, got %v", err)
	}
}