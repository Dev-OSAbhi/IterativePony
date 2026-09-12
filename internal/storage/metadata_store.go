package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MetadataStore handles MongoDB storage for configuration states, optimization reports, and agent metadata.
type MetadataStore struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMetadataStore creates a new MetadataStore with the given MongoDB URI and database name.
func NewMetadataStore(uri string, dbName string) (*MetadataStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)

	store := &MetadataStore{
		client: client,
		db:     db,
	}
	return store, nil
}

// Init ensures required collections exist (MongoDB creates them lazily, but we can validate access).
func (s *MetadataStore) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List collection names to verify we can access the database
	collections, err := s.db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}
	// Optionally, we could create collections here, but MongoDB creates them on first insert.
	_ = collections // unused for now, but we verified access
	return nil
}

// Close closes the MongoDB client connection.
func (s *MetadataStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

// ConfigState represents a configuration state document.
type ConfigState struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Key       string    `bson:"key" json:"key"`
	Value     string    `bson:"value" json:"value"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// InsertConfigState inserts or updates a configuration state.
func (s *MetadataStore) InsertConfigState(cs ConfigState) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cs.UpdatedAt = time.Now()
	filter := bson.M{"key": cs.Key}
	update := bson.M{
		"$set": bson.M{
			"value":     cs.Value,
			"updated_at": cs.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id": cs.ID,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := s.db.Collection("config_states").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to insert config state: %w", err)
	}
	return nil
}

// GetConfigState retrieves a configuration state by key.
func (s *MetadataStore) GetConfigState(key string) (*ConfigState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result ConfigState
	filter := bson.M{"key": key}
	err := s.db.Collection("config_states").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get config state: %w", err)
	}
	return &result, nil
}

// OptimizationReport represents an optimization report document.
type OptimizationReport struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	JobID       string    `bson:"job_id" json:"job_id"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	Recommendations []string `bson:"recommendations" json:"recommendations"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

// InsertOptimizationReport inserts an optimization report.
func (s *MetadataStore) InsertOptimizationReport(orp OptimizationReport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orp.CreatedAt = time.Now()
	_, err := s.db.Collection("optimization_reports").InsertOne(ctx, orp)
	if err != nil {
		return fmt.Errorf("failed to insert optimization report: %w", err)
	}
	return nil
}

// GetOptimizationReports returns reports for a given job ID, sorted by timestamp descending.
func (s *MetadataStore) GetOptimizationReports(jobID string, limit int64) ([]OptimizationReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit)

	filter := bson.M{"job_id": jobID}
	cursor, err := s.db.Collection("optimization_reports").Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find optimization reports: %w", err)
	}
	defer cursor.Close(ctx)

	var reports []OptimizationReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, fmt.Errorf("failed to decode optimization reports: %w", err)
	}
	return reports, nil
}

// AgentMetadata represents metadata for a backup agent/node.
type AgentMetadata struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	AgentID     string    `bson:"agent_id" json:"agent_id"`
	Hostname    string    `bson:"hostname" json:"hostname"`
	IPAddress   string    `bson:"ip_address" json:"ip_address"`
	Status      string    `bson:"status" json:"status"` // e.g., "online", "offline", "maintenance"
	LastSeen    time.Time `bson:"last_seen" json:"last_seen"`
	CapacityGB  int64     `bson:"capacity_gb" json:"capacity_gb"`
	Utilization float64   `bson:"utilization" json:"utilization"` // 0.0 to 1.0
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// UpsertAgentMetadata inserts or updates agent metadata.
func (s *MetadataStore) UpsertAgentMetadata(am AgentMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	am.UpdatedAt = time.Now()
	filter := bson.M{"agent_id": am.AgentID}
	update := bson.M{"$set": am}
	opts := options.Update().SetUpsert(true)

	_, err := s.db.Collection("agent_metadata").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert agent metadata: %w", err)
	}
	return nil
}

// GetAgentMetadata returns metadata for a given agent ID.
func (s *MetadataStore) GetAgentMetadata(agentID string) (*AgentMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result AgentMetadata
	filter := bson.M{"agent_id": agentID}
	err := s.db.Collection("agent_metadata").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get agent metadata: %w", err)
	}
	return &result, nil
}

// NetworkTopology represents the network graph structure for path optimization.
type NetworkTopology struct {
	ID      string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string    `bson:"name" json:"name"` // e.g., "primary_backup_network"
	Version int       `bson:"version" json:"version"`
	Updated time.Time `bson:"updated" json:"updated"`
}

// NetworkNode represents a node in the network topology.
type NetworkNode struct {
	ID   string `bson:"id" json:"id"` // e.g., "agent-001", "storage-node-a", "switch-1"
	Name string `bson:"name" json:"name"` // Human-readable name
}

// NetworkEdge represents a connection between two nodes in the network topology.
type NetworkEdge struct {
	From   string  `bson:"from" json:"from"` // Source node ID
	To     string  `bson:"to" json:"to"`     // Target node ID
	Weight float64 `bson:"weight" json:"weight"` // Latency/cost/distance metric
}

// UpsertNetworkTopology stores or updates the network topology.
func (s *MetadataStore) UpsertNetworkTopology(topology NetworkTopology) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	topology.Updated = time.Now()
	filter := bson.M{"name": topology.Name}
	update := bson.M{
		"$set": bson.M{
			"name":     topology.Name,
			"version":  topology.Version,
			"updated":  topology.Updated,
		},
		"$setOnInsert": bson.M{
			"_id": topology.ID,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := s.db.Collection("network_topologies").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert network topology: %w", err)
	}
	return nil
}

// UpsertNetworkNode stores or updates a network node.
func (s *MetadataStore) UpsertNetworkNode(node NetworkNode) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": node.ID}
	update := bson.M{"$set": node}
	opts := options.Update().SetUpsert(true)

	_, err := s.db.Collection("network_nodes").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert network node: %w", err)
	}
	return nil
}

// UpsertNetworkEdge stores or updates a network edge.
func (s *MetadataStore) UpsertNetworkEdge(edge NetworkEdge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a composite key for the edge
	filter := bson.M{"from": edge.From, "to": edge.To}
	update := bson.M{"$set": edge}
	opts := options.Update().SetUpsert(true)

	_, err := s.db.Collection("network_edges").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert network edge: %w", err)
	}
	return nil
}

// GetNetworkTopology retrieves the named network topology.
func (s *MetadataStore) GetNetworkTopology(name string) (*NetworkTopology, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result NetworkTopology
	filter := bson.M{"name": name}
	err := s.db.Collection("network_topologies").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get network topology: %w", err)
	}
	return &result, nil
}

// GetNetworkNodes retrieves all nodes for a given topology name.
func (s *MetadataStore) GetNetworkNodes(topologyName string) ([]NetworkNode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First get the topology to verify it exists
	topology, err := s.GetNetworkTopology(topologyName)
	if err != nil {
		return nil, err
	}
	if topology == nil {
		return nil, nil
	}

	var nodes []NetworkNode
	cursor, err := s.db.Collection("network_nodes").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get network nodes: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &nodes); err != nil {
		return nil, fmt.Errorf("failed to decode network nodes: %w", err)
	}
	return nodes, nil
}

// GetNetworkEdges retrieves all edges for a given topology name.
func (s *MetadataStore) GetNetworkEdges(topologyName string) ([]NetworkEdge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First get the topology to verify it exists
	topology, err := s.GetNetworkTopology(topologyName)
	if err != nil {
		return nil, err
	}
	if topology == nil {
		return nil, nil
	}

	var edges []NetworkEdge
	cursor, err := s.db.Collection("network_edges").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get network edges: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &edges); err != nil {
		return nil, fmt.Errorf("failed to decode network edges: %w", err)
	}
	return edges, nil
}