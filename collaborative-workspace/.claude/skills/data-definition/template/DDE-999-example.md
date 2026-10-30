# DDE-999 — Customer Support Ticket

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2027-04-24 | - Initial version derived from UC-999 |

## Scope

Covers the data model for the customer support ticketing feature. Supports ticket creation, categorization, priority assignment, attachment storage, and confirmation notification.

## Source Artifacts

| Source ID | Type | Description |
|-----------|------|-------------|
| [UC-999](UC-999-example.md) | Use Case | Create Customer Support Ticket |

## Entity Summary

| Entity | Description | Source |
|--------|-------------|--------|
| [Ticket](#ticket) | A support request submitted by a customer | UC-999 |
| [Customer](#customer) | An authenticated user who submits support tickets | UC-999 |
| [Category](#category) | Classifies the ticket type and determines default priority | UC-999 |
| [Attachment](#attachment) | A file uploaded by the customer as part of a ticket submission | UC-999 |

## Entity Definitions

### Ticket

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| ticket_id | Internal unique identifier for the ticket | Number | Yes | PK |
| ticket_number | Human-readable unique ticket reference shown to the customer | Text | Yes | Unique |
| customer_id | The customer who submitted the ticket | Number | Yes | FK → Customer |
| category_id | The category selected by the customer at submission | Number | Yes | FK → Category |
| subject | Short summary of the support issue | Text | Yes | Max 255 characters |
| description | Full description of the support issue | Text | Yes | |
| status | Current state of the ticket | Text | Yes | Allowed values: OPEN |
| priority | Priority level assigned at creation | Text | Yes | Allowed values: NORMAL, HIGH |
| created_at | Timestamp when the ticket was created | Date | Yes | |

### Customer

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| customer_id | Internal unique identifier for the customer | Number | Yes | PK |
| email | Email address used for ticket confirmation notifications | Text | Yes | Unique |
| account_status | Whether the customer account is active | Text | Yes | Allowed values: ACTIVE, INACTIVE |

### Category

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| category_id | Internal unique identifier for the category | Number | Yes | PK |
| name | Display name of the category | Text | Yes | Unique |
| is_business_critical | Indicates whether this category triggers HIGH priority assignment | Boolean | Yes | |
| default_priority | Default priority assigned to tickets in this category | Text | Yes | Allowed values: NORMAL, HIGH |

### Attachment

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| attachment_id | Internal unique identifier for the attachment | Number | Yes | PK |
| ticket_id | The ticket this attachment belongs to | Number | Yes | FK → Ticket |
| file_name | Original file name as uploaded by the customer | Text | Yes | |
| file_size_bytes | Size of the file in bytes | Number | Yes | Max 20 971 520 (20 MB) |
| file_format | File extension of the uploaded file | Text | Yes | Allowed values: PDF, PNG, JPG, DOCX |
| storage_path | Internal storage location of the file | Text | Yes | |
| uploaded_at | Timestamp when the attachment was stored | Date | Yes | |

## Relationships

| From Entity | Cardinality | To Entity | Description |
|-------------|-------------|-----------|-------------|
| [Customer](#customer) | one-to-many | [Ticket](#ticket) | A customer can submit many tickets; each ticket belongs to one customer |
| [Category](#category) | one-to-many | [Ticket](#ticket) | A category can be applied to many tickets; each ticket has exactly one category |
| [Ticket](#ticket) | one-to-many | [Attachment](#attachment) | A ticket can have many attachments; each attachment belongs to one ticket |

## Entity Relationship Diagrams

| ID |
|----|
| [ERD-999](../entity-relationship-diagrams/ERD-999-example.md) |
