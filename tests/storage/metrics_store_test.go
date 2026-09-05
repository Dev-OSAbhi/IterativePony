package storage_test

import (
	"path/filepath"
	"testing"
	"time"

	"iterative-pony/internal/storage"
)

func TestMetricsStore_InitAndInsertGet(t *testing.T) {
	// Create a temporary file for SQLite database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := storage.NewMetricsStore(dbPath)
	if err != nil {
		t.Fatalf("NewMetricsStore error: %v", err)
	}
	defer store.Close()

	// Initialize schema
	if err := store.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}

	// Insert a metric
	metric := storage.Metric{
		JobID:     "job1",
		Timestamp: time.Now(),
		LatencyMS: 10.5,
		ThroughputMbps: 100.0,
		ErrorRate: 0.01,
		Status:    "success",
	}
	id, err := store.InsertMetric(metric)
	if err != nil {
		t.Fatalf("InsertMetric error: %v", err)
	}
	if id <= 0 {
		t.Errorf("InsertMetric returned invalid ID: %v", id)
	}

	// Retrieve metrics
	metrics, err := store.GetMetrics("job1", 10)
	if err != nil {
		t.Fatalf("GetMetrics error: %v", err)
	}
	if len(metrics) != 1 {
		t.Errorf("GetMetrics len = %v; want 1", len(metrics))
	}
	if metrics[0].JobID != "job1" {
		t.Errorf("GetMetrics[0].JobID = %v; want job1", metrics[0].JobID)
	}
	if metrics[0].LatencyMS != 10.5 {
		t.Errorf("GetMetrics[0].LatencyMS = %v; want 10.5", metrics[0].LatencyMS)
	}
	if metrics[0].ThroughputMbps != 100.0 {
		t.Errorf("GetMetrics[0].ThroughputMbps = %v; want 100.0", metrics[0].ThroughputMbps)
	}
	if metrics[0].ErrorRate != 0.01 {
		t.Errorf("GetMetrics[0].ErrorRate = %v; want 0.01", metrics[0].ErrorRate)
	}
	if metrics[0].Status != "success" {
		t.Errorf("GetMetrics[0].Status = %v; want success", metrics[0].Status)
	}
}

func TestMetricsStore_GetMetrics_Limit(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test2.db")

	store, err := storage.NewMetricsStore(dbPath)
	if err != nil {
		t.Fatalf("NewMetricsStore error: %v", err)
	}
	defer store.Close()
	if err := store.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}

	// Insert three metrics with same jobID but different timestamps
	now := time.Now()
	for i := 0; i < 3; i++ {
		metric := storage.Metric{
			JobID:     "jobX",
			Timestamp: now.Add(time.Duration(-i) * time.Hour), // older as i increases
			LatencyMS: float64(i * 10),
			ThroughputMbps: float64(100 - i*10),
			ErrorRate: 0.0,
			Status:    "success",
		}
		if _, err := store.InsertMetric(metric); err != nil {
			t.Fatalf("InsertMetric error: %v", err)
		}
	}

	// Get with limit 2
	metrics, err := store.GetMetrics("jobX", 2)
	if err != nil {
		t.Fatalf("GetMetrics error: %v", err)
	}
	if len(metrics) != 2 {
		t.Errorf("GetMetrics limit 2 len = %v; want 2", len(metrics))
	}
	// Should be most recent first (i=0 then i=1)
	if metrics[0].LatencyMS != 0.0 { // i=0
		t.Errorf("First metric latency = %v; want 0.0", metrics[0].LatencyMS)
	}
	if metrics[1].LatencyMS != 10.0 { // i=1
		t.Errorf("Second metric latency = %v; want 10.0", metrics[1].LatencyMS)
	}
}

func TestMetricsStore_GetMetrics_UnknownJob(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test3.db")

	store, err := storage.NewMetricsStore(dbPath)
	if err != nil {
		t.Fatalf("NewMetricsStore error: %v", err)
	}
	defer store.Close()
	if err := store.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}

	metrics, err := store.GetMetrics("unknown", 10)
	if err != nil {
		t.Fatalf("GetMetrics error: %v", err)
	}
	if len(metrics) != 0 {
		t.Errorf("GetMetrics for unknown job returned %v metrics; want 0", len(metrics))
	}
}