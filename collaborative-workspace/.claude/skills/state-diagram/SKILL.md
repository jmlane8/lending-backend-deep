---
name: state-diagram
description: Create a State Diagram (STA-NNN) using Mermaid stateDiagram-v2 syntax. Use when the user explicitly asks to draw or generate a state diagram and has provided states and transitions.
disable-model-invocation: true
---

# Skill: State Diagram

## Mermaid Type
`stateDiagram-v2`  
Syntax reference: https://mermaid.js.org/syntax/stateDiagram.html

## ID / Prefix
Prefix: `STA`  
Format: `STA-[3 digits]` — e.g., `STA-001`

## Sequencing
- Start from `STA-001` if no state diagram exists yet.
- Scan `state-diagrams/` and continue from the highest existing ID (ignore `STA-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`state-diagrams/STA-NNN-descriptive-name.md`  
Example: `state-diagrams/STA-003-order-lifecycle.md`  
The output only contains diagram code in Mermaid syntax, without any additional text or explanation.

## Process
- Draw a state diagram that models the complete lifecycle of a single entity (e.g., a loan application, a payment).  
- Include exactly one initial state and at least one final state.  
- Name each state as a noun phrase describing the entity's condition (e.g., "Under Review", "Approved"), not a verb.  
- Label every transition with the event that triggers it; if a guard condition applies, include it in brackets (e.g., [credit score >= 650]).  
- Every state must be reachable from the initial state; remove any unreachable states.  
- Every non-final state must have at least one outgoing transition; avoid dead-end states unless they are explicitly terminal.  
- Include entry or exit actions on a state only when they are part of the user-supplied requirements, not assumed behavior.
