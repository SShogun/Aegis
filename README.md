# Aegis

**Aegis monitors user and system activity, detects threats, assesses risk, and maintains tamper-evident audit logs.**

Aegis is an engineering-focused project for turning security events—such as logins, API requests, resource access, and privilege changes—into explainable risk signals and verifiable audit history. The repository is being built incrementally: today it contains the Go and PostgreSQL foundation; event processing, risk intelligence, and the dashboard are planned layers.

## The problem

Systems produce more activity than people can investigate manually. Suspicious behavior can be difficult to distinguish from normal activity, its severity can be difficult to explain, and historical audit records are only useful when unauthorized changes can be detected. A multi-tenant system must also keep each organization’s data isolated.

Aegis is designed to address these problems through event processing, contextual risk analysis, layered authorization, and cryptographically verifiable audit history.

## What Aegis does

The intended end-to-end flow is:

```text
User/System Activity
        ↓
Event Ingestion and Processing
        ↓
Threat / Anomaly Detection
        ↓
Risk Assessment
        ↓
Dashboard and Alerts
        ↓
Tamper-Evident Audit Trail
```

Only the Phase 1 service foundation is implemented at present. The event, detection, risk, alerting, and audit stages above are planned.

## Example scenario

The following illustrates the intended behavior; it is not an implemented workflow yet.

```text
User logs in
      ↓
Accesses a sensitive resource at an unusual time
      ↓
Downloads far more data than their behavioral baseline
      ↓
Aegis combines anomaly and contextual signals
      ↓
Risk is assessed with contributing factors
      ↓
The decision and security action are recorded for audit
```

## How it works

The long-term architecture separates event handling, intelligence, risk aggregation, and audit integrity so that no single model or component owns the whole decision:

```text
                         Aegis
                           |
        +------------------+------------------+
        |                  |                  |
        v                  v                  v
   Event System       ML / Behavior      Audit System
        |                  |                  |
        +------------------+------------------+
                           |
                           v
                     Risk Engine
                           |
                           v
                      Dashboard
```

This is the documented target architecture. The current repository implements the backend lifecycle and platform foundation only.

## ML and risk intelligence

The planned intelligence layer is hybrid rather than dependent on one opaque model:

```text
Anomaly Detection  → Is this unusual?
Classification      → Does it resemble known malicious behavior?
Behavioral Analysis → Is it unusual for this user or entity?
Risk Engine         → How dangerous is the overall situation?
```

The project overview identifies these planned signal families:

- Isolation Forest for global anomaly detection
- Autoencoder-based anomaly detection
- Behavioral baselines and deviation scoring
- Supervised classification such as XGBoost or LightGBM
- Contextual security signals
- Risk aggregation and explainability

These models, training workflows, inference service, and explainability layer are not present in the current repository.

## Cryptographic audit integrity

Aegis is designed to record important security actions in a tamper-evident audit trail so unauthorized modifications to historical records can be detected.

```text
Audit Event
    ↓
Canonical Representation
    ↓
Cryptographic Hash
    ↓
Hash Chain
    ↓
Integrity Verification
```

Hash-chained audit records are planned. Merkle-tree batches and inclusion proofs are later planned extensions; they are not implemented in Phase 1.

## Security model

The target security model treats identity, authorization, and isolation as separate concerns:

- OAuth/OIDC establishes user identity.
- RBAC provides explicit permissions, with contextual authorization planned later.
- Each organization owns its users, events, policies, models, and audit records.
- Tenant isolation is intended to be enforced in both the application and PostgreSQL Row-Level Security.

OAuth/OIDC, RBAC, multi-tenancy, and PostgreSQL RLS are planned. The current database contains only foundation tables and does not yet expose tenant-aware application data.

## Architecture

```text
                         React Dashboard
                                |
                                v
                           Go Backend
                         /     |      \\
                        /      |       \\
                       v       v        v
                 PostgreSQL  Event     Python ML
                             System     Service
                                          |
                                          v
                                  Risk Intelligence
```

The diagram describes the intended system boundary. The checked-in implementation currently contains the Go backend and PostgreSQL integration; the dashboard, event infrastructure, and Python service are future components.

## Technology stack

| Layer | Technology | Status and purpose |
|---|---|---|
| Backend | Go 1.27 | Implemented API and service lifecycle |
| Database | PostgreSQL 15+ | Implemented persistence and readiness checks |
| Database access | pgx/v5 | Implemented connection pooling |
| Migrations | golang-migrate v4 | Implemented schema migration and rollback support |
| ML service | Python; scikit-learn; XGBoost or LightGBM | Planned model training and inference |
| Frontend | React + TypeScript | Planned analytics dashboard |
| Event infrastructure | Event streams and asynchronous consumers | Planned event processing |
| Identity | OAuth/OIDC | Planned authentication |
| Authorization | RBAC, contextual rules, PostgreSQL RLS | Planned access control and tenant isolation |
| Audit integrity | Hash chains; Merkle trees | Planned tamper evidence and inclusion proofs |

## Current status

### Implemented — Phase 1 foundation

- Environment configuration with validation and safe database URL redaction in logs
- Fail-fast startup when the database URL, environment, or server address is invalid
- Bounded PostgreSQL connection pool and readiness connectivity checks
- `GET /health/live`, independent of PostgreSQL availability
- `GET /health/ready`, reporting PostgreSQL readiness
- Graceful HTTP shutdown on `SIGINT` and `SIGTERM`
- Initial SQL migration and rollback using `golang-migrate`
- Unit tests and opt-in PostgreSQL integration tests

### Planned

- Authentication and tenant-aware authorization
- Event ingestion and asynchronous processing
- Behavioral profiling, trained ML models, and explainable risk aggregation
- Tamper-evident audit records, then Merkle-tree verification
- React analytics dashboard and real-time updates
- Production hardening, observability, deployment, and CI/CD

## Roadmap

The roadmap follows the project’s documented principle that complexity should be earned one architectural responsibility at a time. Phases after Phase 1 are planned, not commitments to a completed implementation.

| Phase | Focus | Status |
|---|---|---|
| Phase 0 | Project definition and architecture | Planned foundation work |
| Phase 1 | Go service, configuration, database, health, and migrations | Implemented |
| Phase 2 | Identity and authentication | Planned |
| Phase 3 | Authorization and permissions | Planned |
| Phase 4 | Multi-tenancy and database isolation | Planned |
| Phase 5 | Event ingestion and asynchronous processing | Planned |
| Phase 6 | Machine learning models and inference | Planned |
| Phase 7 | Risk intelligence and explainability | Planned |
| Phase 8 | Cryptographic audit and verification | Planned |
| Phase 9 | Real-time analytics dashboard | Planned |
| Phase 10 | Production hardening and operations | Planned |

## Engineering highlights

Aegis is technically interesting because its target design brings several responsibilities together while keeping their roles distinct:

- Hybrid anomaly, behavioral, classification, and contextual risk signals
- Separation of ML inference from application-level risk aggregation
- Defense-in-depth tenant isolation and authorization
- Audit history designed for independent integrity verification
- Incremental architecture: each phase introduces a concrete system responsibility

The last four items describe planned architecture beyond the current foundation.

## Project structure

```text
.
├── backend/
│   ├── cmd/api/                 # Backend entrypoint
│   ├── internal/health/         # Liveness and readiness handlers
│   ├── internal/platform/       # Configuration, server, and database lifecycle
│   ├── migrations/              # Up/down SQL migrations
│   └── README.md                # Backend development notes
├── docs/
│   └── 01_OVERVIEW.md           # Product vision and target architecture
└── README.md
```

## Quick start

### Prerequisites

- Go 1.27.0 or newer
- PostgreSQL 15 or newer
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) v4.x
- Docker, if running PostgreSQL locally in a container

### Start PostgreSQL

```bash
docker run -d \
  --name aegis-postgres \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=aegis_dev \
  -p 5432:5432 \
  postgres:15-alpine
```

### Configure and start the backend

```bash
export AEGIS_DATABASE_URL="postgres://user:password@localhost:5432/aegis_dev?sslmode=disable"
export AEGIS_SERVER_ADDR=":8080"
export AEGIS_ENV="development"

cd backend
migrate -database "$AEGIS_DATABASE_URL" -path migrations up
go run ./cmd/api
```

The API listens on `http://localhost:8080`.

### Check health

```bash
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

Both return `{"status":"UP"}` when the corresponding check succeeds. Liveness does not depend on PostgreSQL; readiness returns `503` and `{"status":"DOWN"}` when PostgreSQL is unavailable.

## Configuration

| Variable | Required | Default | Description |
|---|---:|---|---|
| `AEGIS_DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `AEGIS_SERVER_ADDR` | No | `:8080` | HTTP bind address and port |
| `AEGIS_ENV` | No | `development` | `development`, `test`, or `production` |

The backend fails fast when the database URL is missing, the environment is unsupported, or the server address does not include a port.

## Testing

From `backend/`:

```bash
# Unit tests; skips tests requiring PostgreSQL
go test -short ./...

# All tests; requires AEGIS_DATABASE_URL to point to a test database
go test -v ./...

# Formatting and dependency consistency
go fmt ./...
go mod tidy
```

## Documentation

- [Project overview and long-term architecture](docs/01_OVERVIEW.md)
- [Backend setup, health checks, migrations, and testing](backend/README.md)

## License

No license has been declared yet.
