package config

import (
	"os"
)

// Config holds all environmental configuration needed by the application.
type Config struct {
	Port        string
	DatabaseURL string
	LogLevel    string
}

// LoadConfig fetches env variables, applying sensible defaults where necessary.
func LoadConfig() *Config {
	port := getEnv("PORT", "3000")
	// Standard PostgreSQL connection URL: postgres://<username>:<password>@<host>:<port>/<dbname>?sslmode=disable
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable")
	logLevel := getEnv("LOG_LEVEL", "info")

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		LogLevel:    logLevel,
	}
}

// getEnv helper retrieves environmental variable or returns defaultValue if not present.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
