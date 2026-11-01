package pg

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorClass describes err without its message, for workers that handle personal
// data: PostgreSQL error text can quote the values it failed on.
func ErrorClass(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return "postgres:" + pgErr.Code + ":" + pgErr.ConstraintName
	}
	return "internal"
}
