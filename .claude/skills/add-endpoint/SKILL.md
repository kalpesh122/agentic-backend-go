---
name: add-endpoint
description: Add a new REST resource or route to this Go chi + sqlc API (SQL → generate → store → service → handler → route → tests). Use for any new endpoint.
argument-hint: [resource-name]
---

# Add an endpoint

Reference implementation: `internal/notes/`. Copy its shape exactly.

1. **Spec** — new resource → `specs/NNN-<resource>/spec.md` first (`brainstorm-spec`).
2. **Table** — `just migration create_<resource>` → write `-- +goose Up/Down` SQL with an `owner_id text` column and an index on `(owner_id, created_at DESC, id DESC)`. Run the `add-migration` skill.
3. **Queries** — `db/queries/<resource>.sql`: `List` (keyset: `(created_at, id) < (@cursor_created_at, @cursor_id)` with `sqlc.narg`), `Get`, `Create`, `Update` (`coalesce(sqlc.narg(...), col)`), `Delete :execrows`. Every query has `owner_id` in its `WHERE`.
4. **Generate** — `just generate`; read the new `internal/store/<resource>.sql.go` params/rows.
5. **Types** — `internal/<resource>/types.go`: API structs with JSON tags, `CreateInput`, `UpdateInput` (pointer fields), reuse `notes.Cursor` or copy it.
6. **Store** — `store.go`: `Store` interface + `PGStore` mapping sqlc rows to API types; map `pgx.ErrNoRows` → `ErrNotFound`.
7. **Service** — `service.go`: validation (trim, lengths), pagination (`limit+1` → `NextCursor`), `apperr` mapping.
8. **Handler** — `handler.go`: `Routes(chi.Router)`, `respond.Decode` for bodies, `parseID`, `respond.JSON/Error`; read the caller with `auth.Principal(ctx)`.
9. **Mount** — `internal/httpx/router.go`: add the store to `Deps`, then `api.Route("/<resource>", ...)` inside the authenticated `/api` group.
10. **Tests** — `handler_test.go` with a fake store: 201, 422 (+unknown field), 401, list + cursor, cross-principal 404, patch, delete 204, bad/unknown id. `store_integration_test.go` on Postgres for the SQL.
11. **Verify** — `just check` (includes generated-code drift). Curl the route via `just dev`.
