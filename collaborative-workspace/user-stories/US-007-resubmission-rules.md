# US-007 - Enforce Resubmission Rules

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** lending company business representative <br>
**I want** resubmissions blocked for 30 days after a rejection and while the collateral backs an unsettled loan <br>
**So that** rejected applications are not reprocessed too early and no vehicle backs two loans at once <br>

## Acceptance Criterias

### AC-007-001 - First Submission Accepted

**Given** there is no previous application for the customer or the license plate <br>
**When** the partner submits the loan application <br>
**Then** the application passes the resubmission check <br>

### AC-007-002 - Resubmission Within 30 Days of Rejection

**Given** the same application was rejected less than 30 days ago <br>
**When** the partner submits the loan application again <br>
**Then** the submission is rejected <br>

### AC-007-003 - Resubmission After 30 Days

**Given** the same application was rejected 30 or more days ago <br>
**When** the partner submits the loan application again <br>
**Then** the application passes the resubmission check <br>

### AC-007-004 - Collateral Tied to Active Loan

**Given** the license plate is tied to an approved loan that is not fully settled <br>
**When** the partner submits a new application with that license plate <br>
**Then** the submission is rejected <br>

### AC-007-005 - Collateral Loan Fully Settled

**Given** the license plate was tied to an approved loan that is now fully settled <br>
**When** the partner submits a new application with that license plate <br>
**Then** the application passes the resubmission check <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-007](../sequence-diagrams/SEQ-007-resubmission-rules.md) | Sequence Diagram | Resubmission checks |
| [FLW-002](../flowcharts/FLW-002-resubmission-check.md) | Flowchart | Resubmission decision logic |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
