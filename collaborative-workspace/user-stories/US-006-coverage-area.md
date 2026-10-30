# US-006 - Reject Submissions Outside the Coverage Area

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** lending company business representative <br>
**I want** applications from locations outside the defined coverage areas to be rejected <br>
**So that** the company only lends where it is able to operate and serve the loan <br>

## Acceptance Criterias

### AC-006-001 - Location Within Coverage Area

**Given** the customer's domicile address is within a coverage area defined in master data <br>
**When** the partner submits the loan application <br>
**Then** the location passes validation <br>

### AC-006-002 - Location Outside Coverage Area

**Given** the customer's domicile address is outside every coverage area defined in master data <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-006](../sequence-diagrams/SEQ-006-coverage-area.md) | Sequence Diagram | Coverage area check |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
