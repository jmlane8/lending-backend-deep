---
name: entity-relationship-diagram
description: Create an Entity Relationship Diagram (ERD-NNN) using Mermaid erDiagram syntax. Use when the user explicitly asks to draw or generate an ERD and has provided entity or data model details, or when triggered by the data-definition skill.
disable-model-invocation: true
---

# Skill: Entity Relationship Diagram (ERD)

## Mermaid Type
`erDiagram`  
Syntax reference: https://mermaid.js.org/syntax/entityRelationshipDiagram.html

## ID / Prefix
Prefix: `ERD`  
Format: `ERD-[3 digits]` — e.g., `ERD-001`

## Sequencing
- Start from `ERD-001` if no ERD exists yet.
- Scan `entity-relationship-diagrams/` and continue from the highest existing ID (ignore `ERD-NNN-template.md`, if any).
- Never reuse or change an ID.

## Output Path
`entity-relationship-diagrams/ERD-NNN-descriptive-name.md`  
Example: `entity-relationship-diagrams/ERD-001-user-management.md`  
The output only contains diagram code in Mermaid syntax, without any additional text or explanation.

## Process
- Draw an entity relationship diagram for the feature described by the user, or — when triggered by the `data-definition` skill — for the entities and fields in the generated data definition document (use it as the sole input).  
- Use singular entity names (e.g., Loan, Applicant, Vehicle).  
- Show crow's foot cardinality on every relationship (one-to-one, one-to-many, or many-to-many).  
- Replace every many-to-many relationship with a junction entity with its own attribute.  
- Include a primary key for every entity and foreign keys where needed.  
- Model only entities and attributes relevant to the feature. Do not model the entire system in one diagram.
- Use consistent UPPER_SNAKE_CASE entity naming throughout the diagram.  
