package main

import (
	"log"
	"urlshortener/db"
	"urlshortener/handlers"
	"urlshortener/repository"
	"urlshortener/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	_ = godotenv.Load()

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	urlRepo := repository.NewURLRepository(database)
	urlService := services.NewURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	r.POST("/shorten", urlHandler.ShortenURL)

	// /stats/:code and /ping must be registered before /:code so Gin's
	// radix tree prefers specific static prefixes over the wildcard param.
	r.GET("/stats/:code", urlHandler.GetStats)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong", "status": "alive"})
	})
	r.GET("/:code", urlHandler.Redirect)

	r.Run(":8080")
}