---
paths:
  - "internal/**/*.go"
  - "cmd/**/*.go"
  - "db/**"
---

# Go rules

- Package layout is by domain (`internal/notes`), not by layer. Inside a package: `types.go`, `handler.go`, `service.go`, `store.go`.
- Errors: wrap with `fmt.Errorf("doing x: %w", err)`; sentinel errors for store-level conditions (`ErrNotFound`); `*apperr.Error` for anything the client sees.
- Context first, no globals, no `init()` with side effects, no panics for expected failures.
- Handlers: decode with `respond.Decode`, validate in the service, render with `respond.JSON/Error`. Explicit status codes (`201`, `204`).
- Tests: table tests where there are >2 cases; `t.Context()`, `t.Setenv`, `t.Cleanup`; `require` for preconditions, `assert` for the claim. Integration tests skip under `-short`.
- SQL lives in `db/queries`; `internal/store` is generated. `just generate-check` fails on drift.
- Migrations: goose SQL with `Up` and `Down`; additive by default.
