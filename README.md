# Private Markets API

## Structure

Layered Architecture
- `handler/` -> presentation layer (HTTP request/response)
- `service/` -> business logic
- `repository/` -> data access / persistence layer
- `models/` -> domain models
- `dto/` -> data transfer objects (HTTP request/response)

This application utilizes `testcontainers` for integration testing. Dependencies such as the Postgresql database is
spun up in a docker container for a single or multiple tests.

## Getting Started

Install prerequisites:
- [Go](https://go.dev/dl/)
- [Docker](https://www.docker.com/) & Docker Compose
- [sqlc](https://sqlc.dev/) (Type-safe data-layer code generation)

---

### 1. Start infrastructure

Spin up Postgres with Docker Compose:

```bash
docker compose up -d
```

### 2. Setup environment variables
Copy `.env.example` values into a `.env`

### 3. Running the API

```bash
go run ./cmd/api
```

Visit the OpenAPI documentation at http://localhost:7000/docs

All migrations are located under `./internal/database/migrations`. These will be automatically applied on api start through [Golang Migrate](https://github.com/golang-migrate/migrate)

To create a new database migration run:
```bash
migrate create -ext sql -dir ./apps/api/internal/database/migrations -seq <migration_name>
```

### 4. Testing

Run unit tests:

```bash
go test -tags=unit -count=1 -cover ./...
```

Run integration tests:

```bash
go test -tags=integration -count=1 -cover ./...
```
