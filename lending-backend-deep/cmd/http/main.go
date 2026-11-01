package main

import (
	"context"
	"time"
	_ "time/tzdata" // keeps BUSINESS_TIMEZONE working in minimal containers

	"lending-backend-deep/internal/config"
	dbmigrate "lending-backend-deep/internal/migrate"
	"lending-backend-deep/internal/rest/server"
	pgstore "lending-backend-deep/internal/store/pg"
)

func main() {
	// -----------------------------------------------------------------
	// Configuration — all values read from environment variables.
	// -----------------------------------------------------------------
	cfg := config.LoadFromEnv()
	logger := cfg.Logger()

	if _, err := time.LoadLocation(cfg.BusinessTimezone); err != nil {
		logger.Fatal().Err(err).Msg("invalid BUSINESS_TIMEZONE")
		return // unreachable: Fatal calls os.Exit(1)
	}

	// -----------------------------------------------------------------
	// Database migrations — applied before the connection pool or the
	// HTTP server are initialised so the schema is always current when
	// the first request arrives.  ErrNoChange (already up to date) is
	// treated as a success by Run, so this is safe on every startup.
	// -----------------------------------------------------------------
	if err := dbmigrate.Run(context.Background(), cfg.EnvPostgreSQL, logger); err != nil {
		logger.Fatal().Err(err).Msg("failed to run database migrations")
		return // unreachable: Fatal calls os.Exit(1)
	}

	// -----------------------------------------------------------------
	// PostgreSQL connection pool.
	// -----------------------------------------------------------------
	pool, err := pgstore.NewPool(context.Background(), pgstore.Config{
		Host:        cfg.PostgresHost,
		Port:        cfg.PostgresPort,
		Database:    cfg.PostgresDatabase,
		Schema:      cfg.PostgresSchema,
		Username:    cfg.PostgresUsername,
		Password:    cfg.PostgresPassword,
		SSLMode:     cfg.PostgresSSLMode,
		ConnMaxOpen: cfg.PostgresConnMaxOpen,
		ConnMinIdle: cfg.PostgresConnMinIdle,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to postgresql")
		return // unreachable: Fatal calls os.Exit(1)
	}
	defer pool.Close()

	logger.Info().
		Str("host", cfg.PostgresHost).
		Int("port", cfg.PostgresPort).
		Str("db", cfg.PostgresDatabase).
		Msg("connected to postgresql")

	// -----------------------------------------------------------------
	// HTTP server — blocks until SIGINT / SIGTERM.
	// -----------------------------------------------------------------
	if err := server.Serve(cfg, logger, pool); err != nil {
		logger.Fatal().Err(err).Msg("http server error")
	}
}
