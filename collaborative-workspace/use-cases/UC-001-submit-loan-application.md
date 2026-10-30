# UC-001 — Submit Loan Application

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.2` | 2026-10-29 | - Confirmed that vehicle ownership is valid only with the customer or the spouse; added it to field validation (step 5, E-003) and business rules.<br>- Removed the open assumption about "another party" ownership. |
| `0.0.1` | 2026-10-28 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes and OB-001. |

## Overview

| Use Case Element | Description |
|---|---|
| **ID** | UC-001 |
| **Name** | Submit Loan Application |
| **Primary Actor** | Partner System |
| **Secondary Actors** | Lending Platform, Master Data (maintained by the business team), Risk Scoring Service (managed by the risk team) |
| **Description** | A contracted partner submits a car-collateral loan application on behalf of its customer through the Loan Submission API. The platform validates the partner, the customer, spouse, collateral, and loan data against business rules and master data, records the application, acknowledges it to the partner, and hands it over to risk scoring asynchronously. |
| **Trigger** | The partner system sends a loan submission request to the Loan Submission API. |
| **Preconditions** | 1. The partner has a signed legal partnership agreement with the lending company.<br>2. The partner holds valid API credentials, provisioned after the agreement was signed.<br>3. Master data (vehicle brands, models, manufacturing years, loan amounts, loan tenures, coverage areas) is maintained and available. |
| **Postconditions** | 1. The loan application is recorded on the platform.<br>2. The partner has received a submission acknowledgement within 2 seconds.<br>3. The application has been forwarded to the Risk Scoring Service for scoring. |

## Normal Flow

| Step | Actor | Action |
|---|---|---|
| 1 | Partner System | Sends a loan submission request containing customer, collateral, and loan data. |
| 2 | Lending Platform | Authenticates the partner's API credentials. |
| 3 | Lending Platform | Verifies that the partner has an active partnership agreement. |
| 4 | Lending Platform | Validates that all mandatory customer, collateral, and loan fields are present. |
| 5 | Lending Platform | Validates marital status, occupation, and vehicle ownership holder against their accepted values. |
| 6 | Lending Platform | Verifies that the customer is at least 18 years old. |
| 7 | Lending Platform | Verifies that the customer's government-issued ID is not expired. |
| 8 | Lending Platform | Verifies that the customer's domicile address is within a defined coverage area. |
| 9 | Lending Platform | Validates the vehicle brand, model, and manufacturing year against master data. |
| 10 | Lending Platform | Validates the loan amount and loan tenure against master data and the allowed amount range. |
| 11 | Lending Platform | Verifies that the license plate is not tied to an active loan. |
| 12 | Lending Platform | Verifies that the same application was not rejected within the last 30 days. |
| 13 | Lending Platform | Records the loan application. |
| 14 | Lending Platform | Returns a submission acknowledgement to the partner. |
| 15 | Partner System | Receives the submission acknowledgement. |
| 16 | Lending Platform | Forwards the application data to the Risk Scoring Service asynchronously. |
| 17 | Risk Scoring Service | Receives the application for risk scoring. |

## Alternative Flows

**AF-001: Customer Is Married**
- At step 6, the customer's marital status is Married.
- The platform validates that the spouse's mandatory fields (national ID number, full legal name, date of birth, phone number) are present, and verifies that the spouse is at least 18 years old.
- The flow continues at step 7.

## Exceptions

| Exception | Condition | System Response |
|---|---|---|
| E-001: Partner Not Authorized | The credentials are invalid, or the caller is not a legally contracted partner. | Reject the request. No application is recorded. |
| E-002: Mandatory Data Missing | One or more mandatory customer, spouse (when Married), collateral, or loan fields are missing. | Reject the submission and identify the missing fields. |
| E-003: Invalid Field Value | Marital status, occupation, or vehicle ownership holder is not one of the accepted values. | Reject the submission and identify the invalid fields. |
| E-004: Applicant Under Age | The customer, or the spouse when Married, is younger than 18. | Reject the submission. |
| E-005: Expired Government ID | The customer's government-issued ID is expired. | Reject the submission. |
| E-006: Outside Coverage Area | The submission originates from a location outside the defined coverage areas. | Reject the submission. |
| E-007: Collateral Not in Master Data | Vehicle brand, model, or manufacturing year does not match master data. | Reject the submission. |
| E-008: Invalid Loan Terms | Loan amount or tenure is not in master data, or the amount is outside 500–200,000. | Reject the submission. |
| E-009: Collateral Tied to Active Loan | The license plate is tied to an active loan that is not fully settled. | Reject the submission. |
| E-010: Resubmission Too Early | The same application was rejected less than 30 days ago. | Reject the submission. |

## Priority

High — loan submission is the core capability behind the objective of 30,000 partner-sourced applications within 12 months of launch (OB-001 O1).

## Business Rules

- Only legally contracted partners may submit loan applications.
- Mandatory customer data: government-issued ID photo (not expired), national ID number, full legal name, selfie photo, phone number, gender, date of birth, place of birth, marital status, occupation, work duration (years), monthly income, number of dependents, and domicile address (street, district, city, province). Email address is optional.
- Accepted marital status values: Single, Married, Divorced, Widowed.
- Accepted occupation values: Employee, Entrepreneur, Civil Servant.
- The customer must be at least 18 years old.
- Spouse data is required when marital status is Married. Mandatory: national ID number, full legal name, date of birth, phone number. Optional: email address, government-issued ID photo. The spouse must be at least 18 years old.
- Mandatory collateral data: vehicle brand, vehicle model, license plate number, ownership holder, manufacturing year. Collateral values must match master data.
- Vehicle ownership is valid only when held by the customer or the spouse. Any other ownership holder is rejected.
- Mandatory loan data: loan amount (minimum 500, maximum 200,000) and loan tenure, both selected from master data.
- Submissions from locations outside the defined coverage areas are rejected.
- A license plate is valid for submission only if it is not tied to an active loan.
- After a rejection, the same application may be resubmitted only after 30 days.
- After an approval, a new application for the same collateral may be submitted only after the existing loan is fully settled.
- Partners may submit an unlimited number of loan applications.
- A loan submission must be acknowledged within 2 seconds; long-running operations such as risk scoring run asynchronously (OB-001 O3).
- Every loan application goes through risk scoring. Scores are recommendations only and do not approve or reject an application.

## Assumptions

- The customer's "monthly income" satisfies the risk team's "monthly household income – customer". Spouse monthly income is not currently collected at submission.
- Ownership/legal documents (existence and validity) and passport number, requested by the risk team, are not yet part of the submission data.
- The coverage area check is based on the customer's domicile address.
- "Same application" for the 30-day resubmission rule is identified by the customer's national ID number.
- Age is calculated as of the submission date.
- The acknowledgement returns a unique application reference to the partner.

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [FLW-001](../flowcharts/FLW-001-submit-loan-application.md) | Flowchart | Normal flow of loan application submission, including AF-001 |
| [OB-001](../objectives-and-boundaries/OB-001-loan-submission.md) | Business Objectives & Boundaries | Objectives and boundaries for loan submission |
