---
name: my-speckit-to-technical-design
description: Convert a completed Spec-Kit implementation plan (specs/[###-feature]/plan.md) into a Technical Design document (my-technical-design-from-plan.md) in the same feature folder. Use whenever the user asks to turn a Spec-Kit plan into a technical design, generate a technical design from plan.md, or convert /speckit-plan output into design documentation — even if they just say "create the technical design" after running /speckit-plan. Not part of Spec-Kit itself; this is a custom post-plan step.
disable-model-invocation: true
---

# Skill: Spec-Kit Plan → Technical Design

## Purpose
Spec-Kit's `/speckit-plan` produces an implementation plan (`plan.md`) plus supporting design artifacts. This skill converts that plan into a single, self-contained Technical Design document following the team's standard template, so engineers get the familiar TD format (API contracts, database design, basic flow) without re-deriving it by hand from scattered Spec-Kit files.

## Input (strict)
- The user MUST explicitly name which `specs/` subfolder to process (e.g., `specs/001-submit-loan-application`). Never pick a folder yourself — even if only one exists — because the user may have multiple features in flight and the wrong guess silently produces a design for the wrong feature. If the folder is not specified, ask for it and stop until answered.
- **Primary input**: `specs/[###-feature]/plan.md`. This is the authoritative source; the design must not contain anything that contradicts it.
- **Optional enrichment** — use these files when they exist in the same folder to make the corresponding template sections concrete; skip silently when absent:

  | File | Enriches |
  |------|----------|
  | `research.md` | Scope, technology rationale (Phase 0 output) |
  | `data-model.md` | Database Design (Phase 1 output) |
  | `quickstart.md` | Basic Flow, request/response examples (Phase 1 output) |
  | `contracts/` | API Design — endpoints, schemas, examples (Phase 1 output) |

## Precondition Check (before anything else)
Validate, in order:
1. `specs/[###-feature]/` exists and is a directory.
2. `specs/[###-feature]/plan.md` exists.

If either fails, the Spec-Kit plan step has not been completed for this feature. Show an error message stating exactly that — e.g., *"`specs/001-foo/plan.md` not found. Run `/speckit-plan` for this feature first."* — and stop immediately. Do not generate a partial design, do not fall back to `spec.md` or conversation context.

## Template
- Template: `.claude/skills/my-speckit-to-technical-design/template/TD-999-template.md`

Follow its structure exactly; do not invent sections that are not in the template. The `API Design`, `Scheduled Task Design`, and `Database Design` sections are conditional — include a section only when the plan actually covers that part, and keep the inline `None` placeholders where the template instructs.

## Section Mapping
Convert Spec-Kit content into template sections as follows:

| Template section | Built from |
|------------------|-----------|
| Changelog | New document — start at `0.0.1` with a single "Initial conversion from plan.md" entry |
| Source Requirement | The feature's Spec-Kit artifacts: link `spec.md` (and `plan.md`) with relative paths |
| Scope | plan.md `Summary`, `Technical Context`, and explicit out-of-scope statements; enrich rationale from `research.md` |
| API Design | Endpoints named in plan.md, with schemas/examples taken from `contracts/` OpenAPI files. JSON field names in `lower-kebab-case` |
| Scheduled Task Design | Only if plan.md defines a background/scheduled task; otherwise omit |
| Database Design | Entities/fields from `data-model.md` (do not redefine them); `lower_snake_case` fields, standard audit + soft-delete columns, indices, PostgreSQL DDL |
| Basic Flow | One entry point per API endpoint / scheduled task: Mermaid sequence diagram with concrete service names (never "System" or "Backend"), pseudocode with base SQL. Soft-delete only — no hard deletes |

Where plan.md and an enrichment file disagree, prefer the enrichment file for its own domain (e.g., `contracts/` wins for schemas, `data-model.md` wins for fields) — those are the Phase 1 refinements of the plan. Never fabricate details found in neither; write `None` or a `TODO` note instead, and flag it in the summary table.

## Output
- Path: `specs/[###-feature]/my-technical-design-from-plan.md` — the same folder as the plan it was derived from, so the design travels with its feature.
- If the file already exists, overwrite it without asking; the plan is the source of truth and regeneration is the expected way to refresh the design.

## Success Report
After writing the file, display a summary table in the terminal showing the conversion result of each template section, e.g.:

| Template Section | Source | Result |
|------------------|--------|--------|
| Scope | plan.md, research.md | Converted |
| API Design | contracts/submit-loan-application.yaml | Converted (1 endpoint) |
| Scheduled Task Design | — | Skipped (not in plan) |
| Database Design | data-model.md | Converted (5 tables) |
| Basic Flow | plan.md, quickstart.md | Converted (1 entry point) |

Mention the output file path and any sections that needed `TODO` placeholders so the user knows what to review.
Give me concise report. Sacrifice grammar for concision.

## Process
1. Confirm the target `specs/` subfolder from the user's request (ask if not given).
2. Run the Precondition Check; on failure, show the error and stop.
3. Read `plan.md`, then read whichever enrichment files exist (`research.md`, `data-model.md`, `quickstart.md`, `contracts/`).
4. Read the template and build each section per the Section Mapping.
5. Write `my-technical-design-from-plan.md` into the feature folder (overwrite if present).
6. Print the Success Report table.
