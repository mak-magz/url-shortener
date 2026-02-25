package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerHost  string
	ServerPort  string
}

func LoadConfig() *Config {
	// Load .env file
	_ = godotenv.Load()

	config := &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ServerHost:  getEnv("SERVER_HOST", "localhost"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}

	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
