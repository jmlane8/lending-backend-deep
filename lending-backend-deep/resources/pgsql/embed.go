// Package pgsql provides the embedded SQL migration files for the Lending
// Backend database schema.
//
// Migration files live under the migrations/ sub-directory and follow the
// golang-migrate naming convention:
//
//	<version>_<description>.up.sql   ← applied in ascending version order
//	<version>_<description>.down.sql ← applied in descending version order (rollback)
//
// The embedded filesystem is consumed at runtime by internal/migrate to apply
// pending migrations before the HTTP server starts accepting traffic.
package pgsql

import "embed"

// Migrations is an embedded filesystem containing all SQL migration files.
//
//go:embed migrations
var Migrations embed.FS
