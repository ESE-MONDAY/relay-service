# Architecture

## Design Philosophy

The project follows a layered architecture that separates HTTP concerns, business logic, and data persistence.

```text
HTTP Request

        │

        ▼

Gin Router

        │

        ▼

Middleware

        │

        ▼

Handler

        │

        ▼

Service

        │

        ▼

Repository

        │

        ▼

PostgreSQL
```

Each layer has a single responsibility.

### Handler

Responsible for:

* Parsing HTTP requests.
* Returning HTTP responses.
* Mapping application errors into HTTP status codes.

Handlers do **not** contain business logic.

### Service

Responsible for:

* Validation.
* Business rules.
* Entity construction.
* Coordination between repositories.

Services know nothing about HTTP.

### Repository

Responsible for:

* Database queries.
* Persistence.
* Retrieval.

Repositories know nothing about business rules.

## Dependency Injection

The application is bootstrapped inside `internal/app`.

```text
Config

↓

Logger

↓

Database

↓

Repositories

↓

Services

↓

Handlers

↓

Router

↓

HTTP Server
```

Keeping dependency construction in one location makes the system easier to understand and extend.

## Request Lifecycle

A request to `POST /v1/emails` follows this path:

```text
Client

↓

Gin Router

↓

Middleware

↓

EmailHandler

↓

EmailService

↓

EmailRepository

↓

PostgreSQL

↓

HTTP Response
```

## Production Decisions

The project intentionally includes several production-oriented practices:

* Environment-based configuration.
* Structured logging.
* Dependency injection.
* Graceful shutdown.
* Versioned API routes.
* Health and readiness endpoints.
* Consistent response format.

These decisions make the project easier to maintain, test, and scale as new infrastructure components are introduced.

## Future Evolution

Phase 2 introduces asynchronous processing through worker pools and channels.

Later phases will add:

* Job queues.
* Retry policies.
* SMTP provider abstraction.
* Metrics.
* Distributed tracing.
* Rate limiting.
* Multi-provider failover.
* Dead-letter queues.
* Horizontal scalability.
