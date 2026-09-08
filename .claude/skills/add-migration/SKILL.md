---
name: add-migration
description: Change the database schema safely with goose + sqlc (create migration → write SQL → apply → regenerate store → test). Use for any table or column change.
---

# Add a migration

1. `just migration <snake_case_name>` creates `db/migrations/NNNNN_<name>.sql` with `-- +goose Up` / `-- +goose Down` sections.
2. Write additive SQL by default (new nullable column or new table). `NOT NULL` on an existing table needs a `DEFAULT` or a backfill step. Destructive changes need a spec and a two-step rollout (add → backfill → drop). `Down` must be a true inverse.
3. `just docker-up && just migrate` locally. Check `SELECT * FROM goose_db_version`.
4. Update `db/queries/*.sql` for the new shape, then `just generate`. sqlc validates the SQL against the migrations at generation time, so type errors show up here.
5. Update store/service/handler and tests. `just check` (the integration test applies all migrations on a fresh container, proving the migration).
6. Mention the migration in the PR. Production applies it on boot (`MIGRATE_ON_START=true` in the image) or via `bin/migrate`.
