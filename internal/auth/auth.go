// Package auth resolves the caller from a bearer API key.
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/kalpesh122/agentic-backend-go/internal/apperr"
	"github.com/kalpesh122/agentic-backend-go/internal/httpx/respond"
)

type ctxKey struct{}

// Principal returns the authenticated principal stored by RequireAPIKey, or "" if none.
func Principal(ctx context.Context) string {
	p, _ := ctx.Value(ctxKey{}).(string)
	return p
}

// WithPrincipal is used by tests and internal callers to inject a principal.
func WithPrincipal(ctx context.Context, principal string) context.Context {
	return context.WithValue(ctx, ctxKey{}, principal)
}

// RequireAPIKey rejects requests without a valid `Authorization: Bearer <key>` header.
// keys maps API key → principal; the principal scopes the caller's data.
func RequireAPIKey(keys map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(raw, "Bearer ")
			if !ok || token == "" {
				respond.Error(w, r, apperr.Unauthorized(""))
				return
			}
			principal, found := keys[token]
			if !found {
				respond.Error(w, r, apperr.Unauthorized("Invalid API key"))
				return
			}
			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), principal)))
		})
	}
}
