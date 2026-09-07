package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"iterative-pony/internal/simulation"
	"iterative-pony/internal/storage"
)

// Handler holds dependencies for the HTTP handlers.
type Handler struct {
	simulator    *simulation.BackupJobSimulator
	metricsStore *storage.MetricsStore
	metaStore    *storage.MetadataStore
}

// NewHandler creates a new Handler with the given dependencies.
func NewHandler(sim *simulation.BackupJobSimulator, ms *storage.MetricsStore, mds *storage.MetadataStore) *Handler {
	return &Handler{
		simulator:    sim,
		metricsStore: ms,
		metaStore:    mds,
	}
}

// StartSimulation handles POST /simulate/start to start the backup job simulation.
func (h *Handler) StartSimulation(c *gin.Context) {
	h.simulator.Start()
	c.JSON(http.StatusOK, gin.H{"status": "simulation started"})
}

// StopSimulation handles POST /simulate/stop to stop the backup job simulation.
func (h *Handler) StopSimulation(c *gin.Context) {
	log.Println("StopSimulation called")
	h.simulator.Stop()
	c.JSON(http.StatusOK, gin.H{"status": "simulation stopped"})
}

// SimulationStatus handles GET /simulate/status to get the current simulation status.
func (h *Handler) SimulationStatus(c *gin.Context) {
	if h.simulator != nil {
		c.JSON(http.StatusOK, gin.H{"running": h.simulator.IsRunning()})
	} else {
		c.JSON(http.StatusOK, gin.H{"running": false})
	}
}

// GetMetrics handles GET /metrics to retrieve metrics for a given job ID.
func (h *Handler) GetMetrics(c *gin.Context) {
	jobID := c.Query("jobID")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jobID is required"})
		return
	}
	limitStr := c.DefaultQuery("limit", "100")
	limit := 100
	if _, err := fmt.Sscan(limitStr, &limit); err != nil || limit <= 0 {
		limit = 100
	}
	metrics, err := h.metricsStore.GetMetrics(jobID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}

// GetAgentMetadata handles GET /agents/:agentID to retrieve metadata for a specific agent.
func (h *Handler) GetAgentMetadata(c *gin.Context) {
	agentID := c.Param("agentID")
	if agentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agentID is required"})
		return
	}
	metadata, err := h.metaStore.GetAgentMetadata(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if metadata == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agent": metadata})
}