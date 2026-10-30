---
name: mindmap
description: Create a Mindmap diagram (MND-NNN) using Mermaid mindmap syntax. Use when the user explicitly asks to draw or generate a mindmap and has provided topics or concept structure to map.
disable-model-invocation: true
---

# Skill: Mindmap

## Mermaid Type
`mindmap`  
Syntax reference: https://mermaid.js.org/syntax/mindmap.html

## ID / Prefix
Prefix: `MND`  
Format: `MND-[3 digits]` — e.g., `MND-001`

## Sequencing
- Start from `MND-001` if no mindmap exists yet.
- Scan `mindmaps/` and continue from the highest existing ID (ignore `MND-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`mindmaps/MND-NNN-descriptive-name.md`  
Example: `mindmaps/MND-007-project-planning.md`  
The output only contains diagram code in Mermaid syntax, without any additional text or explanation.

## Process
- Draw a mindmap that organizes user-provided topics around a single central concept.  
- Place the central concept at the root, with first-level branches representing distinct, non-overlapping dimensions.  
- Keep node labels short, ideally two to four words, and avoid full sentences.  
- Limit the mindmap to three levels of depth.  
- Split the content into separate mindmaps if a fourth level is needed.  
- Use consistent grammatical structure across sibling nodes.  
- Include only topics explicitly provided by the user.  

