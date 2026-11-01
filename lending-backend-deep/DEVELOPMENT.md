# Development Guide

How to run the service on localhost, with PostgreSQL in Docker.

## Prerequisites

- Docker, with Docker running
- Go 1.26 or later
- `make`. On Windows, install it with `choco install make` or `scoop install make`
- `curl`. In Windows PowerShell, type `curl.exe`, because `curl` is an alias for something else

## Settings

All local settings are in `env.local`, including the database connection. The Makefile and Docker Compose both read it, so you don't need to export anything.

If port `5432` or `8080` is already in use on your machine, change `POSTGRES_PORT` or `HTTP_SERVER_PORT` in `env.local`, and use the new port in the commands below.

## Start the server

Run these from the `lending-backend-deep` folder.

1. Start PostgreSQL and wait until it's ready:

   ```bash
   make db-up
   ```

2. Start the server. It applies the database migrations, then keeps running in this terminal:

   ```bash
   make run
   ```

3. In a second terminal, load the test partners:

   ```bash
   make db-seed
   ```

   You only need to do this once for each database, and it's safe to repeat. It must run after the server's first start, because the server creates the tables.

The server is now at `http://localhost:8080`.

## Check that it's healthy

```bash
curl http://localhost:8080/healthcheck
```

Expected response:

```json
{"status":"ok"}
```

## Call the API

These examples use the test partner key `sandbox-partner-key`, created by `make db-seed`, and the sample application in `resources/sample/loan-application.json`.

| API key | Partner |
|---|---|
| `sandbox-partner-key` | Active |
| `sandbox-partner-2-key` | Active, for checking that partners can't see each other's applications |
| `inactive-partner-key` | Inactive agreement, returns `403` |

1. List the accepted customer values:

   ```bash
   curl -H "api-key: sandbox-partner-key" http://localhost:8080/master-data/v1/enumerations
   ```

2. Submit the sample application:

   ```bash
   curl -X POST http://localhost:8080/loan-applications/v1 -H "api-key: sandbox-partner-key" -H "Content-Type: application/json" -d "@resources/sample/loan-application.json"
   ```

   A `202` response means the application was recorded, not that it's valid:

   ```json
   {"data":{"application-reference":"LA-01a19539-3e93-780f-98c0-74245e6a4e7a","submitted-at":"2026-11-01T03:48:39Z"}}
   ```

3. Read the validation result, using the reference from your own response. Validation usually finishes within a couple of seconds:

   ```bash
   curl -H "api-key: sandbox-partner-key" http://localhost:8080/loan-applications/v1/LA-01a19539-3e93-780f-98c0-74245e6a4e7a/validation-result
   ```

   ```json
   {"data":{"application-reference":"LA-01a19539-3e93-780f-98c0-74245e6a4e7a","status":"PASSED","submitted-at":"2026-11-01T03:48:39Z","validated-at":"2026-11-01T03:48:39Z"}}
   ```

To see a rejection, edit the sample file, for example by removing `gender`, and submit it again. The result then shows `REJECTED` with the affected fields.

## API documentation

While the server is running with the `env.local` settings:

| Path | Content |
|---|---|
| `http://localhost:8080/redoc` | Interactive API documentation, including the rejection webhook. It loads Redoc from a CDN, so the browser needs internet access |
| `http://localhost:8080/openapi.json` | OpenAPI specification |

The source is `resources/openapi/openapi.yml`. After editing it, run `make openapi` to regenerate the bundled spec, which needs Node.js.

## Stop the server

Press `Ctrl+C` in the terminal running `make run`, then choose one:

```bash
make db-down    # stops PostgreSQL and keeps the data
make db-reset   # stops PostgreSQL and deletes the data
```

After `make db-reset`, repeat all the steps in [Start the server](#start-the-server).
