# Backup Job Monitor and Optimizer

A Go-based backend service that simulates backup job execution, collects metrics, detects bottlenecks, and provides optimization recommendations.

## Features

- Simulates backup job execution with configurable agents
- Collects metrics (latency, throughput, error rates)
- Implements bottleneck detection engine
- Provides optimization recommendations (load redistribution, network path optimization)
  - Includes basic recommendations based on metric thresholds
  - Includes advanced, topology-aware recommendations when network topology data is available
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
- `Dockerfile` - Docker image build instructions
- `docker-compose.yml` - Docker Compose file for easy setup of the application and its dependencies (MongoDB)

## Getting Started

### Prerequisites

- Go 1.22+ (for local development)
- Docker and Docker Compose (for containerized deployment)
- OR SQLite and MongoDB (for local development without Docker)

### Installation

```bash
go mod download
```

### Running

#### Local Development (without Docker)

```bash
go run ./cmd/simulator
```

#### Containerized Deployment (with Docker)

```bash
# Build and start all services
docker-compose up --build

# To run in detached mode
docker-compose up -d --build

# To stop and remove containers
docker-compose down

# To stop containers but keep volumes
docker-compose stop
```

## Testing

Tests are located in the `./tests` directory, organized by package. To run all tests:

```bash
go test ./tests/... -v
```

