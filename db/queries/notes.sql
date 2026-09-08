-- name: ListNotes :many
SELECT * FROM notes
WHERE owner_id = @owner_id
  AND (
    sqlc.narg('cursor_created_at')::timestamptz IS NULL
    OR (created_at, id) < (sqlc.narg('cursor_created_at')::timestamptz, sqlc.narg('cursor_id')::uuid)
  )
ORDER BY created_at DESC, id DESC
LIMIT @page_limit;

-- name: GetNote :one
SELECT * FROM notes WHERE owner_id = $1 AND id = $2;

-- name: CreateNote :one
INSERT INTO notes (owner_id, title, body) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateNote :one
UPDATE notes
SET title = coalesce(sqlc.narg('title'), title),
    body = coalesce(sqlc.narg('body'), body),
    updated_at = now()
WHERE owner_id = @owner_id AND id = @id
RETURNING *;

-- name: DeleteNote :execrows
DELETE FROM notes WHERE owner_id = $1 AND id = $2;
