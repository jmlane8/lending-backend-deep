# Implementation Plan: Submit Loan Application with Customer Data

**Branch**: `us-001-loan-submission` (feature `001-submit-loan-customer-data`) | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-submit-loan-customer-data/spec.md`

## Summary

A contracted partner submits a loan application with customer data. The API records the raw submission and acknowledges it with an application reference within 2 seconds. A background worker in the same service then validates the customer data. On failure it deletes the personal data and photos, keeps a minimal Rejected record, and queues a signed webhook. Partners can also read the validation result by reference.

PostgreSQL is both the record store and the work queue (`FOR UPDATE SKIP LOCKED`), so the feature adds no new infrastructure or libraries. The repository has no service code yet, so this feature also bootstraps it from `jmlane8/rest-api-template`.

## Technical Context

**Language/Version**: Go 1.26 (template `go.mod`)

**Primary Dependencies**: chi v5.3.1, pgx v5.10.0, zerolog v1.35.1, golang-migrate v4.19.1, caarlos0/env v11.4.1, google/uuid v1.6.0. No additions.

**Storage**: PostgreSQL 16+. Raw submissions and photos are stored as `BYTEA` ([research R-003](research.md)).

**Testing**: `go test` through `make test` (race detector on). Pure unit tests for validation rules; integration tests against a real PostgreSQL instance named by `TEST_DATABASE_URL`, using `net/http/httptest`.

**Target Platform**: Linux container or any Go 1.26 host. One binary serves HTTP and runs the workers; several replicas can run safely.

**Project Type**: Web service (REST API)

**Performance Goals**: Acknowledgment within 2 seconds for 100% of recorded submissions (SC-001). Validation finished within 1 minute of acknowledgment for 95% of submissions (SC-006).

**Constraints**: No business validation in the request-response cycle (Principle III). Personal data and photos are hard-deleted on rejection (FR-011). Rejection reasons and logs must never contain submitted values. No submission volume limit (AC-001-004).

**Scale/Scope**: 30,000 applications in 12 months (about 82 a day on average) from 20 partners. Request body limit 16 MiB; each decoded photo limit 5 MiB.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Principle | Pre-design | Post-design | Evidence |
|---|---|---|---|
| I. Objective traceability | PASS | PASS | Every FR and SC traces to O1–O4 in the spec. |
| II. Partner system as client | PASS | PASS | Partner-only API. No consumer accounts or UI. |
| III. Asynchronous submission | PASS | PASS with note | All business validation runs in the worker. **Note:** the request path still enforces authentication, a body size limit, and "body is a JSON object", because an unreadable body cannot become a record. See [R-002](research.md). |
| IV. API key security | PASS | PASS with deviation | Per-partner keys replace the template's static key. Only the minimum needed to attribute a request is built here; the lifecycle stays in US-001. See Complexity Tracking. |
| V. Partner self-sufficiency | PASS | PASS | REST, Redoc docs, sandbox deployment, and an enumerations endpoint. |
| VI. Scale without headcount | PASS | PASS | Fully automated; webhook retries and lookup replace support contact. |
| VII. Context7 | PASS | PASS | Library IDs and versions in [research.md](research.md), including noted version gaps. |
| VIII. Template conformance | PASS | PASS with deviations | Commit `147b01094598a24d01e59a214f538707c5369b03` read through GitHub MCP. File map and deviations in [research.md](research.md). |

**Gate result**: PASS. No `NEEDS CLARIFICATION` items remain. All deviations are justified below.

## Project Structure

### Documentation (this feature)

```text
specs/001-submit-loan-customer-data/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── loan-applications.openapi.yaml
│   └── rejection-webhook.md
└── tasks.md             # /speckit-tasks output, not created here
```

### Source Code (repository root)

Files marked `(new)` have no direct template counterpart. Everything else follows the template file of the same path or role.

```text
cmd/http/main.go
internal/
├── config/                       # config.go, http.go, pgsql.go, worker.go (new)
├── migrate/migrate.go
├── openapi/
├── rest/
│   ├── server/server.go          # routes, starts and stops workers
│   ├── response/envelope.go
│   ├── middleware/
│   │   ├── partnerauth.go        # replaces apikey.go
│   │   ├── logger.go
│   │   └── recover.go
│   └── handler/
│       ├── healthcheck/healthcheck.go
│       ├── loanapplication/
│       │   ├── submit.go
│       │   └── validationresult.go
│       └── masterdata/enumerations.go
├── store/pg/
│   ├── pool.go
│   ├── loan_application.go
│   ├── customer.go
│   ├── partner.go
│   ├── master_data.go
│   └── rejection_notification.go
├── validation/                   # (new) pure customer data rules, no I/O
│   ├── customer.go
│   └── customer_test.go
└── worker/                       # (new) background loops
    ├── validator.go
    └── notifier.go
resources/
├── openapi/
│   ├── openapi.yml
│   └── schema/{common,healthcheck,loanapplication,masterdata}/
├── gen/
└── pgsql/
    ├── migrations/001_create_loan_submission.{up,down}.sql
    └── seed/sandbox_partner.sql   # (new) local and sandbox only
tests/integration/                 # (new) needs TEST_DATABASE_URL
```

**Structure Decision**: Single Go service using the template's layout, with the module named `lending-backend-deep`. The calculator sample is not carried over. `internal/validation` and `internal/worker` are the only new packages, because the template has no background processing.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| Minimal `partner` and `partner_credential` tables built ahead of US-001 | FR-011, FR-013 and FR-014 need the submitting partner; Principle IV forbids the template's single static key | A static key cannot attribute a request to a partner |
| New `internal/worker` and `internal/validation` packages | Principle III moves validation out of handlers | Validating in a handler is prohibited; putting rules in `store/pg` would mix them with I/O |
| Template `deleted` soft-delete columns not used on submission data | FR-011 needs real deletion of personal data | A flag would keep the personal data |
| Handler returns `202`, not the template's `200` | The submission is accepted but not yet validated | `200` would imply the work is finished |
| Migration `seed/` folder outside golang-migrate | Sandbox partner data must not run in production | A migration runs in every environment |
| Database trigger for the `issued_at` rule | PostgreSQL `CHECK` constraints can't read another table | Enforcing it only in code would leave direct inserts unchecked |

The data model was reconciled against DDE-001 and ERD-001 on 2026-10-31. Agreed differences are listed in [data-model.md](data-model.md).
