package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mak-magz/url-shortener/internal/url/handler"
	"github.com/mak-magz/url-shortener/internal/url/repository"
	"github.com/mak-magz/url-shortener/internal/url/service"
	"github.com/mak-magz/url-shortener/platform/config"
	"github.com/mak-magz/url-shortener/platform/db"
)

func main() {
	fmt.Println("🚀 URL Shortener API starting...")

	// Load environment variables
	cfg := config.LoadConfig()

	// Connect to database
	pool := db.Connect(cfg.DatabaseURL)
	defer pool.Close()

	// Run database migrations
	db.Migrate(pool)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	urlRepo := repository.NewPostgresRepository(pool)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	router.POST("/api/v1/shorten", urlHandler.CreateShortURL)

	router.GET("/:shortCode", urlHandler.RedirectToOriginalURL)

	router.Run(":8080")
}
