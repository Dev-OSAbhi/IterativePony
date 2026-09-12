package optimizer

import (
	"iterative-pony/internal/analysis"
	"iterative-pony/internal/network"
	"iterative-pony/internal/storage"
)

// AdvancedRecommendations generates optimization suggestions based on detected bottlenecks
// and network topology information for specific, actionable advice.
type AdvancedRecommendations struct {
	metaStore *storage.MetadataStore
	pathOpt   *network.PathOptimizer
}

// NewAdvancedRecommendations creates a new AdvancedRecommendations generator.
func NewAdvancedRecommendations(metaStore *storage.MetadataStore) *AdvancedRecommendations {
	return &AdvancedRecommendations{
		metaStore: metaStore,
		pathOpt:   network.NewPathOptimizer(),
	}
}

// Generate creates optimization recommendations based on the given bottlenecks.
// Unlike the basic version, this provides specific, topology-aware advice.
func (r *AdvancedRecommendations) Generate(bottlenecks []analysis.Bottleneck) []string {
	var recs []string

	// We'll group bottlenecks by type to avoid duplicate recommendations
	seenLatency := false
	seenThroughput := false
	seenErrorRate := false

	for _, b := range bottlenecks {
		switch b.Type {
		case "latency":
			if !seenLatency {
				// Get specific path-based advice for latency issues
				pathAdvice := r.getLatencyRecommendations()
				if pathAdvice != "" {
					recs = append(recs, pathAdvice)
				} else {
					// Fallback to generic advice if no topology data
					recs = append(recs, "Consider increasing network bandwidth or optimizing data transfer paths to reduce latency.")
				}
				seenLatency = true
			}
		case "throughput":
			if !seenThroughput {
				// Get specific parallelism/storage advice for throughput issues
				throughputAdvice := r.getThroughputRecommendations()
				if throughputAdvice != "" {
					recs = append(recs, throughputAdvice)
				} else {
					// Fallback to generic advice
					recs = append(recs, "Consider increasing parallelism or adding more storage nodes to improve throughput.")
				}
				seenThroughput = true
			}
		case "error_rate":
			if !seenErrorRate {
				// Get specific stability advice for error rate issues
				errorAdvice := r.getErrorRateRecommendations()
				if errorAdvice != "" {
					recs = append(recs, errorAdvice)
				} else {
					// Fallback to generic advice
					recs = append(recs, "Consider checking system stability or reducing load to decrease error rates.")
				}
				seenErrorRate = true
			}
		}
	}

	// If no specific bottlenecks were found but we have general advice
	if len(recs) == 0 && len(bottlenecks) > 0 {
		recs = append(recs, "System is operating within normal parameters. Consider monitoring for changes.")
	} else if len(recs) == 0 {
		recs = append(recs, "No data available for analysis. Ensure metrics are being collected.")
	}

	return recs
}

// getLatencyRecommendations provides specific advice for reducing latency based on network topology.
func (r *AdvancedRecommendations) getLatencyRecommendations() string {
	// Try to get network topology for path optimization advice
	topology, err := r.metaStore.GetNetworkTopology("backup_network")
	if err != nil || topology == nil {
		return "" // No topology data available
	}

	// Get nodes and edges to understand the current topology
	nodes, err := r.metaStore.GetNetworkNodes(topology.Name)
	if err != nil {
		return ""
	}
	edges, err := r.metaStore.GetNetworkEdges(topology.Name)
	if err != nil {
		return ""
	}

	if len(nodes) == 0 || len(edges) == 0 {
		return "" // Not enough data for path analysis
	}

	// Build the network graph for path optimization
	graph := make(map[string]*network.Node)
	nodeMap := make(map[string]*network.Node)

	// Create nodes
	for _, n := range nodes {
		node := &network.Node{ID: n.ID}
		graph[n.ID] = node
		nodeMap[n.ID] = node
	}

	// Create edges
	for _, e := range edges {
		if source, exists := nodeMap[e.From]; exists {
			if target, exists := nodeMap[e.To]; exists {
				edge := &network.Edge{Target: target, Weight: e.Weight}
				source.Edges = append(source.Edges, edge)
			}
		}
	}

	// For demonstration, let's assume we want to find paths between agents and storage
	// In a real implementation, we'd analyze actual bottleneck metrics to determine specific source/destination
	// For now, we'll provide a general topology-based suggestion

	// Find if there are multiple paths between any node pair (indicating redundancy opportunities)
	if len(nodes) >= 2 {
		// Take first two nodes as example - in reality we'd analyze actual traffic patterns
		sourceID := nodes[0].ID
		destID := nodes[len(nodes)-1].ID

		path, err := r.pathOpt.FindShortestPath(graph, sourceID, destID)
		if err != nil {
			return ""
		}

		// Provide specific path advice
		if len(path.Nodes) > 2 { // More than direct connection
			return "Consider optimizing your backup path: Current path " +
				" -> " + path.Nodes[0] + " -> " + path.Nodes[len(path.Nodes)-1] +
				" may be suboptimal. Alternative routing could reduce latency by utilizing intermediate nodes."
		}
	}

	return ""
}

// getThroughputRecommendations provides specific advice for improving throughput.
func (r *AdvancedRecommendations) getThroughputRecommendations() string {
	// Check if we have storage nodes that could be added for parallelism
	topology, err := r.metaStore.GetNetworkTopology("backup_network")
	if err != nil || topology == nil {
		return ""
	}

	nodes, err := r.metaStore.GetNetworkNodes(topology.Name)
	if err != nil {
		return ""
	}

	// Count storage-like nodes vs agent nodes
	storageCount := 0
	agentCount := 0

	for _, node := range nodes {
		switch node.Name {
		case "storage", "backup-storage":
			storageCount++
		case "agent", "backup-agent":
			agentCount++
		}
	}

	// If we have fewer than 2 storage nodes, suggest adding more for parallelism
	if storageCount < 2 {
		return "Consider adding additional storage nodes to enable parallel backup operations and improve throughput."
	}

	// If we have many agents but few storage nodes, suggest load balancing
	if agentCount > storageCount*2 {
		return "Consider implementing load balancing across your storage nodes to better distribute backup throughput."
	}

	return ""
}

// getErrorRateRecommendations provides specific advice for improving system stability.
func (r *AdvancedRecommendations) getErrorRateRecommendations() string {
	// Check agent metadata for potential stability issues
	// In a full implementation, we would query recent agent status updates
	// For now, we'll provide general advice based on common causes

	return "Consider checking agent health logs and network connectivity issues that may be causing backup failures."
}
