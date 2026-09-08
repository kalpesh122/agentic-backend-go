package notes

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/kalpesh122/agentic-backend-go/db/migrations"
)

// databaseURL returns TEST_DATABASE_URL when set (CI service container), else starts a Postgres container.
func databaseURL(t *testing.T) string {
	t.Helper()
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	ctx := context.Background()
	pg, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase("app"), tcpostgres.WithUsername("app"), tcpostgres.WithPassword("app"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	url, err := pg.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	return url
}

func TestPGStoreRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	ctx := t.Context()
	pool, err := pgxpool.New(ctx, databaseURL(t))
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	require.NoError(t, migrations.Up(ctx, pool))

	st := NewPGStore(pool)
	owner := "alice-" + uuid.NewString()

	_, err = st.Create(ctx, "bob", "not mine", "")
	require.NoError(t, err)
	var ids []uuid.UUID
	for _, title := range []string{"a", "b", "c"} {
		n, err := st.Create(ctx, owner, title, "body")
		require.NoError(t, err)
		ids = append(ids, n.ID)
	}

	page, err := st.List(ctx, owner, 2, nil)
	require.NoError(t, err)
	require.Len(t, page, 2)
	assert.Equal(t, "c", page[0].Title)
	assert.Equal(t, "b", page[1].Title)

	cursor := &Cursor{CreatedAt: page[1].CreatedAt, ID: page[1].ID}
	rest, err := st.List(ctx, owner, 2, cursor)
	require.NoError(t, err)
	require.Len(t, rest, 1)
	assert.Equal(t, "a", rest[0].Title)

	_, err = st.Get(ctx, "bob", ids[0])
	assert.ErrorIs(t, err, ErrNotFound)

	title := "renamed"
	updated, err := st.Update(ctx, owner, ids[0], &title, nil)
	require.NoError(t, err)
	assert.Equal(t, "renamed", updated.Title)
	assert.Equal(t, "body", updated.Body)

	require.NoError(t, st.Delete(ctx, owner, ids[0]))
	assert.ErrorIs(t, st.Delete(ctx, owner, ids[0]), ErrNotFound)
}
