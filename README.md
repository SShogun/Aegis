# Aegis

### Event-driven risk intelligence with audit trails you can verify

Aegis is an engineering-focused platform for turning organizational activity into explainable risk signals. The long-term design combines event processing, behavioral analysis, machine learning, tenant isolation, and cryptographically verifiable audit history in one system.

> **Current status — Phase 1: Foundation complete**
>
> The repository currently contains the Go backend foundation: validated environment configuration, PostgreSQL connection pooling, migration support, health probes, graceful shutdown, and unit/integration test coverage. Authentication, tenant-aware authorization, event ingestion, ML inference, real-time analytics, and cryptographic audit proofs are planned—not implemented yet.

## Why Aegis?

Modern systems produce more activity than people can investigate manually. Aegis is being built to answer a practical question:

> Can a system learn what normal organizational behavior looks like, explain what makes an event risky, keep tenants isolated, and prove that its audit history has not been silently changed?

The eventual platform will combine multiple signals instead of trusting one opaque model:

~~~text
Organizational event
        |
        v
Authentication + authorization + tenant context
        |
        v
Event processing
        |
        +----------------------+----------------------+
        |                      |
        v                      v
Risk intelligence       Tamper-evident audit
        |                      |
        +----------+-----------+
                   v
          Explainable risk view
~~~

## What is working today

The current backend can:

- Load and validate configuration from environment variables.
- Require a PostgreSQL connection string before startup.
- Create and manage a bounded pgx connection pool.
- Serve GET /health/live without depending on PostgreSQL.
- Serve GET /health/ready based on PostgreSQL connectivity.
- Run the HTTP server with graceful shutdown on SIGINT and SIGTERM.
- Apply and roll back the initial database migration with golang-migrate.
- Redact the database URL from configuration logs.
- Run isolated unit tests and optional PostgreSQL integration tests.

## Current roadmap

| Area | Status |
|---|---|
| Go service lifecycle | ✅ Implemented |
| Environment configuration and validation | ✅ Implemented |
| PostgreSQL pool and health checks | ✅ Implemented |
| Initial migration and rollback scripts | ✅ Implemented |
| Unit tests | ✅ Implemented |
| PostgreSQL integration tests | ✅ Implemented, opt-in locally |
| Authentication / OAuth / OIDC | ⏳ Planned |
| RBAC and fine-grained authorization | ⏳ Planned |
| Multi-tenancy and database RLS | ⏳ Planned |
| Event ingestion and asynchronous processing | ⏳ Planned |
| Behavioral profiling and ML inference | ⏳ Planned |
| Explainable risk aggregation | ⏳ Planned |
| Encrypted, hash-chained audit logs | ⏳ Planned |
| Merkle-tree proofs | ⏳ Planned |
| React analytics dashboard | ⏳ Planned |
| Production deployment and CI/CD | ⏳ Planned |

## Quick start

### Prerequisites

- Go 1.27.0 or newer
- PostgreSQL 15 or newer
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) v4.x
- Docker, if you want to run PostgreSQL locally in a container

### 1. Start PostgreSQL

~~~bash
docker run -d \
  --name aegis-postgres \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=aegis_dev \
  -p 5432:5432 \
  postgres:15-alpine
~~~

### 2. Configure the backend

~~~bash
export AEGIS_DATABASE_URL="postgres://user:password@localhost:5432/aegis_dev?sslmode=disable"
export AEGIS_SERVER_ADDR=":8080"
export AEGIS_ENV="development"
~~~

### 3. Apply the schema and start Aegis

~~~bash
cd backend
migrate -database "$AEGIS_DATABASE_URL" -path migrations up
go run ./cmd/api
~~~

The API is now available at http://localhost:8080.

### 4. Check service health

~~~bash
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
~~~

Expected response when the service and database are healthy:

~~~json
{"status":"UP"}
~~~

/health/live reports whether the process is running. /health/ready reports whether the application can reach PostgreSQL. If PostgreSQL becomes unavailable, liveness remains 200 OK while readiness returns 503 Service Unavailable with {"status":"DOWN"}.

## Configuration

| Variable | Required | Default | Description |
|---|---:|---|---|
| AEGIS_DATABASE_URL | Yes | — | PostgreSQL connection string |
| AEGIS_SERVER_ADDR | No | :8080 | HTTP bind address and port |
| AEGIS_ENV | No | development | Runtime environment: development, test, or production |

The backend fails fast when the database URL is missing, the environment is unsupported, or the server address does not include a port.

## Testing

From backend/:

~~~bash
# Unit tests; skips tests that require a real PostgreSQL instance
go test -short ./...

# All tests; requires AEGIS_DATABASE_URL to point to a test database
go test -v ./...

# Format and inspect the module
go fmt ./...
go mod tidy
~~~

The current short test suite passes. Integration tests are skipped unless AEGIS_DATABASE_URL is set.

## Database migrations

Migration files live in [backend/migrations](backend/migrations). Apply the latest schema with:

~~~bash
cd backend
migrate -database "$AEGIS_DATABASE_URL" -path migrations up
~~~

Roll back the latest migration with:

~~~bash
cd backend
migrate -database "$AEGIS_DATABASE_URL" -path migrations down 1
~~~

Phase 1 creates the migration history proof table and a minimal system_settings table. Domain tables are intentionally deferred until the identity and tenancy phases are defined.

## Architecture direction

Aegis is deliberately being built in layers:

~~~text
Foundation
    |
    v
Identity -> Authorization -> Multi-tenancy
    |
    v
Events -> Risk intelligence -> Cryptographic audit
    |
    v
Real-time analytics -> Production hardening
~~~

The intended final system includes:

- A Go API and event-processing backend.
- A Python service for real, trained ML models and inference.
- A React + TypeScript analytics interface.
- PostgreSQL as the operational datastore.
- OAuth/OIDC, RBAC, and database-level tenant isolation.
- Explainable risk decisions built from rules, behavioral baselines, and model signals.
- Authenticated encryption, hash chains, and Merkle-tree inclusion proofs for audit verification.

These are design targets from the project overview, not claims about the current implementation.

## Repository layout

~~~text
.
├── backend/
│   ├── cmd/api/                 # Backend entrypoint
│   ├── internal/health/         # Liveness and readiness handlers
│   ├── internal/platform/       # Config, server, and PostgreSQL lifecycle
│   ├── migrations/              # Up/down SQL migrations
│   └── README.md                # Backend-specific development notes
├── docs/
│   └── 01_OVERVIEW.md           # Product vision and phased architecture
└── README.md
~~~

## Project principles

- Complexity must be earned one phase at a time.
- Every component should solve a real problem.
- Security is layered: authentication, authorization, isolation, encryption, and auditability each have distinct jobs.
- ML must be trained, evaluated, and monitored—not used as a decorative API wrapper.
- Cryptography should use established primitives, with the interesting work in key handling, integrity, and verification flows.

## Documentation

- [Project overview and long-term architecture](docs/01_OVERVIEW.md)
- [Backend setup, health checks, migrations, and testing](backend/README.md)

## Contributing

The project is being developed phase by phase. Before adding a feature, check the overview and current implementation boundary, keep changes focused on the active phase, and include tests for behavior that can be exercised locally.

## License

No license has been declared yet.
