---
name: user-journey
description: Create a User Journey diagram (UJR-NNN) using Mermaid journey syntax. Use when the user explicitly asks to draw or generate a user journey and has provided steps, actors, and satisfaction scores.
disable-model-invocation: true
---

# Skill: User Journey

## Mermaid Type
`journey`  
Syntax reference: https://mermaid.js.org/syntax/userJourney.html

## ID / Prefix
Prefix: `UJR`  
Format: `UJR-[3 digits]` — e.g., `UJR-001`

## Sequencing
- Start from `UJR-001` if no user journey exists yet.
- Scan `user-journeys/` and continue from the highest existing ID (ignore `UJR-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`user-journeys/UJR-NNN-descriptive-name.md`  
Example: `user-journeys/UJR-002-checkout-process.md`  
The output only contains diagram code in Mermaid syntax, without any additional text or explanation.

## Process
- Draw a user journey that maps a specific actor's experience through a process from start to goal completion.  
- Group steps into named sections or phases that reflect natural stages in the experience (e.g., Discovery, Application, Decision).  
- Each step must represent a single action or touchpoint; assign it to the correct actor.  
- Assign the user-supplied satisfaction score to each step; if any score is missing, ask the user for it — do not invent scores.  
- Cover the full journey including failure paths and recovery steps where the user-supplied input describes them.  
- Use the same actor names as those defined in the related use case or user story, if one exists; otherwise use the actor names from the user's input.  
- Keep step labels short and action-oriented (e.g., "Submit application form", "Receive approval notification").