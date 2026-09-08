# agentic-backend-go command surface. `just check` is the quality gate.
set dotenv-load := true

export GOTOOLCHAIN := "auto"

default:
    @just --list

# Download modules and tools (Go 1.27 toolchain is fetched automatically)
setup:
    go mod download
    go tool sqlc version >/dev/null

# Run the API with live reload (needs `just docker-up` + `just migrate`)
dev:
    go tool air

test *ARGS:
    go test ./... -race -count=1 {{ARGS}}

# Unit tests only (no Docker)
test-short:
    go test ./... -short -race -count=1

lint:
    go vet ./...
    go tool golangci-lint run

fmt:
    go tool golangci-lint fmt

fmt-file FILE:
    gofmt -w "{{FILE}}" >/dev/null 2>&1 || true

typecheck:
    go build ./...

build:
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/migrate ./cmd/migrate

# Regenerate the sqlc store from db/queries + db/migrations
generate:
    go tool sqlc generate

# Fail if generated code is stale
generate-check: generate
    git diff --exit-code -- internal/store

# lint + typecheck + generated-code drift + tests (unit + integration)
check: lint typecheck generate-check test
    @echo "check: green"

council BASE="origin/main":
    scripts/council-review.sh {{BASE}}

# Apply embedded migrations to DATABASE_URL
migrate:
    go run ./cmd/migrate

# Create a new goose migration: just migration add_widgets
migration NAME:
    go tool goose -dir db/migrations create {{NAME}} sql

# Start Postgres (host port 5435) in the background
docker-up:
    docker compose up -d db

docker-down:
    docker compose down

docker-build:
    docker build -t agentic-backend-go:local .
