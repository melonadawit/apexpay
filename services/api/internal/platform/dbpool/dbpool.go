package dbpool

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New builds a tuned pgx connection pool. Sensible production defaults:
//   - prepared-statement cache (reduces planning round-trips per query)
//   - explicit min/max connections to avoid connection churn and pool exhaustion
//   - sane acquire/idle timeouts so a slow upstream never hangs a request
//
// These can be overridden via env for the deployment environment.
func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 20
	cfg.MinConns = 4
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetimeJitter = time.Minute
	cfg.HealthCheckPeriod = time.Minute
	cfg.ConnConfig.ConnectTimeout = 10 * time.Second
	return pgxpool.NewWithConfig(ctx, cfg)
}
