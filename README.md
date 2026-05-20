# Aegis

<p align="center">
  <strong>Production-minded access control and traffic protection in Go.</strong>
</p>

Aegis is a server-rendered (or rather, gRPC-first backend) application built with intent, not with scaffolding. It is an authorization gateway foundation for answering one hot-path question: **should this request be allowed right now?**

The project is intentionally shaped like a real production backend rather than a small demo: configuration is environment-driven, the service starts through a dedicated command, health endpoints are in place, and the next layers are designed around Redis, PostgreSQL, gRPC, and operational telemetry.

## Table of Contents
- [Why This Project Exists](#why-this-project-exists)
- [Product Vision](#product-vision)
- [Current Status](#current-status)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Routes](#api-routes)
- [Configuration](#configuration)
- [Local Development](#local-development)
- [Development Roadmap](#development-roadmap)
- [Docs](#docs)

## Why This Project Exists

Most systems spread access decisions across several places: JWT validation in middleware, policy checks in application code, rate limiting in an API gateway, and audit logs somewhere else. That makes the actual decision path difficult to reason about and harder to test under failure.

Aegis puts the decision pipeline behind one contract. It is built as a serious backend engineering project to demonstrate atomic operations, robust caching, and transactional integrity on the hot path, providing a centralized and observable location for all access-control logic.

## Product Vision

Aegis acts as a unified gatekeeper. The intended workflow looks like this:
1. A Client or Demo Service calls Aegis through gRPC or the REST gateway.
2. Aegis validates the JWT signature, issuer, audience, expiry, and required claims.
3. It checks Redis for token or subject revocation.
4. It loads the matching policy from a memory cache or PostgreSQL.
5. It calculates any adaptive tightening factor for the caller.
6. It applies the configured rate-limit algorithm atomically in Redis.
7. It writes the decision and audit metadata to PostgreSQL (potentially via a transactional outbox).
8. The client receives an allow/deny decision plus rate-limit metadata.

## Current Status

This repository currently contains the service bootstrap:
- Go command entrypoint at `cmd/aegis`
- Environment-based configuration with validation
- Structured logging through `log/slog`
- HTTP server with graceful shutdown
- Liveness and readiness endpoints
- Docker Compose definitions for Redis, PostgreSQL, Aegis, Prometheus, and Grafana

## Architecture

The architecture is intentionally conservative and explicit, designed around high-throughput and failure-resiliency.

```mermaid
flowchart TD
    Client["Client or Demo Service"]
    REST["REST Gateway"]
    GRPC["gRPC API"]
    
    subgraph aegis_core [Aegis Core]
        Auth["JWT and Revocation<br/>(Identity checks)"]
        Policy["Policy Engine<br/>(Access Rules)"]
        Limit["Rate Limiter<br/>(Token Bucket / Sliding Window)"]
        Adaptive["Adaptive EWMA<br/>(Traffic Protection)"]
        Audit["Audit Writer<br/>(Durable Logging)"]
        Worker["Outbox Worker<br/>(Async tasks)"]
    end
    
    subgraph storage_layer [Storage]
        Redis[("(Redis)<br/>Acceleration & Ephemeral State")]
        Postgres[("(PostgreSQL)<br/>Durable Truth")]
    end
    
    Webhook["Webhook Sink"]

    Client --> REST
    Client --> GRPC
    REST --> GRPC
    GRPC --> Auth
    Auth --> Redis
    Auth --> Policy
    Policy --> Postgres
    Policy --> Limit
    Limit --> Redis
    Limit --> Adaptive
    Adaptive --> Redis
    Policy --> Audit
    Audit --> Postgres
    Audit --> Worker
    Worker --> Postgres
    Worker --> Webhook
```

### Request Lifecycle
```
Client Request → gRPC / CheckAccess → Validate JWT → Check Revocations (Redis) → Load Policy (Cache/Postgres) → Apply Rate Limits (Redis Lua) → Record Audit (Postgres) → Return Decision
```

### Two Kinds of State

Aegis explicitly separates its state into **durable truth** and **acceleration**.
- **PostgreSQL** is the durable truth: Policies, audit events, outbox messages, attempts, and dead letters belong here.
- **Redis** is acceleration, not truth: Counters, token buckets, revocations, idempotency keys, and adaptive baselines live here. Hot-path mutations are atomic using Lua scripts or single commands.

## Tech Stack

- **Go 1.25+**: The core service language
- **gRPC / Protocol Buffers**: The primary service-to-service contract
- **PostgreSQL**: Durable storage for policies and audit logs
- **Redis**: Fast, atomic storage for rate limits and revocations
- **log/slog**: Structured application logging

## Project Structure

```text
Aegis/
├── cmd/aegis/                  # Service entrypoint
├── deployments/
│   └── docker-compose.yml      # Local infrastructure configuration
├── docs/
│   └── Aegis_Reading_Map.md    # Build sequence, references, and implementation roadmap
├── internal/
│   ├── app/                    # HTTP server, health handlers, app wiring
│   └── config/                 # Environment config and validation
├── go.mod
├── go.sum
└── README.md
```

## Database Schema

The relational model is designed around atomic policies, robust auditability, and outbox reliability.

```mermaid
erDiagram
    policies ||--o{ audit_events : "evaluated in"
    
    policies {
        bigserial id PK
        text name UK
        text effect
        text resource
        text action
        timestamptz created_at
    }
    
    audit_events {
        bigserial id PK
        text request_id
        text subject
        text decision
        text reason
        timestamptz timestamp
    }
    
    revocations {
        text token_id PK
        timestamptz expires_at
    }
```

## API Routes

| Method | Path / RPC | Status | Description |
|---|---|---:|---|
| `gRPC` | `aegis.v1.Access/CheckAccess` | Planned | Primary hot-path authorization check. |
| `GET` | `/healthz` | `200` | Process liveness check. |
| `GET` | `/readyz` | `200` | Readiness check. Currently returns ready after app wiring succeeds. |

## Configuration

Aegis is configured entirely through environment variables, designed for easy deployment in containerized environments.

| Variable | Default | Required | Purpose |
|---|---|:---:|---|
| `POSTGRES_DSN` | empty | Yes | PostgreSQL connection string. |
| `JWT_ISSUER` | empty | Yes | Expected JWT issuer. |
| `JWT_AUDIENCE` | empty | Yes | Expected JWT audience. |
| `JWKS_URL` | empty | Yes* | Remote JWKS endpoint. |
| `LOCAL_JWKS_PATH` | empty | Yes* | Local JWKS file path for development. |
| `RATE_LIMIT_FAILURE_MODE` | `open` | No | `open` or `closed` on dependency failure. |

## Local Development

Clone the repository:
```bash
git clone https://github.com/SShogun/Aegis.git
cd Aegis
```

Start the required infrastructure (Redis and PostgreSQL):
```bash
docker compose -f deployments/docker-compose.yml up -d redis postgres
```

Run Aegis locally:
```bash
POSTGRES_DSN="postgres://aegis:aegis@localhost:5432/aegis?sslmode=disable" \
JWT_ISSUER="http://localhost:9000" \
JWT_AUDIENCE="aegis" \
LOCAL_JWKS_PATH="./dev/jwks.json" \
go run ./cmd/aegis
```

For Windows PowerShell:
```powershell
$env:POSTGRES_DSN="postgres://aegis:aegis@localhost:5432/aegis?sslmode=disable"
$env:JWT_ISSUER="http://localhost:9000"
$env:JWT_AUDIENCE="aegis"
$env:LOCAL_JWKS_PATH="./dev/jwks.json"
go run ./cmd/aegis
```

## Development Roadmap

| Phase | Status | Scope |
|---:|---|---|
| **0** | Done | Go module, config, logging, app wiring, health endpoints. |
| **1** | Next | gRPC API, protobuf contract, grpc-gateway REST mapping. |
| **2** | Planned | JWT/OIDC validation, JWKS loading, claims extraction. |
| **3** | Planned | Redis revocation checks and atomic rate limiter scripts. |
| **4** | Planned | PostgreSQL policies, audit events, migrations, transaction boundaries. |
| **5** | Planned | Transactional outbox, retries, dead-letter handling. |
| **6** | Planned | Prometheus metrics, OpenTelemetry traces, Grafana dashboard. |
| **7** | Planned | Load tests, failure-mode docs, deployment guide. |

## Docs

- [`docs/Aegis_Reading_Map.md`](docs/Aegis_Reading_Map.md): Implementation sequence, references, and build calendar.

---
**License:** MIT
