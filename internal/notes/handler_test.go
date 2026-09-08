package notes

import (
	"encoding/json/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kalpesh122/agentic-backend-go/internal/auth"
	"github.com/kalpesh122/agentic-backend-go/internal/httpx/respond"
)

// newTestServer mounts the handler behind a fixed principal, mirroring the production auth middleware.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	st := newFakeStore()
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			principal := req.Header.Get("X-Test-Principal")
			if principal == "" {
				respond.Error(w, req, auth.ErrTestUnauthorized())
				return
			}
			next.ServeHTTP(w, req.WithContext(auth.WithPrincipal(req.Context(), principal)))
		})
	})
	r.Route("/api/notes", NewHandler(NewService(st)).Routes)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func do(t *testing.T, srv *httptest.Server, method, path, principal, body string) (*http.Response, map[string]any) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), method, srv.URL+path, strings.NewReader(body))
	require.NoError(t, err)
	if principal != "" {
		req.Header.Set("X-Test-Principal", principal)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := srv.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	var parsed map[string]any
	if len(raw) > 0 {
		require.NoError(t, json.Unmarshal(raw, &parsed), string(raw))
	}
	return res, parsed
}

func TestCreateTrimsTitleAndReturns201(t *testing.T) {
	srv := newTestServer(t)
	res, body := do(t, srv, http.MethodPost, "/api/notes", "alice", `{"title":"  First  ","body":"hello"}`)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "First", body["title"])
	assert.Equal(t, "hello", body["body"])
	_, err := uuid.Parse(body["id"].(string))
	require.NoError(t, err)
}

func TestCreateValidation(t *testing.T) {
	srv := newTestServer(t)
	res, body := do(t, srv, http.MethodPost, "/api/notes", "alice", `{"title":""}`)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	errObj := body["error"].(map[string]any)
	assert.Equal(t, "validation_error", errObj["code"])
	assert.Contains(t, errObj["details"].(map[string]any), "title")
	assert.NotEmpty(t, errObj["request_id"])

	res, body = do(t, srv, http.MethodPost, "/api/notes", "alice", `{"title":"x","unknown":1}`)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	assert.Equal(t, "validation_error", body["error"].(map[string]any)["code"])
}

func TestRequiresPrincipal(t *testing.T) {
	srv := newTestServer(t)
	res, body := do(t, srv, http.MethodGet, "/api/notes", "", "")
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.Equal(t, "unauthorized", body["error"].(map[string]any)["code"])
}

func TestListIsNewestFirstAndPaginates(t *testing.T) {
	srv := newTestServer(t)
	do(t, srv, http.MethodPost, "/api/notes", "bob", `{"title":"not mine"}`)
	for _, title := range []string{"a", "b", "c"} {
		do(t, srv, http.MethodPost, "/api/notes", "alice", `{"title":"`+title+`"}`)
	}
	res, page1 := do(t, srv, http.MethodGet, "/api/notes?limit=2", "alice", "")
	require.Equal(t, http.StatusOK, res.StatusCode)
	items := page1["items"].([]any)
	assert.Equal(t, "c", items[0].(map[string]any)["title"])
	assert.Equal(t, "b", items[1].(map[string]any)["title"])
	next, _ := page1["next_cursor"].(string)
	require.NotEmpty(t, next)

	_, page2 := do(t, srv, http.MethodGet, "/api/notes?limit=2&cursor="+next, "alice", "")
	items2 := page2["items"].([]any)
	require.Len(t, items2, 1)
	assert.Equal(t, "a", items2[0].(map[string]any)["title"])
	assert.Nil(t, page2["next_cursor"])
}

func TestMalformedCursorAndLimit(t *testing.T) {
	srv := newTestServer(t)
	res, _ := do(t, srv, http.MethodGet, "/api/notes?cursor=not-a-cursor", "alice", "")
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	res, _ = do(t, srv, http.MethodGet, "/api/notes?limit=500", "alice", "")
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	res, _ = do(t, srv, http.MethodGet, "/api/notes?limit=abc", "alice", "")
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
}

func TestGetPatchDeleteAndIsolation(t *testing.T) {
	srv := newTestServer(t)
	_, created := do(t, srv, http.MethodPost, "/api/notes", "alice", `{"title":"mine","body":"b"}`)
	url := "/api/notes/" + created["id"].(string)

	res, _ := do(t, srv, http.MethodGet, url, "alice", "")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res, _ = do(t, srv, http.MethodGet, url, "bob", "")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	res, patched := do(t, srv, http.MethodPatch, url, "alice", `{"body":"changed"}`)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "changed", patched["body"])

	res, _ = do(t, srv, http.MethodPatch, url, "alice", `{}`)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

	res, _ = do(t, srv, http.MethodDelete, url, "alice", "")
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	res, _ = do(t, srv, http.MethodGet, url, "alice", "")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestInvalidAndUnknownIDs(t *testing.T) {
	srv := newTestServer(t)
	res, _ := do(t, srv, http.MethodGet, "/api/notes/not-a-uuid", "alice", "")
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	res, _ = do(t, srv, http.MethodGet, "/api/notes/"+uuid.NewString(), "alice", "")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCursorRoundTrip(t *testing.T) {
	c := Cursor{CreatedAt: newFakeStore().clock, ID: uuid.New()}
	got, err := DecodeCursor(c.Encode())
	require.NoError(t, err)
	assert.True(t, got.CreatedAt.Equal(c.CreatedAt))
	assert.Equal(t, c.ID, got.ID)
	_, err = DecodeCursor("zzz")
	require.Error(t, err)
}
