# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o backup-job-monitor ./cmd/simulator

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/backup-job-monitor .

# Copy configuration and environment files
COPY .env .env

# Create directory for SQLite database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8081

# Run the application
CMD ["./backup-job-monitor"]