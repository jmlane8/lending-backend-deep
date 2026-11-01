package integration

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const webhookSecret = "test-notification-secret"

func (e *env) registerWebhook(t *testing.T, url string) {
	t.Helper()
	if _, err := e.pool.Exec(context.Background(), `
		UPDATE partner SET notification_url = $1, notification_secret = $2 WHERE partner_id = $3
	`, url, webhookSecret, partnerID); err != nil {
		t.Fatal(err)
	}
}

// deliverDue makes the pending notification due now and runs one attempt.
func (e *env) deliverDue(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	if _, err := e.pool.Exec(ctx, `UPDATE rejection_notification SET next_attempt_at = NOW()`); err != nil {
		t.Fatal(err)
	}
	if _, err := e.notifier.ProcessNext(ctx); err != nil {
		t.Fatalf("ProcessNext: %v", err)
	}
}

func (e *env) notification(t *testing.T) (status string, attempts int, delay time.Duration) {
	t.Helper()
	if err := e.pool.QueryRow(context.Background(), `
		SELECT status, attempts, next_attempt_at - NOW() FROM rejection_notification
	`).Scan(&status, &attempts, &delay); err != nil {
		t.Fatal(err)
	}
	return status, attempts, delay
}

func TestRejectionWebhook(t *testing.T) {
	t.Run("signed delivery after a retry", func(t *testing.T) {
		e := setup(t)
		var calls atomic.Int32
		var gotBody []byte
		var gotHeaders http.Header
		partner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if calls.Add(1) == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			gotBody, _ = io.ReadAll(r.Body)
			gotHeaders = r.Header.Clone()
			w.WriteHeader(http.StatusNoContent)
		}))
		defer partner.Close()
		e.registerWebhook(t, partner.URL)

		customer := validCustomer(t)
		delete(customer, "gender")
		customer["phone-number"] = "SECRETPHONE"
		reference := e.submit(t, submission(t, customer))
		e.validateAll(t)

		// First attempt fails and is rescheduled about a minute later.
		e.deliverDue(t)
		status, attempts, delay := e.notification(t)
		if status != "PENDING" || attempts != 1 || delay < 50*time.Second || delay > time.Minute {
			t.Fatalf("after failure: %s, %d, %s", status, attempts, delay)
		}

		// It isn't due yet, so nothing is sent.
		if processed, err := e.notifier.ProcessNext(context.Background()); err != nil || processed {
			t.Fatalf("processed = %v, err = %v", processed, err)
		}

		e.deliverDue(t)
		status, attempts, _ = e.notification(t)
		if status != "DELIVERED" || attempts != 2 {
			t.Fatalf("after success: %s, %d", status, attempts)
		}

		timestamp := gotHeaders.Get("x-timestamp")
		mac := hmac.New(sha256.New, []byte(webhookSecret))
		mac.Write([]byte(timestamp + "." + string(gotBody)))
		want := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if gotHeaders.Get("x-signature") != want {
			t.Errorf("signature = %q, want %q", gotHeaders.Get("x-signature"), want)
		}
		if gotHeaders.Get("Content-Type") != "application/json" {
			t.Errorf("content type = %q", gotHeaders.Get("Content-Type"))
		}

		var event struct {
			EventID   string `json:"event-id"`
			EventType string `json:"event-type"`
			Reference string `json:"application-reference"`
			Rejected  string `json:"rejected-at"`
			Reasons   []struct {
				Field string `json:"field"`
				Code  string `json:"code"`
			} `json:"reasons"`
		}
		if err := json.Unmarshal(gotBody, &event); err != nil {
			t.Fatal(err)
		}
		if event.EventType != "loan-application.rejected" || event.Reference != reference ||
			event.EventID != gotHeaders.Get("x-event-id") || len(event.Reasons) != 2 {
			t.Errorf("event = %+v", event)
		}
		if _, err := time.Parse(time.RFC3339, event.Rejected); err != nil {
			t.Errorf("rejected-at: %v", err)
		}
		if strings.Contains(string(gotBody), "SECRETPHONE") {
			t.Errorf("webhook leaks the value: %s", gotBody)
		}
	})

	t.Run("exhausted after 7 attempts", func(t *testing.T) {
		e := setup(t)
		var calls atomic.Int32
		partner := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			calls.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer partner.Close()
		e.registerWebhook(t, partner.URL)

		customer := validCustomer(t)
		delete(customer, "gender")
		reference := e.submit(t, submission(t, customer))
		e.validateAll(t)

		wantDelays := []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute,
			2 * time.Hour, 6 * time.Hour, 12 * time.Hour}
		for i, want := range wantDelays {
			e.deliverDue(t)
			status, attempts, delay := e.notification(t)
			if status != "PENDING" || attempts != i+1 || delay > want || delay < want-10*time.Second {
				t.Fatalf("attempt %d: %s, %d, %s", i+1, status, attempts, delay)
			}
		}
		e.deliverDue(t)

		status, attempts, _ := e.notification(t)
		if status != "EXHAUSTED" || attempts != 7 || calls.Load() != 7 {
			t.Fatalf("%s, attempts = %d, calls = %d", status, attempts, calls.Load())
		}
		if got := e.lookup(t, reference); got.Data.Status != "REJECTED" {
			t.Errorf("lookup after exhaustion: %s", got.Data.Status)
		}
	})

	t.Run("a redirect is a failed attempt", func(t *testing.T) {
		e := setup(t)
		var reached atomic.Bool
		target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			reached.Store(true)
		}))
		defer target.Close()
		partner := httptest.NewServer(http.RedirectHandler(target.URL, http.StatusTemporaryRedirect))
		defer partner.Close()
		e.registerWebhook(t, partner.URL)

		customer := validCustomer(t)
		delete(customer, "gender")
		e.submit(t, submission(t, customer))
		e.validateAll(t)
		e.deliverDue(t)

		if status, attempts, _ := e.notification(t); status != "PENDING" || attempts != 1 || reached.Load() {
			t.Errorf("%s, attempts = %d, reached = %v", status, attempts, reached.Load())
		}
	})

	t.Run("no registered address sends nothing", func(t *testing.T) {
		e := setup(t)
		customer := validCustomer(t)
		delete(customer, "gender")
		reference := e.submit(t, submission(t, customer))
		e.validateAll(t)

		if n := e.count(t, "rejection_notification"); n != 0 {
			t.Errorf("notifications = %d", n)
		}
		if got := e.lookup(t, reference); got.Data.Status != "REJECTED" {
			t.Errorf("status = %s", got.Data.Status)
		}
	})
}
