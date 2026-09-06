package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"iterative-pony/internal/simulation"
)

// Handler holds dependencies for the HTTP handlers.
type Handler struct {
	simulator *simulation.BackupJobSimulator
}

// NewHandler creates a new Handler with the given simulator.
func NewHandler(sim *simulation.BackupJobSimulator) *Handler {
	return &Handler{simulator: sim}
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