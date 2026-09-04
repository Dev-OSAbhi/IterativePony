package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	// Server configuration
	Port string

	// SQLite database configuration
	SQLiteDBPath string

	// MongoDB database configuration
	MongoDBURI string
	MongoDBDatabase string

	 // Environment (development, production, etc.)
	Environment string

	// CORS allowed origins (comma-separated)
	AllowedOrigins []string
}

// Load reads configuration from environment variables and returns a Config struct.
// It also loads a .env file if present.
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Port: getEnv("PORT", "9090"),

		SQLiteDBPath: getEnv("SQLITE_DB_PATH", "./metrics.db"),

		MongoDBURI:   getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBDatabase: getEnv("MONGODB_DATABASE", "backup_monitor"),

		Environment: getEnv("ENVIRONMENT", "development"),

		AllowedOrigins: parseEnvList(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
	}

	return config, nil
}

// getEnv gets an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// parseEnvList parses a comma-separated environment variable into a slice.
// Trims spaces and ignores empty entries.
func parseEnvList(list string) []string {
	if list == "" {
		return []string{}
	}
	result := []string{}
	for _, item := range strings.Split(list, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}