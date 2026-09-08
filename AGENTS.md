# AGENTS.md

Single source of truth for every AI coding agent in this repository. `CLAUDE.md` imports it; `GEMINI.md`
and `.github/copilot-instructions.md` link to it. Keep it under 150 lines: a map, not an encyclopedia.

## What this repository is

A production-grade REST API boilerplate: **Go 1.27 · chi v5 · pgx v5 · sqlc · goose · slog · OpenTelemetry ·
golangci-lint v2**. Standard-library HTTP handlers, SQL you can read (sqlc generates the typed store), embedded
migrations, bearer API-key auth mapped to principals, cursor-paginated CRUD, one error body everywhere, and
tests at two levels: handler tests on an in-memory store (no Docker) and a store integration test on a real
Postgres (Testcontainers).

The `notes` package is the reference vertical slice. Copy its shape for every new resource.

## Command surface

| Command | What it does |
|---------|--------------|
| `just setup` | `go mod download` (Go 1.27.1 toolchain auto-downloads via `GOTOOLCHAIN=auto`) |
| `just docker-up` | Postgres 17 on host port 5435 (`just docker-down` stops it) |
| `just migrate` | Apply embedded goose migrations to `DATABASE_URL` |
| `just dev` | `air` live reload on `PORT` (8080) |
| `just test` / `just test-short` | All tests (integration needs Docker or `TEST_DATABASE_URL`) / unit only |
| `just lint` / `just fmt` | `go vet` + golangci-lint v2 / gofmt + goimports |
| `just generate` | Regenerate `internal/store` from `db/queries/*.sql` |
| `just check` | **The gate**: lint + build + generated-code drift + tests |
| `just migration <name>` | Create a new goose SQL migration |
| `just council` | Local multi-model code review vs main |

## Layout

```
cmd/api/main.go           bootstrap: config → logger → tracing → pool → (migrate) → router → serve → graceful stop
cmd/migrate/main.go       apply migrations and exit
internal/config/          Config from env (caarlos0/env), validated, APIKeyMap()
internal/apperr/          Error{Code, Message, Status, Details}; NotFound/Unauthorized/Validation helpers
internal/auth/            RequireAPIKey middleware; Principal(ctx)
internal/httpx/           router.go (middleware order, /health, /ready, /api mount), server.go (timeouts, drain)
internal/httpx/middleware ClientIP (trusted proxies), AccessLog (slog), SecureHeaders
internal/httpx/respond    JSON(), Error() → standard body, Decode() (strict JSON, size-limited)
internal/notes/           types.go → handler.go → service.go → store.go (Store interface + PGStore); *_test.go
internal/store/           sqlc OUTPUT. Never edit; run `just generate`
internal/observability/   slog setup, optional OTLP tracing
db/migrations/            goose SQL + migrations.go (embed + Up)      db/queries/  sqlc SQL
```

## Workflow (non-negotiable)

1. Non-trivial change → `brainstorm-spec` skill first (`specs/NNN-slug/`).
2. New resource → `add-endpoint` skill. Schema change → `add-migration` skill.
3. TDD: handler tests against a fake `Store` in `internal/<pkg>/handler_test.go`; store tests against Postgres.
4. `just check` green before "done" (`verify-before-done` skill). Paste the tail.
5. Conventional commits; `just council` before pushing anything non-trivial.

## Hard rules

- Handlers only translate HTTP ↔ service types. Rules live in the service; SQL lives in `db/queries`.
- Every store method takes `ownerID`; every query filters by it. Never expose another principal's rows.
- Return `*apperr.Error` for expected failures and let `respond.Error` render it. Wrap causes with `%w`.
- `context.Context` is the first parameter of anything that does IO. No package-level mutable state.
- Never edit `internal/store`; edit SQL and run `just generate`. `just check` fails on drift.
- Migrations are additive by default; destructive changes need a spec and a two-step rollout.
- Use `encoding/json/v2` with `RejectUnknownMembers` for request bodies (see `respond.Decode`).
- Do not use `chi/middleware.RealIP` (deprecated, spoofable); `middleware.ClientIP` handles trusted proxies.
- Secrets only via env; `.env` is gitignored and hook-protected. Never log API keys.
- No new dependency without a one-line justification in the PR.

## Skills

| Skill | When |
|-------|------|
| `add-endpoint` | New resource or route: SQL → generate → store → service → handler → route → tests |
| `add-migration` | Any schema change: goose file → migrate → regenerate |
| `brainstorm-spec`, `tdd`, `debug`, `code-review`, `council-review`, `verify-before-done`, `adr`, `git-hygiene` | Kit workflow skills |

## Subagents

`planner`, `tdd-implementer`, `reviewer`, `security-reviewer`, `explorer`, `council-synthesizer` in `.claude/agents/`.

## Definition of done

Spec acceptance criteria met · handler tests cover happy + failure paths · `just check` green (output pasted) ·
no placeholders · README/ADR updated if behaviour or stack changed · conventional commit.
