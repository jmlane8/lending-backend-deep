# Lending Backend

REST API that lets contracted partners submit car-collateral loan applications with customer data. Built from `jmlane8/rest-api-template`, using Go, chi, pgx, zerolog and golang-migrate, with an embedded OpenAPI 3.0 spec.

## How submission works

1. The partner sends `POST /loan-applications/v1`. The service records the exact body and returns `202` with an application reference. A `202` means recorded, not valid.
2. A background worker in the same process validates the customer data.
3. If every rule passes, the application becomes `PASSED`. The customer is stored and the payload is reduced to its `collateral` and `loan` sections.
4. If any rule fails, the application becomes `REJECTED` and the raw payload is hard-deleted. Only the reference, partner, dates and reasons remain, and a signed webhook is queued.

Partners can read the result at any time through the lookup. PostgreSQL is both the record store and the work queue (`FOR UPDATE SKIP LOCKED`), so several replicas can run safely.

## Prerequisites

| Tool | Purpose |
|---|---|
| Go 1.26+ | Build and run |
| PostgreSQL 16+ | Data and queues |
| make | Makefile targets (optional on Windows) |
| Node.js | `make openapi` only |

Migrations are embedded and applied at startup, so no `migrate` CLI is needed.

## Quick start

To run on localhost with PostgreSQL in Docker, follow [DEVELOPMENT.md](DEVELOPMENT.md). The steps below are for connecting to your own PostgreSQL.

```bash
make env                       # copies .env.example to .env; edit it to match your database
export $(grep -v '^#' .env | xargs)
go run ./cmd/http    # make run would use env.local instead
```

In Windows PowerShell, load `.env` and start the service without `make`:

```powershell
Get-Content .env |
  Where-Object { $_ -notmatch '^\s*#' -and $_ -match '=' } |
  ForEach-Object {
    $name, $value = $_ -split '=', 2
    Set-Item "env:$($name.Trim())" $value.Trim()
  }
go run ./cmd/http
```

The server listens on `http://localhost:8080`. Set `OPENAPI_UI_ENABLE=true` to serve the docs at `/redoc` and the spec at `/openapi.json`.

## Sandbox use

`resources/pgsql/seed/sandbox_partner.sql` creates test partners for local and sandbox environments. Never run it in production.

```bash
psql "$DATABASE_URL" -f resources/pgsql/seed/sandbox_partner.sql
```

| API key | Partner |
|---|---|
| `sandbox-partner-key` | Active |
| `sandbox-partner-2-key` | Active, for checking lookup isolation |
| `inactive-partner-key` | Inactive agreement, returns `403` |

The seed registers no notification address, so these partners get no webhooks until one is set on the `partner` row.

## Endpoints

All routes except `/healthcheck` require the `api-key` header.

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthcheck` | Service health probe |
| `POST` | `/loan-applications/v1` | Submit an application |
| `GET` | `/loan-applications/v1/{application-reference}/validation-result` | Read the validation result. Another partner's reference returns `404` |
| `GET` | `/master-data/v1/enumerations` | List accepted gender, marital status and occupation values |

Marital status, occupation and gender must be sent as the exact API values, such as `SINGLE`, `CIVIL_SERVANT` and `FEMALE`.

`resources/openapi/openapi.yml` is the source of truth for the API. After editing it, run `make openapi` and commit `resources/gen/openapi/openapi.json`.

## Webhook verification

Rejections are sent to the partner's registered address as a signed `loan-application.rejected` event. Only rejections are sent.

1. Compute the lowercase hex HMAC-SHA256 of `<x-timestamp>.<raw body>`, using the partner's notification secret.
2. Compare `sha256=<hex>` with the `x-signature` header, using a constant-time comparison.
3. Reject an old `x-timestamp`. Five minutes is a sensible limit.
4. De-duplicate on `x-event-id`, because delivery is at-least-once.
5. Return any `2xx` within 5 seconds.

Any other result is a failed attempt, including a redirect. Retries run 1 minute, 5 minutes, 30 minutes, 2 hours, 6 hours and 12 hours after the previous attempt, for 7 attempts in total. After that, the result stays available through the lookup. Outside `APP_ENV=local`, an address that isn't `https` counts as a failed attempt and nothing is sent.

## Environment variables

| Variable | Default | Purpose |
|---|---|---|
| `APP_NAME` | `lending-backend-deep` | Name in structured logs |
| `APP_ENV` | `local` | `local`, `dev`, `uat`, `sit`, `sandbox` or `prod` |
| `LOGGER_LEVEL` | `info` | `trace`, `debug`, `info`, `warn`, `error` or `disabled` |
| `LOGGER_OUTPUT` | `stdout` | `console` for pretty logs, `stdout` for JSON |
| `OPENAPI_UI_ENABLE` | `false` | Serve `/openapi.json` and the Redoc UI |
| `OPENAPI_UI_PATH` | `/redoc` | Redoc UI path |
| `HTTP_SERVER_HOST` | `0.0.0.0` | Bind address |
| `HTTP_SERVER_PORT` | `8080` | Bind port |
| `HTTP_SERVER_READ_TIMEOUT` | `30` | Seconds |
| `HTTP_SERVER_WRITE_TIMEOUT` | `30` | Seconds |
| `HTTP_SERVER_SHUTDOWN_TIMEOUT` | `5` | Seconds |
| `POSTGRES_HOST` | `localhost` | Database host |
| `POSTGRES_PORT` | `5432` | Database port |
| `POSTGRES_DATABASE` | `lending` | Database name |
| `POSTGRES_SCHEMA` | `public` | Schema |
| `POSTGRES_USERNAME` | `postgres` | Database user |
| `POSTGRES_PASSWORD` | None | Database password |
| `POSTGRES_SSL_MODE` | `disable` | `sslmode` value |
| `POSTGRES_CONN_MAX_OPEN` | `25` | Maximum open connections |
| `POSTGRES_CONN_MIN_IDLE` | `1` | Minimum idle connections |
| `BUSINESS_TIMEZONE` | None, **required** | IANA time zone that fixes the submission date. The service won't start without it |
| `VALIDATOR_POLL_INTERVAL_SECONDS` | `1` | Wait when no application is due |
| `NOTIFIER_POLL_INTERVAL_SECONDS` | `5` | Wait when no notification is due |

## Tests

```bash
make test                                    # unit tests
TEST_DATABASE_URL=postgres://... make test   # adds integration tests
```

`make test` runs the race detector, which needs a C compiler. Integration tests truncate the partner and application tables, so point `TEST_DATABASE_URL` at a dedicated database.
