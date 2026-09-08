// Package middleware holds HTTP middleware that is not domain specific.
package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// ClientIP resolves the caller's IP. X-Forwarded-For is trusted only when the direct peer is a trusted proxy;
// this avoids the spoofing issues that led chi to deprecate middleware.RealIP.
func ClientIP(trustedProxies []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			peer, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				peer = r.RemoteAddr
			}
			if slices.Contains(trustedProxies, peer) {
				if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
					first, _, _ := strings.Cut(xff, ",")
					if ip := strings.TrimSpace(first); ip != "" {
						r.RemoteAddr = net.JoinHostPort(ip, "0")
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AccessLog writes one structured line per request with the request id, status and latency.
func AccessLog(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.InfoContext(r.Context(), "request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"ms", time.Since(start).Milliseconds(),
				"request_id", chimw.GetReqID(r.Context()),
			)
		})
	}
}

// SecureHeaders sets conservative defaults for an API.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
