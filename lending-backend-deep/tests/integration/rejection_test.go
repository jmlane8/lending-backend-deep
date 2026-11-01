package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestRejectionKeepsOnlyAMinimalRecord(t *testing.T) {
	e := setup(t)
	customer := validCustomer(t)
	delete(customer, "gender")
	delete(customer, "phone-number")

	reference := e.submit(t, submission(t, customer))
	e.validateAll(t)

	got := e.lookup(t, reference)
	if got.Data.Status != "REJECTED" || got.Data.Validated == nil || len(got.Data.Reasons) != 2 {
		t.Fatalf("result = %+v", got.Data)
	}
	for _, r := range got.Data.Reasons {
		if r.Code != "MISSING" {
			t.Errorf("reason = %+v", r)
		}
	}

	if n := e.count(t, "application_payload"); n != 0 {
		t.Errorf("payload rows = %d", n)
	}
	if n := e.count(t, "customer"); n != 0 {
		t.Errorf("customers = %d", n)
	}
	if n := e.count(t, "stored_photo"); n != 0 {
		t.Errorf("photos = %d", n)
	}

	var stage string
	var customerIsNull, rejectedSet bool
	if err := e.pool.QueryRow(context.Background(), `
		SELECT rejection_stage, customer_id IS NULL, rejected_at IS NOT NULL FROM loan_application
	`).Scan(&stage, &customerIsNull, &rejectedSet); err != nil {
		t.Fatal(err)
	}
	if stage != "CUSTOMER_DATA" || !customerIsNull || !rejectedSet {
		t.Errorf("stage = %q, customer NULL = %v, rejected_at set = %v", stage, customerIsNull, rejectedSet)
	}
}

func TestWrongJSONTypeIsAcknowledgedThenInvalid(t *testing.T) {
	e := setup(t)
	customer := validCustomer(t)
	customer["monthly-income"] = "lots"

	reference := e.submit(t, submission(t, customer))
	e.validateAll(t)

	got := e.lookup(t, reference)
	if got.Data.Status != "REJECTED" || len(got.Data.Reasons) != 1 ||
		got.Data.Reasons[0].Field != "customer.monthly-income" || got.Data.Reasons[0].Code != "INVALID" {
		t.Errorf("result = %+v", got.Data)
	}
}

func TestRejectedApplicationCanBeResubmittedImmediately(t *testing.T) {
	e := setup(t)
	bad := validCustomer(t)
	delete(bad, "gender")

	rejected := e.submit(t, submission(t, bad))
	e.validateAll(t)
	if got := e.lookup(t, rejected); got.Data.Status != "REJECTED" {
		t.Fatalf("status = %s", got.Data.Status)
	}

	corrected := e.submit(t, submission(t, validCustomer(t)))
	e.validateAll(t)
	if got := e.lookup(t, corrected); got.Data.Status != "PASSED" {
		t.Errorf("status = %s, reasons = %+v", got.Data.Status, got.Data.Reasons)
	}
}

func TestLookupReasonsHoldNoSubmittedValues(t *testing.T) {
	e := setup(t)
	customer := validCustomer(t)
	customer["phone-number"] = "SECRETPHONE"

	reference := e.submit(t, submission(t, customer))
	e.validateAll(t)

	rec := e.do(t, http.MethodGet, "/loan-applications/v1/"+reference+"/validation-result", partnerKey, nil)
	if strings.Contains(rec.Body.String(), "SECRETPHONE") {
		t.Errorf("response leaks the value: %s", rec.Body)
	}
}

func TestEnumerationsListAcceptedValues(t *testing.T) {
	e := setup(t)
	rec := e.do(t, http.MethodGet, "/master-data/v1/enumerations", partnerKey, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var resp struct {
		Data map[string][]string `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Data["gender"]) != 2 || len(resp.Data["marital-status"]) != 4 || len(resp.Data["occupation"]) != 3 {
		t.Errorf("data = %+v", resp.Data)
	}
	if rec := e.do(t, http.MethodGet, "/master-data/v1/enumerations", "", nil); rec.Code != http.StatusUnauthorized {
		t.Errorf("no key: %d", rec.Code)
	}
}

func TestGenderValueAddedWithoutCodeChange(t *testing.T) {
	e := setup(t)
	ctx := context.Background()
	if _, err := e.pool.Exec(ctx, `
		INSERT INTO master_data_value (category, value, label) VALUES ('GENDER', 'TEST_VALUE', 'Test')
	`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = e.pool.Exec(context.Background(), `DELETE FROM master_data_value WHERE value = 'TEST_VALUE'`)
	})

	customer := validCustomer(t)
	customer["gender"] = "TEST_VALUE"
	reference := e.submit(t, submission(t, customer))
	e.validateAll(t)

	if got := e.lookup(t, reference); got.Data.Status != "PASSED" {
		t.Errorf("status = %s, reasons = %+v", got.Data.Status, got.Data.Reasons)
	}
}
