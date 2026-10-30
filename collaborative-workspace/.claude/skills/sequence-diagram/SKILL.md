---
name: sequence-diagram
description: Create a Sequence Diagram (SEQ-NNN) using Mermaid sequenceDiagram syntax. Use when the user explicitly asks to draw or generate a sequence diagram and has provided actors and interaction steps.
disable-model-invocation: true
---

# Skill: Sequence Diagram

## Mermaid Type
`sequenceDiagram`  
Syntax reference: https://mermaid.js.org/syntax/sequenceDiagram.html

## ID / Prefix
Prefix: `SEQ`  
Format: `SEQ-[3 digits]` — e.g., `SEQ-001`

## Sequencing
- Start from `SEQ-001` if no sequence diagram exists yet.
- Scan `sequence-diagrams/` and continue from the highest existing ID (ignore `SEQ-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`sequence-diagrams/SEQ-NNN-descriptive-name.md`  
Example: `sequence-diagrams/SEQ-007-order-processing.md`  
The output only contains diagram code in Mermaid syntax, without any additional text or explanation.

## Process
- Draw a sequence diagram showing how actors and systems interact over time for a specific scenario.  
- Declare all participants at the top; use the same actor names as the related use case or user story, if one exists; otherwise use the actor names from the user's input.  
- Show synchronous calls with solid arrows and responses with dashed arrows.  
- Label every arrow with its message or return value.  
- Use activation boxes to indicate when a participant is processing.  
- Model conditional or optional behavior with alt or opt fragments.  
- Include only participants directly involved in the scenario.  
- Keep message labels concise and action-oriented (e.g., "POST /loans", "Return acknowledgment").