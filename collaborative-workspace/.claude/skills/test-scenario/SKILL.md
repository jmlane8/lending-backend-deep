---
name: test-scenario
description: Create a Test Scenario document (TS-NNN) — test cases covering positive (happy) and negative (not-happy) paths, with a coverage matrix — derived from an existing use case or user story. Use when the user explicitly asks to create a test scenario, test cases, or test coverage and references one or more files from the use-cases/ or user-stories/ folders as input.
disable-model-invocation: true
---

# Skill: Test Scenario

## Purpose
Turn an approved functional requirement (use case or user story) into a set of executable test cases that a tester — human or automated — can run to verify the feature behaves as specified. Every test case must be concrete and traceable to its source requirement, and the document as a whole must exercise both the happy path and the not-happy path. A test scenario that only proves the success case gives false confidence; the negative cases are where defects hide.

## Input (strict)
- The source MUST be one or more files from the `use-cases/` or `user-stories/` folders. This anchors the tests to an approved requirement instead of an ad-hoc idea.
- Do NOT accept on-the-fly feature descriptions typed by the user as the source, or any source outside `use-cases/` or `user-stories/`.
- If the user does not specify a source file, ask which file(s) from `use-cases/` or `user-stories/` to use.
- If the referenced file does not exist in those folders, stop and inform the user. Do not improvise tests from conversation context.

## Clarifying Questions (mandatory, before generating)
1. **Source requirement** — which file(s) from `use-cases/` or `user-stories/`.
2. **Scope** — full coverage of the source, or a specific subset of flows / acceptance criteria.
3. **Test data constraints** — any concrete limits the source does not state (e.g., exact field min/max, valid formats). If the source lacks a detail needed to write a meaningful test, flag it rather than inventing a value.

Confirm these before writing. Do not guess answers to fill gaps in the source.

## Template
- Template: `.claude/skills/test-scenario/template/TS-999-template.md`

Follow its structure exactly. Do not invent sections that are not in the template. The document opens with a `Changelog` (start a new document at `0.0.1`), followed by the mandatory `Source Requirement` table naming the use case / user story it derives from.

## ID / Prefix
Prefix: `TS`
Format: `TS-[3 digits]` — e.g., `TS-001`

Test case IDs are nested under the document: `TC-[TS digits]-[3 digits]` — e.g., for `TS-007`: `TC-007-001`, `TC-007-002`, …
- The middle segment is fixed and must match the parent document.
- The last segment starts at `001` and increments. Never reuse or change a test case ID.

## Sequencing
- Start from `TS-001` if no test scenario exists yet.
- Scan `test-scenarios/` and continue from the highest existing ID (ignore `TS-999-template.md`).
- Never reuse or change an ID.

## Output Path
`test-scenarios/TS-NNN-descriptive-name.md`
Example: `test-scenarios/TS-001-loan-disbursement.md`

## Process
1. **Read the source.** Read the selected file(s) from `use-cases/` or `user-stories/`.
2. **Confirm scope** and any missing test-data constraints (clarifying questions above).
3. **Source Requirement** — record the source use case / user story ID(s) in the table so the origin of the tests is clear and traceable back to the requirement.
4. **Enumerate the source elements to cover.** This is the backbone of full coverage:
   - From a **use case**: each Normal Flow (happy path), each Alternative Flow (valid variation that still succeeds), and each Exception (failure the system must handle).
   - From a **user story**: each acceptance criterion — both the success ones and the negative/edge ones.
5. **Derive test cases for full-path coverage.** Generate at least one **positive** test case for every Normal/Alternative Flow and success criterion, and at least one **negative** test case for every Exception and negative criterion. The finished document MUST contain both positive and negative cases — never only the happy path. This is the core value of the document.
6. **Make each test case concrete and atomic.** Replace abstract phrases like "a valid amount" with concrete, representative test data — that concreteness is the main thing a test scenario adds over the acceptance criteria it derives from. Give each test case exactly one Expected Result so its outcome is an unambiguous binary pass/fail; if you find yourself listing several independent checks, split them into separate test cases.
7. **Apply boundary and equivalence analysis to negatives.** For any rule with a limit, range, or format, write an in-range case that passes and an out-of-range / invalid-format case that fails — bugs cluster at boundaries. Where the source implies them, cover the standard negative techniques: missing mandatory field, invalid format, unauthorized or expired session, duplicate / idempotency, and system-unavailable.
8. **Keep test cases independent and repeatable.** Each case states its own preconditions and test data and does not rely on leftover state from a previous case, so they can run in any order or in isolation.
9. **Build the Coverage Matrix.** Map every enumerated source element to the test case(s) covering it, and mark each row Positive or Negative. Every source element must appear — an uncovered element is a gap to surface to the user, not to hide.
10. **Do not invent requirements.** If the source does not state a detail needed to write a test (e.g., the exact maximum value or valid format), flag it to the user instead of fabricating it. Consistent test data must come from the requirement, not from assumption.
