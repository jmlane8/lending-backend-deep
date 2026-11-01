package worker

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	pgstore "lending-backend-deep/internal/store/pg"

	"github.com/rs/zerolog"
)

const (
	eventTypeRejected = "loan-application.rejected"
	deliveryTimeout   = 5 * time.Second
	maxAttempts       = 7
)

// retryDelays holds the wait after each failed attempt, by attempt number.
var retryDelays = []time.Duration{
	time.Minute,
	5 * time.Minute,
	30 * time.Minute,
	2 * time.Hour,
	6 * time.Hour,
	12 * time.Hour,
}

// rejectionEvent is the webhook body.  See contracts/rejection-webhook.md.
type rejectionEvent struct {
	EventID              string          `json:"event-id"`
	EventType            string          `json:"event-type"`
	ApplicationReference string          `json:"application-reference"`
	RejectedAt           string          `json:"rejected-at"`
	Reasons              json.RawMessage `json:"reasons"`
}

// Notifier delivers rejections to the partners' registered addresses.
type Notifier struct {
	logger       zerolog.Logger
	store        *pgstore.RejectionNotificationStore
	client       *http.Client
	requireHTTPS bool
}

// NewNotifier returns a Notifier.  With requireHTTPS set, an address that isn't
// https counts as a failed attempt and nothing is sent to it.
func NewNotifier(logger zerolog.Logger, store *pgstore.RejectionNotificationStore, requireHTTPS bool) *Notifier {
	return &Notifier{
		logger: logger,
		store:  store,
		client: &http.Client{
			Timeout: deliveryTimeout,
			// A redirect isn't a 2xx, so it is a failed attempt, not a new target.
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		requireHTTPS: requireHTTPS,
	}
}

// Run delivers due notifications until ctx is cancelled, waiting interval
// whenever none is due.
func (n *Notifier) Run(ctx context.Context, interval time.Duration) {
	poll(ctx, interval, n.logger, "notifier", n.ProcessNext)
}

// ProcessNext claims one due notification and attempts delivery.  It returns
// false when none was due.  Delivery is at-least-once: a crash after the POST
// repeats it with the same event-id.
func (n *Notifier) ProcessNext(ctx context.Context) (bool, error) {
	tx, err := n.store.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(context.Background())

	notification, err := n.store.ClaimNext(ctx, tx)
	if errors.Is(err, pgstore.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	log := n.logger.With().
		Str("application_reference", notification.Reference).
		Str("event_id", notification.EventID.String()).
		Int("attempt", notification.Attempts+1).
		Logger()

	switch err := n.deliver(ctx, notification); {
	case err == nil:
		err = n.store.MarkDelivered(ctx, tx, notification.EventID)
		if err != nil {
			return false, err
		}
		log.Info().Msg("rejection delivered")
	case ctx.Err() != nil:
		// Shutting down: the attempt doesn't count, and the rollback leaves it due.
		return false, ctx.Err()
	case notification.Attempts+1 >= maxAttempts:
		if err := n.store.MarkExhausted(ctx, tx, notification.EventID); err != nil {
			return false, err
		}
		log.Error().Msg("rejection delivery exhausted; partner must use the lookup")
	default:
		delay := retryDelays[notification.Attempts]
		if err := n.store.ScheduleRetry(ctx, tx, notification.EventID, delay); err != nil {
			return false, err
		}
		log.Warn().Dur("retry_in", delay).Msg("rejection delivery failed")
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("Notifier.ProcessNext: commit: %w", err)
	}
	return true, nil
}

// deliver signs and POSTs one event.  Any error is a failed attempt.
func (n *Notifier) deliver(ctx context.Context, notification pgstore.ClaimedNotification) error {
	if n.requireHTTPS {
		parsed, err := url.Parse(notification.NotificationURL)
		if err != nil || parsed.Scheme != "https" {
			return errors.New("notification address is not https")
		}
	}

	body, err := json.Marshal(rejectionEvent{
		EventID:              notification.EventID.String(),
		EventType:            eventTypeRejected,
		ApplicationReference: notification.Reference,
		RejectedAt:           notification.RejectedAt.UTC().Format(time.RFC3339),
		Reasons:              notification.Reasons,
	})
	if err != nil {
		return fmt.Errorf("Notifier.deliver: marshal: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, notification.NotificationURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Notifier.deliver: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-event-id", notification.EventID.String())
	req.Header.Set("x-timestamp", timestamp)
	req.Header.Set("x-signature", sign(notification.NotificationSecret, timestamp, body))

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("Notifier.deliver: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Notifier.deliver: status %d", resp.StatusCode)
	}
	return nil
}

// sign returns the x-signature value for one attempt.
func sign(secret, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
