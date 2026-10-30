# US-008 - Send Applications to Risk Scoring

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version, compiled from the business stakeholder (2027-01-02) and risk stakeholder (2027-01-03) meeting notes. |

## User Story

**As a** risk officer <br>
**I want** every recorded loan application sent to the risk scoring process <br>
**So that** the regional risk committee has a risk recommendation for every application it decides on <br>

## Acceptance Criterias

### AC-008-001 - Recorded Application Forwarded

**Given** a loan application has been recorded and acknowledged <br>
**When** the acknowledgement has been returned to the partner <br>
**Then** the application data required for risk scoring is forwarded to the Risk Scoring Service <br>

### AC-008-002 - Acknowledgement Not Delayed by Scoring

**Given** a valid loan application is submitted <br>
**When** the platform forwards it to risk scoring <br>
**Then** the acknowledgement to the partner is still returned within 2 seconds, without waiting for the risk score <br>

### AC-008-003 - Rejected Submission Not Forwarded

**Given** a submission was rejected during validation <br>
**When** validation completes <br>
**Then** nothing is forwarded to the Risk Scoring Service <br>

### AC-008-004 - Score Is a Recommendation Only

**Given** the Risk Scoring Service returns a risk category (VERY_LOW, LOW, HIGH, or VERY_HIGH) <br>
**When** the platform receives the category <br>
**Then** the application is neither approved nor rejected automatically <br>

## Related Artifacts

| ID | Type | Description |
|----|------|-------------|
| [SEQ-008](../sequence-diagrams/SEQ-008-risk-scoring-handoff.md) | Sequence Diagram | Asynchronous hand-off to risk scoring |
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |
