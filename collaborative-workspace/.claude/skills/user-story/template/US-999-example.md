# US-999 - Create Customer Support Ticket

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `1.0.0` | 2027-04-24 | - Initial version |

## User Story

**As a** customer <br>
**I want** to submit a support request through the support portal <br>
**So that** I can get assistance with my issue and receive a unique ticket number to track its progress <br>

## Acceptance Criterias

### AC-999-001 - Successful Ticket Submission

**Given** the customer is authenticated and the support portal is operational <br>
**When** the customer submits the support request form with all required fields (category, subject, and description) <br>
**Then** a ticket is created with a unique ticket number, the ticket status is set to OPEN, an initial priority is assigned based on the selected category, a confirmation email containing the ticket number is sent to the customer, and a success message is displayed on the portal <br>

### AC-999-002 - Ticket Submitted with Attachments

**Given** the customer is authenticated and is filling out the support request form <br>
**When** the customer uploads one or more files that are within the 20 MB size limit and use a supported format (PDF, PNG, JPG, or DOCX), then submits the form <br>
**Then** the ticket is created successfully and the attachments are stored alongside the ticket <br>

### AC-999-003 - High-Priority Category Auto-Assignment

**Given** the customer is authenticated and is filling out the support request form <br>
**When** the customer selects a category marked as business-critical and submits the form <br>
**Then** the ticket is created with HIGH priority assigned automatically instead of the default NORMAL priority <br>

### AC-999-004 - Required Fields Validation

**Given** the customer is on the support ticket submission form <br>
**When** the customer attempts to submit the form with one or more required fields (category, subject, or description) left empty <br>
**Then** validation errors are displayed identifying the missing fields and the ticket is not created <br>

### AC-999-005 - Invalid Attachment Rejected

**Given** the customer is on the support ticket submission form <br>
**When** the customer uploads a file that exceeds 20 MB or uses an unsupported format, then attempts to submit <br>
**Then** an error message identifying the problematic file is displayed and the ticket is not created <br>

### AC-999-006 - Unauthenticated Access Blocked

**Given** the customer's session is missing or has expired <br>
**When** the customer attempts to access the support portal to submit a request <br>
**Then** the customer is redirected to the login page and no ticket is created <br>

### AC-999-007 - System Unavailable

**Given** the ticketing database or support application is unavailable <br>
**When** the customer attempts to submit a support request <br>
**Then** an error message is displayed informing the customer to try again later and no ticket is created <br>

| ID | Type | Description |
|----|------|-------------|
| [SEQ-999](../sequence-diagrams/SEQ-999-example.md) | Sequence Diagram | Submit Loan Application |
