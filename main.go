package main

import(
	"github.com/gin-gonic/gin"
	"urlshortener/handlers"
	"urlshortener/db"
)

func main() {
	router := gin.Default()

	router.POST("/shorten", handlers.ShortenURL)
	router.GET("/:shortcode", handlers.RedirectURL)
	database:=db.InitDB()
	_ = database

	router.Run()
}