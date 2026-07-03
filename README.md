# Relay Engine

Relay Engine is a production-inspired email relay service written in Go.

The project is designed as a learning exercise in building backend infrastructure from the ground up while following the architectural principles used by companies such as Resend, Stripe, and Moniepoint.

Rather than focusing on sending emails immediately, the project emphasizes the infrastructure required to operate a reliable, scalable, and observable message relay

## Project Goals
- Build a production-quality Go service.
- Learn idiomatic Go architecture.
- Practice dependency injection and clean layering.
- Build a concurrent worker-based relay engine.
- Explore distributed systems concepts through incremental implementation.

The service is implemented with:
- `Gin` for HTTP routing
- `pgx` for PostgreSQL connection pooling
- `Zap` for structured logging
- `github.com/google/uuid` for email IDs

## What it does

The current implementation includes:
- `GET /ping` - returns `pong` for liveness checks
- `GET /test-db` - inserts a sample email record into PostgreSQL and returns the created ID

Email data is persisted into the `emails` table with fields such as sender, recipient, subject, text, HTML, status, and created timestamp.

## Repository layout

- `cmd/api/main.go` - application entrypoint
- `internal/config` - environment configuration loader
- `internal/database` - PostgreSQL connection pool implementation
- `internal/repository` - email repository with save and lookup functions
- `internal/models` - email domain model
- `internal/router` - Gin router and middleware setup
- `internal/middleware` - request logging and recovery middleware
- `migrations` - SQL migration files for the database schema
- `docker-compose.yml` - local PostgreSQL service definition

## Requirements

- Go 1.25
- PostgreSQL 16 (via Docker Compose or local install)
- `migrate` CLI if you want to run migrations from the included Makefile targets

## Environment variables

The service loads configuration from environment variables. The following variables are used:

- `APP_NAME` - application name
- `APP_ENV` - application environment (development, production, etc.)
- `PORT` - HTTP port for the API
- `DB_HOST` - PostgreSQL host (not used directly by the current code, but available in config)
- `DB_PORT` - PostgreSQL port (not used directly by the current code, but available in config)
- `DB_USER` - PostgreSQL user (not used directly by the current code, but available in config)
- `DB_PASSWORD` - PostgreSQL password (not used directly by the current code, but available in config)
- `DB_NAME` - PostgreSQL database name (not used directly by the current code, but available in config)
- `DATABASE_URL` - PostgreSQL connection URL used by the service

Example `DATABASE_URL`:

```bash
postgres://postgres:password@localhost:5432/relay?sslmode=disable
```

## Quick start

1. Start PostgreSQL with Docker Compose:

```bash
docker compose up -d
```

2. Create the database schema:

```bash
make migrate-up
```

3. Start the API:

```bash
make run
```

4. Verify the service:

```bash
curl http://localhost:8080/ping
```

You may need to export `PORT` and `DATABASE_URL` before running the service if they are not already set:

```bash
export PORT=8080
export DATABASE_URL="postgres://postgres:password@localhost:5432/relay?sslmode=disable"
```

Then run:

```bash
go run cmd/api/main.go
```

## Docker Compose

The provided `docker-compose.yml` starts a PostgreSQL container with:

- user: `postgres`
- password: `password`
- database: `relay`
- exposed port: `5432`

## Database migrations

The repository includes a migration to create the `emails` table:

- `migrations/000001_create_emails_up.sql`
- `migrations/000001_create_emails.down.sql`

Run migrations with:

```bash
make migrate-up
```

Rollback the most recent migration with:

```bash
make migrate-down
```

## API endpoints

- `GET /ping` - health check
- `GET /test-db` - saves a sample email record and returns the generated UUID

## Useful Makefile targets

- `make run` - runs the service via `go run cmd/api/main.go`
- `make test` - runs all Go tests
- `make fmt` - formats Go code
- `make vet` - runs `go vet`
- `make lint` - runs `golangci-lint`
- `make up` - brings up Docker Compose services
- `make down` - tears down Docker Compose services
- `make logs` - follows Docker Compose logs
- `make migrate-up` - apply migrations
- `make migrate-down` - rollback the last migration

## Notes

- The service reads `DATABASE_URL` directly and uses it to open a PostgreSQL connection pool.
- The current `/test-db` endpoint is intended as a temporary verification helper.
- The application does not yet expose a full email sending API; it stores email records for relay testing.
