package main

import (
	"log"
	"urlshortener/db"
	"urlshortener/handlers"
	"urlshortener/repository"
	"urlshortener/services"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database := db.InitDB()

	urlRepo := repository.NewURLRepository(database)
	urlService := services.NewURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	r.POST("/shorten", urlHandler.ShortenURL)
	r.GET("/:code", urlHandler.Redirect)

	r.Run(":8080")
}