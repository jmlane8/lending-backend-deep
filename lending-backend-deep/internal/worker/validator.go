// Package worker holds the background loops that run inside the service.  Each
// loop polls PostgreSQL with FOR UPDATE SKIP LOCKED, so several replicas can run
// safely and a crash rolls back and the row is picked up again.
package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pgstore "lending-backend-deep/internal/store/pg"
	"lending-backend-deep/internal/validation"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// Validator runs the customer data rules on recorded submissions.  Log lines must
// never contain submitted values (FR-011).
type Validator struct {
	logger        zerolog.Logger
	applications  *pgstore.LoanApplicationStore
	customers     *pgstore.CustomerStore
	masterData    *pgstore.MasterDataStore
	notifications *pgstore.RejectionNotificationStore
	location      *time.Location
}

// NewValidator returns a Validator that reads submission dates in location.
func NewValidator(
	logger zerolog.Logger,
	applications *pgstore.LoanApplicationStore,
	customers *pgstore.CustomerStore,
	masterData *pgstore.MasterDataStore,
	notifications *pgstore.RejectionNotificationStore,
	location *time.Location,
) *Validator {
	return &Validator{
		logger:        logger,
		applications:  applications,
		customers:     customers,
		masterData:    masterData,
		notifications: notifications,
		location:      location,
	}
}

// Run processes due applications until ctx is cancelled, waiting interval
// whenever none is due.
func (v *Validator) Run(ctx context.Context, interval time.Duration) {
	poll(ctx, interval, v.logger, "validator", v.ProcessNext)
}

// ProcessNext claims and validates one due application in one transaction.  It
// returns false when none was due.  An internal failure rolls back and delays
// the application, so it counts as processed.
func (v *Validator) ProcessNext(ctx context.Context) (bool, error) {
	tx, err := v.applications.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(context.Background())

	application, err := v.applications.ClaimNext(ctx, tx)
	if errors.Is(err, pgstore.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	log := v.logger.With().
		Str("application_reference", application.Reference).
		Str("partner_id", application.PartnerID.String()).
		Logger()

	reasons, err := v.process(ctx, tx, application, log)
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		v.recordFailure(ctx, tx, application, log, err)
		return true, nil
	}

	delay := time.Since(application.SubmittedAt)
	if len(reasons) > 0 {
		log.Info().Str("outcome", "REJECTED").Strs("reason_codes", reasonCodes(reasons)).
			Dur("validation_delay", delay).Msg("application validated")
		return true, nil
	}
	log.Info().Str("outcome", "PASSED").Dur("validation_delay", delay).Msg("application validated")
	return true, nil
}

func (v *Validator) process(ctx context.Context, tx pgx.Tx, application pgstore.ClaimedApplication, log zerolog.Logger) ([]validation.Reason, error) {
	genders, err := v.masterData.ActiveGenders(ctx, tx)
	if err != nil {
		return nil, err
	}

	submissionDate := application.SubmittedAt.In(v.location)
	customer, reasons, err := validation.Validate(application.Body, submissionDate, genders)
	if err != nil {
		return nil, err
	}

	if len(reasons) > 0 {
		return reasons, v.reject(ctx, tx, application, reasons, log)
	}
	return nil, v.pass(ctx, tx, application, customer)
}

func (v *Validator) pass(ctx context.Context, tx pgx.Tx, application pgstore.ClaimedApplication, customer validation.Customer) error {
	photos, err := v.customers.InsertPhotos(ctx, tx, customer.GovernmentIDPhoto, customer.SelfiePhoto)
	if err != nil {
		return err
	}
	customerID, err := v.customers.UpsertCustomer(ctx, tx, customer, photos[0], photos[1])
	if err != nil {
		return err
	}
	retained, err := validation.RetainedSections(application.Body)
	if err != nil {
		return err
	}
	return v.applications.MarkPassed(ctx, tx, application.ID, customerID, retained)
}

func (v *Validator) reject(ctx context.Context, tx pgx.Tx, application pgstore.ClaimedApplication, reasons []validation.Reason, log zerolog.Logger) error {
	encoded, err := json.Marshal(reasons)
	if err != nil {
		return fmt.Errorf("Validator.reject: marshal reasons: %w", err)
	}
	if err := v.applications.MarkRejected(ctx, tx, application.ID, encoded); err != nil {
		return err
	}
	queued, err := v.notifications.Enqueue(ctx, tx, application.ID, application.PartnerID)
	if err != nil {
		return err
	}
	if !queued {
		log.Warn().Msg("partner has no notification address; rejection is available by lookup only")
	}
	return nil
}

// recordFailure rolls back and delays the application.  On shutdown it does
// neither delay nor log an error: the rollback leaves the row due for the next start.
func (v *Validator) recordFailure(ctx context.Context, tx pgx.Tx, application pgstore.ClaimedApplication, log zerolog.Logger, cause error) {
	_ = tx.Rollback(context.Background())
	if ctx.Err() != nil {
		return
	}
	log.Error().Str("error_class", pgstore.ErrorClass(cause)).Msg("application validation failed")
	if err := v.applications.RecordFailure(ctx, application.ID); err != nil {
		log.Error().Str("error_class", pgstore.ErrorClass(err)).Msg("failed to delay application")
	}
}

func reasonCodes(reasons []validation.Reason) []string {
	codes := make([]string, 0, len(reasons))
	for _, r := range reasons {
		codes = append(codes, r.Code)
	}
	return codes
}

// poll calls process until ctx is cancelled, waiting interval whenever there was
// nothing to do or the call failed.
func poll(ctx context.Context, interval time.Duration, logger zerolog.Logger, name string, process func(context.Context) (bool, error)) {
	for ctx.Err() == nil {
		processed, err := process(ctx)
		if err != nil && ctx.Err() == nil {
			logger.Error().Str("worker", name).Str("error_class", pgstore.ErrorClass(err)).Msg("poll failed")
		}
		if processed && err == nil {
			continue
		}
		select {
		case <-ctx.Done():
		case <-time.After(interval):
		}
	}
}
