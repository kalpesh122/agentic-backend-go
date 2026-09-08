// Command api runs the HTTP API.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kalpesh122/agentic-backend-go/db/migrations"
	"github.com/kalpesh122/agentic-backend-go/internal/config"
	"github.com/kalpesh122/agentic-backend-go/internal/httpx"
	"github.com/kalpesh122/agentic-backend-go/internal/notes"
	"github.com/kalpesh122/agentic-backend-go/internal/observability"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := observability.NewLogger(cfg.Env, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdownTracing, err := observability.SetupTracing(ctx, cfg.OTLPEndpoint, "agentic-backend-go")
	if err != nil {
		return fmt.Errorf("tracing: %w", err)
	}
	defer func() {
		flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracing(flushCtx)
	}()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("db pool: %w", err)
	}
	defer pool.Close()
	if cfg.MigrateOnStart {
		if err := migrations.Up(ctx, pool); err != nil {
			return err
		}
		log.Info("migrations applied")
	}

	handler := httpx.NewRouter(httpx.Deps{Config: cfg, Log: log, DB: pool, Notes: notes.NewPGStore(pool)})
	return httpx.Serve(ctx, log, cfg.Port, handler, cfg.ShutdownTimeout)
}
