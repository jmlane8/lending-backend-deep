# US-001 - Restrict Loan Submission to Contracted Partners

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** lending company business representative <br>
**I want** only legally contracted partners to be able to submit loan applications <br>
**So that** every application comes through a channel bound by a partnership agreement and commission terms <br>

## Acceptance Criterias

### AC-001-001 - Contracted Partner Accepted

**Given** a partner has a signed partnership agreement and valid API credentials <br>
**When** the partner submits a loan application <br>
**Then** the submission is accepted for data validation <br>

### AC-001-002 - Invalid Credentials Rejected

**Given** a request carries missing or invalid API credentials <br>
**When** the request submits a loan application <br>
**Then** the submission is rejected and no application is recorded <br>

### AC-001-003 - Non-Contracted Party Rejected

**Given** a caller has no active partnership agreement <br>
**When** the caller submits a loan application <br>
**Then** the submission is rejected and no application is recorded <br>

### AC-001-004 - No Submission Volume Limit

**Given** a contracted partner has already submitted many loan applications <br>
**When** the partner submits another loan application <br>
**Then** the submission is accepted for data validation regardless of how many applications the partner has submitted before <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-001](../sequence-diagrams/SEQ-001-partner-access-control.md) | Sequence Diagram | Partner authentication and agreement check |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
