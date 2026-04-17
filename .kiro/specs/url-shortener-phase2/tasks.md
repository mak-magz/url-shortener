# Implementation Plan: SwiftLink Phase 2 — Docker Infrastructure

## Overview

Containerise the full SwiftLink stack so it can be started with a single `docker compose up` command. Tasks are grouped into five logical areas:

1. **Backend Dockerfile** — harden the existing Dockerfile to meet all requirements
2. **Frontend Dockerfile** — create a new multi-stage Dockerfile for the Nuxt app
3. **Docker Compose** — replace the minimal dev Compose file with a full four-service stack
4. **Nginx** — add reverse-proxy configuration
5. **Hygiene & documentation** — `.dockerignore` files, `.env.example`, verification

No application logic changes are required. All work is infrastructure configuration.

---

## Tasks

- [ ] 1. Harden the backend Dockerfile
  - Keep the builder base image as `golang:1.26-alpine` (already up to date; no downgrade needed — `go.mod` declares `go 1.24.2` as the minimum language version, not the builder image version)
  - Keep the runner base image as `alpine:3.22` (already up to date)
  - Add `RUN adduser -D -u 1001 appuser` in the runner stage before the `COPY` instruction
  - Add `USER appuser` after the `COPY` instruction in the runner stage (Requirement 1.8)
  - Add a `HEALTHCHECK` instruction: `HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -qO- http://localhost:8080/ping || exit 1` (Requirement 1.7)
  - Verify `CGO_ENABLED=0 GOOS=linux` are already set in the build command (Requirement 1.3) — they are; no change needed
  - Verify `EXPOSE 8080` is present — it is; no change needed
  - _Requirements: 1.3, 1.5, 1.6, 1.7, 1.8_

- [ ] 2. Update backend `.dockerignore`
  - Edit `backend/.dockerignore` to ensure it excludes: `.git`, `.env`, `tmp/`, `**/*_test.go`, `*.md`, `.agent`, `.idea`, `.vscode`
  - The existing file already excludes `.git`, `.env`, `.idea`, `.vscode`, `**/*.md`, `.agent` — add `tmp/` and `**/*_test.go`
  - _Requirements: 6.3_

- [ ] 3. Create the frontend Dockerfile
  - Create `frontend/Dockerfile` with a two-stage multi-stage build
  - **Builder stage** (`FROM node:22-alpine AS builder`):
    - `WORKDIR /app`
    - `RUN corepack enable && corepack prepare pnpm@latest --activate`
    - `COPY package.json pnpm-lock.yaml ./`
    - `RUN pnpm install --frozen-lockfile`
    - `COPY . .`
    - `RUN pnpm run build`
  - **Runner stage** (`FROM node:22-alpine AS runner`):
    - `WORKDIR /app`
    - `RUN adduser -D -u 1001 appuser`
    - `COPY --from=builder /app/.output ./.output`
    - `USER appuser`
    - `EXPOSE 3000`
    - `HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD wget -qO- http://localhost:3000/ || exit 1`
    - `CMD ["node", ".output/server/index.mjs"]`
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 2.9_

- [ ] 4. Create the frontend `.dockerignore`
  - Create `frontend/.dockerignore` with the following exclusions:
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
  - _Requirements: 6.4_

- [ ] 5. Create the Nginx configuration
  - Create directory `nginx/`
  - Create `nginx/nginx.conf` with the following structure:
    ```nginx
    events {}

    http {
        server {
            listen 80;

            # API routes → backend
            location /api/ {
                proxy_pass http://backend:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }

            # Short-code redirects → backend (6-char alphanumeric paths)
            location ~* ^/[a-zA-Z0-9]{6}$ {
                proxy_pass http://backend:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }

            # All other traffic → frontend
            location / {
                proxy_pass http://frontend:3000;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }
        }
    }
    ```
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

- [ ] 6. Replace the root `docker-compose.yml` with the full four-service stack
  - Replace the existing `docker-compose.yml` (which only runs postgres) with a full stack definition:
    - **`postgres` service**: `image: postgres:16-alpine`, `volumes: pgdata:/var/lib/postgresql/data`, `environment` block with `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `healthcheck` using `pg_isready -U urlshortener`, `ports: "5432:5432"` (dev convenience)
    - **`backend` service**: `build: context: ./backend`, `env_file: ./backend/.env`, `depends_on: postgres: condition: service_healthy`, no `ports` mapping
    - **`frontend` service**: `build: context: ./frontend`, `env_file: ./frontend/.env`, `depends_on: backend: condition: service_healthy`, no `ports` mapping
    - **`nginx` service**: `image: nginx:alpine`, `ports: "80:80"`, `volumes: ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro`, `depends_on: frontend, backend`
    - Top-level `volumes: pgdata:`
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 3.10_

- [ ] 7. Update `backend/.env` for Compose networking
  - In `backend/.env`, change `DATABASE_URL` host from `localhost` to `postgres` so the backend container can reach the postgres service by its Compose service name
  - Before: `DATABASE_URL=postgres://urlshortener:secret123@localhost:5432/urlshortener?sslmode=disable`
  - After: `DATABASE_URL=postgres://urlshortener:secret123@postgres:5432/urlshortener?sslmode=disable`
  - Also set `SERVER_HOST=0.0.0.0` so the backend binds on all interfaces inside the container
  - _Requirements: 5.1, 3.4_

- [ ] 8. Update `frontend/.env` for Compose networking
  - In `frontend/.env`, update `BACKEND_API_BASE_URL` to point to the nginx proxy (the public entry point) and `APP_BASE_URL` to the public URL
  - Set `BACKEND_API_BASE_URL=http://nginx` (internal Compose DNS) or `http://localhost` (for SSR requests that go through nginx)
  - Set `APP_BASE_URL=http://localhost` (the host-accessible URL via nginx on port 80)
  - _Requirements: 5.2, 3.7_

- [ ] 9. Create `.env.example` at the repository root
  - Create `.env.example` documenting all required environment variables for both services:
    ```dotenv
    # ============================================================
    # SwiftLink — Environment Variable Reference
    # Copy the relevant sections to backend/.env and frontend/.env
    # ============================================================

    # --- Backend (backend/.env) ---

    # PostgreSQL connection string
    # Use 'postgres' as host when running via docker compose
    DATABASE_URL=postgres://urlshortener:secret123@postgres:5432/urlshortener?sslmode=disable

    # Server bind host (use 0.0.0.0 inside Docker)
    SERVER_HOST=0.0.0.0

    # Server bind port
    SERVER_PORT=8080

    # --- Frontend (frontend/.env) ---

    # Base URL of the backend API (internal, via nginx when using Compose)
    BACKEND_API_BASE_URL=http://localhost

    # API version segment (e.g. v1)
    BACKEND_API_VERSION=v1

    # API path prefix (e.g. api)
    BACKEND_API_PREFIX=api

    # Public base URL used to construct short links in the UI
    APP_BASE_URL=http://localhost
    ```
  - _Requirements: 5.3_

- [ ] 10. Checkpoint — verify images build successfully
  - Run `docker build -t swiftlink-backend ./backend` and confirm it exits 0
  - Run `docker build -t swiftlink-frontend ./frontend` and confirm it exits 0
  - Fix any build errors before proceeding
  - _Requirements: 1.1–1.8, 2.1–2.9, 6.1–6.4_

- [ ] 11. Checkpoint — verify full stack starts and is healthy
  - Run `docker compose up --build -d` from the repository root
  - Run `docker compose ps` and confirm all four services show `healthy` or `running`
  - Verify `GET http://localhost/` returns 200 (Nuxt frontend served through nginx)
  - Verify `POST http://localhost/api/v1/shorten` with `{"originalUrl":"https://example.com"}` returns 201
  - Verify `GET http://localhost/{shortCode}` returns 302 redirect
  - Run `docker compose down -v` to clean up
  - _Requirements: 3.1–3.10, 4.1–4.6_

---

## Notes

- The existing `backend/Dockerfile` already has the correct layer ordering (deps before source) and `CGO_ENABLED=0 GOOS=linux` — only the base image versions and the non-root user / health check need to be added.
- The frontend Nuxt `.output` directory is produced by `pnpm run build` and is self-contained. The runner stage only needs Node.js to execute `.output/server/index.mjs`.
- `wget` is used in health checks (not `curl`) because it is available on `alpine` images without additional installation.
- The `backend/.env` and `frontend/.env` files are used by both local development (with `localhost` hostnames) and the Compose stack. Task 7 and 8 update them for Compose. If you need to switch back to local dev, restore `localhost` as the database host.
- The `postgres` service exposes port `5432` to the host for local development convenience (e.g. running `psql` or a DB GUI). This can be removed in a production deployment.
- The `nginx` service has no explicit health check — it starts immediately and returns 502 until upstream services are ready, which is acceptable since the `depends_on` chain ensures the other services are healthy first.
