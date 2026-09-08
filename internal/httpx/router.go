// Package httpx assembles the HTTP server: middleware order, routes, health.
package httpx

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	"github.com/kalpesh122/agentic-backend-go/internal/apperr"
	"github.com/kalpesh122/agentic-backend-go/internal/auth"
	"github.com/kalpesh122/agentic-backend-go/internal/config"
	"github.com/kalpesh122/agentic-backend-go/internal/httpx/middleware"
	"github.com/kalpesh122/agentic-backend-go/internal/httpx/respond"
	"github.com/kalpesh122/agentic-backend-go/internal/notes"
)

func keyByClientIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr, nil //nolint:nilerr // fall back to the raw address rather than failing the request
	}
	return host, nil
}

// Pinger is the readiness dependency (a pgx pool satisfies it).
type Pinger interface {
	Ping(ctx context.Context) error
}

// Deps are the injected collaborators; tests pass fakes.
type Deps struct {
	Config config.Config
	Log    *slog.Logger
	DB     Pinger
	Notes  notes.Store
}

// NewRouter builds the full route tree. Middleware order matters; keep it explicit.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(middleware.ClientIP(d.Config.TrustedProxies))
	r.Use(middleware.AccessLog(d.Log))
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecureHeaders)
	r.Use(chimw.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Get("/ready", func(w http.ResponseWriter, req *http.Request) {
		if err := d.DB.Ping(req.Context()); err != nil {
			d.Log.ErrorContext(req.Context(), "readiness check failed", "err", err)
			respond.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready", "db": "unreachable"})
			return
		}
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ready", "db": "ok"})
	})

	r.Route("/api", func(api chi.Router) {
		if len(d.Config.CORSOrigins) > 0 {
			api.Use(cors.Handler(cors.Options{
				AllowedOrigins:   d.Config.CORSOrigins,
				AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Content-Type"},
				AllowCredentials: true,
			}))
		}
		// Keyed by RemoteAddr, which middleware.ClientIP has already resolved through trusted proxies.
		api.Use(httprate.LimitBy(d.Config.RateLimitPerMin, time.Minute, keyByClientIP,
			httprate.WithLimitHandler(func(w http.ResponseWriter, req *http.Request) {
				respond.Error(w, req, apperr.New(apperr.CodeRateLimited, "Too many requests", http.StatusTooManyRequests))
			}),
		))
		api.Use(auth.RequireAPIKey(d.Config.APIKeyMap()))
		api.Route("/notes", notes.NewHandler(notes.NewService(d.Notes)).Routes)
	})

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		respond.Error(w, req, apperr.NotFound("Route"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		respond.Error(w, req, apperr.New(apperr.CodeNotFound, "Method not allowed", http.StatusMethodNotAllowed))
	})
	return r
}
