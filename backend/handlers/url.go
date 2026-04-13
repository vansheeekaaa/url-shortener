package handlers

import (
	"errors"
	"net/http"
	"os"

	"urlshortener/models"
	"urlshortener/services"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *services.URLService
	baseURL string
}

func NewURLHandler(service *services.URLService) *URLHandler {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &URLHandler{
		service: service,
		baseURL: baseURL,
	}
}

// POST /shorten
func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	shortCode, err := h.service.CreateShortURL(req.URL, req.ExpirySeconds)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidURL):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})

		case errors.Is(err, services.ErrExpiryInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": "expiry_seconds must be >= 0"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, models.ShortenResponse{
		ShortURL:  h.baseURL + "/" + shortCode,
		ShortCode: shortCode,
	})
}

// GET /:code
func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.service.GetOriginalURL(code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrExpired):
			c.JSON(http.StatusGone, gin.H{"error": "link expired"})

		case errors.Is(err, services.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}

// GET /stats/:code
func (h *URLHandler) GetStats(c *gin.Context) {
	code := c.Param("code")

	stats, err := h.service.GetStats(code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, stats)
}