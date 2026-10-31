# Technical Design - [Feature Title]

## Changelog
<!--
Use semantic versioning. Start at 0.0.1 and increment only the patch version. Major/minor versions are user-managed. Sort versions newest first.
Document changes as concise bullet points. Each bullet should ideally be no more than two sentences.
-->

| Version | Date | Changes |
|---------|------|---------|
| `x.y.z` | yyyy-mm-dd | bullet points showing changes |

## Source Requirement
<!-- Mandatory, one or more rows. The use case(s) / user story(ies) this test scenario
is derived from. Use the artifact's official ID as a relative hyperlink when possible. -->

| ID | Type | Description |
|----|------|-------------|
| [UC-002](../use-cases/UC-002-create-support-ticket.md) | Use Case | Create Customer Support Ticket |
| [US-007](../user-stories/US-007-submit-support-request.md) | User Story | Submit support request |

## Scope
<!--
Show what this implementation covers as narrative bullet points.
Show any explicit boundary (out-of-scope)
-->

---

## API Design
<!--
This section is required only when REST API is involved.
Each endpoint is defined in its own sub-section.
For each endpoint, define sub-section as the template below.
All JSON field names in request and response bodies use lower-kebab-case (e.g. `error-code`, `estimated-value`).
-->

### `/api/v1/foo/bar`
<!--
Always include Path Variables, Query Parameters, Request Schema, and Response Schema sections. If a section is not applicable, keep it and set its content to "None".
-->
**Verb**: one of standard HTTP method (`GET | POST | PUT | PATCH | DELETE | OPTIONS`) <br>
**Path**: define only the path, not the server address. Must match with sub-section title (e.g. `/api/v1/foo/bar`) <br>
**Authentication**: None | JWT Bearer | Custom HTTP Header | OAuth | Other (specify) <br> 
**Produces**: MIME type of the response <br>
**Consumes**: [optional] MIME type of the request body (e.g. `application/json`) <br>

#### Path variables
<!-- Path variables are always mandatory, so no Mandatory column is needed. -->
| Parameter | Valid value |
|-----------|-------------|
| parameter name | fixed enum, description (e.g. "valid loan submission ID"), regular expression, or simply "Any" if there is no validation |

#### Query parameters
| Parameter | Mandatory | Valid value |
|-----------|-----------|-------------|
| parameter name | YES/NO | fixed enum, description (e.g. "valid loan submission ID"), regular expression, or simply "Any" if there is no validation |

#### Request Schema
**Schema**
```
{
  "name": string,        // [mandatory] Person full name, use regex ^[A-Za-z ]{3,50}$
  "age": number,         // [optional] Person age, numeric, between 18-150
  "address": {
    "street": string,    // [mandatory] Domicile address, use regex ^[A-Za-z0-9 \-/.,&()#:;\\]{10,200}$
    "city": string,      // [mandatory] Domicile city, one of (ATLANTIS | NORTH POLE | SOUTH POLE), case insensitive
    "zip-code": string   // [optional] If present, must be 5-digit numeric, use regex ^\d{5}$
  },
  "abilities": [         // [optional] Contains one or more strings matching ^[A-Za-z0-9 \-]{3,30}$
    string
  ],
  "active": boolean      // [mandatory]
}
```

**Example**
<!--
Structure, field name, data type, and value must match the respective schema in the same section
-->

```json
{
  "name": "John Doe",
  "age": 25,
  "address": {
    "street": "Alpha Street 143C",
    "city": "Atlantis",
    "zip-code": "24892"
  },
  "abilities": [
    "swim", "run", "eat"
  ],
  "active": true
}
```

#### Response Schema
<!--
Must show: at least one happy path (2xx).
Optionally show: one or more not happy path (non-2xx).
-->

**Schema for 200 OK**
```
{
  "exists": boolean,     // true if user found in database
  "estimated-value": number,
  "spouse": {            // [optional] only present if user has a registered spouse in the system
    "spouse-id": string  // UUID v4
  }
}
```

**Example for 200 OK**
```json
{
  "exists": true,
  "estimated-value": 4000,
  "spouse": {
    "spouse-id": "dabb8c4b-cb2d-4fa3-acd3-00f1a87c8345"
  }
}
```

**Schema for 404 Not Found**
```
{
  "error-code": string,    // machine-readable error identifier
  "error-message": string  // human-readable description of the error
}
```

**Example for 404 Not Found**
```json
{
  "error-code": "USER_NOT_FOUND",
  "error-message": "No user found with the given identifier."
}
```

#### Non-2xx Error Registry
<!-- The 404 and 500 rows below are illustrative examples only; replace them with this endpoint's actual errors. -->

| HTTP Status | Error Code | Scenario | Retryable |
|-------------|------------|----------|-----------|
| Relevant HTTP status code, must exist in the Response schema section | Error code must match the `error-code` on response schema | [scenario, when this will happen] | [YES/NO] |
| 404 | USER_NOT_FOUND | No user with given ID on database | NO |
| 500 | ERR_INTERNAL | Internal server error (unknown cause) | YES |

---

## Scheduled Task Design
<!--
This section is required only when a background/scheduled task is involved.
Each task is defined in its own sub-section.
For each task, define sub-section as the template below.
-->

### [Task Name]
<!--
Always include Trigger, Schedule, Timezone, Concurrency Control, and Idempotency. If a section is not applicable, keep it and set its content to "None".
-->
**Trigger**: Cron | Event-driven | Manual | Other (specify) <br>
**Schedule**: cron expression (e.g. `0 2 * * *`) or triggering condition <br>
**Timezone**: timezone used to evaluate the schedule (e.g. `Asia/Jakarta`) <br>
**Concurrency Control**: how overlapping or duplicate runs are prevented (e.g. distributed lock, DB advisory lock, single-instance only) <br>
**Idempotency**: explain how re-running the task (e.g. due to retry, manual re-trigger, or overlap) does not produce duplicate side effects <br>

#### Inputs
| Parameter | Source | Valid value |
|-----------|--------|-------------|
| parameter name | Config / Database / Upstream API / Environment Variable | fixed enum, description, or simply "Any" if there is no validation |

#### Processing Steps
<!--
High-level, numbered list of what the task does end-to-end.
Detailed pseudocode and sequence diagram belong in the `Basic Flow` section, traced back to this sub-section.
-->

#### Outputs / Side Effects
| Output | Description |
|--------|-------------|
| e.g. Updated rows in `loan_installment` | [what changes and under what condition] |

#### Failure & Retry Registry
<!--
Analogous to the Non-2xx Error Registry in API Design.
-->
| Failure Scenario | Handling | Retryable | Alerting |
|-------------------|----------|-----------|----------|
| [scenario, when this will happen] | [e.g. skip record and continue / abort and roll back run] | YES/NO | [who/what gets notified, e.g. Slack channel, on-call] |

---

## Database Design
<!--
This section is required only when database change is involved.
-->

### Table Definitions
<!-- 
One sub-section per table. 
Include all fields identified from the source requirements. 
Use `lower_snake_case` for field name.
Each table must have:
- `created_at`, `updated_at`: timestamp with timezone, not null, default is current timestamp
- `deleted`: boolean, not null, default false (soft-delete flag)
- `deleted_at`: timestamp with timezone, default null
- `created_by`, `updated_by`: text not null, default is `system`
- `deleted_by`: text default null
-->

#### [Table Name]

**Fields**
| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|----------|-------|
| field_name | Description | Text / Number / Date / Timestamp / Boolean | Yes / No | PK / FK → [Entity Name] / Unique |

**Indices**
<!--
List all indices applicable for this table
-->

**DDL**
<!--
SQL for generating or modifying the table.
Use PostgreSQL syntax.
-->

---

## Basic Flow
<!--
For each entry point, e.g. REST API or scheduled task, define the flow containing basic logic.
-->

### Entry Point 1 (e.g. API /api/v1/foo/bar)
<!--
Entry point must be traceable to the `API Design` or `Scheduled Task` section.
-->

#### Sequence Diagram 
<!--
Provide a step-by-step sequence diagram showing all services involved. 
Use concrete service name; don't use generic labels (e.g. "System", "Application", "Backend").
Use mermaid notation.
-->

#### Basic Logic

```
Pseudocode of the basic logic.
If there is any database operation (SELECT / INSERT / UPDATE), the base SQL must be written.
Deleting database is not allowed, unless soft delete (set the `deleted` field as `true` and `deleted_at` as current timestamp. Set the `deleted_by` as the current user or `system` if the current user cannot be deduced).
```