---
name: technical-design
description: Create a Technical Design document (TD-NNN) — scope, API design, scheduled tasks, database design, and basic flow — derived from an existing use case or user story. Use when the user explicitly asks to create a technical design and references one or more files from the user-stories/ or use-cases/ folders as input.
disable-model-invocation: true
---

# Skill: Technical Design

## Purpose
Turn an approved functional requirement (use case or user story) into an implementation-level blueprint: API contracts, scheduled tasks, database changes, and the basic flow that ties them together. The design is what an engineer (human or AI agent) builds from, so every section must be concrete and traceable to its source requirement.

## Input (strict)
- The functional source MUST be one or more files from the `use-cases/` or `user-stories/` folders. This anchors the design to an approved requirement instead of an ad-hoc idea.
- Do NOT accept on-the-fly feature descriptions typed by the user as the functional source, or any source outside `use-cases/` or `user-stories/`.
- If the user does not specify a source file, ask which file(s) from `use-cases/` or `user-stories/` to use.
- If the referenced file does not exist in those folders, stop and inform the user. Do not improvise the design from conversation context.

## Data Definition Input (when database design is involved)
The `Database Design` section needs entity and field definitions. When the design touches the database, those definitions come from exactly one source, and **the user must choose**:
1. **From data definitions** — one or more files in the `data-definitions/` folder. Reuse the entities and fields as defined there; do not redefine them.
2. **From user input** — entities and fields the user provides directly in the conversation.

Never assume the source, even if a previous technical design in the session used one of them. If the user picks data definitions but names no file, ask which file(s) in `data-definitions/` to use; if a named file does not exist, stop and inform the user. If the design has no database change, skip this choice and omit the Database Design section.

## Clarifying Questions (mandatory, before generating)
1. **Functional source** — which file(s) from `use-cases/` or `user-stories/`.
2. **Scope of the design** — which parts apply: REST API, scheduled/background task, database change. Optional template sections are included only when their part is in scope.
3. **Data definition source** — only when database change is in scope (see above).

Confirm these before writing. Do not guess answers to fill gaps in the source.

## Template
- Template: `.claude/skills/technical-design/template/TD-999-template.md`

Follow its structure exactly. Do not invent sections that are not in the template. The document opens with a `Changelog` (start a new design at `0.0.1`). The `API Design`, `Scheduled Task Design`, and `Database Design` sections are conditional — include a section only when that part is in scope, and keep the inline `None` placeholders where the template instructs.

## ID / Prefix
Prefix: `TD`
Format: `TD-[3 digits]` — e.g., `TD-001`

## Sequencing
- Start from `TD-001` if no technical design exists yet.
- Scan `technical-designs/` and continue from the highest existing ID (ignore `TD-999-template.md`).
- Never reuse or change an ID.

## Output Path
`technical-designs/TD-NNN-descriptive-name.md`  
Example: `technical-designs/TD-001-loan-disbursement-api.md`

## Process
1. Read the selected functional source file(s) from `use-cases/` or `user-stories/`.
2. Confirm scope and, if the database is in scope, the data definition source (clarifying questions above). If the source is `data-definitions/`, read the chosen file(s) before continuing.
3. **Scope** — state what the design covers as narrative bullets, plus the explicit out-of-scope boundary. Name the source use case / user story ID here so the origin of the design is clear.
4. **API Design** (if REST is in scope) — one sub-section per endpoint. Define verb, path, authentication, produces/consumes, path/query parameters, request and response schemas with examples, and the non-2xx error registry. Use `lower-kebab-case` for JSON field names. Keep schema and example consistent within each section.
5. **Scheduled Task Design** (if a background task is in scope) — one sub-section per task. Define trigger, schedule, timezone, concurrency control, idempotency, inputs, processing steps, outputs/side effects, and the failure & retry registry.
6. **Database Design** (if database change is in scope) — build the table definitions from the chosen data definition source. Use `lower_snake_case` field names, mark PK/FK and uniqueness, and include the standard audit and soft-delete columns the template requires (`created_at`, `updated_at`, `deleted`, `deleted_at`, `created_by`, `updated_by`, `deleted_by`). Provide indices and PostgreSQL DDL.
7. **Basic Flow** — for every entry point (each API endpoint or scheduled task in scope), provide a Mermaid sequence diagram using concrete service names (never generic labels like "System" or "Backend") and the basic logic as pseudocode. Write the base SQL for any database operation. No hard deletes — soft-delete only, per the template.
8. **Traceability** — each entry point in Basic Flow must trace back to its API Design or Scheduled Task sub-section, and the whole design back to the source use case / user story.
