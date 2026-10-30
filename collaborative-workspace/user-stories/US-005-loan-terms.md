# US-005 - Validate Loan Amount and Tenure

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** lending company business representative <br>
**I want** the loan amount and tenure in every application to come from master data and stay within 500 to 200,000 <br>
**So that** the company only originates loan products it actually offers <br>

## Acceptance Criterias

### AC-005-001 - Valid Loan Terms Accepted

**Given** the loan amount and tenure exist in master data and the amount is between 500 and 200,000 <br>
**When** the partner submits the loan application <br>
**Then** the loan terms pass validation <br>

### AC-005-002 - Boundary Amounts Accepted

**Given** the loan amount is exactly 500 or exactly 200,000 and exists in master data <br>
**When** the partner submits the loan application <br>
**Then** the loan terms pass validation <br>

### AC-005-003 - Amount Below Minimum

**Given** the loan amount is below 500 <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

### AC-005-004 - Amount Above Maximum

**Given** the loan amount is above 200,000 <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

### AC-005-005 - Tenure Not in Master Data

**Given** the loan tenure does not exist in master data <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

### AC-005-006 - Loan Terms Missing

**Given** the loan amount or tenure is missing <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected, the missing fields are identified, and no application is recorded <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-005](../sequence-diagrams/SEQ-005-loan-terms.md) | Sequence Diagram | Loan amount and tenure validation |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
