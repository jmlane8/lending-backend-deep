// Package integration exercises the service against a real PostgreSQL named by
// TEST_DATABASE_URL.  Tests are skipped when it is unset.
package integration

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"lending-backend-deep/internal/config"
	dbmigrate "lending-backend-deep/internal/migrate"
	"lending-backend-deep/internal/rest/server"
	pgstore "lending-backend-deep/internal/store/pg"
	"lending-backend-deep/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	partnerKey   = "sandbox-partner-key"
	partner2Key  = "sandbox-partner-2-key"
	inactiveKey  = "inactive-partner-key"
	partnerID    = "01900000-0000-7000-8000-000000000001"
	seedFilePath = "../../resources/pgsql/seed/sandbox_partner.sql"
)

type env struct {
	pool      *pgxpool.Pool
	handler   http.Handler
	validator *worker.Validator
	notifier  *worker.Notifier
}

func setup(t *testing.T) *env {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()

	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}
	pgCfg := config.EnvPostgreSQL{
		PostgresHost:     poolCfg.ConnConfig.Host,
		PostgresPort:     int(poolCfg.ConnConfig.Port),
		PostgresDatabase: poolCfg.ConnConfig.Database,
		PostgresSchema:   "public",
		PostgresUsername: poolCfg.ConnConfig.User,
		PostgresPassword: poolCfg.ConnConfig.Password,
		PostgresSSLMode:  "disable",
	}
	logger := zerolog.Nop()
	if err := dbmigrate.Run(ctx, pgCfg, logger); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if _, err := pool.Exec(ctx, `
		TRUNCATE rejection_notification, application_payload, loan_application,
		         customer, stored_photo, partner_credential, partner`); err != nil {
		t.Fatalf("clean tables: %v", err)
	}
	seed, err := os.ReadFile(seedFilePath)
	if err != nil {
		t.Fatalf("read seed: %v", err)
	}
	if _, err := pool.Exec(ctx, string(seed)); err != nil {
		t.Fatalf("load seed: %v", err)
	}

	cfg := &config.Env{}
	return &env{
		pool:    pool,
		handler: server.NewRouter(cfg, logger, pool),
		validator: worker.NewValidator(logger,
			pgstore.NewLoanApplicationStore(pool), pgstore.NewCustomerStore(),
			pgstore.NewMasterDataStore(pool), pgstore.NewRejectionNotificationStore(pool), time.UTC),
		notifier: worker.NewNotifier(logger, pgstore.NewRejectionNotificationStore(pool), false),
	}
}

// do sends a request through the router with the given api-key.
func (e *env) do(t *testing.T, method, path, key string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if key != "" {
		req.Header.Set("api-key", key)
	}
	rec := httptest.NewRecorder()
	e.handler.ServeHTTP(rec, req)
	return rec
}

// submit posts a body, expects 202, and returns the application reference.
func (e *env) submit(t *testing.T, body []byte) string {
	t.Helper()
	rec := e.do(t, http.MethodPost, "/loan-applications/v1", partnerKey, body)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("submit status = %d, body = %s", rec.Code, rec.Body)
	}
	var resp struct {
		Data struct {
			Reference   string `json:"application-reference"`
			SubmittedAt string `json:"submitted-at"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(resp.Data.Reference, "LA-") {
		t.Fatalf("reference = %q", resp.Data.Reference)
	}
	if _, err := time.Parse(time.RFC3339, resp.Data.SubmittedAt); err != nil {
		t.Fatalf("submitted-at: %v", err)
	}
	return resp.Data.Reference
}

// validateAll runs the validator until no application is due.
func (e *env) validateAll(t *testing.T) {
	t.Helper()
	for {
		processed, err := e.validator.ProcessNext(context.Background())
		if err != nil {
			t.Fatalf("ProcessNext: %v", err)
		}
		if !processed {
			return
		}
	}
}

type validationResult struct {
	Data struct {
		Reference string  `json:"application-reference"`
		Status    string  `json:"status"`
		Validated *string `json:"validated-at"`
		Reasons   []struct {
			Field string `json:"field"`
			Code  string `json:"code"`
		} `json:"reasons"`
	} `json:"data"`
}

func (e *env) lookup(t *testing.T, reference string) validationResult {
	t.Helper()
	rec := e.do(t, http.MethodGet, "/loan-applications/v1/"+reference+"/validation-result", partnerKey, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("lookup status = %d, body = %s", rec.Code, rec.Body)
	}
	var result validationResult
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func (e *env) count(t *testing.T, table string) int {
	t.Helper()
	var n int
	if err := e.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM "+table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// validCustomer returns customer data that passes every rule.
func validCustomer(t *testing.T) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 4, 4)), nil); err != nil {
		t.Fatal(err)
	}
	photo := base64.StdEncoding.EncodeToString(buf.Bytes())
	return map[string]any{
		"government-id-photo":       photo,
		"government-id-expiry-date": "2099-05-14",
		"national-id-number":        "3201012345678901",
		"full-legal-name":           "Alex Morgan",
		"selfie-photo":              photo,
		"phone-number":              "+15551234567",
		"email":                     "alex@example.com",
		"gender":                    "FEMALE",
		"date-of-birth":             "1990-02-11",
		"place-of-birth":            "Springfield",
		"marital-status":            "SINGLE",
		"occupation":                "EMPLOYEE",
		"work-duration-years":       6,
		"monthly-income":            4500.5,
		"number-of-dependents":      1,
		"domicile-street":           "12 Example Street",
		"domicile-district":         "Central",
		"domicile-city":             "Springfield",
		"domicile-province":         "North",
	}
}

func submission(t *testing.T, customer map[string]any) []byte {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"customer":   customer,
		"collateral": map[string]any{"plate": "AB-1"},
		"loan":       map[string]any{"amount": 1000},
	})
	if err != nil {
		t.Fatal(err)
	}
	return body
}
