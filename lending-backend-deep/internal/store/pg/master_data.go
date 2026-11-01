package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MasterDataStore handles reads of the master_data_value table.
type MasterDataStore struct {
	pool *pgxpool.Pool
}

// NewMasterDataStore returns a new MasterDataStore backed by pool.
func NewMasterDataStore(pool *pgxpool.Pool) *MasterDataStore {
	return &MasterDataStore{pool: pool}
}

// ListActive returns the active values grouped by category.
func (s *MasterDataStore) ListActive(ctx context.Context) (map[string][]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT category, value
		FROM master_data_value
		WHERE is_active = TRUE AND deleted = FALSE
		ORDER BY category, value
	`)
	if err != nil {
		return nil, fmt.Errorf("MasterDataStore.ListActive: %w", err)
	}
	defer rows.Close()

	values := map[string][]string{}
	for rows.Next() {
		var category, value string
		if err := rows.Scan(&category, &value); err != nil {
			return nil, fmt.Errorf("MasterDataStore.ListActive: scan: %w", err)
		}
		values[category] = append(values[category], value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("MasterDataStore.ListActive: rows: %w", err)
	}
	return values, nil
}

// ActiveGenders returns the accepted gender values, inside the caller's transaction.
func (s *MasterDataStore) ActiveGenders(ctx context.Context, tx pgx.Tx) ([]string, error) {
	rows, err := tx.Query(ctx, `
		SELECT value
		FROM master_data_value
		WHERE category = 'GENDER' AND is_active = TRUE AND deleted = FALSE
	`)
	if err != nil {
		return nil, fmt.Errorf("MasterDataStore.ActiveGenders: %w", err)
	}
	genders, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("MasterDataStore.ActiveGenders: rows: %w", err)
	}
	return genders, nil
}
