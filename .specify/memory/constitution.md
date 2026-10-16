==================
Version change: 1.0.0 → 1.1.0
Bump rationale: MINOR — new principles added (VII, VIII) and quality gates materially
expanded. No existing principle removed or redefined. (1.1.0 was not yet committed, so both
additions ship in this single MINOR amendment.)

Modified principles: none renamed.
  - (new) → VII. Current Technical Documentation via Context7
  - (new) → VIII. Technical Template Conformance (jmlane8/rest-api-template)

Modified sections:
  - Specification Workflow & Quality Gates: Constitution Check now covers I–VIII; added
    Context7 documentation gate for /speckit-plan, /speckit-tasks, /speckit-implement,
    /speckit-converge; added template conformance gate (template commit SHA recorded in
    research.md, template read via GitHub MCP).

Added sections: none
Removed sections: none

Follow-up TODOs: none.
-->

# Lending Backend Constitution

This platform lets vetted external partners (digital-native companies or startups with
40,000+ users) submit car-collateral loan applications on behalf of their end consumers.
Partners earn a commission on approved loans; the company performs all financing. The goal is
scalable revenue growth through independent external distribution channels without expanding
marketing spend or operational headcount.

## Core Principles

### I. Business-Objective Traceability

- `docs/constitution-input-001.md` is the authoritative source of truth for business problems,
  objectives (O1–O4), and scope boundaries during all functional specification phases
  (`/speckit-specify`, `/speckit-clarify`, `/speckit-checklist`).
- Every functional requirement and success criterion in a spec MUST trace to at least one
  business objective (O1–O4) or stated business problem. Requirements with no trace MUST be
  removed or escalated as a proposed amendment.
- Where a spec, plan, or stakeholder request conflicts with the source-of-truth document, the
  document prevails until this constitution is amended.

Rationale: Keeps feature work anchored to the revenue, onboarding, performance, and scale
outcomes the platform exists to deliver.

### II. Partner System as the Client

- The direct user of the platform is the external partner system (website, mobile app, or
  partner backend). End consumers MUST NOT call the platform directly.
- The platform identifies and authorizes partners, not individual end users. Specs MUST NOT
  introduce end-consumer accounts, logins, or identity management.
- Consumer-facing UIs and partner-facing dashboards are out of scope and MUST NOT appear in
  specs.

Rationale: The product is a machine-to-machine distribution channel; consumer-facing concerns
belong to partners.

### III. Asynchronous Submission (NON-NEGOTIABLE)

- The submission API MUST return an acknowledgment within 2 seconds (O3).
- All validations — including vehicle identification via the internal loan system — MUST run
  after the submission is accepted, never inside the request-response cycle.
- Credit bureau checks and underwriting (including manual steps) MUST be processed
  asynchronously. Synchronous credit checks or underwriting are prohibited.
- Loan status and decisions MUST be delivered to the partner separately from the submission
  acknowledgment.
- Unavailability or failure of downstream dependencies (internal loan system, credit bureau)
  MUST NOT cause submission API errors; failed validations MUST be retried internally.

Rationale: Downstream response times are variable or unpredictable; decoupling them protects
the 2-second objective and keeps partners able to submit during outages.

### IV. API Key Security & Partner Traceability

- Authentication MUST be API key-based machine-to-machine authentication.
- Every API key MUST be traceable to exactly one partner, and every request MUST be
  attributable to that partner.
- Keys MUST be revocable immediately when compromised.
- The initial credential is issued by the company only after the partner's legal agreement is
  signed; subsequent key rotations MUST be self-service for the partner.
- OAuth, session-based flows, and a dedicated authentication service are out of scope.

Rationale: A static, pre-vetted partner set needs simple, auditable, revocable credentials,
not a general-purpose identity system.

### V. Partner Self-Sufficiency

- A new partner MUST be able to read documentation, test in a sandbox, and go live within
  5 business days with minimal assistance from the company (O2).
- The public interface MUST be a REST API (chosen over gRPC and GraphQL to minimize
  integration friction).
- Every partner-facing capability MUST ship with developer documentation and sandbox support.
- Accepted vehicle brands and valid enumerations (e.g., gender) MUST be exposed through master
  data endpoints so partners submit consistent values.

Rationale: Self-service onboarding is the lever that lets partner count grow without support
cost growing with it.

### VI. Operational Scale Without Headcount Growth

- The platform MUST support at least 20 active partner integrations within 12 months of launch
  without increasing operational team headcount (O4).
- Any spec that introduces a recurring manual operational task (other than underwriting steps
  and initial credential issuance) MUST justify it against O4 or replace it with a
  self-service or automated alternative.

Rationale: Growth must come from partner channels, not proportional operational effort.

### VII. Current Technical Documentation via Context7

- For every coding or technical-framework task in any Spec Kit prompt (technology research
  in `/speckit-plan`, design artifacts, `/speckit-tasks`, `/speckit-implement`,
  `/speckit-converge`, and technical debugging), the agent MUST fetch current documentation for each library, framework, SDK,
  API, CLI tool, or cloud service involved through the Context7 MCP server
  (`resolve-library-id` → `query-docs`) before relying on it.
- This applies even to well-known technologies; training-data knowledge alone MUST NOT be the
  basis for API syntax, configuration, version choices, or setup steps.
- When a version is pinned or named, the version-specific Context7 library ID MUST be used.
- Context7 MUST be preferred over general web search for library documentation. Web search is
  permitted only when Context7 has no matching library or returns insufficient results, and
  that fallback MUST be noted in the plan's research notes.
- Context7 is not required for non-technical work: business-level specification and
  clarification, refactoring, writing scripts from scratch, business-logic debugging, code
  review, or general programming concepts.

Rationale: Framework APIs change faster than model training data; grounding technical
decisions in current docs prevents deprecated or incorrect usage from entering plans and code.

### VIII. Technical Template Conformance

- The technical template `https://github.com/jmlane8/rest-api-template` is the authoritative
  reference for code structure, conventions, and technical stack in all technical phases
  (`/speckit-plan`, `/speckit-tasks`, `/speckit-implement`, `/speckit-converge`).
- Before planning or implementing, the agent MUST read the template through the GitHub MCP
  server (`get_file_contents` on `jmlane8/rest-api-template`) — at minimum `README.md`,
  `go.mod`, `Makefile`, and the template files analogous to the code being written. Local
  clones, copies, or recollection MUST NOT substitute for reading the current template.
- Implementation MUST follow the template's:
  - **Layout**: `cmd/http/` entry point; `internal/config`, `internal/migrate`,
    `internal/rest/{server,response,middleware,handler/<domain>}`, `internal/openapi`,
    `internal/store/pg`; `resources/openapi`, `resources/gen`, `resources/pgsql/migrations`.
  - **Spec-first API**: the hand-authored OpenAPI 3.0 YAML in `resources/openapi/` is the
    API source of truth; one schema file per endpoint under `schema/<domain>/`; the bundled
    `resources/gen/openapi/openapi.json` (via `make openapi`) is committed and embedded.
    Handlers are written to match the spec — never the reverse.
  - **Handler conventions**: one file per endpoint; hand-written request/response structs;
    responses through the shared `response.WriteJSON` / `response.WriteError` envelope;
    routes registered in `internal/rest/server`.
  - **Stack**: the Go version and libraries pinned in the template's `go.mod` (chi, pgx v5,
    zerolog, golang-migrate with embedded migrations applied at startup, caarlos0/env,
    google/uuid). Adding a library the template does not use, or replacing one it does,
    MUST be justified in the plan's Complexity Tracking.
  - **Tooling**: Makefile targets (`env`, `serve`, `openapi`, `test`, `vendor`) and
    `.env.example`-based configuration.
- Business principles I–VI prevail over template behavior. Where the template's sample code
  conflicts with them, follow the template's structure but not the conflicting behavior —
  notably: the template persists synchronously within the request (Principle III requires
  asynchronous processing after acknowledgment), and the template validates one static
  `API_KEY` (Principle IV requires per-partner, revocable, rotatable keys). Each such
  deviation MUST be recorded in the plan.
- Version choices for Principle VII lookups come from the template's `go.mod`; Context7 is
  then queried for those versions.

Rationale: A shared template keeps services consistent, reviewable, and quick to build;
reading it live through GitHub MCP ensures the current conventions are followed.

## Scope Boundaries & Measurable Targets

**Measurable targets** (from the source-of-truth document):

| ID | Objective | Target |
|----|-----------|--------|
| O1 | Revenue growth | 30,000 new loan applications via partner channels within 12 months of launch |
| O2 | Partner self-sufficiency | Docs → sandbox → live within 5 business days |
| O3 | API performance | Submission acknowledgment within 2 seconds |
| O4 | Operational scale | ≥ 20 active partner integrations in 12 months, no added ops headcount |

**In scope**: REST API; asynchronous loan processing with separate decision delivery; vehicle
identification via the existing internal loan system (10–200 ms); asynchronous credit bureau
checks; asynchronous underwriting and credit assessment; partner credential provisioning and
self-service rotation; developer documentation and sandbox; master data endpoints.

**Out of scope**: consumer-facing UI or partner-facing dashboard; dedicated authentication
service; synchronous credit checks or underwriting; OAuth or session-based authentication.

Adding any out-of-scope item to a spec requires a constitution amendment first.

## Specification Workflow & Quality Gates

- Before drafting or clarifying a functional spec, the author MUST read
  `docs/constitution-input-001.md` and cite the relevant objectives (O1–O4) in the spec.
- Clarification questions MUST NOT re-open decisions already fixed by this constitution
  (async model, REST, API keys, scope boundaries).
- Every plan MUST include a Constitution Check confirming compliance with Principles I–VIII;
  any deviation MUST be recorded with justification and an amendment proposal.
- Specs MUST express performance and onboarding criteria in the measurable terms of the
  targets table above.
- Plans (`research.md`) MUST list each technology decision with the Context7 library ID and
  version consulted (Principle VII); tasks, convergence, and implementation work that touches
  a library MUST re-check its Context7 docs for the specific API being used.
- Plans (`research.md`) MUST record the `jmlane8/rest-api-template` commit SHA read via
  GitHub MCP, map each planned package/file to its template counterpart, and list every
  deviation from the template with the principle that requires it (Principle VIII).
  Tasks MUST reference the template file each new file is modeled on.

## Governance

- This constitution supersedes all other project practices. `docs/constitution-input-001.md`
  is its referenced source of truth for business intent; a new revision of that input (e.g.,
  `constitution-input-002.md`) takes effect only through a constitution amendment that
  updates the references here.
- Amendments require: a written change description, review and approval by the project owner,
  an updated Sync Impact Report, and a migration note for any affected specs or plans.
- Versioning follows semantic versioning: MAJOR for removed or redefined principles or scope
  boundaries; MINOR for new principles/sections or materially expanded guidance; PATCH for
  clarifications and wording fixes.
- Compliance review: every spec, plan, and pull request MUST be checked against this
  constitution; `/speckit-analyze` findings that cite a constitutional violation are blocking.

**Version**: 1.1.0 | **Ratified**: 2026-10-15 | **Last Amended**: 2026-10-16
