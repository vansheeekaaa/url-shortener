package main

import (
	"urlshortener/db"
	"urlshortener/handlers"
	"urlshortener/repository"
	"urlshortener/services"
	
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database := db.InitDB()

	urlRepo := repository.NewURLRepository(database)
	urlService := services.NewURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	r.POST("/shorten", urlHandler.ShortenURL)
	r.GET("/:code", urlHandler.Redirect)

	r.Run(":8080")
}