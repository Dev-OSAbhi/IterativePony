# Backup Job Monitor and Optimizer

A Go-based backend service that simulates backup job execution, collects metrics, detects bottlenecks, and provides optimization recommendations.

## Overview

This project aligns with the Rubrik Software Engineer - Intern role, focusing on backend development, bottleneck identification, and system optimization.

## Features

- Simulates backup job execution with configurable agents
- Collects metrics (latency, throughput, error rates)
- Implements bottleneck detection engine
- Provides optimization recommendations (load redistribution, network path optimization)
- Exposes REST API for monitoring and control
- Dual storage: SQLite for metrics, MongoDB for metadata
- Simple CLI to start/stop simulation and configure parameters

## Project Structure

- `cmd/simulator/main.go` - Application entry point
- `internal/api/` - HTTP API handlers
- `internal/simulation/` - Backup job simulation
- `internal/analysis/` - Bottleneck detection
- `internal/storage/` - Database interactions (SQLite & MongoDB)
- `internal/network/` - Path optimization algorithms
- `internal/optimizer/` - Recommendation generation
- `config/` - Configuration loading

## Getting Started

### Prerequisites

- Go 1.22+
- SQLite
- MongoDB

### Installation

```bash
go mod download
```

### Running

```bash
go run ./cmd/simulator
```

## Implementation Plan

The implementation will be broken down into small, incremental steps:

### Phase 1: Basic Server Setup
1. ✅ Implement configuration loading from environment variables
2. ✅ Set up Gin web server with basic middleware (logger, recovery)
3. ✅ Create health check endpoint
4. ✅ Set up SQLite database connection for metrics storage
5. ✅ Set up MongoDB connection for metadata storage

### Phase 2: Basic Simulation Engine
1. ✅ Create simple backup job simulator that generates dummy metrics
2. ✅ Implement metrics collection and storage in SQLite
3. ✅ Create API endpoint to trigger simulation

### Phase 3: Basic API and Storage
1. ✅ Implement basic CRUD operations for metrics in SQLite
2. ✅ Implement basic storage for agent metadata in MongoDB
3. ✅ Create API endpoints to retrieve metrics and agent information

### Phase 4: Initial Analysis and Optimization
1. Implement basic bottleneck detection (simple threshold-based)
2. Generate simple optimization recommendations
3. Create API endpoints for analysis results and recommendations

## Verification Steps

1. ✅ Run the server and verify it starts without errors
2. ✅ Access health check endpoint and verify response
3. ✅ Verify SQLite database connection works (health endpoint shows sqlite: connected)
4. ✅ Verify MongoDB connection works (health endpoint shows mongo: connected)
5. ✅ Trigger simulation via API and verify metrics are stored
6. ✅ Retrieve metrics via API and verify correct data returned
7. Check that basic bottleneck detection and recommendations work

## Testing

Tests are located in the `./tests` directory, organized by package. To run all tests:

```bash
go test ./tests/... -v
```