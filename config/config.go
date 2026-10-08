package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	SQLiteDBPath string
	MongoDBURI      string
	MongoDBDatabase string
	Environment string
	AllowedOrigins []string
}

func Load() (*Config, error) {
	godotenv.Load()

	config := &Config{
		Port: getEnv("PORT", "9090"),
		SQLiteDBPath: getEnv("SQLITE_DB_PATH", "./metrics.db"),
		MongoDBURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBDatabase: getEnv("MONGODB_DATABASE", "backup_monitor"),
		Environment: getEnv("ENVIRONMENT", "development"),
		AllowedOrigins: parseEnvList(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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