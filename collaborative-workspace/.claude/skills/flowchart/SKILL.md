---
name: flowchart
description: Create a Flowchart diagram (FLW-NNN) using Mermaid flowchart syntax. Use when the user explicitly asks to draw or generate a flowchart and has provided process or decision flow details.
disable-model-invocation: true
---

# Skill: Flowchart

## Mermaid Type
`flowchart` 
Syntax reference: https://mermaid.js.org/syntax/flowchart.html

## ID / Prefix
Prefix: `FLW`  
Format: `FLW-[3 digits]` — e.g., `FLW-001`

## Sequencing
- Start from `FLW-001` if no flowchart exists yet.
- Scan `flowcharts/` and continue from the highest existing ID (ignore `FLW-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`flowcharts/FLW-NNN-descriptive-name.md`  
Example: `flowcharts/FLW-004-customer-support-process.md`  
The output only contains diagram code in Mermaid syntax, without any additional text or explanation.

## Process
- Draw a flowchart for the process or decision flow described by the user.  
- Use a single start node and one or more clearly labeled end nodes.  
- Give each decision node exactly two outgoing branches with condition labels (e.g., Yes/No, Approved/Rejected).  
- Represent only one action per process step.  
- Split compound actions into separate process steps.  
- When multiple actors are involved, group steps by actor using subgraphs (Mermaid flowchart has no native swim lanes).  
- Arrange the layout to avoid crossing lines.  
- Keep node labels short and in active verb form (e.g., "Validate application", "Send notification").  
