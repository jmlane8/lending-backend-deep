package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds the parameters needed to open a pgx connection pool.
type Config struct {
	Host        string
	Port        int
	Database    string
	Schema      string
	Username    string
	Password    string
	SSLMode     string
	ConnMaxOpen int
	ConnMinIdle int
}

// NewPool creates and validates a *pgxpool.Pool from the given Config.
// It pings the database before returning so startup fails fast on bad credentials.
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
		cfg.Schema,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pg.NewPool: parse config: %w", err)
	}

	// MaxConns caps the total number of open connections (active + idle).
	poolCfg.MaxConns = int32(cfg.ConnMaxOpen)
	// MinIdleConns ensures the pool keeps at least N idle connections ready,
	// reducing tail-latency on bursts.  Per pgx/v5 docs, MinIdleConns is
	// superior to MinConns for this purpose because MinConns eagerly creates
	// connections regardless of load, whereas MinIdleConns only refills the
	// idle tier after connections are used or closed.
	poolCfg.MinIdleConns = int32(cfg.ConnMinIdle)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pg.NewPool: create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pg.NewPool: ping: %w", err)
	}

	return pool, nil
}
