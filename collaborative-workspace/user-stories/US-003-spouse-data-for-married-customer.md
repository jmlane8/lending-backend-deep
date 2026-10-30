# US-003 - Provide Spouse Data for Married Customers

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** lending company business representative <br>
**I want** partners to provide spouse data whenever the customer is married <br>
**So that** the spouse's identity is on record before the company assesses a married customer's loan <br>

## Acceptance Criterias

### AC-003-001 - Married Customer with Complete Spouse Data

**Given** the customer's marital status is Married and the spouse's national ID number, full legal name, date of birth, and phone number are provided <br>
**When** the partner submits the loan application <br>
**Then** the spouse data passes validation <br>

### AC-003-002 - Spouse Optional Fields Omitted

**Given** the customer is Married and the spouse's email address and government-issued ID photo are not provided <br>
**When** the partner submits the loan application <br>
**Then** the spouse data passes validation <br>

### AC-003-003 - Spouse Mandatory Field Missing

**Given** the customer is Married and one or more mandatory spouse fields are missing <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected, the missing spouse fields are identified, and no application is recorded <br>

### AC-003-004 - Spouse Under 18

**Given** the customer is Married and the spouse is younger than 18 years old on the submission date <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

### AC-003-005 - Unmarried Customer Without Spouse Data

**Given** the customer's marital status is Single, Divorced, or Widowed and no spouse data is provided <br>
**When** the partner submits the loan application <br>
**Then** the application passes validation without spouse data <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-003](../sequence-diagrams/SEQ-003-spouse-data-for-married-customer.md) | Sequence Diagram | Spouse data validation for married customers |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
