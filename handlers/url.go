package handlers

import(
	"net/http"
	"github.com/gin-gonic/gin"

	"urlshortener/models"
	"urlshortener/services"
)

func ShortenURL(c *gin.Context) {
	//mapping
	request := models.ShortenRequest{}
	err := c.BindJSON(&request)

	if(err!=nil) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "invalid request",
		})
		return 
	}

	//validate
	if(request.URL == "") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "no url found",
		})
		return
	}

	shortCode := services.CreateShortURL(request.URL)

	c.JSON(http.StatusOK, models.ShortenResponse{
		ShortURL : shortCode,
	})
}

func RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortcode")

	url, exists := services.GetOriginalURL(shortCode)

	if(!exists) {
		c.JSON(http.StatusNotFound, gin.H{
			"error" : "not found",
		})
		return
	}
	c.Redirect(http.StatusFound, url)
}