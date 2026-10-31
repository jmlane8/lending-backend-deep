---

description: "Task list for Submit Loan Application with Customer Data"
---

# Tasks: Submit Loan Application with Customer Data

**Input**: Design documents from `specs/001-submit-loan-customer-data/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md), [my-technical-design-from-plan.md](my-technical-design-from-plan.md) (TD)

**Tests**: Included. The spec doesn't request TDD, but the plan names `internal/validation/customer_test.go` and `tests/integration/`, and quickstart scenario 13 runs `TestRejectionWebhook`. Test tasks are limited to those.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Paths are relative to the repository root, `lending-backend-deep/`

## Rules for Every Task

- **Template**: "Model:" names the `jmlane8/rest-api-template` file at commit `147b01094598a24d01e59a214f538707c5369b03`. Read it through GitHub MCP `get_file_contents` before writing the new file. Local copies and recollection don't count (Principle VIII). Replace the import prefix `rest-api-template` with `lending-backend-deep`.
- **Libraries**: Before using a library API, check it through Context7 with the IDs in [research.md](research.md). Exact versions were unavailable for some libraries (Principle VII).
- **Responses**: Use `response.WriteJSONWithEnvelope` and `response.WriteErrorWithEnvelope`. JSON fields use `lower-kebab-case`.
- **Personal data**: Never log, or put in rejection reasons or webhooks, any submitted value (FR-011).
- **DDL and SQL**: Use the DDL and base SQL in the TD as written. Don't redefine them.

## Decisions Made for This Task List

The plan and TD left these open. They're set here so tasks can run, and T045 copies them into the TD.

| Open item | Decision |
|---|---|
| `403` error code | `PARTNER_AGREEMENT_INACTIVE` |
| `404` error code | `APPLICATION_NOT_FOUND` |
| `413` error code | `REQUEST_BODY_TOO_LARGE` |
| Polling interval | `VALIDATOR_POLL_INTERVAL_SECONDS` default `1`; `NOTIFIER_POLL_INTERVAL_SECONDS` default `5` |
| Internal failure backoff | `next_validation_at = NOW() + LEAST(2^validation_attempts, 60)` minutes, retrying indefinitely |
| Alert target | Error-level log only, because no alerting system exists |
| Earlier `stored_photo` rows after a customer upsert | Soft-delete them per the template. Their content stays until retention rules are defined |

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the Go service from the template, without the calculator sample

- [ ] T001 Create `go.mod` with module `lending-backend-deep`, `go 1.26`, and the template's exact requirements, then run `go mod tidy` and commit `go.sum`. Model: `go.mod`
- [ ] T002 [P] Create `Makefile` with the template's targets unchanged, including `run` (there is no `serve` target). Model: `Makefile`
- [ ] T003 [P] Create `.env.example` and `.gitignore`. Set `APP_NAME=lending-backend-deep`, `POSTGRES_DATABASE=lending`, add `sandbox` to the `APP_ENV` comment, remove `API_KEY` and `POSTGRES_CONN_MAX_IDLE`, and add `BUSINESS_TIMEZONE=UTC`, `VALIDATOR_POLL_INTERVAL_SECONDS=1`, `NOTIFIER_POLL_INTERVAL_SECONDS=5`. Model: `.env.example`, `.gitignore`
- [ ] T004 [P] Create `internal/migrate/migrate.go`, `resources/pgsql/embed.go`, and `internal/store/pg/pool.go`, changing only imports and the package comment. Model: same paths
- [ ] T005 [P] Create `internal/rest/response/envelope.go`, `internal/rest/middleware/logger.go`, `internal/rest/middleware/recover.go`, and `internal/rest/handler/healthcheck/healthcheck.go`, changing only imports. Model: same paths
- [ ] T006 [P] Create `internal/openapi/ui.go` and `resources/gen/embed.go`. Model: same paths
- [ ] T007 Create `internal/config/config.go`, `internal/config/http.go`, and `internal/config/pgsql.go`, removing `APIKey` and defaulting `APP_NAME` to `lending-backend-deep`. Model: same paths

**Checkpoint**: The module has its shared packages, ready for the schema and routing

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, configuration, partner authentication, and server wiring that every story needs

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T008 Create `resources/pgsql/migrations/001_create_loan_submission.up.sql` with the TD DDL for all 8 tables in dependency order: `partner`, `partner_credential` (with the `issued_at` trigger), `stored_photo`, `customer`, `loan_application`, `application_payload`, `rejection_notification`, `master_data_value`, including its seed rows. Model: `resources/pgsql/migrations/001_create_calculator_transactions.up.sql`
- [ ] T009 Create `resources/pgsql/migrations/001_create_loan_submission.down.sql`, dropping everything from T008 in reverse order, including the trigger and function. Model: `resources/pgsql/migrations/001_create_calculator_transactions.down.sql`
- [ ] T010 Create `internal/config/worker.go` with `EnvWorker`: required `BusinessTimezone` (`env:"BUSINESS_TIMEZONE,required"`, no default), `ValidatorPollIntervalSeconds` default `1`, `NotifierPollIntervalSeconds` default `5`, using integer seconds with duration helper methods. Embed it in `Env` in `internal/config/config.go`. Model: `internal/config/http.go`
- [ ] T011 [P] Create `internal/store/pg/partner.go` with `PartnerStore.FindByKeyHash(ctx, hash []byte)`, returning partner ID, agreement status, or not found, using the TD Entry Point 1 SQL. Treat `pgx.ErrNoRows` as not found. Model: `internal/store/pg/calculator_tx.go`
- [ ] T012 Create `internal/rest/middleware/partnerauth.go` with `PartnerAuth`. Read header `api-key`, hash it with SHA-256, call `FindByKeyHash`, return `401 UNAUTHORIZED` for a missing or unknown key and `403 PARTNER_AGREEMENT_INACTIVE` when status isn't `ACTIVE`, and put the partner ID in the request context with an accessor function. Model: `internal/rest/middleware/apikey.go`
- [ ] T013 Create `internal/rest/server/server.go` with the template's router, middlewares, healthcheck, OpenAPI routes and graceful shutdown, plus an empty `PartnerAuth`-protected route group. Add a cancelable context for workers, stopped after `srv.Shutdown`. Model: `internal/rest/server/server.go`
- [ ] T014 Create `cmd/http/main.go`, adding a blank import of `time/tzdata` and exiting at startup if `BUSINESS_TIMEZONE` fails `time.LoadLocation`. Model: `cmd/http/main.go`
- [ ] T015 [P] Create `resources/openapi/openapi.yml`, `resources/openapi/schema/common/error.yaml`, and `resources/openapi/schema/healthcheck/healthcheck.yaml`, with title `Loan Submission API`, the `ApiKeyAuth` scheme from `contracts/loan-applications.openapi.yaml`, and the healthcheck path. Model: same paths
- [ ] T016 [P] Create `resources/pgsql/seed/sandbox_partner.sql`, with an `ACTIVE` partner keyed by the SHA-256 of `sandbox-partner-key`, a second `ACTIVE` partner keyed by `sandbox-partner-2-key` for lookup isolation, and an `INACTIVE` partner keyed by `inactive-partner-key`. Set `issued_at` after `agreement_signed_date`. No template model
- [ ] T017 Create `internal/validation/customer.go` with exported types only: `Reason{Field, Code, Message}`, the reason codes `MISSING`, `INVALID`, `CUSTOMER_UNDER_AGE`, `GOVERNMENT_ID_EXPIRED`, and a `Customer` struct holding parsed fields and decoded photos. No template model
- [ ] T018 [P] Create `tests/integration/helper_test.go`, which skips tests when `TEST_DATABASE_URL` is unset, applies migrations, loads the seed, cleans tables between tests, and builds a valid submission body with tiny JPEG photos. No template model
- [ ] T019 Run `make openapi`, then confirm `go build ./...` passes and the service starts with migrations applied, generating `resources/gen/openapi/openapi.json`

**Checkpoint**: Foundation ready. The service starts, authenticates partners, and has its schema

---

## Phase 3: User Story 1 - Submit a valid application and receive an acknowledgment (Priority: P1) 🎯 MVP

**Goal**: A partner submits valid data, gets a reference within 2 seconds, and sees the application pass validation

**Independent Test**: Submit complete, valid data once with and once without `email`. Each returns `202` within 2 seconds, and the lookup then shows `PASSED`.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T020 [P] [US1] Create `tests/integration/submit_test.go` covering `202` with an `LA-` reference, a missing `email` passing, `401` for a wrong key, `403` for an inactive partner, `400` for a non-object body, `413` over 16 MiB, nothing recorded for those four errors, and `404` for another partner's reference (AC-002-001, AC-002-002, FR-014)

### Implementation for User Story 1

- [ ] T021 [P] [US1] Create `resources/openapi/schema/loanapplication/submit.yaml` and `resources/openapi/schema/loanapplication/validationresult.yaml` from `contracts/loan-applications.openapi.yaml`. Mark nothing as `required` in the request. Add both paths to `resources/openapi/openapi.yml`. Model: `resources/openapi/schema/calculator/add.yaml`
- [ ] T022 [P] [US1] Create `internal/store/pg/loan_application.go` with `Create` (one transaction inserting `loan_application` and `application_payload`), `FindByReference` scoped to the partner, `ClaimNext`, `MarkPassed`, and `RecordFailure`, using the TD SQL. `application_status` is "`PENDING`, `PASSED`, or `REJECTED`". `RecordFailure` runs outside the rolled-back transaction with `next_validation_at = NOW() + LEAST(2^validation_attempts, 60)` minutes. Model: `internal/store/pg/calculator_tx.go`
- [ ] T023 [P] [US1] Create `internal/store/pg/customer.go` with `InsertPhotos` and `UpsertCustomer`, using `ON CONFLICT (national_id_number) DO UPDATE`, inside the caller's transaction. Soft-delete the customer's previous `stored_photo` rows by setting `deleted`, `deleted_at`, and `deleted_by = 'system'`. `monthly_income` is `NUMERIC(20,2)`. Model: `internal/store/pg/calculator_tx.go`
- [ ] T024 [US1] Add `Validate(body []byte, submissionDate time.Time, genders []string) (Customer, []Reason, error)` to `internal/validation/customer.go`, parsing the `customer` object with `json.Number` for numbers and decoding base64 photos into `Customer`. In this phase it returns no reasons. A parse failure returns an error, which the worker retries until US2 adds rules
- [ ] T025 [US1] Create `internal/rest/handler/loanapplication/submit.go`. Limit the body to 16 MiB with `http.MaxBytesReader` and return `413 REQUEST_BODY_TOO_LARGE` on `*http.MaxBytesError`. Decode into a struct with `Customer`, `Collateral`, and `Loan` as `json.RawMessage`, and return `400 INVALID_REQUEST_BODY` if the body isn't a JSON object. Don't type individual fields (R-002). Generate the ID with `uuid.NewV7`, set the reference to `LA-<uuid>`, store the exact body, and return `202`. Storage failure returns `500 INTERNAL_SERVER_ERROR`. Model: `internal/rest/handler/calculator/add.go`
- [ ] T026 [US1] Create `internal/rest/handler/loanapplication/validationresult.go` for `GET /loan-applications/v1/{application-reference}/validation-result`, using `chi.URLParam` and the context partner ID. Return `404 APPLICATION_NOT_FOUND` when absent or owned by another partner. Omit `reasons` unless `REJECTED`, and never report loan decisions. Model: `internal/rest/handler/calculator/add.go`
- [ ] T027 [US1] Create `internal/worker/validator.go` with a polling loop that stops on context cancel. Each cycle runs one transaction: claim with `FOR UPDATE OF la SKIP LOCKED`, derive the submission date from `submitted_at` with `BUSINESS_TIMEZONE`, load active `GENDER` values, call `Validate`, then upsert the customer, set `PASSED`, and rewrite the payload to hold only `collateral` and `loan`. On error, roll back and call `RecordFailure`. Log the reference, outcome, and delay since `submitted_at` only. No template model
- [ ] T028 [US1] Register `POST /loan-applications/v1` and the lookup route in the `PartnerAuth` group, and start the validator goroutine, in `internal/rest/server/server.go`
- [ ] T029 [US1] Run `make openapi` to update `resources/gen/openapi/openapi.json`, then run the US1 tests

**Checkpoint**: User Story 1 works on its own. Don't release it to partners yet, because invalid data isn't rejected until US2

---

## Phase 4: User Story 2 - Reject incomplete or invalid customer data with the affected fields identified (Priority: P2)

**Goal**: Invalid data is rejected after acknowledgment, every affected field is identified, personal data is deleted, and the partner is told

**Independent Test**: Omit each mandatory field in turn and send unaccepted values. Each gets `202`, then `REJECTED` with every affected field, no `application_payload` row, a NULL `customer_id`, and a signed webhook.

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T030 [P] [US2] Create `internal/validation/customer_test.go` with table-driven tests: each mandatory field absent, null, empty, and whitespace-only; several missing at once; wrong JSON types; both enums invalid; each FR-015 rule; a truncated image; a non-text `email` treated as absent; and NUL characters (AC-002-003 to AC-002-005, FR-015)
- [ ] T031 [P] [US2] Create `tests/integration/rejection_test.go` covering `REJECTED` with reasons, no `application_payload` row, NULL `customer_id`, `rejection_stage = 'CUSTOMER_DATA'`, wrong-typed `monthly-income` giving `202` then `INVALID`, and an immediate corrected resubmission passing (FR-011, FR-012, FR-016)
- [ ] T032 [P] [US2] Create `tests/integration/webhook_test.go` with `TestRejectionWebhook`, using `httptest`: a valid signature and headers, retry after a failure, `DELIVERED` on `2xx`, and `EXHAUSTED` after 7 attempts (FR-013)

### Implementation for User Story 2

- [ ] T033 [US2] Add the field rules to `Validate` in `internal/validation/customer.go`, running every rule so all problems appear in one rejection, with messages that never contain submitted values. Absent, null, or empty or whitespace-only text is `MISSING`. A wrong JSON type is `INVALID`. Apply these rules with code `INVALID`: `marital-status` one of `SINGLE`, `MARRIED`, `DIVORCED`, `WIDOWED` (case-sensitive); `occupation` one of `EMPLOYEE`, `ENTREPRENEUR`, `CIVIL_SERVANT` (case-sensitive); `national-id-number` "Letters and digits only"; `phone-number` "Optional `+`, then 8 to 15 digits"; `gender` "Active `GENDER` master data value"; `date-of-birth` "Real date in `YYYY-MM-DD`, not after the submission date"; `government-id-expiry-date` "Real date in `YYYY-MM-DD`"; `work-duration-years` and `number-of-dependents` "Whole number, 0 or more, at most 32,767"; `monthly-income` "Greater than 0, plain decimal, up to 18 whole digits and 2 decimals"; photos "Base64 JPEG or PNG that fully decodes, at most 5 MiB and 50 megapixels", using `image.DecodeConfig` for the pixel limit before a full decode; any text field "No NUL characters". A non-text `email` is treated as absent and never causes rejection
- [ ] T034 [US2] Add `MarkRejected` to `internal/store/pg/loan_application.go`. In the caller's transaction, hard-delete the `application_payload` row and set `REJECTED`, `rejection_reasons`, `rejection_stage = 'CUSTOMER_DATA'`, `validated_at`, and `rejected_at`, leaving `customer_id` NULL (FR-011, FR-016). This is the only hard delete in the system
- [ ] T035 [P] [US2] Create `internal/store/pg/rejection_notification.go` with `Enqueue` (the TD `INSERT ... SELECT` that adds no row when the partner has no `notification_url`), `ClaimNext` with `FOR UPDATE OF rn SKIP LOCKED`, `MarkDelivered`, `MarkExhausted`, and `ScheduleRetry`. `status` is "`PENDING`, `DELIVERED`, or `EXHAUSTED`". Model: `internal/store/pg/calculator_tx.go`
- [ ] T036 [US2] Update `internal/worker/validator.go` so non-empty reasons call `MarkRejected` and `Enqueue` in the same transaction, logging a warning with the reference if no notification row is added
- [ ] T037 [US2] Create `internal/worker/notifier.go` per `contracts/rejection-webhook.md`: event type `loan-application.rejected`; headers `x-event-id`, `x-timestamp`, `x-signature`; HMAC-SHA256 of `<x-timestamp>.<raw body>`; a 5-second timeout; retries after 1 minute, 5 minutes, 30 minutes, 2 hours, 6 hours and 12 hours; `EXHAUSTED` after 7 attempts. No template model
- [ ] T038 [US2] Create `internal/store/pg/master_data.go` with `ListActive` and `ActiveGenders`, `internal/rest/handler/masterdata/enumerations.go`, and `resources/openapi/schema/masterdata/enumerations.yaml` for `GET /master-data/v1/enumerations`, and add the path to `resources/openapi/openapi.yml`. Replace the validator's gender lookup with `ActiveGenders`. Model: `internal/store/pg/calculator_tx.go`, `internal/rest/handler/calculator/add.go`, `resources/openapi/schema/calculator/add.yaml`
- [ ] T039 [US2] Register the enumerations route in the `PartnerAuth` group and start the notifier goroutine in `internal/rest/server/server.go`, then run `make openapi` to update `resources/gen/openapi/openapi.json`

**Checkpoint**: User Stories 1 and 2 work together, and invalid data no longer proceeds

---

## Phase 5: User Story 3 - Enforce customer eligibility: minimum age and valid government ID (Priority: P3)

**Goal**: Applications pass only for customers aged 18 or over on the submission date, with an unexpired ID

**Independent Test**: Submit a 17-year-old, a customer turning 18 on the submission date, an expired ID, and an ID expiring on the submission date. The first and third are rejected, and the others pass.

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T040 [US3] Add table-driven tests to `internal/validation/customer_test.go`: under 18 by one day, 18 today, a 29 February birth in a non-leap year, ID expired yesterday, and ID expiring today (AC-002-006 to AC-002-008)

### Implementation for User Story 3

- [ ] T041 [US3] Add eligibility rules to `Validate` in `internal/validation/customer.go`. `date-of-birth`: "Customer is 18 on the submission date. 29 February births reach 18 on 1 March in non-leap years", with code `CUSTOMER_UNDER_AGE`. `government-id-expiry-date`: "Not before the submission date", with code `GOVERNMENT_ID_EXPIRED`. Skip each rule if its date already failed as `MISSING` or `INVALID`

**Checkpoint**: All user stories work

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, verification, and keeping the design documents in step

- [ ] T042 [P] Create `README.md` for this service, covering setup with `make run`, endpoints, the asynchronous flow, webhook verification, sandbox use, and every environment variable. Model: `README.md`
- [ ] T043 Check that no log line in `internal/worker/` or `internal/rest/` can contain a submitted value, and that request bodies are never logged (FR-011)
- [ ] T044 Measure SC-001 and SC-006 against a running service with maximum-size bodies, using `submitted_at`, `acknowledged_at`, and `validated_at`, and record the results in `specs/001-submit-loan-customer-data/quickstart.md`
- [ ] T045 [P] Update `specs/001-submit-loan-customer-data/my-technical-design-from-plan.md` to replace every `TODO` with the decisions at the top of this file, and update the first paragraph of `specs/001-submit-loan-customer-data/data-model.md` to say every table has the audit and soft-delete columns
- [ ] T046 Run all 20 scenarios in `specs/001-submit-loan-customer-data/quickstart.md` and `make test` with `TEST_DATABASE_URL` set, then fix any failures

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup. Blocks all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational
- **User Story 2 (Phase 4)**: Depends on US1
- **User Story 3 (Phase 5)**: Depends on US1 and US2
- **Polish (Phase 6)**: Depends on all stories

### User Story Dependencies

- **US1 (P1)**: Provides `Validate`, the worker, and the lookup
- **US2 (P2)**: Extends `Validate` and the worker with the rejection path. It can't run in parallel with US1
- **US3 (P3)**: Adds rules that use the US2 rejection path. It can't run in parallel with US2

### Within Each User Story

- Tests are written first and must fail before implementation
- Stores before handlers and workers
- Handlers and workers before route registration

### Parallel Opportunities

- Setup: T002 to T006 after T001
- Foundational: T011, T015, T016, and T018 after T008 to T010
- US1: T020 to T023
- US2: T030 to T032, then T035 alongside T033 and T034
- Polish: T042 and T045

---

## Parallel Example: User Story 1

```bash
Task: "Create tests/integration/submit_test.go"
Task: "Create OpenAPI schemas in resources/openapi/schema/loanapplication/"
Task: "Create internal/store/pg/loan_application.go"
Task: "Create internal/store/pg/customer.go"
```

---

## Implementation Strategy

### MVP First

1. Complete Phases 1 and 2.
2. Complete Phase 3.
3. **STOP and VALIDATE** with the US1 independent test.
4. Treat this as an internal demo only: it accepts data without rejecting it.

### Incremental Delivery

1. Setup and Foundational give a running, authenticated service.
2. US1 gives submission, acknowledgment, and lookup.
3. US2 gives rejection, deletion, and webhooks. Sandbox access for partners can start here.
4. US3 gives eligibility rules. This is the earliest production-ready point.

---

## Notes

- [P] tasks touch different files and have no dependencies on incomplete tasks
- Commit after each task or logical group
- Stop at any checkpoint to validate the story independently
- The questions for the business stakeholder remain open: country and currency, the narrowed lookup, gender values, and the webhook
