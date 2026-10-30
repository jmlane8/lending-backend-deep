# Grounding k6 scripts in the OpenAPI contract

The test scenario says *what* to verify; the OpenAPI document is the authority for *how the API actually looks*. Read this before generating so the endpoints, payloads, auth, and status codes come from real contract — not from guessing.

## 1. Fetch & validate
Resolve the user's OpenAPI source:
- **URL** → fetch it (e.g. WebFetch). 
- **Local file** (anywhere on disk, not necessarily the working folder) → read it.
- Parse as **JSON** or **YAML** depending on the extension/content.

Then confirm it is genuinely **OpenAPI 3.x** — all of:
- the root has `openapi:` with a `3.` major version (e.g. `3.0.3`, `3.1.0`); and
- there is a non-empty `paths` object.

If you cannot fetch/read it, it doesn't parse, or that signature is missing → **STOP** and tell the user. Do not infer the schema from the test scenario's prose; the whole point of this input is to avoid guessing. Swagger 2.0 (`swagger: 2.0`) is out of scope: if you receive one, stop and tell the user this skill targets OpenAPI 3.x.

## 2. Extract the schema (what to read from where)
OpenAPI 3.x layout:

| Need | Where in the OpenAPI |
|---|---|
| `BASE_URL` | `servers[0].url` |
| Endpoint + method | a key under `paths` + the HTTP verb under it (`get`/`post`/…) |
| Operation | `paths[path][method]` — its `operationId`, `parameters`, `requestBody`, `responses` |
| Request body shape | `requestBody.content["application/json"].schema` — resolve `$ref` into `components.schemas` |
| Field names / types / required | the resolved schema's `properties` and `required` array (recurse into nested objects/`$ref`) |
| Content type | the media-type key under `content` (usually `application/json`) |
| Auth (how the key is sent) | `components.securitySchemes` + the operation/global `security` — e.g. `type: http, scheme: bearer` → `Authorization: Bearer <key>`; `type: apiKey, in: header, name: X-API-Key` → that header |
| Valid status codes | the keys of the operation's `responses` (e.g. `200`, `202`, `400`, `401`, `413`) |

Use these to build `config.js` (BASE_URL, endpoint, `headers()`), the valid payload's **shape**, and the status codes the checks assert. The scenario supplies the concrete *values*; the contract supplies the *structure* and the *legal responses*.

## 3. Reconcile against the test scenario
Before writing files, cross-check the scenario against the contract:

- **Endpoint/path** — does the path the scenario exercises exist in `paths`?
- **Method** — is the verb the scenario uses defined on that path?
- **Fields** — are the fields named in the scenario's Test Data present in the request schema? Are required schema fields all covered by the scenario's valid case?
- **Status codes** — are the scenario's expected results (e.g. 401, 413, 422) among the operation's documented `responses`?
- **Domain sanity** — is the contract obviously for the same API (not a different service entirely)?

If everything lines up, generate normally.

## 4. On mismatch — always ask, never assume
A mismatch is any of: a path/field/method/status the scenario uses that the contract doesn't have (or vice-versa for required fields), or a contract that is plainly for a different API. When you find one:

1. **Inform** the user of the **specific** mismatches — name the field/path/status, and which side has what. Example: *"TC-001-030 expects HTTP 413, but the OpenAPI for `POST /loan-applications` only documents 200/202/400/401. Field `selfie_photo` in the scenario is not in the request schema."*
2. **Explain the risk** — you can force-generate, but since the two sources disagree the output will have to favour one over the other, so it may be biased or unpredictable (e.g. asserting a status the API never returns, or sending a field it rejects).
3. **Ask** the user to choose: **force-generate** or **stop**. Do not pick for them, and do not silently reconcile by dropping or inventing things. Proceed only on their explicit choice; if they force, note in the report which mismatches were forced through.

## Example (abridged)
OpenAPI:
```yaml
servers: [{ url: https://api.test.lendingplatform.internal/v1 }]
paths:
  /loan-applications:
    post:
      security: [{ partnerToken: [] }]
      requestBody:
        content: { application/json: { schema: { $ref: '#/components/schemas/LoanApplication' } } }
      responses: { '202': {}, '400': {}, '401': {} }
components:
  securitySchemes:
    partnerToken: { type: http, scheme: bearer }
```
→ `config.js`: `BASE_URL` = that server, `ENDPOINT = ${BASE_URL}/loan-applications`, `headers(k) => { Authorization: 'Bearer ' + k }`.
→ valid payload shape = the `LoanApplication` schema's properties; concrete values from the scenario's reference case.
→ checks assert `202` for the happy path and `400`/`401` for the negatives — **not** a `413`, because the contract doesn't document it (that would be a mismatch to raise with the user).
