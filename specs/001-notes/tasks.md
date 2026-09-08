# Tasks: Notes resource

- [x] Migration, queries, generated store — test: `TestPGStoreRoundTrip` (Testcontainers, ~10 s)
- [x] Types and cursor — test: `TestCursorRoundTrip`
- [x] Store — test: `TestPGStoreRoundTrip` (pagination, isolation, update, delete)
- [x] Service validation + pagination — test: `TestCreateValidation`, `TestMalformedCursorAndLimit`, `TestListIsNewestFirstAndPaginates`
- [x] Handler, auth, error shape — test: `TestRequiresPrincipal`, `TestGetPatchDeleteAndIsolation`, `TestInvalidAndUnknownIDs`
- [x] Docs / README / ADR updated
- [x] `just check` green — lint 0 issues, no sqlc drift, all packages ok
