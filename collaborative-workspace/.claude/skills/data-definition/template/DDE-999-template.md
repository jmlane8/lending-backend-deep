# DDE-999 — [Data Model Name]

## Changelog
<!--
Use semantic versioning. Start at 0.0.1 and increment only the patch version. Major/minor versions are user-managed. Sort versions newest first.
Document changes as concise bullet points. Each bullet should ideally be no more than two sentences.
-->

| Version | Date | Changes |
|---------|------|---------|
| `x.y.z` | yyyy-mm-dd | Summary of changes |

## Scope
<!-- 
  Briefly describe which feature or business domain this data model covers. 
-->

Short description of the business domain and the features this model supports. 

## Source Artifacts
<!-- 
List every user story or use case that contributed entities or fields to this model. 
Use the source's official ID as a relative hyperlink. 
-->

| Source ID | Type | Description |
|-----------|------|-------------|
| ID | User Story / Use Case | [Brief description] |

## Entity Summary
<!--
Each entity must exist on `Entity Definitions` section.
Use the entity name as a relative hyperlink.
-->

| Entity | Description |
|--------|-------------|
| Entity Name | Brief description of the entity's role |

## Entity Definitions
<!-- 
One sub-section per entity. 
Include all fields identified from the source requirements. 
Fields may be extended by the functional or technical team. 
Use `lower_snake_case` for entity field name.
-->

### [Entity Name]

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|----------|-------|
| field_name | Description | Text / Number / Date / Boolean | Yes / No | PK / FK → [Entity Name] / Unique |

## Relationships
<!--
Each entity (from and to) must exist on `Entity Definitions` section.
Use the entity name as a relative hyperlink.
-->

| From Entity | Cardinality | To Entity | Description |
|-------------|------------|-----------|-------------|
| [Entity A] | one-to-many, one-to-one, etc | [Entity B] | Brief description |

## Entity Relationship Diagrams
<!-- 
List of generated entity relationship diagram, if any.
Use the ERD ID as a relative hyperlink where possible. 
-->

| ID | 
|----|
| ERD ID |