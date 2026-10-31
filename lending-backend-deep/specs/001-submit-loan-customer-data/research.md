# Research: Submit Loan Application with Customer Data

## Template Baseline (Principle VIII)

**Template**: `jmlane8/rest-api-template`, commit `147b01094598a24d01e59a214f538707c5369b03`, read through GitHub MCP on 2026-10-31.

**Files read**: `README.md`, `go.mod`, `Makefile`, `.env.example`, `cmd/http/main.go`, `internal/config/config.go`, `internal/migrate/migrate.go`, `internal/rest/server/server.go`, `internal/rest/response/envelope.go`, `internal/rest/middleware/apikey.go`, `internal/rest/middleware/logger.go`, `internal/rest/handler/calculator/add.go`, `internal/store/pg/pool.go`, `internal/store/pg/calculator_tx.go`, `resources/openapi/openapi.yml`, `resources/openapi/schema/common/error.yaml`, and `resources/pgsql/migrations/001_create_calculator_transactions.up.sql`.

### Template mapping

| Planned file | Template counterpart |
|---|---|
| `cmd/http/main.go` | `cmd/http/main.go` |
| `internal/config/*` | `internal/config/*` |
| `internal/config/worker.go` | `internal/config/http.go` (same style) |
| `internal/migrate/migrate.go` | same path |
| `internal/openapi/*` | same path |
| `internal/rest/server/server.go` | same path |
| `internal/rest/response/envelope.go` | same path |
| `internal/rest/middleware/partnerauth.go` | `internal/rest/middleware/apikey.go` |
| `internal/rest/middleware/logger.go`, `recover.go` | same paths |
| `internal/rest/handler/healthcheck/*` | same path |
| `internal/rest/handler/loanapplication/*.go`, `masterdata/*.go` | `internal/rest/handler/calculator/add.go` |
| `internal/store/pg/pool.go` | same path |
| `internal/store/pg/<entity>.go` | `internal/store/pg/calculator_tx.go` |
| `resources/openapi/schema/<domain>/*.yaml` | `resources/openapi/schema/calculator/add.yaml` |
| `resources/pgsql/migrations/001_*` | `resources/pgsql/migrations/001_*` |
| `Makefile`, `.env.example`, `go.mod` | same paths |
| `internal/validation/*`, `internal/worker/*` | none (see plan Complexity Tracking) |

### Deviations from the template

| Deviation | Required by |
|---|---|
| Validation runs in a worker after recording, not in the handler | Principle III |
| Per-partner hashed keys replace the static `API_KEY` | Principle IV |
| `202 Accepted` for submission instead of `200` | Principle III, HTTP meaning |
| Submission data is hard-deleted rather than soft-deleted | FR-011 |
| Request logging must never add body fields; the template logger is already safe and is kept as is | FR-011 |
| New `internal/worker` and `internal/validation` packages | Principle III |

### Template and constitution disagreements

These don't affect the design, but the constitution should be corrected in a PATCH amendment.

| Constitution or README says | Template code does | This plan follows |
|---|---|---|
| Makefile target `serve` | Target is `run`; there is no `serve` | `make run` |
| `response.WriteJSON` / `response.WriteError` | `WriteJSONWithEnvelope` / `WriteErrorWithEnvelope` | Template code |
| README errors look like `{"code":400,"message":...}` | `{"error-code":...,"error-message":...}` | Template code |
| `POSTGRES_CONN_MAX_IDLE` in `.env.example` | Not passed to `pg.Config` by `main.go` | Leave out of the new `.env.example` |

## Technology References (Principle VII)

| Technology | Version (source) | Context7 library ID | Topics checked |
|---|---|---|---|
| pgx | 5.10.0 (`go.mod`) | `/jackc/pgx/v5.10.0` | `Pool.Begin`, `pgx.ErrNoRows`, bytea and JSON type mapping |
| PostgreSQL | 16 | `/websites/postgresql/16` | `SKIP LOCKED`, row locks, `VACUUM` |
| chi | 5.3.1 (`go.mod`) | `/go-chi/docs` (unversioned) | URL parameters |
| golang-migrate | 4.19.1 (`go.mod`) | `/golang-migrate/migrate/v4_18_3` (nearest available) | `iofs` source; template's `migrate.go` is the reference for the pgx5 driver |
| zerolog | 1.35.1 (`go.mod`) | `/rs/zerolog` (unversioned) | Sub-loggers with context fields |
| caarlos0/env | 11.4.1 (`go.mod`) | `/caarlos0/env` (unversioned) | `envDefault`, `required` behaviour |
| google/uuid | 1.6.0 (`go.mod`) | `/google/uuid` (unversioned) | `NewV7` |
| Go standard library | 1.26 (`go.mod`) | `/websites/pkg_go_dev_go1_25_3` (nearest full library) | `http.MaxBytesReader`, `image`, `time/tzdata`, `crypto/hmac` |

No web search fallback was used. Where Context7 has no exact version, re-check the specific API during tasks and implementation.

## Decisions

### R-001: Background processing inside the service

- **Decision**: Goroutine loops in the same binary, polling PostgreSQL with `SELECT ... FOR UPDATE SKIP LOCKED`.
- **Rationale**: PostgreSQL documents `SKIP LOCKED` for queue-like tables with multiple consumers. The row is locked for the length of one transaction, so a crash rolls back and the row is picked up again without any lease logic. The submission and its queue entry are the same row, so they can't disagree.
- **Alternatives considered**: A message broker adds infrastructure the template lacks. `LISTEN/NOTIFY` still needs polling to recover missed signals.

### R-002: What the request path still checks

- **Decision**: Only authentication, an active agreement, the 16 MiB body limit, and "the body is a JSON object" run in the request. Everything else, including a missing `customer` object or a field with the wrong JSON type, is recorded and rejected by the worker.
- **Rationale**: Principle III prohibits validations in the request, but a body that can't be read can't be recorded. `400`, `401`, `403` and `413` are protocol failures with no application created. The OpenAPI request schema therefore marks nothing as `required` and documents mandatory fields in descriptions, so the contract doesn't promise a synchronous check.
- **Alternatives considered**: A typed Go request struct would fail on wrong types inside the handler and turn field problems into `400`s, which contradicts FR-012.

### R-003: Store the raw body as `BYTEA`

- **Decision**: Keep the exact request body in `application_payload.body`, with photos as base64 strings inside the JSON. On pass, photos move to `stored_photo`; on rejection, the payload row is deleted.
- **Rationale**: Nothing queries inside the payload, and `BYTEA` avoids parsing up to 16 MiB as `JSONB` while the partner waits. One row delete removes all personal data.
- **Alternatives considered**: Multipart upload is more friction than a single JSON document (O2). Object storage is new infrastructure and is only worth adding when storage volume requires it. At most about 10 MB per passed application, 30,000 applications is under 300 GB in year one, which is manageable.

### R-004: Photo checks

- **Decision**: Accept JPEG and PNG, up to 5 MiB decoded. Run `image.DecodeConfig` first to reject images over 50 megapixels, then fully decode.
- **Rationale**: `DecodeConfig` reads only the header, so on its own it would pass a truncated file. The size and pixel limits stop a small file expanding into a very large decode.
- **Alternatives considered**: HEIC and WebP need libraries the template doesn't use. These limits are not in the spec and should be confirmed during sandbox testing.

### R-005: Webhook delivery

- **Decision**: Notification rows are written in the same transaction as the rejection. The notifier POSTs JSON to the partner's `notification_url`, signed with HMAC-SHA256, with a 5-second timeout. Retries happen after 1 minute, 5 minutes, 30 minutes, 2 hours, 6 hours and 12 hours, about 21 hours in total, then the row is marked `EXHAUSTED`.
- **Rationale**: This is at-least-once delivery, so payloads carry a stable `event-id` for de-duplication. After `EXHAUSTED`, the FR-014 lookup is the recovery route. Only rejections are sent (FR-013), and payloads contain no submitted values.
- **Alternatives considered**: Sending inside the worker's validation transaction would tie validation speed to partner availability.

### R-006: Submission date and time zone

- **Decision**: The submission date comes from `submitted_at`, converted using the required `BUSINESS_TIMEZONE` setting. The binary imports `time/tzdata`.
- **Rationale**: The country is unconfirmed, so a default could silently produce wrong age and expiry results around midnight. A required setting makes startup fail until someone decides. `time/tzdata` keeps minimal containers working.

### R-007: Identifiers

- **Decision**: `loan_application_id` is a UUIDv7 and `application_reference` is `LA-` followed by that UUID.
- **Rationale**: UUIDv7 keeps primary key inserts roughly ordered. The reference reveals nothing beyond the time, and lookups are scoped to the owning partner, so another partner's reference returns `404`.

### R-008: Partner authentication

- **Decision**: Look up the SHA-256 hash of the `api-key` header in `partner_credential`. A missing or unknown key gets `401`; a valid key with an inactive agreement gets `403`.
- **Rationale**: Keys are high-entropy, so a fast hash fits, and storing only hashes means a database leak doesn't expose keys. Issuing, revoking and rotating keys stay in US-001. There is no rate limiting because AC-001-004 forbids a volume limit.

### R-009: Observability

- **Decision**: Use zerolog with `application_reference`, `partner_id`, outcome, reason codes, and validation delay. No metrics library.
- **Rationale**: SC-001 and SC-006 can be measured from `submitted_at`, `acknowledged_at`, and `validated_at`. Adding Prometheus would be a new dependency this story doesn't need.

## Risks and Open Items

| Item | Impact | Owner |
|---|---|---|
| A retried `POST` after a lost response creates a duplicate application; no idempotency key is in scope | Duplicate customer data checks and possibly duplicate rejections | Decide with US-007 |
| Deleted rows stay in PostgreSQL dead tuples until `VACUUM`, and in backups until they expire | "Deleted" is not instantly physical; matters for compliance | Business and compliance, once the country is known |
| `partner.notification_secret` must be readable to sign webhooks, so it can't be hashed | A database leak exposes signing secrets | US-001 |
| Retention of passed customers' data is undefined | Compliance gap | Business and compliance |
| Customer rows are keyed by national ID (DDE-001), so a later application updates shared data, even from another partner | Earlier values are overwritten | Review when loan decisions are designed |
| Collateral and loan sections are stored but not validated | "Passed" means customer data only until US-004 and US-005 are built | Those stories |
| Principle III also applies to US-003 to US-007 checks | Their stories and sequence diagrams still show synchronous rejection | Workspace owner |
