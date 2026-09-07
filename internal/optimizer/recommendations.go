package optimizer

import (
	"iterative-pony/internal/analysis"
)

// Recommendations generates optimization suggestions based on detected bottlenecks.
type Recommendations struct{}

// NewRecommendations creates a new Recommendations generator.
func NewRecommendations() *Recommendations {
	return &Recommendations{}
}

// Generate creates optimization recommendations based on the given bottlenecks.
func (r *Recommendations) Generate(bottlenecks []analysis.Bottleneck) []string {
	var recs []string

	// We'll group bottlenecks by type to avoid duplicate recommendations
	seenLatency := false
	seenThroughput := false
	seenErrorRate := false

	for _, b := range bottlenecks {
		switch b.Type {
		case "latency":
			if !seenLatency {
				recs = append(recs, "Consider increasing network bandwidth or optimizing data transfer paths to reduce latency.")
				seenLatency = true
			}
		case "throughput":
			if !seenThroughput {
				recs = append(recs, "Consider increasing parallelism or adding more storage nodes to improve throughput.")
				seenThroughput = true
			}
		case "error_rate":
			if !seenErrorRate {
				recs = append(recs, "Consider checking system stability or reducing load to decrease error rates.")
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