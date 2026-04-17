# Implementation Plan: URL Shortener Phase 1

## Overview

Bring the SwiftLink prototype to production-ready state. Tasks are grouped in three logical phases:

1. **Backend fixes & hardening** — bug fix, validation, crypto short-code, collision retry, rate limiter, health check, error type
2. **New backend features** — wiring rate limiter into the router
3. **Frontend completion** — Pinia store, updated Link_Input, Dashboard page, AppHeader nav link

Each task builds on the previous ones. Property-based tests use `pgregory.net/rapid` (backend) and Vitest (frontend).

---

## Tasks

- [ ] 1. Fix `GetURLByShortCode` SELECT query bug
  - In `backend/internal/url/repository/url_repo.go`, add `short_code` to the `SELECT` column list in `GetURLByShortCode`
  - Update the corresponding `Scan` call to include `&u.ShortCode` in the correct position
  - _Requirements: 3.1, 3.2_

  - [ ]* 1.1 Write property test for short-code round-trip (Property 3)
    - Add `pgregory.net/rapid` to `go.mod` via `go get pgregory.net/rapid`
    - Create `backend/internal/url/repository/url_repo_test.go`
    - Use `rapid.StringMatching` to generate random 6-char alphanumeric short codes
    - For each generated code, configure a mock repo that returns a `model.URL` with `ShortCode` set to that code; assert the returned `ShortCode` equals the queried value
    - Tag: `// Feature: url-shortener-phase1, Property 3: Short code round-trip`
    - **Property 3: Short code round-trip**
    - **Validates: Requirements 3.2, 3.3**

- [ ] 2. Add URL validation to the service layer
  - In `backend/internal/url/model/url.go`, add the `url` binding tag to `CreateURLRequest.OriginalURL`: `binding:"required,url"`
  - In `backend/internal/url/service/url_service.go`, add a `net/url.Parse` check inside `CreateShortURL`: reject any URL whose scheme is not `http` or `https`, or whose host is empty; return `appErrors.NewBadRequestError` on failure
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

  - [ ]* 2.1 Write unit tests for URL validation
    - Extend `backend/internal/url/service/url_service_test.go` with table-driven tests
    - Cases: empty string, `"not-a-url"`, `"ftp://example.com"`, `"javascript:alert(1)"`, `"file:///etc/passwd"` → expect error, mock `CreateShortURL` never called
    - Cases: `"http://example.com"`, `"https://example.com/path?q=1"` → expect success
    - _Requirements: 1.1, 1.2, 1.3, 1.4_

  - [ ]* 2.2 Write property test for invalid URL rejection (Property 1)
    - In `backend/internal/url/service/url_service_test.go` (or a new `url_service_prop_test.go`)
    - Use `rapid` to generate strings that are NOT valid http/https URLs (non-http/https schemes, random strings, empty)
    - For each, call `URLService.CreateShortURL` with a counting mock repo; assert error returned and mock `CreateShortURL` call count is 0
    - Tag: `// Feature: url-shortener-phase1, Property 1: Invalid URL rejection`
    - **Property 1: Invalid URL rejection**
    - **Validates: Requirements 1.2, 1.3**

  - [ ]* 2.3 Write property test for valid URL acceptance (Property 2)
    - Use `rapid` to generate valid `http`/`https` URLs (random hosts, paths, query strings)
    - For each, call `URLService.CreateShortURL` with the existing `MockRepo`; assert no error, `ShortCode` length is 6, `OriginalURL` in response equals input
    - Tag: `// Feature: url-shortener-phase1, Property 2: Valid URL acceptance`
    - **Property 2: Valid URL acceptance**
    - **Validates: Requirements 1.4**

- [ ] 3. Switch short-code generation to `crypto/rand`
  - In `backend/internal/url/service/url_service.go`, replace the `math/rand` import with `crypto/rand`
  - Rewrite `generateShortCode` to use `crypto/rand.Read` to fill a byte slice, then map each byte to the `charset` using modulo
  - _Requirements: 2.4_

- [ ] 4. Add collision retry logic to `CreateShortURL`
  - In `backend/internal/url/service/url_service.go`, wrap the `repo.CreateShortURL` call in a retry loop (up to 5 attempts)
  - On a unique-constraint violation (detect via `pgconn.PgError` with code `"23505"`), generate a new short code and retry
  - After 5 failed attempts, return `appErrors.NewInternalError("Failed to generate unique short code after retries", err)`
  - _Requirements: 2.1, 2.2, 2.3_

  - [ ]* 4.1 Write unit tests for collision retry
    - Add a `MockRepoWithCollisions` in `url_service_test.go` that returns a unique-constraint error for the first N calls then succeeds
    - Test: N=1 (succeeds on 2nd attempt), N=4 (succeeds on 5th attempt), N=5 (all fail → error returned after exactly 5 attempts)
    - _Requirements: 2.1, 2.2, 2.3_

- [ ] 5. Add `NewTooManyRequestsError` to the errors package
  - In `backend/platform/errors/errors.go`, add:
    ```go
    func NewTooManyRequestsError(message string, err error) *AppError {
        return NewAppError(http.StatusTooManyRequests, message, err)
    }
    ```
  - _Requirements: 7.2, 10.1_

- [ ] 6. Implement rate limiter middleware
  - Create `backend/platform/middleware/rate_limiter.go`
  - Add `golang.org/x/time/rate` to `go.mod` if not already present (it is transitively available; add it as a direct dependency)
  - Define `RateLimiterMiddleware` with a `sync.Map` of `*rate.Limiter` keyed by client IP, limit `rate.Every(3 * time.Second)`, burst 20
  - Implement `getLimiter(ip string) *rate.Limiter` using `sync.Map.LoadOrStore`
  - Implement `Handler() gin.HandlerFunc`: call `limiter.Allow()`; if false, return `NewTooManyRequestsError` with `Retry-After` header; on internal error, log and call `c.Next()`
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [ ]* 6.1 Write unit tests for rate limiter
    - Create `backend/platform/middleware/rate_limiter_test.go`
    - Test: single IP — 20 requests succeed, 21st returns 429 with `Retry-After` header
    - Test: two different IPs — limits are independent (20 each)
    - Test: limiter internal error path — request proceeds (fail-open)
    - _Requirements: 7.1, 7.2, 7.4, 7.5_

  - [ ]* 6.2 Write property test for rate limit enforcement (Property 6)
    - Use `rapid` to generate random IP address strings
    - For each IP, send 21 requests through the middleware via `httptest`; assert first 20 return non-429, 21st returns 429 with `Retry-After` header
    - Tag: `// Feature: url-shortener-phase1, Property 6: Rate limit enforcement`
    - **Property 6: Rate limit enforcement**
    - **Validates: Requirements 7.1**

- [ ] 7. Wire rate limiter into the router and update health check
  - In `backend/cmd/api/main.go`, instantiate `RateLimiterMiddleware` and apply its `Handler()` only to the `POST /api/v1/shorten` route (not to `GET /:shortCode` or `GET /ping`)
  - Update the `GET /ping` handler to call `pool.Ping(ctx)` and return `{"status": "ok"}` on success or `{"status": "ok", "db": "unavailable"}` on failure — both with HTTP 200
  - _Requirements: 7.4, 9.1, 9.2, 9.3_

- [ ] 8. Checkpoint — verify backend builds and all tests pass
  - Run `go build ./...` from `backend/` to confirm no compile errors
  - Run `go test ./...` from `backend/` to confirm all unit and property tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. Add property test for error response structure (Property 7)
  - Create `backend/internal/url/handler/url_handler_prop_test.go`
  - Use `rapid` to generate diverse error-producing inputs: invalid URLs, missing `originalUrl` field, non-existent short codes
  - For each, call the full Gin handler stack via `httptest.NewRecorder`; assert response body is valid JSON with `code` (int) == HTTP status and `message` (non-empty string), and `Content-Type` is `application/json`
  - Tag: `// Feature: url-shortener-phase1, Property 7: Error response structure`
  - **Property 7: Error response structure**
  - **Validates: Requirements 10.1, 10.3, 10.4**
  - _Requirements: 10.1, 10.3, 10.4_

- [ ] 10. Create `UrlRecord` TypeScript type
  - Create `frontend/app/types/url.ts` with the `UrlRecord` interface:
    ```typescript
    export interface UrlRecord {
      id: number
      originalUrl: string
      shortCode: string
      clicks: number
      createdAt: string
    }
    ```
  - _Requirements: 5.1, 6.2_

- [ ] 11. Implement `useUrlStore` Pinia store
  - Create `frontend/app/stores/useUrlStore.ts`
  - Define a `setup` store with `urls: ref<UrlRecord[]>([])`, `addUrl(record: UrlRecord): void`, `clearUrls(): void`, and `latestUrl: computed<UrlRecord | null>`
  - Store is session-scoped (no `localStorage` persistence)
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

  - [ ]* 11.1 Write unit tests for `useUrlStore` (Property 4)
    - Create `frontend/test/unit/useUrlStore.test.ts`
    - Test: initial `urls` is empty, `latestUrl` is null
    - Test: `addUrl` appends a record and `urls.length` increases by exactly 1
    - Test: `latestUrl` returns the most recently added record
    - Test: `clearUrls` empties the list
    - Use `rapid`-style property approach via Vitest: generate N random `UrlRecord` objects, call `addUrl` for each, assert `urls.length === N` and each record is present
    - Tag: `// Feature: url-shortener-phase1, Property 4: URL store append`
    - **Property 4: URL store append**
    - **Validates: Requirements 5.1, 5.2**

- [ ] 12. Update `HeroAppLinkInput` component
  - In `frontend/app/components/Hero/AppLinkInput.vue`:
    - In `onSuccess`: call `urlStore.addUrl(response.data)` (import and use `useUrlStore`)
    - In `onError`: replace `console.error` with `toast.add` — show the API `message` field for 400 responses, `"Too many requests. Please wait before trying again."` for 429, and `"Service unavailable. Please try again later."` for 500 / network errors
    - Clear error state when a new submission begins (before calling `shortenUrl`)
    - Remove all `console.log` / `console.error` watch calls
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 5.2, 5.3_

  - [ ]* 12.1 Write component tests for `HeroAppLinkInput`
    - Create `frontend/test/nuxt/AppLinkInput.test.ts`
    - Use `mountSuspended` from `@nuxt/test-utils/runtime`
    - Mock `$fetch` to return a 400 error; assert a toast with the API error message is shown
    - Mock `$fetch` to return a 500 error; assert a toast with the generic message is shown
    - Mock `$fetch` to return a 429 error; assert the rate-limit toast is shown
    - Mock `$fetch` to return success; assert `urlStore.addUrl` was called with the response data
    - Assert error is cleared when a new submission starts
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 5.2_

- [ ] 13. Create Dashboard page
  - Create `frontend/app/pages/dashboard.vue`
  - Import and use `useUrlStore`; read `urlStore.urls` reactively
  - Render a table/list with columns: Original URL, Short URL (`appBaseUrl + '/' + shortCode`), Clicks, Created At, Copy button
  - Copy button: write full short URL to clipboard via `navigator.clipboard.writeText`, show a success toast
  - Empty state: display a message with a `<NuxtLink to="/">` link when `urlStore.urls` is empty
  - _Requirements: 6.1, 6.2, 6.3, 6.4_

  - [ ]* 13.1 Write component tests for Dashboard (Property 5)
    - Create `frontend/test/nuxt/dashboard.test.ts`
    - Use `mountSuspended`; seed the store with N records (N ≥ 1); assert each record's `originalUrl`, `shortCode`, `clicks`, and `createdAt` appear in the rendered DOM
    - Test empty state: empty store → empty-state message and home link rendered
    - Test copy button: clicking copy calls `navigator.clipboard.writeText` with the correct full short URL and shows a toast
    - Use Vitest's `fc`-style loop to test with multiple record sets
    - Tag: `// Feature: url-shortener-phase1, Property 5: Dashboard record display`
    - **Property 5: Dashboard record display**
    - **Validates: Requirements 6.2**

- [ ] 14. Add Dashboard link to `AppHeader`
  - In `frontend/app/components/AppHeader.vue`, add a `{ label: 'Dashboard', to: '/dashboard', active: route.path.startsWith('/dashboard') }` entry to the `items` computed array
  - _Requirements: 6.5_

- [ ] 15. Checkpoint — verify frontend builds and all tests pass
  - Run `pnpm run typecheck` from `frontend/` to confirm no TypeScript errors
  - Run `pnpm test --run` from `frontend/` to confirm all Vitest tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 16. Final checkpoint — verify all tests pass
  - Confirm `go test ./...` from `backend/` still passes
  - Confirm `pnpm test --run` from `frontend/` still passes
  - Ensure all tests pass, ask the user if questions arise.

---

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Property tests require `pgregory.net/rapid` — add it with `go get pgregory.net/rapid` from `backend/`
- Each property test must be tagged with `// Feature: url-shortener-phase1, Property N: <text>` and run with at least 100 iterations
- Backend property tests run via `go test -count=1 -run TestProperty ./...`
- Frontend tests run via `pnpm test --run` from `frontend/`
- The `url` binding tag on `CreateURLRequest` handles the `required` + well-formed URL check at the handler layer; the service layer adds the `http`/`https` scheme check on top
- The rate limiter is fail-open: internal errors log and allow the request through (Requirement 7.5)
- The health check always returns HTTP 200 regardless of DB state (Requirement 9.3)
