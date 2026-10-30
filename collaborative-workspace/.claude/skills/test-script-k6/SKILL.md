---
name: test-script-k6
description: Convert an existing Test Scenario document (test-scenarios/TS-NNN-*.md) into runnable Grafana k6 test scripts, grounded in the API's real OpenAPI contract. Use when the user explicitly asks to generate k6 scripts, k6 tests, "Grafana k6", or performance/load/functional test scripts from a test scenario, and references a file in the test-scenarios/ folder as the input. Requires a second input — the API's OpenAPI document (URL, JSON, or YAML) — which is fetched, validated as real OpenAPI 3.x, and reconciled against the scenario before generating. Produces one runnable .js per test case under test-scripts-k6/<scenario-name>/, each named k6-TC-NNN-NNN-<short-kebab-title>.js and tagged for traceability, plus a shared minimal config.js (connection settings) and a README for the folder.
disable-model-invocation: true
---

# Skill: Test Script (Grafana k6)

## Purpose
Turn an approved Test Scenario (TS-NNN) into runnable Grafana k6 scripts, grounded in the API's real OpenAPI 3.x contract. The scenario states *what* to verify — test data, steps, one expected result per case; this skill renders that into k6 without re-deciding coverage. The value is fidelity and traceability: every check traces to a `TC-NNN-NNN`, and nothing — endpoint, payload field, or status code — is invented. A script that quietly fabricates a value gives false confidence, exactly like a test scenario that only proves the happy path.

## Input (strict)
Two required, never generate with only one
  
1. **Test Scenario** — one or more files from `test-scenarios/` (e.g. `TS-001-…md`). Not a use case, user story, feature description, or ad-hoc prompt. If none is named, ask; if the named file is absent, stop.
2. **OpenAPI 3.x document** — a URL, `.json`, or `.yaml`/`.yml` the user points to (need not be in the working folder). It is the schema authority. If none is given, ask. See *Ground the schema*.

## Ground the schema (OpenAPI 3.x)
Full how-to: `references/openapi-grounding.md`. The rules:
- **Fetch & validate.** Fetch/read and parse (JSON or YAML); confirm `openapi: 3.x` + a `paths` object. If you can't fetch/parse it or that signature is missing, **STOP** — never infer the schema from the scenario's prose. Swagger 2.0 is out of scope (stop, say so).
- **Use as the authority.** From the contract, not memory: `servers`→`BASE_URL`; `paths`+method→endpoint & `http.<method>`; `requestBody` schema (`$ref`→`components.schemas`)→payload field names/types/required; `securitySchemes`→the `headers()` auth; `responses`→the status codes you may assert.
- **Reconcile, then ask.** Cross-check the scenario's path/method/fields/expected-status against the contract. If they align, generate. On any mismatch (a path/field/method/status one side has and the other lacks, or a contract for a different API): inform the user of the **specific** mismatches, warn that force-generating may be biased/unpredictable, and **ask** force vs stop. Always ask — never assume.

## Output
**Folder** `test-scripts-k6/<scenario-filename>/`, containing:
- **One `.js` per test case**, named `k6-<TC-id>-<short-kebab-title>.js` (lowercase kebab from the case title), e.g. `k6-TC-001-002-non-contracted-party-submission-is-rejected.js`. A multi-step case keeps all its calls in its own file — the boundary is the case, not the endpoint.
- **One shared `config.js`** (uniform name) — minimal, **connection settings only**: `BASE_URL` (servers), shared endpoint path(s) (paths), credentials, a `headers()` helper (auth per `securitySchemes`), each via `__ENV.NAME || 'default'`. Per-case values — `vus`/`iterations`, `thresholds`, payload, checks — must **not** go here.
- **One `README.md`** telling the user to set `BASE_URL` and the key(s) in `config.js` (or via env) before running. Keep it brief.

Each case file runs standalone (`k6 run k6-…js`) and may import **only** `./config.js` — never another case file.

## Generating each case
Start from the skeletons in `references/`: `config-skeleton.js` (shared config) and `k6-skeleton.js` (each case — fixed order: imports incl. `./config.js` → header → payload → custom metrics → `options` → `default` function). `references/k6-patterns.md` has the cheat sheet and tricky-case mappings. **k6 syntax:** if the Context7 MCP is available, confirm constructs against `/grafana/k6-docs` first (the k6 API evolves); else rely on `k6-patterns.md`. Map the scenario like so:

| Test Scenario element | k6 |
|---|---|
| Endpoint (OpenAPI `paths`, cross-checked vs the case's Steps) | the `http.<method>` target |
| HTTP-expressible precondition (e.g. valid API key) | a value from `config.js` + request header |
| Environment/state precondition (seed data, pre-existing record, service down) | a `// PRECONDITION (manual/external): …` comment — flag, never fake |
| Test Data | request body: field names/types from the OpenAPI schema, values from the case; clone the valid payload, mutate only what changes |
| Steps | the HTTP call(s) inside `group('TC-NNN-NNN — title', …)` |
| Timestamp / SLA step | no manual timer — a tagged `http_req_duration` threshold |
| Expected Result | `check(res, {…})` (status codes from OpenAPI `responses`), plus a threshold when it's about time |

Bake into every file:
- **Traceability** — tag each request and check `{ tc: 'TC-NNN-NNN', type: 'positive'|'negative' }`; keep the TC id in the file name and `group()` label.
- **Pass/fail gate** — `thresholds: { checks: ['rate==1.0'] }`, so the run exits non-zero on any failed check.
- **Isolation** — one `default` function, `options = { vus: 1, iterations: 1, … }`. A concurrency case (simultaneous requests) uses `vus: 2, iterations: 1` and asserts the cross-request rule with a `Counter` + `count==N` threshold, since a VU can't see another's response. Ordering between state-sensitive cases is the runner's job — run files in `TC` order.

## Gaps
The contract resolves most would-be gaps (wire formats, content types, valid status codes) — consult it first. A gap remains only where **both** the scenario and the OpenAPI are silent. There, generate the runnable skeleton and mark the unknown with a precise `// TODO(gap): …` quoting the source gap — never fabricate a value. Typical: a TBD numeric limit neither bounds; a wire format neither defines (use a named placeholder + TODO); an external precondition (script the request, assert the expected response, flag the manual setup).

## Process
1. Confirm both inputs exist (scenario in `test-scenarios/`; an OpenAPI source). Missing → stop/ask.
2. Fetch & validate the OpenAPI (`openapi: 3.x` + `paths`). Invalid → stop.
3. Read the whole scenario — Test Basis, Coverage Matrix, every case.
4. Reconcile scenario vs contract; on mismatch, inform + ask force/stop before proceeding.
5. Write `config.js` from the contract (BASE_URL, endpoint, `headers()`).
6. Encode the canonical valid payload — shape from the OpenAPI schema, values from the scenario's reference case.
7. Render one file per `TC-NNN-NNN` per the table above; set per-case `options` + thresholds; assert only documented status codes.
8. Mark gaps (`TODO(gap)` / manual-precondition) only where both sources are silent.
9. Write the folder: `config.js`, the case files (header comment naming source scenario + TC id + OpenAPI), and `README.md`.
10. Report: files written, TC→file map, the OpenAPI used (+ any forced mismatches), ordering-sensitive cases, all TODOs / manual preconditions.
