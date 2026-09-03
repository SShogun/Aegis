# Aegis

**Aegis monitors user and system activity, detects threats, assesses risk, and maintains tamper-evident audit logs.**

Aegis is an engineering-focused security platform for turning security events—such as logins, API requests, resource access, and privilege changes—into explainable risk signals and verifiable audit history.

The project is being built incrementally. Each phase introduces one major architectural capability and leaves behind a testable result before the next layer is added.

## The Problem

Systems produce more activity than people can investigate manually. Suspicious behavior can be difficult to distinguish from normal activity, its severity can be difficult to explain, and historical audit records are only useful when unauthorized changes can be detected.

A multi-tenant security platform must also ensure that each organization's data and operations remain isolated.

Aegis is designed to address these problems through:

* Event ingestion and validation
* Hybrid anomaly and behavioral detection
* Contextual risk analysis
* Layered authorization and tenant isolation
* Explainable risk signals
* Cryptographically verifiable audit history
* Real-time security intelligence

## What Aegis Does

The intended end-to-end flow is:

```text
User / System Activity
          ↓
    Event Ingestion
          ↓
   Event Validation
          ↓
   ML / Anomaly Signals
          ↓
 Contextual Risk Analysis
          ↓
    Risk Assessment
          ↓
 Dashboard / Alerts
          ↓
 Tamper-Evident Audit Trail
```

The project is being implemented layer by layer rather than building this entire pipeline at once.

**Current implementation status: Phase 5 and Phase 2 (backend) are complete.**

The independent hybrid ML system is implemented and can train models, persist artifacts, perform inference, and combine multiple detection signals. The Go risk-intelligence integration begins in Phase 6.

## Example Scenario

The following illustrates the intended long-term behavior:

```text
User logs in
      ↓
Accesses sensitive resources at an unusual time
      ↓
Downloads significantly more data than normal
      ↓
Multiple detection methods identify unusual behavior
      ↓
ML signals are combined with security context
      ↓
Aegis produces an explainable risk assessment
      ↓
The event and resulting security action are recorded
      ↓
Audit integrity can later be verified
```

The complete workflow is a target architecture. Not every stage is currently implemented.

---

## How It Works

Aegis separates detection, risk aggregation, and audit integrity so that no single model or component owns the entire security decision.

```text
                         Aegis
                           |
        +------------------+------------------+
        |                  |                  |
        v                  v                  v
   Event System       ML / Behavior      Audit System
        |                  |                  |
        |                  v                  |
        |             ML Signals             |
        |                  |                  |
        +------------------+------------------+
                           |
                           v
                     Risk Engine
                           |
                           v
                      Dashboard
```

The architecture is intentionally incremental.

The current implementation has established the backend foundation and independent ML capability. Future phases connect those capabilities into the complete risk-intelligence platform.


---

# Authentication Architecture

Aegis uses OpenID Connect (OIDC) to authenticate users and issues its own persistent, HTTP-only sessions to secure the backend API.

```mermaid
flowchart TD
    User([User]) -->|GET /auth/login| Login[Login Handler]
    Login -->|Redirect| OIDC[OIDC Provider]
    OIDC -->|Redirect w/ Code| Callback[Callback Handler]
    Callback --> ID[Verify ID Token]
    ID --> Mapping[Find/Create Internal User]
    Mapping --> Session[Create Persistent Session]
    Session --> Cookie[Set HttpOnly Cookie]
    Cookie --> Auth[Authenticated Requests]
    Auth --> Middleware[Auth Middleware]
    Middleware --> Endpoint[Protected Endpoint]
```

### Public and Protected Boundaries

The backend enforces strict boundaries between public authentication routes and protected data.

* `/auth/login` — **Public**
* `/auth/callback` — **Public**
* `/auth/logout` — **Public** (invalidates supplied session if present)
* `/auth/me` — **Authenticated** (requires valid session)

No OAuth tokens or session IDs are exposed to client JavaScript.

---

# ML and Risk Intelligence

Aegis uses a **hybrid detection approach** rather than depending on a single model.

The current Phase 5 ML pipeline contains four complementary signal families:

```text
                    Event Features
                         |
        +----------------+----------------+
        |                |                |
        v                v                v
  Supervised       Isolation Forest   Behavioral
  Classifier       Anomaly Detection   Baseline
        |                |                |
        +----------------+----------------+
                         |
                    Autoencoder
                         |
                         v
                  ML Signal Layer
                         |
                         v
                   Hybrid Score
```

### Implemented ML Components

#### 1. Supervised Classifier

A Decision Tree classifier is trained using labeled events:

```text
Features
   ↓
Decision Tree
   ↓
Probability of suspicious class
```

The classifier learns from known normal and suspicious examples.

#### 2. Isolation Forest

Isolation Forest performs unsupervised anomaly detection.

It does not require the `is_suspicious` label to determine whether an event looks unusual.

```text
Event
  ↓
Isolation Forest
  ↓
Anomaly Score
```

#### 3. Behavioral Baseline

The behavioral detector establishes a statistical baseline from observed feature values and measures how far an event deviates from that baseline.

```text
Observed Event
      ↓
Compare Against Baseline
      ↓
Normalized Deviation
      ↓
Behavioral Signal
```

The current implementation uses a global statistical baseline. Per-user or per-entity behavioral profiling is planned for later risk-intelligence work.

#### 4. Autoencoder

The autoencoder learns to reconstruct event features.

Large reconstruction errors indicate that an event is harder for the model to represent.

```text
Input Event
     ↓
Autoencoder
     ↓
Reconstructed Event
     ↓
Reconstruction Error
     ↓
Anomaly Signal
```

#### 5. Hybrid Signal

The four signals are combined into a single ML anomaly signal:

```text
Classifier       ──┐
Isolation Forest ──┤
Behavioral       ──┼──→ Hybrid ML Score
Autoencoder      ──┘
```

The current implementation uses initial engineering weights:

```text
Classifier       30%
Isolation Forest 25%
Behavioral       25%
Autoencoder      20%
```

These weights are an initial prototype configuration, not a scientifically validated production calibration.

### Signal Semantics

The individual signals are normalized to a 0–1 range where applicable:

```text
0 ───────────────────────────── 1
│                                │
Less suspicious              More suspicious
```

These values should generally be interpreted as **ML/anomaly signals**, not probabilities.

The supervised classifier produces a probability-like suspicious-class value through `predict_proba()`. The other detection methods produce anomaly or deviation scores that are normalized for comparison and ensemble use.

The hybrid score is an ML signal and is **not yet the final Aegis security risk score**.

Final application-level risk aggregation is introduced in Phase 6.

---

# Phase 5 Evaluation

The current synthetic dataset contains both normal and suspicious events.

The completed evaluation showed:

| Signal           |    Normal | Suspicious |
| ---------------- | --------: | ---------: |
| Classifier       |     0.000 |      0.900 |
| Isolation Forest |     0.142 |      0.684 |
| Behavioral       |     0.070 |      0.462 |
| Autoencoder      |     0.147 |      0.416 |
| **Hybrid**       | **0.082** |  **0.640** |

The important result is that the different detection approaches produce **higher anomaly/suspicion signals for the known suspicious records than for the normal records**.

The classifier provides the strongest separation, while Isolation Forest, behavioral analysis, and the autoencoder provide complementary signals.

The hybrid signal combines those perspectives into a single ML-level assessment.

### Evaluation Limitation

The current dataset is intentionally small and synthetic.

The evaluation demonstrates that the pipeline works coherently and that the implemented detection methods distinguish the provided normal and suspicious examples.

It should **not** be interpreted as proof of real-world attack-detection accuracy.

Production-quality evaluation will require larger and more representative datasets, proper validation/test separation, realistic attack distributions, calibration, and model monitoring.

---

# Cryptographic Audit Integrity

Aegis is designed to record important security actions in a tamper-evident audit trail so unauthorized modifications to historical records can be detected.

The planned model is:

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

Hash-chained audit records are planned for Phase 7.

Merkle-tree batching and inclusion proofs are planned as later extensions where they provide useful verification capabilities.

**Cryptographic audit integrity is not currently implemented.**

---

# Security Model

The target security model treats identity, authorization, and tenant isolation as separate concerns.

* OAuth/OIDC establishes user identity.
* RBAC provides explicit permissions.
* Contextual authorization can provide additional policy decisions.
* Organizations act as tenant boundaries.
* Tenant-scoped access is enforced at the application and database layers.
* PostgreSQL Row-Level Security is planned for defense-in-depth tenant isolation.

These capabilities are introduced incrementally through Phases 2 and 3.

---

# Architecture

The intended architecture separates application responsibilities:

```text
                         React Dashboard
                                |
                                v
                           Go Backend
                         /     |      \
                        /      |       \
                       v       v        v
                 PostgreSQL  Event     Python ML
                             System     Service
                                          |
                                          v
                                  ML / Risk Signals
                                          |
                                          v
                                   Risk Intelligence
```

The architectural ownership is:

| Component            | Responsibility                                                 |
| -------------------- | -------------------------------------------------------------- |
| React                | Presentation and dashboard                                     |
| Go                   | API, policies, orchestration, application-level risk decisions |
| Python               | ML training, artifacts, inference, ML signals                  |
| PostgreSQL           | Relational persistence and tenant isolation                    |
| Event infrastructure | Asynchronous event delivery and processing                     |
| Audit subsystem      | Cryptographic integrity and verification                       |

The Python ML system is deliberately independent from the Go backend in Phase 5.

Phase 6 introduces the integration between ML signals and the Go-owned risk engine.

---

# Technology Stack

| Layer                | Technology                             | Status             |
| -------------------- | -------------------------------------- | ------------------ |
| Backend              | Go                                     | Implemented        |
| Database             | PostgreSQL                             | Implemented        |
| Database access      | pgx/v5                                 | Implemented        |
| Migrations           | golang-migrate                         | Implemented        |
| ML                   | Python                                 | Implemented        |
| ML libraries         | NumPy, pandas, scikit-learn, joblib    | Implemented        |
| Frontend             | React + TypeScript                     | Planned            |
| Event infrastructure | Event streams / asynchronous consumers | Planned            |
| Identity             | OAuth/OIDC                             | Planned            |
| Authorization        | RBAC, contextual rules, PostgreSQL RLS | Planned            |
| Risk intelligence    | Go + Python ML signals                 | Planned in Phase 6 |
| Audit integrity      | Hash chains / Merkle trees             | Planned            |
| Real-time delivery   | WebSockets                             | Planned            |

---

# Current Status

## Phase 0 — Foundation and Development Environment

**Status: Complete**

The project foundation and development structure have been established.

## Phase 1 — Core Go Backend and PostgreSQL Data Layer

**Status: Complete**

Implemented:

* Environment configuration and validation
* Safe database URL redaction
* PostgreSQL connection pooling
* Database readiness checks
* `GET /health/live`
* `GET /health/ready`
* Graceful HTTP shutdown
* SQL migrations
* Unit tests
* PostgreSQL integration tests

## Phase 2 — Identity and Authentication

**Status: Backend Complete, Frontend Deferred**

Implemented (Backend):

* OAuth 2.0 / OpenID Connect provider configuration and discovery
* Internal user identity mapping
* Persistent PostgreSQL sessions
* Login, callback handling, and HTTP-only cookie management
* Authentication middleware and authenticated request context
* Safe identity retrieval (`GET /auth/me`)
* Idempotent logout

*Note: Frontend authentication state is explicitly DEFERRED until a frontend application framework is added.*

## Phase 3 — Authorization and Multi-Tenancy

**Status: Planned**

Planned:

* Organizations and tenants
* Memberships
* Roles and permissions
* RBAC
* Tenant context
* Tenant-scoped data access
* PostgreSQL Row-Level Security

## Phase 4 — Event Domain and Ingestion

**Status: Planned**

Planned:

* Canonical event schema
* Event validation
* Event creation
* Event persistence
* Event retrieval
* Event filtering
* Event lifecycle states

## Phase 5 — Independent Hybrid Machine Learning System

**Status: Complete**

Implemented:

* Dataset loading
* Feature selection and preprocessing
* Supervised Decision Tree classifier
* Isolation Forest anomaly detection
* Behavioral baseline scoring
* Autoencoder anomaly detection
* Signal calibration
* Hybrid signal aggregation
* Model artifact persistence
* Artifact loading
* Independent event inference
* Unified model training
* Evaluation across normal and suspicious events

The ML service is independent of the Go backend.

## Phase 6 — Adaptive Risk Intelligence

**Status: Next**

Planned:

* Event-to-feature mapping
* ML signal integration
* Behavioral context
* Contextual security signals
* Go-owned risk aggregation
* Risk scoring
* Explainable risk factors
* Backend risk retrieval

## Phase 7 — Cryptographic Audit Integrity

**Status: Planned**

Planned:

* Canonical audit records
* Cryptographic hashing
* Hash chains
* Chain verification
* Integrity reporting
* Optional Merkle-tree verification

## Phase 8 — Asynchronous Event-Driven Processing

**Status: Planned**

Planned:

* Event publishing
* Background consumers
* Retry handling
* Idempotent processing
* Deferred risk processing
* Queue or stream integration

## Phase 9 — Real-Time Analytics and Dashboard Delivery

**Status: Planned**

Planned:

* WebSocket delivery
* Live risk updates
* Event analytics
* Risk trend visualization
* Model metrics
* Audit integrity visibility

## Phase 10 — Observability, Security Hardening, and Release Readiness

**Status: Planned**

Planned:

* Structured logging
* Metrics
* Tracing
* Rate limiting
* Security hardening
* CI/CD
* Test consolidation
* Deployment readiness
* Production documentation

---

# Roadmap

The project follows the principle that complexity should be earned one architectural responsibility at a time.

| Phase | Focus                                                    | Status       |
| ----- | -------------------------------------------------------- | ------------ |
| 0     | Foundation and Development Environment                   | Complete     |
| 1     | Core Go Backend and PostgreSQL Data Layer                | Complete     |
| 2     | Identity and Authentication                              | Backend Complete |
| 3     | Authorization and Multi-Tenancy                          | Planned      |
| 4     | Event Domain and Ingestion                               | Planned      |
| 5     | Independent Hybrid Machine Learning System               | **Complete** |
| 6     | Adaptive Risk Intelligence                               | Next         |
| 7     | Cryptographic Audit Integrity                            | Planned      |
| 8     | Asynchronous Event-Driven Processing                     | Planned      |
| 9     | Real-Time Analytics and Dashboard Delivery               | Planned      |
| 10    | Observability, Security Hardening, and Release Readiness | Planned      |

### Dependency Chain

```text
Phase 0
   ↓
Phase 1
   ↓
Phase 2
   ↓
Phase 3
   ↓
Phase 4
   ↓
Phase 5 ✓
   ↓
Phase 6 ← Next
   ↓
Phase 7
   ↓
Phase 8
   ↓
Phase 9
   ↓
Phase 10
```

Each phase introduces one dominant capability and creates the minimum dependency required for the next phase.

---

# Project Structure

```text
.
├── backend/
│   ├── cmd/api/                 # Backend entrypoint
│   ├── internal/health/         # Liveness and readiness handlers
│   ├── internal/platform/       # Configuration, server, database lifecycle
│   ├── migrations/              # Up/down SQL migrations
│   └── README.md                # Backend development notes
│
├── ml-service/
│   ├── data/
│   │   └── events.csv           # ML training/evaluation dataset
│   │
│   ├── models/                  # Saved ML artifacts
│   │
│   ├── training/
│   │   ├── preprocessing.py
│   │   ├── train_classifier.py
│   │   ├── evaluate_classifier.py
│   │   ├── train_isolation_forest.py
│   │   ├── behavioral_baseline.py
│   │   ├── train_autoencoder.py
│   │   ├── train_all.py
│   │   └── artifacts.py
│   │
│   ├── inference/
│   │   ├── signal_calibration.py
│   │   ├── signals.py
│   │   ├── ensemble.py
│   │   ├── ml_signal.py
│   │   ├── predict.py
│   │   └── test_inference.py
│   │
│   ├── evaluation/
│   │   └── evaluate_all.py
│   │
│   └── requirements.txt
│
├── docs/
│   └── 01_OVERVIEW.md           # Project vision and architecture
│
└── README.md
```

The Python virtual environment is intentionally excluded from version control.

---

# Quick Start

## Backend

### Prerequisites

* Go
* PostgreSQL
* golang-migrate
* Docker, if running PostgreSQL locally in a container

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

The API listens on:

```text
http://localhost:8080
```

### Check health

```bash
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

---

# ML Service

The ML service is independently executable from `ml-service/`.

## Install dependencies

```bash
cd ml-service

python -m venv .venv
source .venv/bin/activate

pip install -r requirements.txt
```

## Train all models

```bash
cd training
python train_all.py
```

This trains the four implemented detection components and saves their artifacts into:

```text
ml-service/models/
```

## Test independent inference

```bash
cd ../inference
python test_inference.py
```

Inference accepts an event containing:

```text
hour
failed_logins
requests_per_minute
files_downloaded
```

and produces:

```text
Classifier Signal
Isolation Forest Signal
Behavioral Signal
Autoencoder Signal
Hybrid Score
Interpretation
```

## Evaluate the complete ML pipeline

```bash
cd ../evaluation
python evaluate_all.py
```

This evaluates the saved inference pipeline against the available dataset and compares normal and suspicious events across all implemented signal families.

---

# Configuration

The backend currently supports:

| Variable             | Required | Default       | Description                  |
| -------------------- | -------: | ------------- | ---------------------------- |
| `AEGIS_DATABASE_URL` |      Yes | —             | PostgreSQL connection string |
| `AEGIS_SERVER_ADDR`  |       No | `:8080`       | HTTP bind address and port   |
| `AEGIS_ENV`          |       No | `development` | Runtime environment          |

The backend fails fast when required configuration is missing or invalid.

---

# Testing

## Backend

From `backend/`:

```bash
go test -short ./...
```

Runs unit tests while skipping tests requiring PostgreSQL.

```bash
go test -v ./...
```

Runs the full test suite when a suitable PostgreSQL test database is available.

```bash
go fmt ./...
go mod tidy
```

## ML

From `ml-service/`:

```bash
cd training
python train_all.py
```

Then:

```bash
cd ../inference
python test_inference.py
```

And:

```bash
cd ../evaluation
python evaluate_all.py
```

---

# Engineering Principles

Aegis is being built around several architectural principles:

### Independent capabilities first

Each major capability is implemented and tested independently before integration.

### ML does not own the final security decision

The Python service produces ML signals.

The eventual Go risk engine owns the application-level risk decision.

```text
Python ML
   ↓
ML Signals
   ↓
Go Risk Aggregator
   ↓
Final Risk Assessment
```

### Multiple detection perspectives

Aegis does not depend on one anomaly detector.

The ML layer combines:

* Supervised classification
* Global anomaly detection
* Behavioral deviation
* Reconstruction-based anomaly detection

### Defense in depth

Identity, authorization, tenant isolation, risk intelligence, and audit integrity remain separate responsibilities.

### Earn complexity gradually

Queues, WebSockets, cryptographic verification, contextual risk, and production infrastructure are introduced only after their prerequisite capabilities are stable.

---

# Documentation

* Project overview and long-term architecture: `docs/01_OVERVIEW.md`
* Backend setup, health checks, migrations, and testing: `backend/README.md`
* Phase roadmap: `03_PHASE_PLANNING.md`

---

# License

No license has been declared yet.
