# Aegis Backend - Phase 1

This document outlines how to run and validate the Aegis Phase 1 backend. 

> **Note:** The current Phase 1 implementation focuses on the core foundational platform (server lifecycle, configuration, health checks, and database migrations). Future capabilities like authentication, RBAC, multi-tenancy, ML, risk intelligence, asynchronous processing, WebSockets, Merkle trees, and production deployment are **planned for later phases** and are not included in this document.

## Prerequisites

- **Go version:** 1.27.0 or higher
- **PostgreSQL version:** 15.0 or higher (required for relational data and health checks)
- **Golang Migrate CLI:** `migrate` (v4.x) for running database migrations

## Configuration (Environment Variables)

The backend is configured entirely via environment variables.

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `AEGIS_DATABASE_URL` | PostgreSQL connection string | Yes | None |
| `AEGIS_SERVER_ADDR` | HTTP server address/port | No | `:8080` |
| `AEGIS_ENV` | Application environment (`development`, `test`, `production`) | No | `development` |

### Development Defaults Example

For local development, you might use:
```bash
export AEGIS_DATABASE_URL="postgres://user:password@localhost:5432/aegis_dev?sslmode=disable"
export AEGIS_SERVER_ADDR=":8080"
export AEGIS_ENV="development"
```

## Running the Database Locally

To start a local PostgreSQL instance for development using Docker:

```bash
docker run -d \
  --name aegis-postgres \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=aegis_dev \
  -p 5432:5432 \
  postgres:15-alpine
```

## Database Migrations

Aegis uses `golang-migrate` for schema versioning. Migrations are located in the `migrations/` directory.

To apply database migrations (Up):
```bash
migrate -database "${AEGIS_DATABASE_URL}" -path migrations up
```

To rollback database migrations (Down):
```bash
migrate -database "${AEGIS_DATABASE_URL}" -path migrations down 1
```

## Starting the Backend

Ensure `AEGIS_DATABASE_URL` is exported, then run:

```bash
cd backend
go run ./cmd/api
```

## Health Endpoints


The backend exposes health endpoints for monitoring and orchestration:

- **Liveness Endpoint:** `GET /health/live`
  - Validates that the Go HTTP server process is running.
  - Does not rely on database connectivity.
  - Expected success response: `200 OK` with `{"status": "UP"}`.

- **Readiness Endpoint:** `GET /health/ready`
  - Validates that the application is fully ready to accept traffic, which includes successfully pinging the PostgreSQL database.
  - Expected success response: `200 OK` with `{"status": "UP"}`.

### Expected Behavior When PostgreSQL is Unavailable
If the database goes down or becomes unreachable:
- The **Liveness** endpoint (`/health/live`) will continue to return `200 OK` (`{"status": "UP"}`).
- The **Readiness** endpoint (`/health/ready`) will return `503 Service Unavailable` (`{"status": "DOWN"}`).

## Testing and Validation

### Running Backend Tests

To run the unit tests (excluding database integration tests):
```bash
go test -short ./...
```

To run all tests (including integration tests):
1. Ensure your test database is running.
2. Export the connection string: `export AEGIS_DATABASE_URL="postgres://user:password@localhost:5432/aegis_test?sslmode=disable"`
3. Run the tests:
```bash
go test -v ./...
```

### Formatting and Static Checks

To ensure code quality and proper formatting, run:

```bash
# Format code
go fmt ./...

# Tidy module dependencies
go mod tidy

# Run static analysis (requires golangci-lint installed)
golangci-lint run
```
