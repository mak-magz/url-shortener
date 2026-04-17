# Requirements Document

## Introduction

Phase 2 of SwiftLink containerises the entire application stack. The goal is to produce Docker images for both the Go backend and the Nuxt 4 frontend using multi-stage builds, wire them together with PostgreSQL in a single Docker Compose file, and expose the stack through an Nginx reverse proxy. After this phase, the full application can be started with a single `docker compose up` command with no local Go or Node.js toolchain required.

## Glossary

- **Backend_Image**: The Docker image built from `backend/Dockerfile` that runs the Go/Gin API server.
- **Frontend_Image**: The Docker image built from `frontend/Dockerfile` that runs the Nuxt 4 SSR server.
- **Compose_Stack**: The set of services defined in the root `docker-compose.yml` file (frontend, backend, postgres, nginx).
- **Nginx_Proxy**: The Nginx reverse-proxy container that routes HTTP traffic to the frontend or backend service.
- **Builder_Stage**: The intermediate Docker build stage that compiles or bundles application code; its artefacts are copied into the final image.
- **Runner_Stage**: The final, minimal Docker image stage that contains only the compiled binary or built assets needed to run the service.
- **Health_Check**: A Docker-native `HEALTHCHECK` instruction or Compose `healthcheck` block that periodically tests whether a container is ready to serve traffic.
- **Postgres_Service**: The `postgres:16-alpine` container that provides the database for the Compose stack.
- **pgdata_Volume**: The named Docker volume that persists PostgreSQL data across container restarts.
- **env_file**: A `.env` file (or Compose `environment` block) that supplies runtime configuration to a container without baking secrets into the image.

---

## Requirements

### Requirement 1: Backend Multi-Stage Dockerfile

**User Story:** As a developer, I want a multi-stage Dockerfile for the backend, so that the production image is minimal and contains only the compiled binary.

#### Acceptance Criteria

1. THE `Backend_Image` SHALL use a `golang:1.26-alpine` base image in the `Builder_Stage`.
2. WHEN the `Builder_Stage` runs, THE `Backend_Image` SHALL download Go module dependencies before copying application source code, so that the dependency layer is cached independently.
3. WHEN the `Builder_Stage` compiles the binary, THE `Backend_Image` SHALL set `CGO_ENABLED=0` and `GOOS=linux` to produce a statically linked binary.
4. THE `Runner_Stage` SHALL use `alpine:3.22` as its base image.
5. THE `Runner_Stage` SHALL copy only the compiled binary from the `Builder_Stage`; no Go toolchain or source files SHALL be present in the final image.
6. THE `Backend_Image` SHALL expose port `8080`.
7. THE `Backend_Image` SHALL define a `HEALTHCHECK` that calls `GET /ping` on `localhost:8080` with a 5-second interval, 3-second timeout, and 3 retries before marking the container unhealthy.
8. THE `Backend_Image` SHALL set a non-root user (`appuser`) as the runtime user in the `Runner_Stage`.

### Requirement 2: Frontend Multi-Stage Dockerfile

**User Story:** As a developer, I want a multi-stage Dockerfile for the frontend, so that the production image contains only the Nuxt build output and its Node.js runtime dependencies.

#### Acceptance Criteria

1. THE `Frontend_Image` SHALL use a `node:22-alpine` base image in the `Builder_Stage`.
2. THE `Frontend_Image` SHALL install `pnpm` via `corepack enable && corepack prepare pnpm@latest --activate` in the `Builder_Stage`.
3. WHEN the `Builder_Stage` installs dependencies, THE `Frontend_Image` SHALL copy `package.json` and `pnpm-lock.yaml` before copying the rest of the source, so that the dependency layer is cached independently.
4. WHEN the `Builder_Stage` builds the application, THE `Frontend_Image` SHALL run `pnpm run build` to produce the Nuxt output bundle.
5. THE `Runner_Stage` SHALL use `node:22-alpine` as its base image.
6. THE `Runner_Stage` SHALL copy only the `.output` directory from the `Builder_Stage`; no build toolchain or source files SHALL be present in the final image.
7. THE `Frontend_Image` SHALL expose port `3000`.
8. THE `Frontend_Image` SHALL set a non-root user (`appuser`) as the runtime user in the `Runner_Stage`.
9. THE `Frontend_Image` SHALL define a `HEALTHCHECK` that calls `GET /` on `localhost:3000` with a 10-second interval, 5-second timeout, and 3 retries before marking the container unhealthy.

### Requirement 3: Docker Compose Full-Stack Deployment

**User Story:** As a developer, I want a single `docker compose up` command to start the entire SwiftLink stack, so that I can run the full application locally without installing Go or Node.js.

#### Acceptance Criteria

1. THE `Compose_Stack` SHALL define four services: `postgres`, `backend`, `frontend`, and `nginx`.
2. THE `Compose_Stack` SHALL build the `backend` service from `./backend/Dockerfile` and the `frontend` service from `./frontend/Dockerfile`.
3. THE `Postgres_Service` SHALL use the `postgres:16-alpine` image and mount the `pgdata_Volume` to persist data.
4. WHEN the `backend` service starts, THE `Compose_Stack` SHALL ensure the `Postgres_Service` is healthy before the `backend` container begins accepting traffic, using a `depends_on` condition of `service_healthy`.
5. WHEN the `frontend` service starts, THE `Compose_Stack` SHALL ensure the `backend` service is healthy before the `frontend` container begins accepting traffic, using a `depends_on` condition of `service_healthy`.
6. THE `Compose_Stack` SHALL supply runtime environment variables to the `backend` service via an `env_file` reference to `./backend/.env`.
7. THE `Compose_Stack` SHALL supply runtime environment variables to the `frontend` service via an `env_file` reference to `./frontend/.env`.
8. THE `Compose_Stack` SHALL NOT expose the `backend` or `frontend` service ports directly to the host; all external traffic SHALL be routed through the `Nginx_Proxy`.
9. THE `Nginx_Proxy` service SHALL expose port `80` on the host.
10. THE `Compose_Stack` SHALL define a named `pgdata_Volume` so that database data persists across `docker compose down` and `docker compose up` cycles.

### Requirement 4: Nginx Reverse Proxy Configuration

**User Story:** As a developer, I want Nginx to route requests to the correct service, so that the frontend and backend are accessible through a single port without CORS issues.

#### Acceptance Criteria

1. THE `Nginx_Proxy` SHALL route all requests to `/api/` to the `backend` service on port `8080`.
2. THE `Nginx_Proxy` SHALL route all requests to `/:shortCode` (i.e. paths that are not `/api/`) to the `backend` service on port `8080`, so that short-link redirects work through the proxy.
3. THE `Nginx_Proxy` SHALL route all remaining requests to the `frontend` service on port `3000`.
4. THE `Nginx_Proxy` SHALL pass the `Host`, `X-Real-IP`, and `X-Forwarded-For` headers to upstream services on every proxied request.
5. THE `Nginx_Proxy` SHALL use the Docker Compose service names (`backend`, `frontend`) as upstream hostnames, relying on Docker's internal DNS resolution.
6. WHEN an upstream service is unavailable, THE `Nginx_Proxy` SHALL return an HTTP `502 Bad Gateway` response.

### Requirement 5: Environment Variable Configuration

**User Story:** As a developer, I want all service configuration to come from environment variables, so that the same images can be used in different environments without rebuilding.

#### Acceptance Criteria

1. THE `Backend_Image` SHALL read `DATABASE_URL`, `SERVER_HOST`, and `SERVER_PORT` from environment variables at runtime; no default values SHALL be baked into the image.
2. THE `Frontend_Image` SHALL read `BACKEND_API_BASE_URL`, `BACKEND_API_VERSION`, `BACKEND_API_PREFIX`, and `APP_BASE_URL` from environment variables at runtime.
3. THE `Compose_Stack` SHALL provide a `.env.example` file at the repository root documenting all required environment variables for each service.
4. IF a required environment variable is missing at container startup, THEN THE `Backend_Image` SHALL exit with a non-zero status code and a descriptive error message.

### Requirement 6: Image Hygiene and Security

**User Story:** As a developer, I want Docker images to be minimal and secure, so that the attack surface and image size are kept small.

#### Acceptance Criteria

1. THE `Backend_Image` SHALL NOT include the Go toolchain, source code, test files, or `.env` files in the `Runner_Stage`.
2. THE `Frontend_Image` SHALL NOT include the Node.js build toolchain, source code, or `.env` files in the `Runner_Stage`.
3. THE `Backend_Image` SHALL include a `.dockerignore` file that excludes `.git`, `.env`, `tmp/`, and test files from the build context.
4. THE `Frontend_Image` SHALL include a `.dockerignore` file that excludes `.git`, `.env`, `.nuxt/`, `node_modules/`, and test directories from the build context.
5. WHEN the `Runner_Stage` executes, THE `Backend_Image` SHALL run as a non-root user.
6. WHEN the `Runner_Stage` executes, THE `Frontend_Image` SHALL run as a non-root user.
