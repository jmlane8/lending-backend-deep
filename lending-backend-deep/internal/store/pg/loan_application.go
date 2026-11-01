package pg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LoanApplicationStore handles persistence for the loan_application and
// application_payload tables.
type LoanApplicationStore struct {
	pool *pgxpool.Pool
}

// NewLoanApplicationStore returns a new LoanApplicationStore backed by pool.
func NewLoanApplicationStore(pool *pgxpool.Pool) *LoanApplicationStore {
	return &LoanApplicationStore{pool: pool}
}

// CreateApplicationParams carries the inputs required to record a submission.
type CreateApplicationParams struct {
	ID          uuid.UUID
	Reference   string
	PartnerID   uuid.UUID
	SubmittedAt time.Time
	Body        []byte // exact request body
}

// Create records the application and its raw body in one transaction.
func (s *LoanApplicationStore) Create(ctx context.Context, p CreateApplicationParams) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("LoanApplicationStore.Create: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO loan_application (
			loan_application_id, application_reference, partner_id,
			submitted_at, acknowledged_at
		) VALUES ($1, $2, $3, $4, NOW())
	`, p.ID, p.Reference, p.PartnerID, p.SubmittedAt); err != nil {
		return fmt.Errorf("LoanApplicationStore.Create: application: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO application_payload (loan_application_id, body)
		VALUES ($1, $2)
	`, p.ID, p.Body); err != nil {
		return fmt.Errorf("LoanApplicationStore.Create: payload: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("LoanApplicationStore.Create: commit: %w", err)
	}
	return nil
}

// ValidationResult is an application's current validation state.
type ValidationResult struct {
	Reference   string
	Status      string // "PENDING" | "PASSED" | "REJECTED"
	SubmittedAt time.Time
	ValidatedAt *time.Time
	Reasons     json.RawMessage
}

// FindByReference returns the application with the given reference if it belongs
// to partnerID, or ErrNotFound.
func (s *LoanApplicationStore) FindByReference(ctx context.Context, reference string, partnerID uuid.UUID) (ValidationResult, error) {
	var r ValidationResult
	var reasons []byte
	err := s.pool.QueryRow(ctx, `
		SELECT application_reference, application_status, submitted_at,
		       validated_at, rejection_reasons
		FROM loan_application
		WHERE application_reference = $1
		  AND partner_id = $2
		  AND deleted = FALSE
	`, reference, partnerID).Scan(&r.Reference, &r.Status, &r.SubmittedAt, &r.ValidatedAt, &reasons)
	if errors.Is(err, pgx.ErrNoRows) {
		return ValidationResult{}, ErrNotFound
	}
	if err != nil {
		return ValidationResult{}, fmt.Errorf("LoanApplicationStore.FindByReference: %w", err)
	}
	r.Reasons = reasons
	return r, nil
}

// Begin starts the transaction that holds a claimed application's row lock.
func (s *LoanApplicationStore) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("LoanApplicationStore.Begin: %w", err)
	}
	return tx, nil
}

// ClaimedApplication is a pending application locked by the caller's transaction.
type ClaimedApplication struct {
	ID          uuid.UUID
	Reference   string
	PartnerID   uuid.UUID
	SubmittedAt time.Time
	Body        []byte
}

// ClaimNext locks one due application, skipping rows locked by other workers.
// It returns ErrNotFound when none is due.
func (s *LoanApplicationStore) ClaimNext(ctx context.Context, tx pgx.Tx) (ClaimedApplication, error) {
	var a ClaimedApplication
	err := tx.QueryRow(ctx, `
		SELECT la.loan_application_id, la.application_reference, la.partner_id,
		       la.submitted_at, ap.body
		FROM loan_application la
		JOIN application_payload ap USING (loan_application_id)
		WHERE la.application_status = 'PENDING'
		  AND la.next_validation_at <= NOW()
		  AND la.deleted = FALSE
		ORDER BY la.next_validation_at
		LIMIT 1
		FOR UPDATE OF la SKIP LOCKED
	`).Scan(&a.ID, &a.Reference, &a.PartnerID, &a.SubmittedAt, &a.Body)
	if errors.Is(err, pgx.ErrNoRows) {
		return ClaimedApplication{}, ErrNotFound
	}
	if err != nil {
		return ClaimedApplication{}, fmt.Errorf("LoanApplicationStore.ClaimNext: %w", err)
	}
	return a, nil
}

// MarkPassed sets the application to PASSED and reduces its payload to retained.
func (s *LoanApplicationStore) MarkPassed(ctx context.Context, tx pgx.Tx, id, customerID uuid.UUID, retained []byte) error {
	if _, err := tx.Exec(ctx, `
		UPDATE loan_application
		SET application_status = 'PASSED', customer_id = $2,
		    validated_at = NOW(), updated_at = NOW()
		WHERE loan_application_id = $1
	`, id, customerID); err != nil {
		return fmt.Errorf("LoanApplicationStore.MarkPassed: application: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE application_payload
		SET body = $2, updated_at = NOW()
		WHERE loan_application_id = $1
	`, id, retained); err != nil {
		return fmt.Errorf("LoanApplicationStore.MarkPassed: payload: %w", err)
	}
	return nil
}

// MarkRejected hard-deletes the payload and sets the application to REJECTED,
// leaving customer_id NULL.  This is the only hard delete in the system: FR-011
// requires the personal data to be gone, and a soft delete would keep it.
func (s *LoanApplicationStore) MarkRejected(ctx context.Context, tx pgx.Tx, id uuid.UUID, reasons []byte) error {
	if _, err := tx.Exec(ctx, `
		DELETE FROM application_payload WHERE loan_application_id = $1
	`, id); err != nil {
		return fmt.Errorf("LoanApplicationStore.MarkRejected: payload: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE loan_application
		SET application_status = 'REJECTED', rejection_reasons = $2::text::jsonb,
		    rejection_stage = 'CUSTOMER_DATA', validated_at = NOW(),
		    rejected_at = NOW(), updated_at = NOW()
		WHERE loan_application_id = $1
	`, id, string(reasons)); err != nil {
		return fmt.Errorf("LoanApplicationStore.MarkRejected: application: %w", err)
	}
	return nil
}

// RecordFailure delays the next attempt after an internal failure, by 2^attempts
// minutes up to 60.  It runs outside the rolled-back transaction.
func (s *LoanApplicationStore) RecordFailure(ctx context.Context, id uuid.UUID) error {
	if _, err := s.pool.Exec(ctx, `
		UPDATE loan_application
		SET next_validation_at = NOW() + make_interval(
		        mins => LEAST(power(2, LEAST(validation_attempts, 6)), 60)::int),
		    validation_attempts = validation_attempts + 1,
		    updated_at = NOW()
		WHERE loan_application_id = $1
	`, id); err != nil {
		return fmt.Errorf("LoanApplicationStore.RecordFailure: %w", err)
	}
	return nil
}
