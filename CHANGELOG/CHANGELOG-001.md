# Changelog

All notable changes to the project are documented here.

---

## Phase 1 — Production API Foundation

### Part 1 — Project Bootstrap

* Initialized Go module.
* Created project layout.
* Added Docker Compose.
* Added PostgreSQL.
* Added SQL migrations.
* Added configuration management.

### Part 2 — Logging & Middleware

* Integrated Zap logger.
* Added request ID middleware.
* Added request logging middleware.
* Added panic recovery middleware.

### Part 3 — Database Layer

* Configured pgx connection pool.
* Implemented Email repository.
* Added repository interfaces.
* Introduced Email model.

### Part 4 — Business Layer

* Added DTOs.
* Added validation.
* Implemented Email service.
* Introduced application error types.

### Part 5 — HTTP Layer

* Added Email handler.
* Added Health handler.
* Added centralized handler container.
* Added versioned routing.
* Added consistent JSON responses.

### Part 6 — Application Bootstrap

* Introduced `internal/app`.
* Centralized dependency injection.
* Added HTTP server configuration.
* Added graceful shutdown.
* Added signal handling.
* Added readiness and health endpoints.

---

## Upcoming

### Phase 2

* Worker pools.
* Go channels.
* In-memory queue.
* Background processing.
* Retry logic.
* SMTP abstraction.

### Phase 3

* Redis-backed queue.
* Scheduled retries.
* Dead-letter queue.
* Rate limiting.

### Phase 4

* Metrics.
* Prometheus.
* OpenTelemetry.
* Distributed tracing.
* Grafana dashboards.

### Phase 5

* Multi-provider failover.
* Domain verification.
* DKIM signing.
* Bounce handling.
* Webhooks.
* Horizontal scaling.
