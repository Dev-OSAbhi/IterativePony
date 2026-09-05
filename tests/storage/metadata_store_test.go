package storage_test

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"iterative-pony/internal/storage"
)

func TestMetadataStore_CRUD(t *testing.T) {
	// Use a separate test database to avoid interfering with other data
	mongoURI := "mongodb://localhost:27017"
	testDB := "backup_monitor_test_" + primitive.NewObjectID().Hex()

	store, err := storage.NewMetadataStore(mongoURI, testDB)
	if err != nil {
		t.Fatalf("NewMetadataStore error: %v", err)
	}
	defer store.Close() // Ensure close

	if err := store.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}

	// Test ConfigState
	cs := storage.ConfigState{
		Key:   "test_key",
		Value: "test_value",
	}
	if err := store.InsertConfigState(cs); err != nil {
		t.Fatalf("InsertConfigState error: %v", err)
	}
	fetched, err := store.GetConfigState("test_key")
	if err != nil {
		t.Fatalf("GetConfigState error: %v", err)
	}
	if fetched == nil {
		t.Fatalf("GetConfigState returned nil for existing key")
	}
	if fetched.Value != "test_value" {
		t.Errorf("GetConfigState.Value = %v; want test_value", fetched.Value)
	}
	// Update
	cs.Value = "updated_value"
	if err := store.InsertConfigState(cs); err != nil {
		t.Fatalf("InsertConfigState update error: %v", err)
	}
	fetched2, err := store.GetConfigState("test_key")
	if err != nil {
		t.Fatalf("GetConfigState error after update: %v", err)
	}
	if fetched2 == nil {
		t.Fatalf("GetConfigState returned nil after update")
	}
	if fetched2.Value != "updated_value" {
		t.Errorf("GetConfigState.Value after update = %v; want updated_value", fetched2.Value)
	}
	// Non-existent key
	nonexistent, err := store.GetConfigState("nonexistent")
	if err != nil {
		t.Fatalf("GetConfigState error for nonexistent: %v", err)
	}
	if nonexistent != nil {
		t.Errorf("GetConfigState for nonexistent returned %v; want nil", nonexistent)
	}

	// Test OptimizationReport
	orp := storage.OptimizationReport{
		JobID:         "job123",
		Timestamp:     time.Now(),
		Recommendations: []string{"rec1", "rec2"},
	}
	if err := store.InsertOptimizationReport(orp); err != nil {
		t.Fatalf("InsertOptimizationReport error: %v", err)
	}
	reports, err := store.GetOptimizationReports("job123", 10)
	if err != nil {
		t.Fatalf("GetOptimizationReports error: %v", err)
	}
	if len(reports) != 1 {
		t.Errorf("GetOptimizationReports len = %v; want 1", len(reports))
	}
	if reports[0].JobID != "job123" {
		t.Errorf("GetOptimizationReports[0].JobID = %v; want job123", reports[0].JobID)
	}
	if len(reports[0].Recommendations) != 2 {
		t.Errorf("GetOptimizationReports[0].Recommendations len = %v; want 2", len(reports[0].Recommendations))
	}
	// Ensure timestamp is set (not zero)
	if reports[0].Timestamp.IsZero() {
		t.Errorf("OptimizationReport timestamp is zero")
	}
	// Ensure CreatedAt is set
	if reports[0].CreatedAt.IsZero() {
		t.Errorf("OptimizationReport CreatedAt is zero")
	}

	// Test AgentMetadata
	am := storage.AgentMetadata{
		AgentID:     "agent001",
		Hostname:    "host1",
		IPAddress:   "10.0.0.1",
		Status:      "online",
		CapacityGB:  1000,
		Utilization: 0.5,
	}
	if err := store.UpsertAgentMetadata(am); err != nil {
		t.Fatalf("UpsertAgentMetadata error: %v", err)
	}
	fetchedAM, err := store.GetAgentMetadata("agent001")
	if err != nil {
		t.Fatalf("GetAgentMetadata error: %v", err)
	}
	if fetchedAM == nil {
		t.Fatalf("GetAgentMetadata returned nil for existing agent")
	}
	if fetchedAM.Hostname != "host1" {
		t.Errorf("GetAgentMetadata.Hostname = %v; want host1", fetchedAM.Hostname)
	}
	if fetchedAM.Utilization != 0.5 {
		t.Errorf("GetAgentMetadata.Utilization = %v; want 0.5", fetchedAM.Utilization)
	}
	// Update
	am.Status = "maintenance"
	if err := store.UpsertAgentMetadata(am); err != nil {
		t.Fatalf("UpsertAgentMetadata update error: %v", err)
	}
	fetchedAM2, err := store.GetAgentMetadata("agent001")
	if err != nil {
		t.Fatalf("GetAgentMetadata error after update: %v", err)
	}
	if fetchedAM2 == nil {
		t.Fatalf("GetAgentMetadata returned nil after update")
	}
	if fetchedAM2.Status != "maintenance" {
		t.Errorf("GetAgentMetadata.Status after update = %v; want maintenance", fetchedAM2.Status)
	}
	// Non-existent agent
	nonexistentAM, err := store.GetAgentMetadata("nonexistent")
	if err != nil {
		t.Fatalf("GetAgentMetadata error for nonexistent: %v", err)
	}
	if nonexistentAM != nil {
		t.Errorf("GetAgentMetadata for nonexistent returned %v; want nil", nonexistentAM)
	}
}