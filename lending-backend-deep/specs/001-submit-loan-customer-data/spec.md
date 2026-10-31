# Feature Specification: Submit Loan Application with Customer Data

**Feature Branch**: `us-001-loan-submission`

**Created**: 2026-10-30

**Status**: Draft

**Input**: User story US-002 — Submit Loan Application with Customer Data (`collaborative-workspace/user-stories/US-002-submit-loan-application-customer-data.md`, v0.0.1). "As a partner company I want to submit my customer's loan application together with their personal data so that my customers can obtain car-collateral loans through my channel and I earn commission on approved loans."

**Business objectives** (`docs/constitution-input-001.md`): O1 Revenue growth, O2 Partner self-sufficiency, O3 API performance (acknowledgment within 2 seconds).

## Clarifications

### Session 2026-10-30

- Q: With checks now running after acknowledgment, what should the platform keep for a submission that fails customer data validation? → A: Keep only a minimal record marked Rejected (reference, partner, dates, reasons); delete the personal data and photos when the submission is rejected.
- Q: How should a partner find out that an acknowledged submission was rejected? → A: Both: the platform sends the rejection, and the partner can also look up status and reasons by reference at any time.
- Q: How soon after acknowledgment must the platform finish customer data validation and make any rejection available to the partner? → A: Within 1 minute of acknowledgment, for 95% of submissions.
- Q: Beyond checking that fields are present and that marital status and occupation are accepted values, should the platform also reject values that are clearly wrong? → A: Yes, basic checks too: national ID and phone formats, gender from master data, dates that are real and not impossible, whole or positive numbers where expected, and photos that are readable images.
- Q: Which country's national ID and phone number formats should the platform accept? → A: Indonesia: 16-digit NIK, Indonesian mobile.
- Q: Should the 30-day resubmission wait apply to a submission rejected for customer data problems, such as a missing field or an expired ID? → A: No wait for data problems: the partner can correct and resubmit immediately. The 30 days apply only to other rejections, such as a later decision on the application.
- Q: Given that Ruth put the status API out of scope, how should FR-014's lookup be handled? → A: Keep sending rejections, and narrow the lookup to the validation result only (pending, passed, or rejected with reasons). Tracking of loan decisions stays out of scope.
- Q: Is Indonesia actually confirmed as the market, or should the format checks stay country-neutral until someone confirms it? → A: Not confirmed: use neutral checks for now (national ID must contain only letters and digits, phone must be a plausible international number), and record the country as an open question for Ruth. This replaces the earlier Indonesia answer.
- Q: Which gender values should the platform accept? → A: Male and Female, published through master data as the starting values and maintained by the business team.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Submit a valid application and receive an acknowledgment (Priority: P1)

A contracted partner's system submits a loan application for one of its customers, including the customer's personal data. The platform records the application and acknowledges it to the partner within 2 seconds, then validates the customer data after acknowledgment. When the data is complete and valid, the application passes validation. The customer's email address is optional and its absence does not prevent submission.

**Why this priority**: This is the core value of the story — without an accepted, acknowledged submission no partner-sourced application exists (O1), and the acknowledgment time is a fixed platform objective (O3).

**Independent Test**: Submit an application with complete, valid customer data (once with and once without an email address) and confirm that each is recorded and acknowledged within 2 seconds and then passes validation.

**Acceptance Scenarios**:

1. **Given** a contracted partner has complete and valid customer, collateral, and loan data, **When** the partner submits the loan application, **Then** the application is recorded and an acknowledgment is returned to the partner within 2 seconds. *(AC-002-001)*
2. **Given** a contracted partner has complete mandatory customer data without an email address, **When** the partner submits the loan application, **Then** the application is acknowledged and passes validation. *(AC-002-002)*

---

### User Story 2 - Reject incomplete or invalid customer data with the affected fields identified (Priority: P2)

When the customer data is missing mandatory fields, or the marital status or occupation is not an accepted value, the platform rejects the submission after acknowledgment and tells the partner separately which fields are missing or invalid, so the partner can correct the data without contacting the company.

**Why this priority**: Clear, field-level rejection keeps partner integrations self-sufficient (O2) and ensures only complete applications enter the pipeline (O1).

**Independent Test**: Submit applications that omit each mandatory field in turn, and applications with an unaccepted marital status or occupation; confirm each is acknowledged and then rejected, the affected fields are identified, the application goes no further, and only a minimal Rejected record remains.

**Acceptance Scenarios**:

1. **Given** the customer data is missing one or more mandatory fields (ID photo, ID expiry date, national ID number, full legal name, selfie photo, phone number, gender, date of birth, place of birth, marital status, occupation, work duration, monthly income, number of dependents, or domicile street, district, city, or province), **When** the partner submits the loan application, **Then** the submission is acknowledged, then rejected, the missing fields are identified to the partner separately, and only a minimal Rejected record remains. *(AC-002-003)*
2. **Given** the marital status is not one of Single, Married, Divorced, or Widowed, **When** the partner submits the loan application, **Then** the submission is acknowledged, then rejected, and the invalid field is identified to the partner separately. *(AC-002-004)*
3. **Given** the occupation is not one of Employee, Entrepreneur, or Civil Servant, **When** the partner submits the loan application, **Then** the submission is acknowledged, then rejected, and the invalid field is identified to the partner separately. *(AC-002-005)*
4. **Given** a customer field fails a basic value check in FR-015, **When** the partner submits the loan application, **Then** the submission is acknowledged, then rejected, and every invalid field is identified to the partner separately. *(added by clarification)*

---

### User Story 3 - Enforce customer eligibility: minimum age and valid government ID (Priority: P3)

An application passes validation only for customers who are at least 18 years old on the submission date and whose government-issued ID has not expired. Eligibility is checked after acknowledgment.

**Why this priority**: These eligibility rules prevent ineligible applications from consuming downstream processing, but they build on the submission and field validation delivered by Stories 1 and 2.

**Independent Test**: Submit applications for a customer who is 17, a customer whose 18th birthday is the submission date, and a customer with an expired government ID; confirm the first and third are rejected and the second passes the age check.

**Acceptance Scenarios**:

1. **Given** the customer is younger than 18 years old on the submission date, **When** the partner submits the loan application, **Then** the submission is acknowledged, then rejected. *(AC-002-006)*
2. **Given** the customer's 18th birthday is on the submission date, **When** the partner submits the loan application, **Then** the customer passes the age check. *(AC-002-007)*
3. **Given** the customer's government-issued ID is expired, **When** the partner submits the loan application, **Then** the submission is acknowledged, then rejected. *(AC-002-008)*

---

### Edge Cases

- Several mandatory fields are missing at once: all missing fields are identified in the same rejection, so the partner can fix them in one attempt.
- A mandatory field is present but empty: it is treated as missing.
- Marital status and occupation are both invalid: both fields are identified as invalid.
- The customer turns 18 the day after the submission date: rejected as under 18.
- The government-issued ID expires on the submission date: the ID is still valid that day and is not treated as expired.
- The email address is present: it is accepted along with the application; its absence never causes a rejection.
- Delivery of a rejection to the partner keeps failing: the partner can still get the validation result and reasons by looking up the application reference.
- The partner corrects a rejected submission and sends it again straight away: it is treated as a new submission, and the 30-day resubmission wait does not apply.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The platform MUST allow a contracted partner to submit a loan application for a customer together with the customer's personal data. *(O1)*
- **FR-002**: The platform MUST require these customer fields on every submission: government-issued ID photo, government-issued ID expiry date, national ID number, full legal name, selfie photo, phone number, gender, date of birth, place of birth, marital status, occupation, work duration (years), monthly income, number of dependents, and domicile address (street, district, city, province). *(O1)*
- **FR-003**: The platform MUST treat the customer's email address as optional and MUST NOT reject a submission because it is absent. *(O1)*
- **FR-004**: When one or more mandatory customer fields are missing, the platform MUST reject the submission and identify every missing field to the partner. *(O1, O2)*
- **FR-005**: The platform MUST accept only Single, Married, Divorced, or Widowed as the marital status; any other value MUST cause rejection with the marital status field identified as invalid. *(O1, O2)*
- **FR-006**: The platform MUST accept only Employee, Entrepreneur, or Civil Servant as the occupation; any other value MUST cause rejection with the occupation field identified as invalid. *(O1, O2)*
- **FR-007**: The platform MUST reject a submission when the customer is younger than 18 years old on the submission date, and MUST treat a customer whose 18th birthday falls on the submission date as meeting the age requirement. *(O1)*
- **FR-008**: The platform MUST reject a submission when the customer's government-issued ID is expired on the submission date. *(O1)*
- **FR-009**: The platform MUST record each submission as a loan application and return an acknowledgment to the partner before validating it. *(O1, O3)*
- **FR-010**: The platform MUST return the submission acknowledgment within 2 seconds of receiving the submission. *(O3)*
- **FR-011**: A rejected application MUST NOT proceed to any further processing. When a submission is rejected, the platform MUST delete its personal data and photos and keep only a minimal Rejected record: application reference, submitting partner, dates, and rejection reasons. *(O1)*
- **FR-012**: The platform MUST run the customer data validations in FR-004 to FR-008 and FR-015 after the submission is acknowledged, never inside the submission response, and MUST report any rejection to the partner separately from the acknowledgment (Constitution Principle III). This supersedes the in-response rejection described in US-002 and SEQ-002. *(O3)*
- **FR-013**: The platform MUST send each rejection, with its reasons, to the address the partner has registered for notifications, and MUST retry when delivery fails. *(O2, O4)*
- **FR-014**: The acknowledgment MUST include an application reference, and the partner MUST be able to look up the application's validation result (pending, passed, or rejected with reasons) by that reference. The lookup MUST NOT report loan decisions or any status after validation. *(O2, O4)*
- **FR-015**: The platform MUST reject a submission, identifying every affected field, when: the national ID number contains anything other than letters and digits; the phone number is not a plausible international number (an optional leading + followed by 8 to 15 digits); the gender is not an accepted value published through master data (initially Male or Female); the date of birth or ID expiry date is not a real calendar date, or the date of birth is after the submission date; work duration or number of dependents is not a whole number of 0 or more; monthly income is not greater than 0; or the ID photo or selfie photo is not a readable image. *(O1, O2)*
- **FR-016**: A rejection under FR-004 to FR-008 or FR-015 MUST NOT trigger the 30-day resubmission wait. The partner MUST be able to submit corrected data immediately. *(O1, O2)*

### Key Entities

- **Loan Application**: A partner's request for a car-collateral loan on behalf of one customer. It is recorded when the submission is acknowledged, then either passes validation or becomes a minimal Rejected record holding the application reference, submitting partner, dates, and rejection reasons, with no personal data or photos.
- **Customer**: The partner's end consumer applying for the loan. Holds identity data (government-issued ID photo and expiry, national ID number, full legal name, selfie photo, gender, date of birth, place of birth), contact data (phone number, optional email), personal and financial profile (marital status, occupation, work duration, monthly income, number of dependents), and domicile address (street, district, city, province).
- **Submission Acknowledgment**: The platform's confirmation to the partner that the application was recorded, carrying the application reference. It does not mean the application passed validation.
- **Submission Rejection**: The platform's notice to the partner, separate from the acknowledgment, that a submission failed validation. It is sent to the partner's registered address and is also available by application reference. It identifies the missing or invalid fields, or the eligibility rule that failed (under age, expired ID).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of submissions are acknowledged within 2 seconds of receipt. *(O3)*
- **SC-002**: 100% of submissions missing a mandatory customer field, carrying an unaccepted marital status or occupation, or failing a basic value check in FR-015, are rejected with every affected field identified, and none of them proceeds to further processing or keeps personal data or photos after rejection. *(O1, O2)*
- **SC-003**: 0 applications pass validation for customers under 18 on the submission date or with an expired government-issued ID. *(O1)*
- **SC-004**: A partner developer can correct a rejected submission using only the information in the rejection, without contacting the company, in at least 95% of cases during sandbox testing. *(O2)*
- **SC-005**: 0 valid submissions are rejected solely because the customer's email address is absent. *(O1)*
- **SC-006**: For 95% of submissions, customer data validation finishes and any rejection is available to the partner within 1 minute of acknowledgment. *(O2)*

## Assumptions

- Partner identification and authorization (that the caller is a contracted partner) are covered by a separate story (US-001) and are not part of this specification.
- The validity of collateral data and loan data mentioned in AC-002-001 is covered by separate stories (US-004 car collateral data, US-005 loan terms). This specification covers only the customer data.
- Spouse data for married customers, coverage area, resubmission rules, and the risk-scoring handoff are covered by separate stories (US-003, US-006, US-007, US-008) and are not part of this specification, except that customer data rejections are exempt from the US-007 30-day wait (FR-016).
- The government-issued ID expiry date is supplied with the submission alongside the ID photo, because the expiry check (AC-002-008) needs it, so it is treated as a mandatory field. The story does not say how it is captured.
- "Submission date" means the calendar date on which the platform receives the submission, in the company's business time zone.
- The operating country is not confirmed, so the national ID and phone checks in FR-015 are country-neutral. The business stakeholder needs to confirm the country, after which these checks can be tightened. Format checks confirm shape only; they do not confirm that a national ID or phone number exists or belongs to the customer.
- An ID whose expiry date equals the submission date is not yet expired.
- A customer born on 29 February reaches 18 on 1 March in non-leap years.
- Gender is added to master data, maintained by the business team. Its starting values (Male, Female) still need the business stakeholder's confirmation.
- General loan status tracking through a status API is out of scope for this project (business stakeholder notes, 2027-01-02). The validation result lookup in FR-014 is narrower than that and still needs the business stakeholder's confirmation.
- The acknowledgment's content beyond confirming that the application was recorded and giving the application reference is not defined by this story.
