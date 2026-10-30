---
name: data-definition
description: Create a Data Definition / logical data model (DDE-NNN) in a markdown table, optionally with an ERD. Use when the user explicitly asks to create a data definition and references one or more files from the user-stories/ or use-cases/ folders as input.
disable-model-invocation: true
---

# Skill: Data Definition

## Purpose
Produce a logical data model that defines business entities, their relationships, and their fields, derived from functional requirements.

## Input (strict)
- Input MUST be one or more files from the `user-stories/` or `use-cases/` folders.
- Do NOT accept on-the-fly descriptions typed by the user, or any source outside `user-stories/` or `use-cases/`.
- If the user does not specify a source file, ask which file(s) from `user-stories/` or `use-cases/` to use.
- If the referenced file does not exist in those folders, stop and inform the user. Do not improvise entities from conversation context.

## Clarifying Question (mandatory, every prompt)
Before generating, ALWAYS ask the user to choose one of:
1. **Data model only** — logical data model in markdown tables.
2. **Data model + ERD** — logical data model plus an Entity Relationship Diagram.

Never assume the answer, even if a previous data definition in the session included an ERD.

## Template & Example
- Template: `.claude/skills/data-definition/template/DDE-999-template.md`
- Example: `.claude/skills/data-definition/template/DDE-999-example.md`

Follow its structure exactly. Do not invent sections that are not in the template.

## ID / Prefix
Prefix: `DDE`
Format: `DDE-[3 digits]` — e.g., `DDE-001`

## Sequencing
- Start from `DDE-001` if no data definition exists yet.
- Scan `data-definitions/` and continue from the highest existing ID (ignore `DDE-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`data-definitions/DDE-NNN-descriptive-name.md`
Example: `data-definitions/DDE-002-loan-application-data-model.md`

## Process

### Always (logical data model)
1. Read the selected source file(s) from `user-stories/` or `use-cases/`.
2. Identify business entities and the relationships between them.
3. Define entities & relationships in markdown tables.
4. Add fields per entity. Fields typically come from functional requirements (i.e., derived from the user stories or use cases), but both functional and technical team members may add fields afterwards — leave the tables easy to extend.
5. For each field, capture at minimum: field name, description, data type (logical, e.g., text/number/date/boolean), and required/optional. Mark identifying fields and references between entities.
6. Trace each entity back to its source (user story or use case ID) so the origin of the requirement is clear.
7. If you generate entity identifier field (such as `person_id` for table `people`), make sure to mark it as the primary key of the entity, use a consistent naming convention across entities (e.g., `<entity>_id`), and ensure the data type is UUID.

### If user chose Markdown + ERD (do in sequence)
1. Generate the data definition in markdown tables (as above).
2. Trigger the `entity-relationship-diagram` skill, using the generated data definition as its input.
3. Put the resulting ERD link as a **Related Artifact** in the data definition document.