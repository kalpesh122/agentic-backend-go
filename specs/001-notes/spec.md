# Spec: Notes resource (reference vertical slice)

Status: implemented
Date: 2026-09-09
Owner: Kalpesh Mali

## Problem

Every new resource needs a worked example to copy: SQL, generated store, service rules, handler, auth, pagination, error handling, and both test levels. Without one, each agent invents its own shape.

## Goals

- A complete, principal-scoped CRUD resource with keyset pagination.
- Handler tests on a fake store (no Docker) plus a store test on real Postgres.

## Non-goals

- Sharing notes between principals, full-text search, attachments, OpenAPI generation.

## Acceptance criteria (testable sentences)

1. `POST /api/notes` with a valid body returns 201 and the created note; the title is trimmed.
2. `POST /api/notes` with an empty title, or an unknown JSON field, returns 422 with `details`.
3. Any `/api/notes` request without a valid bearer key returns 401 with `error.code = "unauthorized"`.
4. `GET /api/notes?limit=2` returns the newest two notes and a `next_cursor`; passing it returns the rest and `next_cursor: null`.
5. A malformed cursor, a non-integer `limit`, or `limit > 100` returns 422.
6. `GET /api/notes/{id}` for another principal's note returns 404.
7. `PATCH /api/notes/{id}` with `{}` returns 422; with a field returns 200 and the updated note.
8. `DELETE /api/notes/{id}` returns 204 and a subsequent GET returns 404.
9. A non-UUID id returns 422; an unknown UUID returns 404.
10. The store's keyset query returns the correct second page on real Postgres.

## Design

### Data / schema changes

`notes(id uuid pk default gen_random_uuid(), owner_id text, title text check 1..200, body text default '' check ≤10000, created_at timestamptz, updated_at timestamptz)`; index `(owner_id, created_at DESC, id DESC)`.

### API / interface changes

`internal/notes/handler.go`; bearer API keys via `API_KEYS`; standard error body.

### Behaviour

Keyset pagination with a row-value comparison `(created_at, id) < ($cursor_created_at, $cursor_id)`; cursor is raw-URL base64 of `RFC3339Nano|uuid`.

## Risks and mitigations

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Cross-principal data leak | low | every sqlc query filters by `owner_id`; criteria 6 and 10 |
| Stale generated store | medium | `just generate-check` in the gate |

## Open questions

None.
