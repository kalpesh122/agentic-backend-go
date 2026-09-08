# Plan: Notes resource

Spec: `spec.md`

## Steps

### 1. Migration + queries + generated store
- Files: `db/migrations/00001_create_notes.sql`, `db/queries/notes.sql`, `internal/store/*` (generated)
- Test: `TestPGStoreRoundTrip` applies migrations on a fresh container

### 2. Types + cursor
- Files: `internal/notes/types.go`
- Test: `TestCursorRoundTrip`

### 3. Store (interface + PGStore)
- Files: `internal/notes/store.go`
- Test: `TestPGStoreRoundTrip`

### 4. Service
- Files: `internal/notes/service.go`
- Test: `TestCreateValidation`, `TestMalformedCursorAndLimit`, `TestListIsNewestFirstAndPaginates`

### 5. Handler + router mount + auth
- Files: `internal/notes/handler.go`, `internal/httpx/router.go`, `internal/auth/auth.go`
- Test: `TestRequiresPrincipal`, `TestGetPatchDeleteAndIsolation`, `TestInvalidAndUnknownIDs`

## Rollback

`goose down` (the `Down` section drops the table); remove the route mount in `router.go`.

## Out of scope / follow-ups

- OpenAPI document generation.
