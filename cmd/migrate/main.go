// Command migrate applies pending database migrations and exits.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kalpesh122/agentic-backend-go/db/migrations"
	"github.com/kalpesh122/agentic-backend-go/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: db pool:", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := migrations.Up(ctx, pool); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
	fmt.Println("migrations applied")
}
