package notes

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kalpesh122/agentic-backend-go/internal/store"
)

// ErrNotFound is returned by Store implementations when no row matches.
var ErrNotFound = errors.New("note not found")

// Store is the persistence port. Every method is scoped to an owner.
type Store interface {
	List(ctx context.Context, ownerID string, limit int, cursor *Cursor) ([]Note, error)
	Get(ctx context.Context, ownerID string, id uuid.UUID) (Note, error)
	Create(ctx context.Context, ownerID, title, body string) (Note, error)
	Update(ctx context.Context, ownerID string, id uuid.UUID, title, body *string) (Note, error)
	Delete(ctx context.Context, ownerID string, id uuid.UUID) error
}

// PGStore implements Store on Postgres through the sqlc-generated queries.
type PGStore struct{ q *store.Queries }

// NewPGStore wraps a pgx pool.
func NewPGStore(pool *pgxpool.Pool) *PGStore { return &PGStore{q: store.New(pool)} }

func fromRow(r store.Note) Note {
	return Note{ID: r.ID, Title: r.Title, Body: r.Body, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}

// List returns up to limit notes newest first, starting after cursor when given.
func (s *PGStore) List(ctx context.Context, ownerID string, limit int, cursor *Cursor) ([]Note, error) {
	params := store.ListNotesParams{OwnerID: ownerID, PageLimit: int32(limit)} //nolint:gosec // limit is validated ≤ 100
	if cursor != nil {
		params.CursorCreatedAt = pgtype.Timestamptz{Time: cursor.CreatedAt, Valid: true}
		params.CursorID = pgtype.UUID{Bytes: cursor.ID, Valid: true}
	}
	rows, err := s.q.ListNotes(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	out := make([]Note, 0, len(rows))
	for _, r := range rows {
		out = append(out, fromRow(r))
	}
	return out, nil
}

// Get returns one note or ErrNotFound.
func (s *PGStore) Get(ctx context.Context, ownerID string, id uuid.UUID) (Note, error) {
	r, err := s.q.GetNote(ctx, store.GetNoteParams{OwnerID: ownerID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	if err != nil {
		return Note{}, fmt.Errorf("get note: %w", err)
	}
	return fromRow(r), nil
}

// Create inserts a note.
func (s *PGStore) Create(ctx context.Context, ownerID, title, body string) (Note, error) {
	r, err := s.q.CreateNote(ctx, store.CreateNoteParams{OwnerID: ownerID, Title: title, Body: body})
	if err != nil {
		return Note{}, fmt.Errorf("create note: %w", err)
	}
	return fromRow(r), nil
}

// Update patches the non-nil fields.
func (s *PGStore) Update(ctx context.Context, ownerID string, id uuid.UUID, title, body *string) (Note, error) {
	r, err := s.q.UpdateNote(ctx, store.UpdateNoteParams{OwnerID: ownerID, ID: id, Title: title, Body: body})
	if errors.Is(err, pgx.ErrNoRows) {
		return Note{}, ErrNotFound
	}
	if err != nil {
		return Note{}, fmt.Errorf("update note: %w", err)
	}
	return fromRow(r), nil
}

// Delete removes a note or returns ErrNotFound.
func (s *PGStore) Delete(ctx context.Context, ownerID string, id uuid.UUID) error {
	n, err := s.q.DeleteNote(ctx, store.DeleteNoteParams{OwnerID: ownerID, ID: id})
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
