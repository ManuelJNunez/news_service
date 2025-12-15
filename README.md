# News Service

⚠️ **WARNING: This application is intentionally vulnerable to SQL Injection (SQLi) attacks. It is designed for educational and testing purposes only. DO NOT use in production environments.**

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (optional, for containerized deployment)

## Building the Project

### Option 1: Build with Go

```bash
# Download dependencies
go mod download

# Build the binary
go build -o news-service ./cmd/api

# Run the application
./news-service
```

### Option 2: Build with Docker

```bash
# Build and run using Docker Compose
docker-compose up --build

# Or build the Docker image manually
docker build -t news-service .
```

## Running Tests

```bash
# Run unit tests
go test ./...

# Run end-to-end tests
go test -tags=e2e ./test/e2e/...
```
