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
	r.GET("/:code", urlHandler.Redirect)

	//health check endpoint for UptimeRobot
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong", "status": "alive"})
	})

	r.Run(":8080")
}