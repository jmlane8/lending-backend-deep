---
name: user-story
description: Create a user story document (US-NNN) with acceptance criteria from user-supplied input. Use when the user explicitly asks to write or generate a user story and has provided actual requirements or process detail.
disable-model-invocation: true
---

# Skill: User Story

## ID Format
`US-[3 digits]` — e.g., `US-001`, `US-002`

## Sequencing
- Start from `US-001` if no user story exists yet.
- Scan `user-stories/` and continue from the highest existing ID (ignore `US-999-template.md` and `US-999-example.md`).
- Never reuse or change a user story ID.

## Output Path
`user-stories/US-NNN-descriptive-name.md` — one file per user story.

## Template & Example
- Template: `.claude/skills/user-story/template/US-999-template.md`
- Example: `.claude/skills/user-story/template/US-999-example.md`

Follow its structure exactly. Do not invent sections that are not in the template.

## Acceptance Criteria

Each user story must have one or more acceptance criteria.

### ID Format
`AC-[user story digits]-[3 digits]`  - e.g., for `US-007`: `AC-007-001`, `AC-007-002`, …

- The middle segment (`007`) is fixed and must match the parent user story.
- The last segment starts from `001` if none exist for that story, or continues from the highest existing.
- Never reuse or change an acceptance criteria ID.

## Process
- Generate the user story document based on the user's input, following the structure and format of the template. 
- Ensure the content is clear, concise, and accurately reflects the user's requirements.
- The role must be a specific actor, not a generic term like "user"; the benefit must state the business value, not restate the goal.
- Each story must satisfy the INVEST criteria: Independent, Negotiable, Valuable, Estimable, Small, and Testable.
- Avoid including technical implementation details in the story statement itself.
- Each acceptance criterion must use the Given/When/Then format and must be independently testable with a binary pass or fail outcome.
- A story should represent a single deliverable unit of value; split it if it covers more than one distinct capability.
- Include edge cases and negative scenarios as separate acceptance criteria, not as footnotes in the story.
- If you meet any ambiguity or conflicting information in the user's input, inform the user.  

