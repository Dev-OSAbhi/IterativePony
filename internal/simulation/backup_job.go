package simulation

import (
	"math/rand"
	"time"

	"iterative-pony/internal/storage"
)

// BackupJobSimulator simulates a backup job and generates metrics.
type BackupJobSimulator struct {
	store    *storage.MetricsStore
	jobID    string
	interval time.Duration
	running  bool
	stopCh   chan struct{}
	r        *rand.Rand
}

// NewBackupJobSimulator creates a new simulator with the given store and jobID.
func NewBackupJobSimulator(store *storage.MetricsStore, jobID string) *BackupJobSimulator {
	return &BackupJobSimulator{
		store:    store,
		jobID:    jobID,
		interval: 5 * time.Second, // generate metrics every 5 seconds
		running:  false,
		stopCh:   make(chan struct{}),
		r:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins the simulation. It creates a goroutine that generates and stores metrics until Stop is called.
func (s *BackupJobSimulator) Start() {
	if s.running {
		return
	}
	s.running = true
	go s.run()
}

// Stop signals the simulation to stop and waits for the goroutine to finish.
func (s *BackupJobSimulator) Stop() {
	if !s.running {
		// fmt.Printf("Stop: returning early because running is false\n")
		return
	}
	// fmt.Printf("Stop: about to close stopCh, running=%v\n", s.running)
	close(s.stopCh)
	// fmt.Printf("Stop: closed stopCh\n")
	s.running = false
}

// IsRunning returns true if the simulator is currently running.
func (s *BackupJobSimulator) IsRunning() bool {
	return s.running
}

// run is the main loop of the simulator.
func (s *BackupJobSimulator) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.generateAndStoreMetrics()
		case <-s.stopCh:
			return
		}
	}
}

// generateAndStoreMetrics creates a dummy metric and stores it in the MetricsStore.
func (s *BackupJobSimulator) generateAndStoreMetrics() {
	// Generate dummy values
	latency := s.r.Float64() * 100         // 0-100 ms
	throughput := s.r.Float64()*1000 + 100 // 100-1100 Mbps
	errorRate := s.r.Float64() * 0.1       // 0-10% error
	status := "success"
	if s.r.Float64() < errorRate {
		status = "failed"
	}

	metric := storage.Metric{
		JobID:          s.jobID,
		Timestamp:      time.Now(),
		LatencyMS:      latency,
		ThroughputMbps: throughput,
		ErrorRate:      errorRate,
		Status:         status,
	}

	if _, err := s.store.InsertMetric(metric); err != nil {
		// In a real app, we might log this error.
		// For now, we just ignore it to keep the simulation running.
	}
}
