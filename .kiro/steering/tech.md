# Tech Stack

## Backend (Go)

- **Language**: Go 1.24
- **Framework**: Gin (`github.com/gin-gonic/gin`)
- **Database**: PostgreSQL 16 via `pgx/v5` connection pool (`pgxpool`)
- **Config**: `godotenv` for `.env` loading
- **Validation**: `go-playground/validator/v10`
- **CORS**: `gin-contrib/cors`
- **Testing**: `go.uber.org/mock` for mocks
- **Module**: `github.com/mak-magz/url-shortener`
- **Hot reload** (dev): `air` (`.air.toml`)

## Frontend (Nuxt/Vue)

- **Framework**: Nuxt 4 with Vue 3, SSR enabled
- **Language**: TypeScript
- **UI Library**: Nuxt UI v4 (built on Tailwind CSS v4)
- **State Management**: Pinia + `@pinia/colada` for async queries
- **Validation**: Valibot
- **Package Manager**: pnpm 10
- **Testing**: Vitest (unit + Nuxt component tests), Playwright (e2e)
- **Linting**: ESLint with Nuxt ESLint config

## Infrastructure

- **Database**: PostgreSQL 16 (Docker Compose for local dev)
- **Container**: Docker with multi-stage builds (Phase 2)

---

## Common Commands

### Backend

```bash
# Run from backend/
go run cmd/api/main.go       # start the API server
go build ./...               # build all packages
go test ./...                # run all tests
air                          # hot-reload dev server
```

### Frontend

```bash
# Run from frontend/
pnpm dev                     # start dev server
pnpm build                   # production build
pnpm lint                    # run ESLint
pnpm typecheck               # TypeScript type check
pnpm test                    # run Vitest (watch mode)
pnpm test:unit               # unit tests only
pnpm test:nuxt               # Nuxt component tests only
pnpm test:e2e                # Playwright e2e tests
```

### Infrastructure

```bash
# From repo root
docker compose up -d         # start PostgreSQL locally
docker compose down          # stop services
```

---

## Environment Variables

### Backend (`backend/.env`)

| Variable       | Description                        |
|----------------|------------------------------------|
| `DATABASE_URL` | PostgreSQL connection string (required) |
| `SERVER_HOST`  | Server host (default: `localhost`) |
| `SERVER_PORT`  | Server port (default: `8080`)      |

### Frontend (`frontend/.env`)

| Variable                | Description                        |
|-------------------------|------------------------------------|
| `BACKEND_API_BASE_URL`  | Base URL of the backend API        |
| `BACKEND_API_VERSION`   | API version (e.g. `v1`)            |
| `BACKEND_API_PREFIX`    | API prefix (e.g. `/api`)           |
| `APP_BASE_URL`          | Public base URL of the frontend    |
