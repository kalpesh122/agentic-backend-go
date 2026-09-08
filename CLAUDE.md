@AGENTS.md

## Claude Code specifics

- Use the `add-endpoint` and `add-migration` skills; they encode the package shape and the sqlc/goose loop.
- Hooks enforce the hard rules in AGENTS.md (dangerous commands, `.env`, `go.sum`, gofmt-on-edit, `just check` on stop). Do not work around them.
- Use plan mode for anything touching `db/`, `internal/httpx/router.go`, or `internal/config`.
- Handler tests never need Docker; only `store_integration_test.go` does (or `TEST_DATABASE_URL`).
