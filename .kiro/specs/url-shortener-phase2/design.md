# Design Document: SwiftLink Phase 2 — Docker Infrastructure

## Overview

Phase 2 containerises the entire SwiftLink stack so the full application can be started with a single `docker compose up` command. The deliverables are:

- A production-grade multi-stage `backend/Dockerfile` (Go binary, non-root, health check)
- A new multi-stage `frontend/Dockerfile` (Nuxt `.output`, non-root, health check)
- An updated root `docker-compose.yml` wiring all four services together with proper startup ordering
- An `nginx/` directory containing the reverse-proxy configuration
- An updated `.env.example` documenting every required variable
- Updated `.dockerignore` files for both services

No application logic changes are required. All work is infrastructure configuration.

---

## Architecture

```mermaid
flowchart TD
    Client["Browser / curl"] -->|"port 80"| Nginx["nginx:alpine (Nginx_Proxy)"]
    Nginx -->|"/api/*"| Backend["backend:8080 (Backend_Image)"]
    Nginx -->|"/:shortCode"| Backend
    Nginx -->|"/* default"| Frontend["frontend:3000 (Frontend_Image)"]
    Backend -->|"TCP 5432"| Postgres["postgres:16-alpine (Postgres_Service)"]
```

### Startup Dependency Chain

```mermaid
flowchart LR
    Postgres["postgres (healthy)"] --> Backend["backend (healthy)"]
    Backend --> Frontend["frontend (healthy)"]
    Frontend --> Nginx["nginx"]
```

Docker Compose `depends_on` with `condition: service_healthy` enforces this order. Nginx has no explicit `depends_on` — it starts immediately and relies on its own upstream error handling (502) until the other services are ready.

### Network Topology

All four services share a single user-defined bridge network (`swiftlink-net`) created implicitly by Compose. The `backend` and `frontend` services do **not** publish ports to the host; only `nginx` publishes port `80`. The `postgres` service publishes port `5432` to the host for local development convenience (can be removed in production).

---

## Components and Interfaces

### Backend Dockerfile (`backend/Dockerfile`)

Two-stage build:

| Stage | Base Image | Purpose |
|-------|-----------|---------|
| `builder` | `golang:1.26-alpine` | Download deps, compile static binary |
| `runner` | `alpine:3.22` | Minimal runtime image |

Key decisions:
- `COPY go.mod go.sum ./` then `RUN go mod download` before `COPY . .` — dependency layer is cached independently of source changes.
- `CGO_ENABLED=0 GOOS=linux` produces a fully static binary that runs on `alpine` without glibc.
- `adduser -D -u 1001 appuser` creates a non-root user; `USER appuser` activates it before `CMD`.
- `HEALTHCHECK` uses `wget -qO- http://localhost:8080/ping` (wget is available on alpine; curl is not installed by default).

### Frontend Dockerfile (`frontend/Dockerfile`)

Two-stage build:

| Stage | Base Image | Purpose |
|-------|-----------|---------|
| `builder` | `node:22-alpine` | Install pnpm, install deps, run `pnpm run build` |
| `runner` | `node:22-alpine` | Serve Nuxt SSR output |

Key decisions:
- `corepack enable && corepack prepare pnpm@latest --activate` installs pnpm without a separate `npm install -g pnpm` step.
- `COPY package.json pnpm-lock.yaml ./` then `RUN pnpm install --frozen-lockfile` before `COPY . .` — dependency layer cached independently.
- Only `.output/` is copied into the runner stage. The Nuxt `.output` directory is self-contained and includes the Node.js server entry point at `.output/server/index.mjs`.
- Runner CMD: `node .output/server/index.mjs`
- `HEALTHCHECK` uses `wget -qO- http://localhost:3000/` (wget available on node:22-alpine).
- Non-root user created with `adduser -D -u 1001 appuser`.

### Docker Compose (`docker-compose.yml`)

Replaces the existing minimal Compose file (which only ran postgres) with a full four-service stack.

| Service | Image / Build | Ports (host) | Health Check |
|---------|--------------|-------------|-------------|
| `postgres` | `postgres:16-alpine` | `5432:5432` | `pg_isready` |
| `backend` | `./backend/Dockerfile` | none | `GET /ping` |
| `frontend` | `./frontend/Dockerfile` | none | `GET /` |
| `nginx` | `nginx:alpine` | `80:80` | none |

`env_file` directives point to `./backend/.env` and `./frontend/.env` respectively. The `pgdata` named volume persists database data.

### Nginx Configuration (`nginx/nginx.conf`)

Single `server` block on port 80 with three `location` blocks:

```
/api/          → proxy_pass http://backend:8080
/:shortCode    → proxy_pass http://backend:8080  (regex location)
/              → proxy_pass http://frontend:3000  (default)
```

Routing rationale:
- `/api/` prefix covers all REST endpoints (`POST /api/v1/shorten`).
- Short-code redirects (`GET /:shortCode`) must also hit the backend. A regex location `~* ^/[a-zA-Z0-9]{6}$` matches exactly 6-character alphanumeric paths and routes them to the backend before the default catch-all.
- All other traffic (the Nuxt SSR pages) falls through to the frontend.

Standard proxy headers (`Host`, `X-Real-IP`, `X-Forwarded-For`, `X-Forwarded-Proto`) are set on every proxied request.

### `.dockerignore` Files

**`backend/.dockerignore`** (update existing):
```
.git
.env
tmp/
**/*_test.go
*.md
.agent
.idea
.vscode
```

**`frontend/.dockerignore`** (new):
```
.git
.env
.nuxt/
node_modules/
test/
tests/
*.md
.agent
```

---

## Data Models

This phase introduces no new application data models. The only "data" concerns are:

- **`pgdata` volume**: Named Docker volume at `/var/lib/postgresql/data` inside the postgres container. Persists across `docker compose down` / `up` cycles. Destroyed only by `docker compose down -v`.
- **Environment variables**: Documented in `.env.example`. No new variables are introduced; existing variables are wired into the Compose stack.

### Environment Variable Reference

| Service | Variable | Description | Example |
|---------|----------|-------------|---------|
| backend | `DATABASE_URL` | PostgreSQL connection string | `postgres://urlshortener:secret123@postgres:5432/urlshortener?sslmode=disable` |
| backend | `SERVER_HOST` | Bind host | `0.0.0.0` |
| backend | `SERVER_PORT` | Bind port | `8080` |
| frontend | `BACKEND_API_BASE_URL` | Backend base URL (internal, via nginx) | `http://nginx` |
| frontend | `BACKEND_API_VERSION` | API version segment | `v1` |
| frontend | `BACKEND_API_PREFIX` | API path prefix | `api` |
| frontend | `APP_BASE_URL` | Public base URL for short links | `http://localhost` |

> **Note on `DATABASE_URL` in Compose**: The host must be `postgres` (the Compose service name), not `localhost`. The `backend/.env` file used for local development uses `localhost:5432`; for the Compose stack the env var is overridden in the Compose file's `environment` block or the `.env` file is updated to use the service name.

---

## Correctness Properties

This feature consists entirely of Dockerfile, Docker Compose, and Nginx configuration files. All acceptance criteria are infrastructure configuration checks (SMOKE) or single-example behavioral checks (EXAMPLE). There are no pure functions with wide input spaces, no parsers, no serializers, and no business logic transformations.

Property-based testing is **not applicable** to this feature. The appropriate testing strategy is:
- Smoke tests: inspect configuration files for required directives
- Integration tests: `docker compose up` and verify services start healthy
- Example-based tests: verify the one behavioral criterion (missing env var causes non-zero exit)

---

## Error Handling

### Missing Environment Variables

The backend's `config.LoadConfig()` already calls `log.Fatal("DATABASE_URL is not set")` when `DATABASE_URL` is empty, which exits with status 1. No changes needed to the application code.

The frontend's Nuxt `runtimeConfig` reads from `process.env` at build time and runtime. Missing variables result in `undefined` values; the frontend does not crash on startup but API calls will fail. This is acceptable for Phase 2 — a future phase can add startup validation.

### Nginx Upstream Unavailability

When a proxied upstream is unreachable, Nginx returns `502 Bad Gateway` by default. No custom error page configuration is required for Phase 2.

### Container Startup Ordering

The `depends_on: condition: service_healthy` chain ensures the backend does not start until postgres passes its health check, and the frontend does not start until the backend passes its health check. If a service fails its health check repeatedly, Compose will not start the dependent service, and the operator will see a clear error in `docker compose ps`.

### Build Failures

- If `go build` fails in the backend builder stage, the Docker build exits non-zero and no image is produced.
- If `pnpm run build` fails in the frontend builder stage, the Docker build exits non-zero and no image is produced.

---

## Testing Strategy

PBT is not applicable to this feature (pure infrastructure configuration). The testing approach is:

### Smoke Tests (Configuration Inspection)

These are manual or CI-automated checks that inspect file content:

1. **Dockerfile correctness**: Verify base images, layer ordering, build flags, EXPOSE, HEALTHCHECK, USER instructions match requirements.
2. **docker-compose.yml correctness**: Verify four services defined, correct build contexts, `depends_on` conditions, `env_file` references, no exposed ports for backend/frontend, nginx on port 80, named volume.
3. **nginx.conf correctness**: Verify location blocks, proxy_pass targets, proxy header directives.
4. **.dockerignore correctness**: Verify exclusion patterns for both services.
5. **.env.example completeness**: Verify all required variables are documented.

### Integration Test

Run the full stack and verify end-to-end:

```bash
docker compose up --build -d
# Wait for all services to be healthy
docker compose ps
# Verify: POST http://localhost/api/v1/shorten returns 201
# Verify: GET http://localhost/{shortCode} returns 302
# Verify: GET http://localhost/ returns 200 (Nuxt frontend)
docker compose down -v
```

### Example-Based Unit Test (Requirement 5.4)

The backend's missing-env-var behavior is already covered by the existing `config.LoadConfig()` implementation (`log.Fatal` on empty `DATABASE_URL`). A unit test can verify this:

```go
// In platform/config/config_test.go
func TestLoadConfig_MissingDatabaseURL(t *testing.T) {
    // Unset DATABASE_URL, verify log.Fatal is triggered
    // Use a mock logger or os.Exit interceptor
}
```

This is an example-based test (single scenario: DATABASE_URL is empty), not a property test.

### Build Verification

```bash
# Backend image builds successfully
docker build -t swiftlink-backend ./backend

# Frontend image builds successfully  
docker build -t swiftlink-frontend ./frontend

# Full stack starts and all services reach healthy state
docker compose up --build -d && docker compose ps
```
