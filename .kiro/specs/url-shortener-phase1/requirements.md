# Requirements Document

## Introduction

SwiftLink is a URL shortening service consisting of a Go/Gin backend API and a Nuxt 4 / Vue 3 frontend. Phase 1 documents the current working state of the application and defines the improvements needed to make it production-ready. The scope covers backend URL validation, short-code collision handling, a bug fix for the redirect query, frontend error handling and state management, a URL dashboard, and rate limiting.

## Glossary

- **API**: The Go/Gin backend HTTP service.
- **Frontend**: The Nuxt 4 / Vue 3 web application.
- **Short_Code**: A 6-character alphanumeric string that uniquely identifies a shortened URL.
- **Original_URL**: The full destination URL submitted by a user for shortening.
- **URL_Record**: The database row containing `id`, `original_url`, `short_code`, `clicks`, and `created_at`.
- **URL_Service**: The Go service layer responsible for business logic around URL creation and retrieval.
- **URL_Repository**: The Go data-access layer that reads from and writes to PostgreSQL.
- **Redirect_Handler**: The Go handler that resolves a Short_Code to an Original_URL and issues an HTTP 302 redirect.
- **Shorten_Handler**: The Go handler that accepts an Original_URL and returns a URL_Record.
- **Link_Input**: The `HeroAppLinkInput` Vue component that accepts a URL from the user and calls the shorten API.
- **URL_Store**: A Pinia store that manages the list of shortened URLs in the Frontend.
- **Dashboard**: The `/dashboard` page that lists all URL_Records created in the current browser session.
- **Rate_Limiter**: The middleware that restricts the number of requests a client can make within a time window.
- **Validator**: The component (backend or frontend) responsible for checking that a value conforms to expected rules.

---

## Requirements

### Requirement 1: URL Validation on the Backend

**User Story:** As a developer, I want the API to validate that `originalUrl` is a well-formed HTTP or HTTPS URL before storing it, so that the database does not contain invalid or non-navigable entries.

#### Acceptance Criteria

1. WHEN a `POST /api/v1/shorten` request is received with an `originalUrl` that is empty or missing, THE Shorten_Handler SHALL return HTTP 400 with a structured error response indicating the field is required.
2. WHEN a `POST /api/v1/shorten` request is received with an `originalUrl` that is not a valid URL (e.g., `"not-a-url"`, `"ftp://example.com"`, `"javascript:alert(1)"`), THE Shorten_Handler SHALL return HTTP 400 with a structured error response indicating the URL is invalid.
3. WHEN a `POST /api/v1/shorten` request is received with an `originalUrl` whose scheme is neither `http` nor `https`, THE Shorten_Handler SHALL return HTTP 400 with a structured error response.
4. WHEN a `POST /api/v1/shorten` request is received with a valid `http` or `https` `originalUrl`, THE Shorten_Handler SHALL proceed to create and return the URL_Record.
5. THE Validator SHALL enforce URL validation at the service layer so that it applies regardless of the HTTP transport used.

---

### Requirement 2: Short Code Collision Handling

**User Story:** As a developer, I want the system to handle short-code collisions gracefully, so that two different Original_URLs are never assigned the same Short_Code.

#### Acceptance Criteria

1. WHEN the URL_Service generates a Short_Code that already exists in the database, THE URL_Service SHALL regenerate a new Short_Code and retry the insert.
2. THE URL_Service SHALL retry Short_Code generation up to 5 times before returning an error.
3. IF the URL_Service exhausts all retries without finding a unique Short_Code, THEN THE Shorten_Handler SHALL return HTTP 500 with a structured error response.
4. THE URL_Service SHALL use a cryptographically random source when generating Short_Codes to minimise collision probability.

---

### Requirement 3: Redirect Query Bug Fix

**User Story:** As a developer, I want the `GetURLByShortCode` query to return the `short_code` column, so that the URL_Record returned from the redirect lookup is complete.

#### Acceptance Criteria

1. WHEN the URL_Repository executes `GetURLByShortCode`, THE URL_Repository SHALL include `short_code` in the `SELECT` column list.
2. WHEN the Redirect_Handler resolves a Short_Code, THE URL_Record returned by the URL_Repository SHALL have a non-empty `ShortCode` field equal to the queried Short_Code.
3. FOR ALL valid Short_Codes stored in the database, querying by Short_Code and reading back the `ShortCode` field SHALL return the same value that was stored (round-trip property).

---

### Requirement 4: Frontend Error Handling

**User Story:** As a user, I want to see clear error messages in the UI when URL shortening fails, so that I understand what went wrong and can correct my input.

#### Acceptance Criteria

1. WHEN the Link_Input component receives an error response from the API, THE Link_Input SHALL display the error message to the user in a visible UI element (e.g., a toast notification or inline error).
2. WHEN the API returns HTTP 400 due to an invalid URL, THE Link_Input SHALL display a message indicating the URL is invalid.
3. WHEN the API returns HTTP 500 or a network error occurs, THE Link_Input SHALL display a generic error message indicating the service is unavailable.
4. WHEN an error is displayed, THE Link_Input SHALL NOT log the error only to the browser console as the sole user-facing feedback.
5. WHEN the user corrects the URL and resubmits, THE Link_Input SHALL clear any previously displayed error before making a new request.

---

### Requirement 5: Frontend State Management with Pinia

**User Story:** As a developer, I want URL shortening state managed in a Pinia store, so that shortened URL data is accessible across components and pages without prop-drilling.

#### Acceptance Criteria

1. THE URL_Store SHALL maintain a list of URL_Records created during the current browser session.
2. WHEN the Link_Input successfully shortens a URL, THE URL_Store SHALL append the new URL_Record to its list.
3. WHEN the URL_Store list is updated, THE Link_Input SHALL reactively reflect the most recently created URL_Record.
4. THE URL_Store SHALL expose the list of URL_Records to the Dashboard page.
5. WHILE the URL_Store list is empty, THE Dashboard SHALL display an empty-state message prompting the user to shorten a URL.

---

### Requirement 6: URL Dashboard Page

**User Story:** As a user, I want a dashboard page that lists all URLs I have shortened in the current session, so that I can review and copy my short links without re-shortening them.

#### Acceptance Criteria

1. THE Frontend SHALL provide a `/dashboard` route that renders the Dashboard page.
2. WHEN the Dashboard page loads, THE Dashboard SHALL retrieve the list of URL_Records from the URL_Store and display each record's `originalUrl`, `shortCode`, `clicks`, and `createdAt`.
3. WHEN the user clicks a copy button next to a URL_Record, THE Dashboard SHALL copy the full short URL (base URL + Short_Code) to the clipboard and display a confirmation toast.
4. WHILE the URL_Store list is empty, THE Dashboard SHALL display an empty-state message with a link to the home page.
5. THE Dashboard SHALL be navigable from the AppHeader navigation menu.

---

### Requirement 7: Rate Limiting

**User Story:** As a system operator, I want the API to enforce rate limits on the shorten endpoint, so that the service is protected from abuse and excessive load.

#### Acceptance Criteria

1. THE Rate_Limiter SHALL limit each client IP address to 20 `POST /api/v1/shorten` requests per minute.
2. WHEN a client exceeds the rate limit, THE Rate_Limiter SHALL return HTTP 429 with a structured error response and a `Retry-After` header indicating when the client may retry.
3. WHILE a client is within the rate limit, THE Rate_Limiter SHALL allow requests to proceed to the Shorten_Handler without additional latency beyond 5ms.
4. THE Rate_Limiter SHALL NOT apply to the `GET /:shortCode` redirect endpoint or the `GET /ping` health check endpoint.
5. IF the Rate_Limiter encounters an internal error evaluating a request, THEN THE API SHALL allow the request to proceed and log the error.

---

### Requirement 9: Health Check Endpoint

**User Story:** As a system operator, I want a reliable health check endpoint, so that container orchestration and monitoring tools can verify the API is running.

#### Acceptance Criteria

1. THE API SHALL expose a `GET /ping` endpoint that returns HTTP 200 with the body `"pong"`.
2. WHEN the database connection pool is healthy, THE API SHALL respond to `GET /ping` within 200ms.
3. IF the database connection pool is unavailable, THEN THE API SHALL still respond to `GET /ping` with HTTP 200 to indicate the process is alive, and SHALL include a `"db": "unavailable"` field in the response body.

---

### Requirement 10: Structured API Error Responses

**User Story:** As a frontend developer, I want all API errors to follow a consistent JSON structure, so that the Frontend can reliably parse and display error messages.

#### Acceptance Criteria

1. THE API SHALL return all error responses as JSON objects with at minimum a `code` (integer HTTP status) and `message` (string) field.
2. WHEN a validation error occurs, THE API SHALL include field-level detail in the error response so the Frontend can identify which input was invalid.
3. THE API SHALL return `Content-Type: application/json` on all error responses.
4. FOR ALL error-producing inputs to the API, the response body SHALL be parseable as a JSON object with a `code` field equal to the HTTP status code (round-trip property).
