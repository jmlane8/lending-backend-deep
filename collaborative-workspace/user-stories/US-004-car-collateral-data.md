# US-004 - Validate Car Collateral Data

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** lending company business representative <br>
**I want** the car collateral in every application to match master data and be owned by the customer or the spouse <br>
**So that** every loan is backed by a vehicle the company recognizes and that belongs to the applicant's household <br>

## Acceptance Criterias

### AC-004-001 - Valid Collateral Accepted

**Given** the vehicle brand, model, and manufacturing year exist in master data and the ownership holder is the customer <br>
**When** the partner submits the loan application <br>
**Then** the collateral passes validation <br>

### AC-004-002 - Spouse-Owned Collateral Accepted

**Given** the customer is Married, the vehicle data exists in master data, and the ownership holder is the spouse <br>
**When** the partner submits the loan application <br>
**Then** the collateral passes validation <br>

### AC-004-003 - Mandatory Collateral Field Missing

**Given** one or more of vehicle brand, model, license plate number, ownership holder, or manufacturing year is missing <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected, the missing fields are identified, and no application is recorded <br>

### AC-004-004 - Collateral Not in Master Data

**Given** the vehicle brand, model, or manufacturing year does not exist in master data <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

### AC-004-005 - Invalid Ownership Holder

**Given** the ownership holder is neither the customer nor the spouse <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-004](../sequence-diagrams/SEQ-004-car-collateral-data.md) | Sequence Diagram | Collateral validation against master data |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
