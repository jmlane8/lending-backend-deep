---
name: use-case
description: Create a use case document (UC-NNN) from user-supplied input. Use when the user explicitly asks to write or generate a use case and has provided actual process flow or business objective input.
disable-model-invocation: true
---

# Skill: Use Case

## ID Format
`UC-[3 digits]` — e.g., `UC-001`, `UC-002`

## Sequencing
- Start from `UC-001` if no use case exists yet.
- Scan `use-cases/` and continue from the highest existing ID (ignore `UC-999-template.md` and `UC-999-example.md`).
- Never reuse or change a use case ID.

## Output Path
`use-cases/UC-NNN-descriptive-name.md` — one file per use case.  
Example: `use-cases/UC-001-user-login.md`

## Template & Example
- Template: `.claude/skills/use-case/template/UC-999-template.md`
- Example: `.claude/skills/use-case/template/UC-999-example.md`

Follow its structure exactly. Do not invent sections that are not in the template.

## Process
- Generate the use case document based on the user's input, following the structure and format of the template.  
- Ensure the content is clear, concise, and accurately reflects the user's requirements.  
- Keep the use case focused on a single user goal. Split into multiple use cases if more than one goal is present.  
- If you meet any ambiguity or conflicting information in the user's input, inform the user.  