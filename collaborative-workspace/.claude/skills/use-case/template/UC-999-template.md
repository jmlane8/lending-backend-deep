# UC-999 — [Use Case Name]

## Changelog
<!--
Use semantic versioning. Start at 0.0.1 and increment only the patch version. Major/minor versions are user-managed. Sort versions newest first.
Document changes as concise bullet points. Each bullet should ideally be no more than two sentences.
-->

| Version | Date | Changes |
|---------|------|---------|
| `x.y.z` | yyyy-mm-dd | bullet points showing changes |

## Overview
<!-- 
All fields are mandatory, except Secondary Actors 
-->

| Use Case Element | Description |
|---|---|
| **ID** | Unique identifier for the use case (e.g., UC-001). |
| **Name** | Short, descriptive name of the use case. |
| **Primary Actor** | The main actor that initiates the use case. |
| **Secondary Actors** | Other actors, systems, or services involved in the use case. |
| **Description** | Brief summary of the business goal and what the use case accomplishes. |
| **Trigger** | Event or action that starts the use case. |
| **Preconditions** | Conditions that must be true before the use case can begin. |
| **Postconditions** | Conditions that must be true after the use case completes successfully. |

## Normal Flow
<!-- 
Mandatory. 
Happy path only. The sequence of steps when everything goes as expected. Keep each step atomic: one actor, one action per row. Do not handle errors or variations here. 
-->

| Step | Actor | Action |
|---|---|---|
| 1 | Actor | Describe the first interaction or action. |
| 2 | Actor/System | Describe the next action. |
| N | Actor/System | Continue until the use case is complete. |

## Alternative Flows
<!-- 
Optional, can be zero or more. 
Valid variations that still result in a successful outcome (not errors). Each alternative flow anchors to a specific step in the Normal Flow where the variation begins, and rejoins it at a specific step. 
-->

**AF-001: [Alternative Flow Name 1]**
- At step [X], [describe the condition that triggers this alternative flow].
- [Describe how the flow differs from the normal flow.]
- The flow continues at step [Y].

**AF-00N: [Alternative Flow Name N]**
- At step [P], [describe the condition that triggers this alternative flow].
- [Describe the alternate behavior.]
- The flow continues at step [Q].

## Exceptions
<!-- 
Optional, can be zero or more. 
Failure conditions that prevent the use case from completing successfully. Focus on business-level failures the system must handle gracefully, not implementation-level (coding/technical) errors. 
-->

| Exception | Condition | System Response |
|---|---|---|
| E-001: Exception Name | Describe the condition that causes the exception. | Describe how the system responds. |
| E-00N: Exception Name | Describe the condition that causes the exception. | Describe how the system responds. |

## Priority

[High / Medium / Low] — Explain why this use case is important relative to other requirements.

## Business Rules
<!-- 
Optional, can be zero or more. 
Constraints or policies that govern the behavior described in this use case. Each rule should originating from business or regulatory requirements. Not technical aspect (such as database product) 
-->

- Business rule 1
- Business rule 2
- Business rule N

## Assumptions
<!-- 
Optional, can be zero or more. 
Things believed to be true but not formally confirmed. Each one implicitly flags a risk: if it turns out to be false, the use case may need revision. Do not list confirmed facts or decisions already made. 
-->

- Assumption 1
- Assumption 2
- Assumption N

## Related Artifacts
<!-- 
Optional, can be zero or more. 
Artifacts directly connected to this use case that provide meaningful additional context. 
Use the artifact's official ID as a relative hyperlink where possible. 
-->

| ID | Type | Description |
|----|------|-------------|
| Artifact ID | Artifact type | What the other artifact shown |