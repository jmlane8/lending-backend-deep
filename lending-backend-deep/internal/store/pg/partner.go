package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a requested row does not exist.
var ErrNotFound = errors.New("not found")

// PartnerStore handles reads of the partner and partner_credential tables.
type PartnerStore struct {
	pool *pgxpool.Pool
}

// NewPartnerStore returns a new PartnerStore backed by pool.
func NewPartnerStore(pool *pgxpool.Pool) *PartnerStore {
	return &PartnerStore{pool: pool}
}

// PartnerCredential identifies the partner behind a valid API key.
type PartnerCredential struct {
	PartnerID       uuid.UUID
	AgreementStatus string // "ACTIVE" | "INACTIVE"
}

// FindByKeyHash returns the partner owning a valid credential with the given
// SHA-256 key hash, or ErrNotFound.
func (s *PartnerStore) FindByKeyHash(ctx context.Context, hash []byte) (PartnerCredential, error) {
	var c PartnerCredential
	err := s.pool.QueryRow(ctx, `
		SELECT p.partner_id, p.agreement_status
		FROM partner_credential c
		JOIN partner p ON p.partner_id = c.partner_id
		WHERE c.key_hash = $1 AND c.is_valid = TRUE
		  AND c.deleted = FALSE AND p.deleted = FALSE
	`, hash).Scan(&c.PartnerID, &c.AgreementStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return PartnerCredential{}, ErrNotFound
	}
	if err != nil {
		return PartnerCredential{}, fmt.Errorf("PartnerStore.FindByKeyHash: %w", err)
	}
	return c, nil
}
