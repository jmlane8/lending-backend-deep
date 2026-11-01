package integration

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSubmitAcknowledgesThenPasses(t *testing.T) {
	e := setup(t)

	started := time.Now()
	reference := e.submit(t, submission(t, validCustomer(t)))
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Errorf("acknowledged in %s", elapsed)
	}

	if got := e.lookup(t, reference); got.Data.Status != "PENDING" || got.Data.Validated != nil {
		t.Errorf("before validation: %+v", got.Data)
	}

	e.validateAll(t)

	got := e.lookup(t, reference)
	if got.Data.Status != "PASSED" || got.Data.Validated == nil || got.Data.Reasons != nil {
		t.Errorf("after validation: %+v", got.Data)
	}
	if n := e.count(t, "customer"); n != 1 {
		t.Errorf("customers = %d", n)
	}
	if n := e.count(t, "stored_photo"); n != 2 {
		t.Errorf("photos = %d", n)
	}

	var body []byte
	if err := e.pool.QueryRow(context.Background(), `SELECT body FROM application_payload`).Scan(&body); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(body, []byte("Alex")) || !bytes.Contains(body, []byte(`"collateral"`)) {
		t.Errorf("payload not reduced: %s", body)
	}
}

func TestSubmitWithoutEmailPasses(t *testing.T) {
	e := setup(t)
	customer := validCustomer(t)
	delete(customer, "email")

	reference := e.submit(t, submission(t, customer))
	e.validateAll(t)

	if got := e.lookup(t, reference); got.Data.Status != "PASSED" {
		t.Errorf("status = %s, reasons = %+v", got.Data.Status, got.Data.Reasons)
	}
}

func TestSecondApplicationReplacesCustomerPhotos(t *testing.T) {
	e := setup(t)
	e.submit(t, submission(t, validCustomer(t)))
	e.validateAll(t)

	customer := validCustomer(t)
	customer["full-legal-name"] = "Alex Q Morgan"
	e.submit(t, submission(t, customer))
	e.validateAll(t)

	if n := e.count(t, "customer"); n != 1 {
		t.Errorf("customers = %d", n)
	}
	ctx := context.Background()
	var name string
	if err := e.pool.QueryRow(ctx, `SELECT full_legal_name FROM customer`).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != "Alex Q Morgan" {
		t.Errorf("name = %q", name)
	}
	var deleted int
	if err := e.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM stored_photo WHERE deleted AND deleted_by = 'system' AND deleted_at IS NOT NULL
	`).Scan(&deleted); err != nil {
		t.Fatal(err)
	}
	if deleted != 2 {
		t.Errorf("soft-deleted photos = %d", deleted)
	}
}

func TestProtocolFailuresRecordNothing(t *testing.T) {
	e := setup(t)
	valid := submission(t, validCustomer(t))

	tests := []struct {
		name   string
		key    string
		body   []byte
		status int
		code   string
	}{
		{"missing key", "", valid, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"wrong key", "wrong", valid, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"inactive partner", inactiveKey, valid, http.StatusForbidden, "PARTNER_AGREEMENT_INACTIVE"},
		{"not json", partnerKey, []byte("not json"), http.StatusBadRequest, "INVALID_REQUEST_BODY"},
		{"json array", partnerKey, []byte("[1]"), http.StatusBadRequest, "INVALID_REQUEST_BODY"},
		{"json null", partnerKey, []byte("null"), http.StatusBadRequest, "INVALID_REQUEST_BODY"},
		{"trailing data", partnerKey, []byte("{} x"), http.StatusBadRequest, "INVALID_REQUEST_BODY"},
		{"over 16 MiB", partnerKey, []byte(`{"a":"` + strings.Repeat("x", 16<<20) + `"}`),
			http.StatusRequestEntityTooLarge, "REQUEST_BODY_TOO_LARGE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := e.do(t, http.MethodPost, "/loan-applications/v1", tt.key, tt.body)
			if rec.Code != tt.status || !strings.Contains(rec.Body.String(), tt.code) {
				t.Errorf("got %d %s", rec.Code, rec.Body)
			}
			if n := e.count(t, "loan_application"); n != 0 {
				t.Errorf("applications recorded = %d", n)
			}
		})
	}
}

func TestLookupIsScopedToThePartner(t *testing.T) {
	e := setup(t)
	reference := e.submit(t, submission(t, validCustomer(t)))
	path := "/loan-applications/v1/" + reference + "/validation-result"

	rec := e.do(t, http.MethodGet, path, partner2Key, nil)
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "APPLICATION_NOT_FOUND") {
		t.Errorf("other partner: %d %s", rec.Code, rec.Body)
	}
	rec = e.do(t, http.MethodGet, "/loan-applications/v1/LA-unknown/validation-result", partnerKey, nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown reference: %d", rec.Code)
	}
}

func TestInternalFailureRetriesLater(t *testing.T) {
	e := setup(t)
	reference := e.submit(t, submission(t, validCustomer(t)))

	// A valid submission can't fail on its own, so make the insert fail.
	ctx := context.Background()
	if _, err := e.pool.Exec(ctx, `
		ALTER TABLE customer ADD CONSTRAINT test_block CHECK (full_legal_name = 'blocked') NOT VALID
	`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = e.pool.Exec(context.Background(), `ALTER TABLE customer DROP CONSTRAINT IF EXISTS test_block`)
	})

	e.validateAll(t)

	var attempts int
	var delay time.Duration
	if err := e.pool.QueryRow(ctx, `
		SELECT validation_attempts, next_validation_at - NOW() FROM loan_application
	`).Scan(&attempts, &delay); err != nil {
		t.Fatal(err)
	}
	if attempts != 1 || delay < 50*time.Second {
		t.Errorf("attempts = %d, delay = %s", attempts, delay)
	}
	if got := e.lookup(t, reference); got.Data.Status != "PENDING" {
		t.Errorf("status = %s", got.Data.Status)
	}
	if n := e.count(t, "stored_photo"); n != 0 {
		t.Errorf("photos left by the rolled-back attempt = %d", n)
	}

	if _, err := e.pool.Exec(ctx, `ALTER TABLE customer DROP CONSTRAINT test_block`); err != nil {
		t.Fatal(err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE loan_application SET next_validation_at = NOW()`); err != nil {
		t.Fatal(err)
	}
	e.validateAll(t)
	if got := e.lookup(t, reference); got.Data.Status != "PASSED" {
		t.Errorf("status after recovery = %s", got.Data.Status)
	}
}
