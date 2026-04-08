package main

import(
	"github.com/gin-gonic/gin"
	"urlshortener/handlers"
)

func main() {
	router := gin.Default()

	router.POST("/shorten", handlers.ShortenURL)
	router.GET("/:shortcode", handlers.RedirectURL)

	router.Run()
}