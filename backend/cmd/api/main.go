package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	router.Run(":8080")
}
