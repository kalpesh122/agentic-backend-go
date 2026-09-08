// Package notes is the reference vertical slice: handler → service → store.
package notes

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Note is the API representation of a note.
type Note struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Page is a cursor-paginated list.
type Page struct {
	Items      []Note  `json:"items"`
	NextCursor *string `json:"next_cursor"`
}

// CreateInput is the body of POST /api/notes.
type CreateInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// UpdateInput is the body of PATCH /api/notes/{id}; nil fields are left unchanged.
type UpdateInput struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

// Cursor encodes the keyset position (created_at, id) of the last item in a page.
type Cursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

// Encode produces an opaque, URL-safe token.
func (c Cursor) Encode() string {
	return base64.RawURLEncoding.EncodeToString([]byte(c.CreatedAt.UTC().Format(time.RFC3339Nano) + "|" + c.ID.String()))
}

// DecodeCursor parses a token produced by Encode. It returns an error for anything malformed.
func DecodeCursor(raw string) (Cursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Cursor{}, fmt.Errorf("decode cursor: %w", err)
	}
	ts, id, ok := strings.Cut(string(b), "|")
	if !ok {
		return Cursor{}, fmt.Errorf("decode cursor: missing separator")
	}
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return Cursor{}, fmt.Errorf("decode cursor: %w", err)
	}
	u, err := uuid.Parse(id)
	if err != nil {
		return Cursor{}, fmt.Errorf("decode cursor: %w", err)
	}
	return Cursor{CreatedAt: t, ID: u}, nil
}
