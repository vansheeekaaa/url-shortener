package handlers

import (
	"errors"
	"log"
	"net/http"

	"urlshortener/models"
	"urlshortener/services"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	shortCode, err := h.service.CreateShortURL(req.URL, req.ExpirySeconds)
	if err != nil {
		log.Println("Service error:", err)

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
		ShortURL: shortCode,
	})
}

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

	c.Redirect(http.StatusFound, originalURL)
}