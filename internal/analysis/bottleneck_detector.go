package analysis

import (
	"fmt"
	"time"

	"iterative-pony/internal/storage"
)

// Bottleneck represents a detected bottleneck in the system.
type Bottleneck struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // e.g., "latency", "throughput", "error_rate"
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Message   string    `json:"message"`
}

// BottleneckDetector analyzes metrics to identify bottlenecks.
type BottleneckDetector struct{}

// NewBottleneckDetector creates a new BottleneckDetector.
func NewBottleneckDetector() *BottleneckDetector {
	return &BottleneckDetector{}
}

// DetectBottlenecks analyzes the given metrics and returns any detected bottlenecks.
// It uses simple threshold-based detection.
//
// Thresholds (example values):
//   - Latency > 100 ms
//   - Throughput < 50 Mbps
//   - Error rate > 0.05 (5%)
func (d *BottleneckDetector) DetectBottlenecks(metrics []storage.Metric) []Bottleneck {
	var bottlenecks []Bottleneck

	for _, m := range metrics {
		// Check latency
		if m.LatencyMS > 100.0 {
			bottlenecks = append(bottlenecks, Bottleneck{
				Timestamp: m.Timestamp,
				Type:      "latency",
				Value:     m.LatencyMS,
				Threshold: 100.0,
				Message:   fmt.Sprintf("High latency: %.2f ms > 100 ms threshold", m.LatencyMS),
			})
		}

		// Check throughput
		if m.ThroughputMbps < 50.0 {
			bottlenecks = append(bottlenecks, Bottleneck{
				Timestamp: m.Timestamp,
				Type:      "throughput",
				Value:     m.ThroughputMbps,
				Threshold: 50.0,
				Message:   fmt.Sprintf("Low throughput: %.2f Mbps < 50 Mbps threshold", m.ThroughputMbps),
			})
		}

		// Check error rate
		if m.ErrorRate > 0.05 {
			bottlenecks = append(bottlenecks, Bottleneck{
				Timestamp: m.Timestamp,
				Type:      "error_rate",
				Value:     m.ErrorRate,
				Threshold: 0.05,
				Message:   fmt.Sprintf("High error rate: %.2f%% > 5%% threshold", m.ErrorRate*100),
			})
		}
	}

	return bottlenecks
}