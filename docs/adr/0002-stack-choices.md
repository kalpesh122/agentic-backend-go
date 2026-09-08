# 0002. Stack choices for the Go API boilerplate

Date: 2026-09-09
Status: accepted

## Context

The boilerplate must be productive for AI coding agents: explicit code over frameworks, SQL that can be read and reviewed, a single command gate, and tests that catch the common mistakes (wrong SQL, unscoped queries, stale generated code).

## Decision

Go 1.27 with chi v5 on standard `http.Handler`; pgx v5 with sqlc-generated typed queries; goose SQL migrations embedded in the binary; `log/slog`; `encoding/json/v2`; golangci-lint v2; Testcontainers for the store integration test; tools pinned through the `go.mod` `tool` directive; domain packages under `internal/` following go.dev/doc/modules/layout.

## Alternatives considered

- **Stdlib `ServeMux` only** — fine below ~15 routes, but no sub-routers, grouped middleware, or method-not-allowed handling without boilerplate. chi adds exactly those and nothing else.
- **Echo / Gin / Fiber** — Echo v5 changed every handler signature and lives at a new module path; Gin and Fiber use non-standard context types. All three are things agents hallucinate old versions of.
- **GORM / ent / sqlx** — an ORM hides the SQL an agent needs to review; sqlc keeps SQL visible and generates the boring code, and `just check` fails if the generated code drifts.
- **golang-migrate** — solid, but goose's `-- +goose Up/Down` files plus `embed.FS` are simpler for one binary.
- **`chi/middleware.RealIP`** — deprecated after IP-spoofing advisories; `middleware.ClientIP` trusts `X-Forwarded-For` only from configured proxies.
- **Taskfile / Makefile** — the kit standardises on `just` so every boilerplate exposes the same recipes.
- **A full OpenAPI generator** — deferred; add `oapi-codegen` or `huma` when a typed client is needed.

## Consequences

- Good: handler tests run in under a second with no Docker; the SQL is proven on real Postgres; generated code cannot drift silently; a ~15 MB distroless image.
- Bad: integration tests need Docker or `TEST_DATABASE_URL`; the dev Postgres uses host port 5435 to avoid clashing with local installs; `go tool golangci-lint` compiles a large binary on first use.
- Neutral: OpenTelemetry and CORS are opt-in; API-key auth is deliberately minimal.
