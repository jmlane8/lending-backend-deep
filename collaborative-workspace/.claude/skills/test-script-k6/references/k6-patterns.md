# k6 Patterns & Reference Structure

Idioms for generating the scripts, plus pointers to the skeletons. The k6 API matches current Grafana k6 (v1.x / v2.x). If the Context7 MCP is available, confirm any construct you're unsure about against `/grafana/k6-docs` before relying on it.

## Schema from the OpenAPI contract
Ground endpoints, request fields, auth, and the status codes you assert in the user's OpenAPI document — see [`openapi-grounding.md`](openapi-grounding.md). If it can't be validated, stop; if it disagrees with the scenario, ask the user before generating.

## Skeletons
- [`config-skeleton.js`](config-skeleton.js) → the folder's shared `config.js` (connection only; no per-case values).
- [`k6-skeleton.js`](k6-skeleton.js) → each case file (`TC-NNN-NNN`), section order: imports incl. `./config.js` → header → payload → custom metrics → `options` → `default` function.

Adapt placeholders; don't copy values verbatim. When a `config.js` is produced, also drop a short `README.md` in the folder.

## Construct cheat sheet

| You need to… | Use |
|---|---|
| Send a JSON POST | `http.post(url, JSON.stringify(body), { headers, tags: { tc: 'TC-001-001' } })` |
| Assert outcome (binary pass/fail) | `check(res, { 'TC-…: status is 200/202': (r) => r.status === 200 \|\| r.status === 202 }, { tc: 'TC-…' })` |
| Read a field safely from a JSON body | `(r) => { try { return r.json('application_id') !== undefined; } catch (_) { return false; } }` |
| Make the run fail if any check fails | `thresholds: { checks: ['rate==1.0'] }` |
| Enforce a time SLA on the case | tag the request `{ tc: 'TC-001-027' }`, then `thresholds: { 'http_req_duration{tc:TC-001-027}': ['p(100)<2000'] }` |
| Run the case once | `export const options = { vus: 1, iterations: 1, thresholds: {…} }` + an `export default function` |
| Produce real concurrency | `vus: 2, iterations: 1` (each VU runs the `default` function once → 2 simultaneous requests) |
| Assert a cross-request outcome | a dedicated `Counter`, incremented in the `default` function, with a `count==N` threshold |
| Build a negative payload | clone the valid payload, mutate one field — `delete body.customer.national_id_number`, or set an invalid value |
| Share connection settings | in `config.js`: `export const BASE_URL = __ENV.BASE_URL \|\| '…';` — case files `import { BASE_URL } from './config.js'` |

## Tricky cases
- **Time SLA** — tag the request with the TC id and bound its built-in `http_req_duration` via a tagged threshold; the `check` covers correctness, the threshold covers time.
- **Concurrency** — `options = { vus: 2, iterations: 1 }`; a VU sees only its own response, so assert the aggregate rule with a `Counter` incremented on success + a `count==1` threshold, not a per-VU check.
- **Gaps** — consult the contract first; a gap remains only where **both** the scenario and the OpenAPI are silent. Then generate the skeleton + a precise `// TODO(gap): …`; never fabricate.

## Faithfulness
- **Error-body checks are best-effort.** A substring check (`r.body.includes('national_id_number')`) is reasonable; tighten it (`r.json('errors')`) when the schema is known (e.g. the OpenAPI error-response schema), else add a `TODO(gap)`.
- **Don't seed data from the script.** Environment-state preconditions go in the file header as manual steps; the script asserts the *response*, it doesn't create the state.
- **One expected result per case**, mirroring the source. If a case needs several independent assertions, they should already be separate cases — don't merge them.
