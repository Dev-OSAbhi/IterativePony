package config_test

import (
	"os"
	"testing"

	"iterative-pony/config"
)

func TestLoad_WithEnvVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "1234")
	os.Setenv("SQLITE_DB_PATH", "./test.sqlite")
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DATABASE", "testdb")
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("SQLITE_DB_PATH")
		os.Unsetenv("MONGODB_URI")
		os.Unsetenv("MONGODB_DATABASE")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("ALLOWED_ORIGINS")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "1234" {
		t.Errorf("Load().Port = %v; want %v", cfg.Port, "1234")
	}
	if cfg.SQLiteDBPath != "./test.sqlite" {
		t.Errorf("Load().SQLiteDBPath = %v; want %v", cfg.SQLiteDBPath, "./test.sqlite")
	}
	if cfg.MongoDBURI != "mongodb://localhost:27017" {
		t.Errorf("Load().MongoDBURI = %v; want %v", cfg.MongoDBURI, "mongodb://localhost:27017")
	}
	if cfg.MongoDBDatabase != "testdb" {
		t.Errorf("Load().MongoDBDatabase = %v; want %v", cfg.MongoDBDatabase, "testdb")
	}
	if cfg.Environment != "test" {
		t.Errorf("Load().Environment = %v; want %v", cfg.Environment, "test")
	}
	expectedOrigins := []string{"http://localhost:3000", "http://localhost:8080"}
	if len(cfg.AllowedOrigins) != len(expectedOrigins) {
		t.Errorf("Load().AllowedOrigins = %v; want %v", cfg.AllowedOrigins, expectedOrigins)
	}
	for i, v := range cfg.AllowedOrigins {
		if v != expectedOrigins[i] {
			t.Errorf("Load().AllowedOrigins[%d] = %v; want %v", i, v, expectedOrigins[i])
		}
	}
}

func TestLoad_WithDefaults(t *testing.T) {
	// Unset relevant env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("SQLITE_DB_PATH")
	os.Unsetenv("MONGODB_URI")
	os.Unsetenv("MONGODB_DATABASE")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("ALLOWED_ORIGINS")
	// Also ensure no .env file influences (we can't easily delete, but we assume none)
	defer func() {
		// Reset to original? Not needed for test.
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("Load().Port = %v; want %v", cfg.Port, "9090")
	}
	if cfg.SQLiteDBPath != "./metrics.db" {
		t.Errorf("Load().SQLiteDBPath = %v; want %v", cfg.SQLiteDBPath, "./metrics.db")
	}
	if cfg.MongoDBURI != "mongodb://localhost:27017" {
		t.Errorf("Load().MongoDBURI = %v; want %v", cfg.MongoDBURI, "mongodb://localhost:27017")
	}
	if cfg.MongoDBDatabase != "backup_monitor" {
		t.Errorf("Load().MongoDBDatabase = %v; want %v", cfg.MongoDBDatabase, "backup_monitor")
	}
	if cfg.Environment != "development" {
		t.Errorf("Load().Environment = %v; want %v", cfg.Environment, "development")
	}
	expectedOrigins := []string{"http://localhost:3000"}
	if len(cfg.AllowedOrigins) != len(expectedOrigins) {
		t.Errorf("Load().AllowedOrigins = %v; want %v", cfg.AllowedOrigins, expectedOrigins)
	}
	for i, v := range cfg.AllowedOrigins {
		if v != expectedOrigins[i] {
			t.Errorf("Load().AllowedOrigins[%d] = %v; want %v", i, v, expectedOrigins[i])
		}
	}
}