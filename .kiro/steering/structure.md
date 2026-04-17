# Project Structure

## Repository Root

```
url-shortener/
├── backend/          # Go API server
├── frontend/         # Nuxt 4 frontend
├── docker-compose.yml  # Local PostgreSQL
└── .env.example
```

---

## Backend (`backend/`)

```
backend/
├── cmd/api/main.go           # Entry point — wires up router, DB, middleware, handlers
├── internal/
│   └── url/
│       ├── handler/          # HTTP handlers (Gin) — parse request, call service, return response
│       ├── service/          # Business logic — interfaces + implementations
│       ├── repository/       # DB access — interfaces + pgx implementations
│       └── model/            # Structs: domain models, request/response types
└── platform/
    ├── config/               # Config loading from env
    ├── db/                   # pgxpool setup + inline SQL migrations
    ├── errors/               # AppError type + constructors (NotFound, Internal, BadRequest)
    └── middleware/           # Gin middleware: CORS, error handler, logger
```

### Backend Architecture Patterns

- **Layered architecture**: Handler → Service → Repository. Each layer depends only on the interface of the layer below.
- **Interface-driven**: `Handler`, `Service`, and `URLRepository` are all Go interfaces. Concrete types are unexported; constructors return the interface.
- **Error handling**: Handlers call `c.Error(err)` to attach errors; `ErrorMiddleware` intercepts and serializes `AppError` to JSON. Never write error responses directly in handlers.
- **DB migrations**: Run inline at startup via `db.Migrate(pool)` — no external migration tool.
- **No ORM**: Raw SQL with `pgx`. Queries are written directly in repository methods.

---

## Frontend (`frontend/`)

```
frontend/
├── app/
│   ├── app.vue               # Root component (layout + NuxtPage)
│   ├── app.config.ts         # App-level config (theme, etc.)
│   ├── assets/css/main.css   # Global styles
│   ├── components/           # Vue components, organized by feature/section
│   ├── layouts/              # Nuxt layouts (default, admin, user)
│   ├── pages/                # File-based routing
│   ├── stores/               # Pinia stores
│   └── utils/                # Shared utilities (e.g. validators.ts)
├── test/
│   ├── unit/                 # Vitest unit tests
│   └── nuxt/                 # Vitest Nuxt component tests
├── tests/                    # Playwright e2e tests
├── nuxt.config.ts            # Nuxt configuration
└── vitest.config.ts          # Vitest configuration (two projects: unit, nuxt)
```

### Frontend Architecture Patterns

- **File-based routing**: Pages live in `app/pages/`. Nested folders = nested routes.
- **Nuxt UI components**: Use `U*` components (e.g. `UButton`, `UApp`) from `@nuxt/ui` before writing custom ones.
- **Icons**: Use Iconify via Nuxt Icon — `i-lucide-*` for Lucide icons, `i-simple-icons-*` for brand icons.
- **Async data / queries**: Use `@pinia/colada` (`useQuery`, `useMutation`) for server data fetching, not raw `useFetch` where possible.
- **Validation**: Use Valibot schemas for form validation.
- **ESLint style**: Tabs for indentation, no trailing commas, `1tbs` brace style (configured in `nuxt.config.ts`).
- **Runtime config**: Backend API coordinates come from `useRuntimeConfig().public.*` — never hardcode URLs.
