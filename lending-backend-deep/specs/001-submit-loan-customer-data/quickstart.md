# Quickstart: Validate the Feature End to End

Use this after implementation to prove the feature works. Request and response shapes are in [contracts/](contracts/), and rules are in [data-model.md](data-model.md).

## Prerequisites

- Go 1.26+, PostgreSQL 16+, `make`, and `psql`
- An empty database for the service and another for tests

## Setup

```bash
make env                 # then edit .env
make run                 # applies migrations and starts the API and workers
psql "$DATABASE_URL" -f resources/pgsql/seed/sandbox_partner.sql
```

Set these in `.env`:

| Variable | Purpose |
|---|---|
| `POSTGRES_*` | Database connection, as in the template |
| `BUSINESS_TIMEZONE` | Required. For example `UTC` for local testing |
| `OPENAPI_UI_ENABLE=true` | Serves docs at `/redoc` |

Also export `DATABASE_URL` in your shell, pointing to the same database as `POSTGRES_*`. Only `psql` uses it.

The seed creates an active partner with the API key `sandbox-partner-key` and an inactive partner with `inactive-partner-key`. Never run it in production.

## Automated checks

```bash
make test                                         # unit tests
TEST_DATABASE_URL=postgres://... make test        # adds integration tests
```

## Manual scenarios

Send each request to `http://localhost:8080` with the `api-key` header. Use tiny valid JPEG or PNG files, base64-encoded, for the two photo fields.

| # | Scenario | Action | Expected outcome | Covers |
|---|---|---|---|---|
| 1 | Valid submission | `POST /loan-applications/v1` with all fields | `202` in under 2 seconds with `LA-...`; lookup shows `PASSED` within a minute | AC-002-001, SC-001, SC-006 |
| 2 | Email left out | Same, without `email` | Lookup shows `PASSED` | AC-002-002, SC-005 |
| 3 | Fields missing | Remove `gender` and `phone-number` | `202`; lookup shows `REJECTED` with both fields `MISSING` | AC-002-003 |
| 4 | Accepted values | Send `marital-status: ENGAGED` and `occupation: STUDENT` | Lookup shows both fields `INVALID` | AC-002-004, AC-002-005 |
| 5 | Under 18 | `date-of-birth` one day short of 18 years before today | `CUSTOMER_UNDER_AGE` | AC-002-006 |
| 6 | Turns 18 today | `date-of-birth` exactly 18 years before today | `PASSED` | AC-002-007 |
| 7 | Expired ID | `government-id-expiry-date` yesterday | `GOVERNMENT_ID_EXPIRED` | AC-002-008 |
| 8 | ID expires today | `government-id-expiry-date` today | `PASSED` | Edge case |
| 9 | Basic value checks | Bad phone number, `monthly-income: 0`, and a non-image photo | All three listed as `INVALID` | FR-015 |
| 10 | Wrong JSON type | `monthly-income: "lots"` | `202`, then `INVALID`, not `400` | FR-012 |
| 11 | Rejected data is gone | After scenario 3, query the database | No `application_payload` row; `customer_id` is NULL | FR-011 |
| 12 | Correct and resubmit | Resend scenario 3 immediately with fixes | `202`, then `PASSED`, with no wait | FR-016 |
| 13 | Webhook sent | Run `go test ./tests/integration -run TestRejectionWebhook` | Signed delivery, retry after failure, stops after 7 attempts | FR-013 |
| 14 | Lookup isolation | Fetch a reference with a second partner's key | `404` | FR-014 |
| 15 | Lookup scope | Inspect a `PASSED` response | Validation result only, no loan decision | FR-014 |
| 16 | Access control | Send a wrong key, then `inactive-partner-key` | `401`, then `403`; no application recorded | US-001 |
| 17 | Unreadable body | Send `not json`, then a body over 16 MiB | `400`, then `413`; no application recorded | R-002 |
| 18 | Crash recovery | Stop the service while rows are `PENDING`, then restart | Rows finish normally | R-001 |
| 19 | Accepted values listed | `GET /master-data/v1/enumerations` | Includes `MALE` and `FEMALE` | Principle V |
| 20 | No personal data in logs | Search logs after scenarios 1 to 9 | No submitted values appear | FR-011 |
