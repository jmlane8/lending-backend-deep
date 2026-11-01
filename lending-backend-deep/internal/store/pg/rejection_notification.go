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

// RejectionNotificationStore handles persistence for the rejection_notification table.
type RejectionNotificationStore struct {
	pool *pgxpool.Pool
}

// NewRejectionNotificationStore returns a new RejectionNotificationStore backed by pool.
func NewRejectionNotificationStore(pool *pgxpool.Pool) *RejectionNotificationStore {
	return &RejectionNotificationStore{pool: pool}
}

// Enqueue adds a notification in the caller's transaction.  It adds no row, and
// returns false, when the partner has no notification_url.
func (s *RejectionNotificationStore) Enqueue(ctx context.Context, tx pgx.Tx, applicationID, partnerID uuid.UUID) (bool, error) {
	eventID, err := uuid.NewV7()
	if err != nil {
		return false, fmt.Errorf("RejectionNotificationStore.Enqueue: id: %w", err)
	}
	tag, err := tx.Exec(ctx, `
		INSERT INTO rejection_notification (event_id, loan_application_id)
		SELECT $1, $2
		FROM partner
		WHERE partner_id = $3 AND notification_url IS NOT NULL
	`, eventID, applicationID, partnerID)
	if err != nil {
		return false, fmt.Errorf("RejectionNotificationStore.Enqueue: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// Begin starts the transaction that holds a claimed notification's row lock.
func (s *RejectionNotificationStore) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("RejectionNotificationStore.Begin: %w", err)
	}
	return tx, nil
}

// ClaimedNotification is a pending notification locked by the caller's transaction.
type ClaimedNotification struct {
	EventID            uuid.UUID
	Attempts           int
	Reference          string
	RejectedAt         time.Time
	Reasons            json.RawMessage
	NotificationURL    string
	NotificationSecret string
}

// ClaimNext locks one due notification, skipping rows locked by other workers.
// It returns ErrNotFound when none is due.
func (s *RejectionNotificationStore) ClaimNext(ctx context.Context, tx pgx.Tx) (ClaimedNotification, error) {
	var n ClaimedNotification
	var reasons []byte
	err := tx.QueryRow(ctx, `
		SELECT rn.event_id, rn.attempts, la.application_reference, la.rejected_at,
		       la.rejection_reasons, COALESCE(p.notification_url, ''),
		       COALESCE(p.notification_secret, '')
		FROM rejection_notification rn
		JOIN loan_application la USING (loan_application_id)
		JOIN partner p ON p.partner_id = la.partner_id
		WHERE rn.status = 'PENDING'
		  AND rn.next_attempt_at <= NOW()
		  AND rn.deleted = FALSE
		ORDER BY rn.next_attempt_at
		LIMIT 1
		FOR UPDATE OF rn SKIP LOCKED
	`).Scan(&n.EventID, &n.Attempts, &n.Reference, &n.RejectedAt, &reasons,
		&n.NotificationURL, &n.NotificationSecret)
	if errors.Is(err, pgx.ErrNoRows) {
		return ClaimedNotification{}, ErrNotFound
	}
	if err != nil {
		return ClaimedNotification{}, fmt.Errorf("RejectionNotificationStore.ClaimNext: %w", err)
	}
	n.Reasons = reasons
	return n, nil
}

// MarkDelivered records a successful delivery.
func (s *RejectionNotificationStore) MarkDelivered(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) error {
	if _, err := tx.Exec(ctx, `
		UPDATE rejection_notification
		SET status = 'DELIVERED', attempts = attempts + 1,
		    delivered_at = NOW(), updated_at = NOW()
		WHERE event_id = $1
	`, eventID); err != nil {
		return fmt.Errorf("RejectionNotificationStore.MarkDelivered: %w", err)
	}
	return nil
}

// MarkExhausted records the final failed attempt.
func (s *RejectionNotificationStore) MarkExhausted(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) error {
	if _, err := tx.Exec(ctx, `
		UPDATE rejection_notification
		SET status = 'EXHAUSTED', attempts = attempts + 1, updated_at = NOW()
		WHERE event_id = $1
	`, eventID); err != nil {
		return fmt.Errorf("RejectionNotificationStore.MarkExhausted: %w", err)
	}
	return nil
}

// ScheduleRetry records a failed attempt and sets when the next one may run.
func (s *RejectionNotificationStore) ScheduleRetry(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, delay time.Duration) error {
	if _, err := tx.Exec(ctx, `
		UPDATE rejection_notification
		SET attempts = attempts + 1,
		    next_attempt_at = NOW() + make_interval(secs => $2),
		    updated_at = NOW()
		WHERE event_id = $1
	`, eventID, delay.Seconds()); err != nil {
		return fmt.Errorf("RejectionNotificationStore.ScheduleRetry: %w", err)
	}
	return nil
}
