package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"iterative-pony/config"
	"iterative-pony/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize SQLite metrics store
	sqliteStore, err := storage.NewMetricsStore(cfg.SQLiteDBPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite store: %v", err)
	}
	defer sqliteStore.Close()

	// Initialize database schema
	if err := sqliteStore.Init(); err != nil {
		log.Fatalf("Failed to initialize SQLite schema: %v", err)
	}

	// Initialize MongoDB metadata store
	mongoStore, err := storage.NewMetadataStore(cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB store: %v", err)
	}
	defer mongoStore.Close()

	// Initialize MongoDB (verify access)
	if err := mongoStore.Init(); err != nil {
		log.Fatalf("Failed to initialize MongoDB store: %v", err)
	}

	// Set up Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"environment": cfg.Environment,
			"sqlite":  "connected",
			"mongo":   "connected",
		})
	})

	// Run server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}