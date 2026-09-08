package notes

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/kalpesh122/agentic-backend-go/internal/apperr"
)

const (
	maxTitle = 200
	maxBody  = 10_000
	maxLimit = 100
)

// Service owns the business rules: validation beyond shape, pagination, error mapping.
type Service struct{ store Store }

// NewService wires a Store.
func NewService(store Store) *Service { return &Service{store: store} }

func validateTitle(t string, problems map[string]string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		problems["title"] = "must not be empty"
	} else if utf8.RuneCountInString(t) > maxTitle {
		problems["title"] = "must be at most 200 characters"
	}
	return t
}

func validateBody(b string, problems map[string]string) {
	if utf8.RuneCountInString(b) > maxBody {
		problems["body"] = "must be at most 10000 characters"
	}
}

func mapErr(err error) error {
	if errors.Is(err, ErrNotFound) {
		return apperr.NotFound("Note")
	}
	return err
}

// List returns a page of the owner's notes.
func (s *Service) List(ctx context.Context, ownerID string, limit int, rawCursor string) (Page, error) {
	if limit < 1 || limit > maxLimit {
		return Page{}, apperr.Validation(map[string]string{"limit": "must be between 1 and 100"})
	}
	var cursor *Cursor
	if rawCursor != "" {
		c, err := DecodeCursor(rawCursor)
		if err != nil {
			return Page{}, apperr.Wrap(apperr.Validation(map[string]string{"cursor": "malformed cursor"}), err)
		}
		cursor = &c
	}
	rows, err := s.store.List(ctx, ownerID, limit+1, cursor)
	if err != nil {
		return Page{}, err
	}
	page := Page{Items: rows}
	if len(rows) > limit {
		page.Items = rows[:limit]
		last := page.Items[len(page.Items)-1]
		next := Cursor{CreatedAt: last.CreatedAt, ID: last.ID}.Encode()
		page.NextCursor = &next
	}
	return page, nil
}

// Get returns one note.
func (s *Service) Get(ctx context.Context, ownerID string, id uuid.UUID) (Note, error) {
	n, err := s.store.Get(ctx, ownerID, id)
	return n, mapErr(err)
}

// Create validates and inserts.
func (s *Service) Create(ctx context.Context, ownerID string, in CreateInput) (Note, error) {
	problems := map[string]string{}
	title := validateTitle(in.Title, problems)
	validateBody(in.Body, problems)
	if len(problems) > 0 {
		return Note{}, apperr.Validation(problems)
	}
	return s.store.Create(ctx, ownerID, title, in.Body)
}

// Update validates and patches the provided fields.
func (s *Service) Update(ctx context.Context, ownerID string, id uuid.UUID, in UpdateInput) (Note, error) {
	if in.Title == nil && in.Body == nil {
		return Note{}, apperr.Validation(map[string]string{"body": "provide at least one field"})
	}
	problems := map[string]string{}
	if in.Title != nil {
		t := validateTitle(*in.Title, problems)
		in.Title = &t
	}
	if in.Body != nil {
		validateBody(*in.Body, problems)
	}
	if len(problems) > 0 {
		return Note{}, apperr.Validation(problems)
	}
	n, err := s.store.Update(ctx, ownerID, id, in.Title, in.Body)
	return n, mapErr(err)
}

// Delete removes a note.
func (s *Service) Delete(ctx context.Context, ownerID string, id uuid.UUID) error {
	return mapErr(s.store.Delete(ctx, ownerID, id))
}
