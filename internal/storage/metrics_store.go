package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// MetricsStore handles SQLite storage for time-series metrics.
type MetricsStore struct {
	db *sql.DB
}

// NewMetricsStore creates a new MetricsStore with the given SQLite database path.
func NewMetricsStore(dbPath string) (*MetricsStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	store := &MetricsStore{db: db}
	return store, nil
}

// Init initializes the database schema if it doesn't exist.
func (s *MetricsStore) Init() error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS backup_job_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		latency_ms REAL NOT NULL,
		throughput_mbps REAL NOT NULL,
		error_rate REAL NOT NULL,
		status TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_backup_job_metrics_job_id_timestamp ON backup_job_metrics(job_id, timestamp);
	`
	_, err := s.db.Exec(createTableSQL)
	return err
}

// Metric represents a single backup job metric record.
type Metric struct {
	ID        int64     `json:"id"`
	JobID     string    `json:"job_id"`
	Timestamp time.Time `json:"timestamp"`
	LatencyMS float64   `json:"latency_ms"`
	ThroughputMbps float64 `json:"throughput_mbps"`
	ErrorRate float64   `json:"error_rate"`
	Status    string    `json:"status"`
}

// InsertMetric inserts a new metric record.
func (s *MetricsStore) InsertMetric(m Metric) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO backup_job_metrics (job_id, timestamp, latency_ms, throughput_mbps, error_rate, status) VALUES (?, ?, ?, ?, ?, ?)`,
		m.JobID, m.Timestamp, m.LatencyMS, m.ThroughputMbps, m.ErrorRate, m.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert metric: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return id, nil
}

// GetMetrics retrieves metrics for a given job ID, optionally limited by count and time range.
func (s *MetricsStore) GetMetrics(jobID string, limit int) ([]Metric, error) {
	query := `
	SELECT id, job_id, timestamp, latency_ms, throughput_mbps, error_rate, status
	FROM backup_job_metrics
	WHERE job_id = ?
	ORDER BY timestamp DESC
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := s.db.Query(query, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.ID, &m.JobID, &m.Timestamp, &m.LatencyMS, &m.ThroughputMbps, &m.ErrorRate, &m.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metric row: %w", err)
		}
		metrics = append(metrics, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating metrics rows: %w", err)
	}
	return metrics, nil
}

// Close closes the database connection.
func (s *MetricsStore) Close() error {
	return s.db.Close()
}