package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"iterative-pony/config"
	"iterative-pony/internal/analysis"
	"iterative-pony/internal/optimizer"
	"iterative-pony/internal/api"
	"iterative-pony/internal/simulation"
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

	// Create backup job simulator
	simulator := simulation.NewBackupJobSimulator(sqliteStore, "backup-job-001")

	// Create analysis and optimizer instances
	analyser := analysis.NewBottleneckDetector()
	optimizer := optimizer.NewAdvancedRecommendations(mongoStore)

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

	// Create API handler
	handler := api.NewHandler(simulator, sqliteStore, mongoStore, analyser, optimizer)

	// Simulation endpoints
	simGroup := router.Group("/simulate")
	{
		simGroup.POST("/start", handler.StartSimulation)
		simGroup.POST("/stop", handler.StopSimulation)
		simGroup.GET("/status", handler.SimulationStatus)
	}

	// Metrics endpoints
	metricsGroup := router.Group("/metrics")
	{
		metricsGroup.GET("", handler.GetMetrics)
	}

	// Agents endpoints
	agentsGroup := router.Group("/agents")
	{
		agentsGroup.GET("/:agentID", handler.GetAgentMetadata)
	}

	// Analysis endpoints
	analysisGroup := router.Group("/analysis")
	{
		analysisGroup.GET("/bottlenecks", handler.GetBottlenecks)
	}

	// Optimization endpoints
	optimizationGroup := router.Group("/optimization")
	{
		optimizationGroup.GET("/recommendations", handler.GetRecommendations)
	}

	// Run server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}