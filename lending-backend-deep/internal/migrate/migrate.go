// Package migrate provides database schema migration support for the Lending
// Backend.  It is intentionally kept separate from the internal/rest/* packages
// so that the migration concern does not bleed into HTTP-layer code.
//
// Migrations are embedded at compile time from resources/pgsql/migrations and
// applied in ascending version order using golang-migrate.  The function is
// designed to be called once at application startup, before the HTTP server
// begins accepting traffic.
package migrate

import (
	"context"
	"errors"
	"fmt"

	"lending-backend-deep/internal/config"
	"lending-backend-deep/resources/pgsql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // register pgx5 driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog"
)

// zerologAdapter bridges zerolog to the migrate.Logger interface so that
// golang-migrate's internal log output is routed through the application logger
// instead of going directly to stderr.
type zerologAdapter struct {
	log zerolog.Logger
}

func (a zerologAdapter) Printf(format string, v ...any) {
	a.log.Debug().Msgf(format, v...)
}

func (a zerologAdapter) Verbose() bool {
	return a.log.GetLevel() <= zerolog.DebugLevel
}

// Run applies all pending "up" migrations embedded in resources/pgsql/migrations.
//
// Behaviour:
//   - If the schema is already at the latest version (ErrNoChange), Run returns
//     nil — it is safe to call on every startup.
//   - Any other error is wrapped and returned so the caller can decide whether
//     to abort startup (the recommended behaviour).
//
// The migration state table (schema_migrations) is created automatically by
// golang-migrate on first run.
func Run(ctx context.Context, cfg config.EnvPostgreSQL, logger zerolog.Logger) error {
	// Build the iofs source from the embedded migration files.
	src, err := iofs.New(pgsql.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("migrate: create iofs source: %w", err)
	}

	// pgx5:// scheme is required by the golang-migrate pgx/v5 database driver.
	dsn := fmt.Sprintf(
		"pgx5://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
		cfg.PostgresSSLMode,
		cfg.PostgresSchema,
	)

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		return fmt.Errorf("migrate: init: %w", err)
	}

	m.Log = zerologAdapter{log: logger}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Warn().Err(srcErr).Msg("migrate: failed to close source")
		}
		if dbErr != nil {
			logger.Warn().Err(dbErr).Msg("migrate: failed to close db connection")
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info().Msg("database schema is up to date, no migrations to apply")
			return nil
		}
		return fmt.Errorf("migrate: up: %w", err)
	}

	logger.Info().Msg("database migrations applied successfully")
	return nil
}
