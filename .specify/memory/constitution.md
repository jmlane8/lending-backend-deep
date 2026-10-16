<!--
SYNC IMPACT REPORT
==================
Version change: (unratified template) → 1.0.0
Bump rationale: Initial ratification. All template placeholders replaced with project-specific
governance derived from docs/constitution-input-001.md.

Modified principles (placeholder → new title):
  - [PRINCIPLE_1_NAME] → I. Business-Objective Traceability
  - [PRINCIPLE_2_NAME] → II. Partner System as the Client
  - [PRINCIPLE_3_NAME] → III. Asynchronous Submission (NON-NEGOTIABLE)
  - [PRINCIPLE_4_NAME] → IV. API Key Security & Partner Traceability
  - [PRINCIPLE_5_NAME] → V. Partner Self-Sufficiency
  - (new)              → VI. Operational Scale Without Headcount Growth

Added sections:
  - Preamble (business context)
  - Scope Boundaries & Measurable Targets (was [SECTION_2_NAME])
  - Specification Workflow & Quality Gates (was [SECTION_3_NAME])

Removed sections: none

Follow-up TODOs: none. Ratification date set to the date of first adoption (2026-10-15).
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
- Every plan MUST include a Constitution Check confirming compliance with Principles I–VI;
  any deviation MUST be recorded with justification and an amendment proposal.
- Specs MUST express performance and onboarding criteria in the measurable terms of the
  targets table above.

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

**Version**: 1.0.0 | **Ratified**: 2026-10-15 | **Last Amended**: 2026-10-15
