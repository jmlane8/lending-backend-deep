# k6 Scripts - TS-001 Submit Loan Application with Customer Data

One script per test case in [TS-001](../../test-scenarios/TS-001-submit-loan-application-customer-data.md). The contract is the service's OpenAPI document at `http://localhost:8080/openapi.json`.

## Before you run

1. Start the service and seed the sandbox partners, as described in `DEVELOPMENT.md`. The service must use `BUSINESS_TIMEZONE=UTC`.
2. Set `BASE_URL` and `API_KEY` in `config.js`, or pass them as environment variables. The defaults are `http://localhost:8080` and `sandbox-partner-key`.
3. For `TC-001-003`, set `LARGE_PHOTO_FILE` to the absolute path of a text file holding a base64 image of 5 MiB decoded.

## Run

Run from this folder:

```bash
k6 run k6-TC-001-001-valid-application-is-acknowledged-within-2-seconds.js
k6 run -e LARGE_PHOTO_FILE=/path/to/photo-5mib.b64.txt k6-TC-001-003-largest-allowed-application-is-acknowledged-within-2-seconds.js
```

Cases don't depend on each other, so they can run in any order. A failed check makes `k6 run` exit with a non-zero code. Don't run across midnight UTC, because dates are worked out from the current date.

## Forced mismatches

These cases have no matching API operation. They were generated on request.

| Case | What the script does | What's still needed |
|---|---|---|
| `TC-001-030` | Submits and confirms the rejection, then logs the reference | Check by hand that `customer_id` is NULL and no `application_payload` row exists |
| `TC-001-045` | Fails on purpose | Run against `validation.Validate` with submission date `2026-02-28` |
| `TC-001-048` | Fails on purpose | Run against `validation.Validate` with submission date `2026-03-01` |
