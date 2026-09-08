package notes

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// fakeStore is an in-memory Store used by handler and service tests.
type fakeStore struct {
	mu    sync.Mutex
	rows  []Note
	owner map[uuid.UUID]string
	clock time.Time
}

func newFakeStore() *fakeStore {
	return &fakeStore{owner: map[uuid.UUID]string{}, clock: time.Date(2026, 9, 9, 10, 0, 0, 0, time.UTC)}
}

func (f *fakeStore) List(_ context.Context, ownerID string, limit int, cursor *Cursor) ([]Note, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []Note
	for _, n := range f.rows {
		if f.owner[n.ID] != ownerID {
			continue
		}
		if cursor != nil {
			before := n.CreatedAt.Before(cursor.CreatedAt) || (n.CreatedAt.Equal(cursor.CreatedAt) && n.ID.String() < cursor.ID.String())
			if !before {
				continue
			}
		}
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.After(out[j].CreatedAt)
		}
		return out[i].ID.String() > out[j].ID.String()
	})
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (f *fakeStore) Get(_ context.Context, ownerID string, id uuid.UUID) (Note, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, n := range f.rows {
		if n.ID == id && f.owner[id] == ownerID {
			return n, nil
		}
	}
	return Note{}, ErrNotFound
}

func (f *fakeStore) Create(_ context.Context, ownerID, title, body string) (Note, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.clock = f.clock.Add(time.Second)
	n := Note{ID: uuid.New(), Title: title, Body: body, CreatedAt: f.clock, UpdatedAt: f.clock}
	f.rows = append(f.rows, n)
	f.owner[n.ID] = ownerID
	return n, nil
}

func (f *fakeStore) Update(_ context.Context, ownerID string, id uuid.UUID, title, body *string) (Note, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, n := range f.rows {
		if n.ID == id && f.owner[id] == ownerID {
			if title != nil {
				n.Title = *title
			}
			if body != nil {
				n.Body = *body
			}
			n.UpdatedAt = f.clock.Add(time.Minute)
			f.rows[i] = n
			return n, nil
		}
	}
	return Note{}, ErrNotFound
}

func (f *fakeStore) Delete(_ context.Context, ownerID string, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	idx := slices.IndexFunc(f.rows, func(n Note) bool { return n.ID == id && f.owner[id] == ownerID })
	if idx < 0 {
		return ErrNotFound
	}
	f.rows = slices.Delete(f.rows, idx, idx+1)
	delete(f.owner, id)
	return nil
}
