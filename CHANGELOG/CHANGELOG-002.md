# Changelog

## Phase 2: Asynchronous Email Processing

### Part 1 — Worker Pool

* Added configurable worker pool.
* Implemented concurrent background workers using goroutines.
* Added graceful worker startup and shutdown.
* Integrated `sync.WaitGroup` for worker lifecycle management.
* Added structured logging for worker events.

### Part 2 — In-Memory Queue

* Introduced producer/consumer queue abstraction.
* Implemented in-memory job queue using Go channels.
* Added publisher and consumer interfaces.
* Added typed email jobs.
* Connected API requests to asynchronous job processing.

### Part 3 — Email Processor

* Added processor abstraction.
* Implemented `EmailProcessor`.
* Separated job execution from HTTP request handling.
* Added helper methods for loading and updating email state.
* Improved processor logging.

### Part 4 — Email Lifecycle

* Implemented complete email state transitions:
  * `queued`
  * `processing`
  * `sent`
  * `failed`
* Added repository support for updating email status.
* Added status transition logging.

### Part 5 — Sender Abstraction

* Introduced sender interface.
* Added `NoopSender` for local development.
* Decoupled processing logic from email delivery implementation.
* Prepared the codebase for SMTP providers.

### Part 6 — Idempotency

* Added support for the `Idempotency-Key` HTTP header.
* Added idempotency key persistence.
* Implemented lookup by idempotency key.
* Prevented duplicate email creation on client retries.
* Returned existing email responses for duplicate requests.
* Added automatic idempotency key generation when the client does not provide one.
* Prepared database for unique idempotency constraints.

### Part 7 — Repository Improvements

* Added `FindByIdempotencyKey`.
* Extended repository interfaces.
* Improved repository error handling.
* Refactored persistence layer for idempotent operations.

### Part 8 — Service Improvements

* Refactored email creation workflow.
* Added duplicate request detection.
* Improved separation between validation, persistence, and queue publishing.
* Prevented duplicate job publication.

### Part 9 — Testing

* Verified worker pool concurrency.
* Load-tested API using `hey`.
* Validated asynchronous processing pipeline.
* Verified idempotent request handling.
* Confirmed correct status transitions.
* Tested graceful shutdown behavior.


## Current Architecture

```
HTTP API
    │
    ▼
Gin Handler
    │
    ▼
Email Service
    │
    ├── Validation
    ├── Idempotency
    ├── Repository
    └── Queue Publisher
             │
             ▼
      In-Memory Queue
             │
             ▼
        Worker Pool
             │
             ▼
      Email Processor
             │
             ▼
      Sender Interface
             │
             ▼
        Email Provider
```

## Upcoming

### Phase 3 — Durable Messaging

* Redis Streams
* Consumer groups
* Message acknowledgements
* Durable job storage
* Queue recovery after crashes

### Phase 4 — Delivery Reliability

* SMTP integration
* Retry policies
* Exponential backoff
* Dead-letter queue
* Delivery tracking

### Phase 5 — Observability

* Prometheus metrics
* Grafana dashboards
* OpenTelemetry tracing
* Queue metrics
* Worker metrics

### Phase 6 — Production Readiness

* Rate limiting
* Authentication
* API keys
* Webhooks
* Horizontal scaling
* Multi-provider email failover