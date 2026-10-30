# US-002 - Submit Loan Application with Customer Data

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** partner company <br>
**I want** to submit my customer's loan application together with their personal data <br>
**So that** my customers can obtain car-collateral loans through my channel and I earn commission on approved loans <br>

## Acceptance Criterias

### AC-002-001 - Valid Application Acknowledged

**Given** a contracted partner has complete and valid customer, collateral, and loan data <br>
**When** the partner submits the loan application <br>
**Then** the application is recorded and an acknowledgement is returned to the partner within 2 seconds <br>

### AC-002-002 - Optional Email Omitted

**Given** a contracted partner has complete mandatory customer data without an email address <br>
**When** the partner submits the loan application <br>
**Then** the application is accepted <br>

### AC-002-003 - Mandatory Customer Field Missing

**Given** the customer data is missing one or more mandatory fields (ID photo, national ID number, full legal name, selfie photo, phone number, gender, date of birth, place of birth, marital status, occupation, work duration, monthly income, number of dependents, or domicile street, district, city, or province) <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected, the missing fields are identified, and no application is recorded <br>

### AC-002-004 - Invalid Marital Status

**Given** the marital status is not one of Single, Married, Divorced, or Widowed <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected and the invalid field is identified <br>

### AC-002-005 - Invalid Occupation

**Given** the occupation is not one of Employee, Entrepreneur, or Civil Servant <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected and the invalid field is identified <br>

### AC-002-006 - Customer Under 18

**Given** the customer is younger than 18 years old on the submission date <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

### AC-002-007 - Customer Turns 18 on Submission Date

**Given** the customer's 18th birthday is on the submission date <br>
**When** the partner submits the loan application <br>
**Then** the customer passes the age check <br>

### AC-002-008 - Expired Government ID

**Given** the customer's government-issued ID is expired <br>
**When** the partner submits the loan application <br>
**Then** the submission is rejected <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-002](../sequence-diagrams/SEQ-002-submit-loan-application-customer-data.md) | Sequence Diagram | Customer data validation and acknowledgement |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
