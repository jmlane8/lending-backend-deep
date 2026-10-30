# UC-999 — Create Customer Support Ticket

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `1.0.0` | 2027-04-24 | - Initial version |

## Overview

| Use Case Element | Description |
|---|---|
| **ID** | UC-002 |
| **Name** | Create Customer Support Ticket |
| **Primary Actor** | Customer |
| **Secondary Actors** | Customer Support System, Notification Service |
| **Description** | A customer submits a support request through the company's support portal. The system validates the request, creates a ticket, assigns an initial priority, and notifies the customer that the ticket has been received. |
| **Trigger** | The customer submits the support request form through the web portal. |
| **Preconditions** | 1. The customer has an active account.<br>2. The customer is authenticated.<br>3. The support portal is operational. |
| **Postconditions** | 1. A support ticket is created with a unique ticket number.<br>2. The ticket status is set to OPEN.<br>3. A confirmation notification is sent to the customer.<br>4. The ticket becomes available to the support team. |

## Normal Flow

| Step | Actor | Action |
|---|---|---|
| 1 | Customer | Opens the support portal and selects "Create New Ticket". |
| 2 | Customer | Enters the ticket details, including category, subject, description, and optional attachments. |
| 3 | Customer | Clicks the Submit button. |
| 4 | Support System | Validates the submitted information and attachment formats. |
| 5 | Support System | Generates a unique ticket number. |
| 6 | Support System | Creates the ticket with status OPEN. |
| 7 | Support System | Assigns an initial priority based on the selected category. |
| 8 | Support System | Stores the ticket and related attachments. |
| 9 | Notification Service | Sends a confirmation email containing the ticket number. |
| 10 | Support System | Displays a successful submission message to the customer. |

## Alternative Flows

**AF-001: Ticket Submitted with Attachments**
- At step 2, the customer uploads one or more supporting files.
- The system validates the file size and format.
- The flow continues at step 4.

**AF-002: High-Priority Category Selected**
- At step 2, the customer selects a category marked as business-critical.
- At step 7, the system automatically assigns HIGH priority instead of NORMAL priority.
- The flow continues normally.

## Exceptions

| Exception | Condition | System Response |
|---|---|---|
| E-001: Customer Not Authenticated | The customer session is missing or expired. | Redirect the customer to the login page. No ticket is created. |
| E-002: Required Fields Missing | Subject, category, or description is missing. | Display validation errors and prevent submission. |
| E-003: Invalid Attachment | An uploaded file exceeds the maximum size or uses an unsupported format. | Display an error identifying the problematic file. No ticket is created. |
| E-004: System Unavailable | The ticketing database or application is unavailable. | Display an error message and request the customer to try again later. |

## Priority

Medium — important for customer satisfaction and operational efficiency.

## Business Rules

- Every ticket must have exactly one category.
- Ticket numbers must be unique across the platform.
- Attachments must not exceed 20 MB per file.
- Supported attachment formats are PDF, PNG, JPG, and DOCX.
- Every newly created ticket must start with status OPEN.
- A confirmation notification must be sent after successful ticket creation.

## Assumptions

- Customers have access to a valid email address for notifications.
- Support agents monitor newly created tickets during business hours.
- The notification service is available and capable of delivering emails.
- File scanning for malware is performed by an external security service.

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-999](../sequence-diagrams/SEQ-999-example.md) | Sequence Diagram | Submit Loan Application |
