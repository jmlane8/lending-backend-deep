# Test Scenario - Submit Loan Application with Customer Data

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-11-01 | - Initial version, covering all 8 acceptance criteria of US-002 with expected results that follow the implemented asynchronous behaviour. |

## Source Requirement

| ID | Type | Description |
|----|------|-------------|
| [US-002](../user-stories/US-002-submit-loan-application-customer-data.md) | User Story | Submit Loan Application with Customer Data |

---

## Test Basis & Assumptions

- **Superseded behaviour**: US-002 v0.0.1 says a failed submission is rejected and no application is recorded. These cases follow the later spec (`specs/001-submit-loan-customer-data/spec.md`, FR-011 and FR-012) as built: the submission is acknowledged with `202`, then rejected, and a minimal Rejected record remains. US-002 needs updating to match.
- **Test data source**: US-002 states no formats or limits. API values, date formats and size limits come from the spec, `data-model.md` and `contracts/loan-applications.openapi.yaml`.
- **Mandatory field added**: ID expiry date is mandatory under FR-002, although AC-002-003 doesn't list it. AC-002-008 can't be checked without it.
- **Open value conflict**: US-002 names values such as `Single`, but the API accepts only exact values such as `SINGLE`. TC-001-036 and TC-001-041 record this current behaviour, which the business stakeholder hasn't confirmed.
- **Not covered**: rejection webhooks, value checks beyond these criteria (FR-015), partner authentication (US-001), and collateral and loan validity (US-004, US-005).
- **Environment**: the service runs as described in `DEVELOPMENT.md`, with `make db-seed` loaded and `BUSINESS_TIMEZONE=UTC`. The active partner key is `sandbox-partner-key`.
- **Dates**: each date is set relative to the day of the test run. The values in brackets assume a submission date of `2026-11-01`. Avoid running across midnight in the business time zone.
- **Lookup timing**: validation finishes within 1 minute for 95% of submissions, so a lookup that is still `PENDING` after 1 minute needs investigating.
- **Independence**: each case submits its own application, so cases can run in any order.
- **Baseline customer record**: `resources/sample/loan-application.json`, with these values:

| Field | Value |
|---|---|
| government-id-photo | Valid base64 image from the sample file |
| government-id-expiry-date | `2099-05-14` |
| national-id-number | `3201012345678901` |
| full-legal-name | `Alex Morgan` |
| selfie-photo | Valid base64 image from the sample file |
| phone-number | `+15551234567` |
| email | `alex@example.com` |
| gender | `FEMALE` |
| date-of-birth | `1990-02-11` |
| place-of-birth | `Springfield` |
| marital-status | `SINGLE` |
| occupation | `EMPLOYEE` |
| work-duration-years | `6` |
| monthly-income | `4500.5` |
| number-of-dependents | `1` |
| domicile-street | `12 Example Street` |
| domicile-district | `Central` |
| domicile-city | `Springfield` |
| domicile-province | `North` |

## Coverage Matrix

| Source Element | Path | Covered By |
|---|---|---|
| AC-002-001 - Valid Application Acknowledged | Positive | TC-001-001 to TC-001-003 |
| AC-002-002 - Optional Email Omitted | Positive | TC-001-004 to TC-001-006 |
| AC-002-003 - Mandatory Customer Field Missing | Negative | TC-001-007 to TC-001-030 |
| AC-002-004 - Invalid Marital Status (accepted values) | Positive | TC-001-031 to TC-001-034 |
| AC-002-004 - Invalid Marital Status | Negative | TC-001-035, TC-001-036, TC-001-042 |
| AC-002-005 - Invalid Occupation (accepted values) | Positive | TC-001-037 to TC-001-039 |
| AC-002-005 - Invalid Occupation | Negative | TC-001-040 to TC-001-042 |
| AC-002-006 - Customer Under 18 | Negative | TC-001-043 to TC-001-045, TC-001-052 |
| AC-002-007 - Customer Turns 18 on Submission Date | Positive | TC-001-046 to TC-001-048 |
| AC-002-008 - Expired Government ID (unexpired boundary) | Positive | TC-001-050, TC-001-051 |
| AC-002-008 - Expired Government ID | Negative | TC-001-049, TC-001-052 |

---

## Test Scenarios

### Scenario 1: Valid application is acknowledged and passes validation

#### TC-001-001 - Valid application is acknowledged within 2 seconds
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-001 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| customer | Baseline record, unchanged |
| collateral | `{}` |
| loan | `{}` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Measure the time from sending the request to receiving the full response. |

**Expected Result**
- The response is `202` with an `application-reference` starting `LA-`, received within 2 seconds.

#### TC-001-002 - Valid application passes validation
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-001 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| customer | Baseline record, unchanged |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-003 - Largest allowed application is acknowledged within 2 seconds
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-001 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| government-id-photo | Base64 JPEG, 5 MiB decoded |
| selfie-photo | Base64 JPEG, 5 MiB decoded |
| Request body | At most 16 MiB |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Measure the time from sending the request to receiving the full response. |

**Expected Result**
- The response is `202` with an `application-reference`, received within 2 seconds.

### Scenario 2: Email address is optional

#### TC-001-004 - Application without email passes validation
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-002 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| email | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-005 - Application with null email passes validation
**Type**: Positive <br>
**Priority**: Medium <br>
**Traces To**: AC-002-002 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| email | `null` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-006 - Application with empty email passes validation
**Type**: Positive <br>
**Priority**: Medium <br>
**Traces To**: AC-002-002 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| email | `""` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

### Scenario 3: Missing mandatory customer field is rejected with the field identified

#### TC-001-007 - Missing ID photo is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| government-id-photo | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.government-id-photo`, `code` `MISSING`.

#### TC-001-008 - Missing ID expiry date is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| government-id-expiry-date | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.government-id-expiry-date`, `code` `MISSING`.

#### TC-001-009 - Missing national ID number is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| national-id-number | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.national-id-number`, `code` `MISSING`.

#### TC-001-010 - Missing full legal name is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| full-legal-name | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.full-legal-name`, `code` `MISSING`.

#### TC-001-011 - Missing selfie photo is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| selfie-photo | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.selfie-photo`, `code` `MISSING`.

#### TC-001-012 - Missing phone number is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| phone-number | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.phone-number`, `code` `MISSING`.

#### TC-001-013 - Missing gender is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| gender | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.gender`, `code` `MISSING`.

#### TC-001-014 - Missing date of birth is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.date-of-birth`, `code` `MISSING`.

#### TC-001-015 - Missing place of birth is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| place-of-birth | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.place-of-birth`, `code` `MISSING`.

#### TC-001-016 - Missing marital status is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.marital-status`, `code` `MISSING`.

#### TC-001-017 - Missing occupation is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.occupation`, `code` `MISSING`.

#### TC-001-018 - Missing work duration is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| work-duration-years | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.work-duration-years`, `code` `MISSING`.

#### TC-001-019 - Missing monthly income is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| monthly-income | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.monthly-income`, `code` `MISSING`.

#### TC-001-020 - Missing number of dependents is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| number-of-dependents | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.number-of-dependents`, `code` `MISSING`.

#### TC-001-021 - Missing domicile street is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| domicile-street | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.domicile-street`, `code` `MISSING`.

#### TC-001-022 - Missing domicile district is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| domicile-district | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.domicile-district`, `code` `MISSING`.

#### TC-001-023 - Missing domicile city is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| domicile-city | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.domicile-city`, `code` `MISSING`.

#### TC-001-024 - Missing domicile province is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| domicile-province | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.domicile-province`, `code` `MISSING`.

#### TC-001-025 - Several missing fields are identified together
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| gender | Field left out |
| phone-number | Field left out |
| domicile-city | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly three reasons, for `customer.gender`, `customer.phone-number` and `customer.domicile-city`, each with `code` `MISSING`.

#### TC-001-026 - Empty mandatory field is treated as missing
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| full-legal-name | `""` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.full-legal-name`, `code` `MISSING`.

#### TC-001-027 - Whitespace-only mandatory field is treated as missing
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| place-of-birth | `"   "` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.place-of-birth`, `code` `MISSING`.

#### TC-001-028 - Null mandatory field is treated as missing
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | `null` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.occupation`, `code` `MISSING`.

#### TC-001-029 - Submission with a missing field is still acknowledged
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| gender | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |

**Expected Result**
- The response is `202` with an `application-reference`, and no field problem appears in the response.

#### TC-001-030 - Rejected submission keeps no customer data
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-003 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed
- The tester can read the test database

**Test Data**
| Field | Value |
|---|---|
| gender | Field left out |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |
| 3 | Query `loan_application` and `application_payload` for the application. |

**Expected Result**
- The `loan_application` row has a NULL `customer_id`, and no `application_payload` row exists.

### Scenario 4: Marital status must be an accepted value

#### TC-001-031 - Marital status SINGLE is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-004 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `SINGLE` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-032 - Marital status MARRIED is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-004 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `MARRIED` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-033 - Marital status DIVORCED is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-004 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `DIVORCED` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-034 - Marital status WIDOWED is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-004 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `WIDOWED` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-035 - Unlisted marital status is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-004 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `SEPARATED` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.marital-status`, `code` `INVALID`.

#### TC-001-036 - Marital status in the wrong case is rejected
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-004 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `Single` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.marital-status`, `code` `INVALID`.

### Scenario 5: Occupation must be an accepted value

#### TC-001-037 - Occupation EMPLOYEE is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-005 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | `EMPLOYEE` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-038 - Occupation ENTREPRENEUR is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-005 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | `ENTREPRENEUR` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-039 - Occupation CIVIL_SERVANT is accepted
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-005 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | `CIVIL_SERVANT` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-040 - Unlisted occupation is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-005 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | `STUDENT` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.occupation`, `code` `INVALID`.

#### TC-001-041 - Occupation in display form is rejected
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-005 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| occupation | `Civil Servant` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.occupation`, `code` `INVALID`.

#### TC-001-042 - Invalid marital status and occupation are identified together
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-004, AC-002-005 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| marital-status | `SEPARATED` |
| occupation | `STUDENT` |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly two reasons, for `customer.marital-status` and `customer.occupation`, each with `code` `INVALID`.

### Scenario 6: Customer must be at least 18 on the submission date

#### TC-001-043 - Customer one day under 18 is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-006 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | Submission date minus 18 years, plus 1 day (`2008-11-02`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.date-of-birth`, `code` `CUSTOMER_UNDER_AGE`.

#### TC-001-044 - Customer well under 18 is rejected
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-006 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | Submission date minus 10 years (`2016-11-01`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.date-of-birth`, `code` `CUSTOMER_UNDER_AGE`.

#### TC-001-045 - Customer born 29 February is under 18 on 28 February of a non-leap year
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-006 <br>

**Preconditions**
- Run against `validation.Validate` with the fixed submission date below, because the API always uses the current date

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | `2008-02-29` |
| Submission date | `2026-02-28` |

**Steps**
| Step | Action |
|---|---|
| 1 | Call `validation.Validate` with the baseline record and the Test Data. |

**Expected Result**
- One reason is returned: `field` `customer.date-of-birth`, `code` `CUSTOMER_UNDER_AGE`.

#### TC-001-046 - Customer whose 18th birthday is the submission date passes
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-007 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | Submission date minus 18 years (`2008-11-01`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-047 - Customer one day over 18 passes
**Type**: Positive <br>
**Priority**: Medium <br>
**Traces To**: AC-002-007 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | Submission date minus 18 years, minus 1 day (`2008-10-31`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-048 - Customer born 29 February reaches 18 on 1 March of a non-leap year
**Type**: Positive <br>
**Priority**: Medium <br>
**Traces To**: AC-002-007 <br>

**Preconditions**
- Run against `validation.Validate` with the fixed submission date below, because the API always uses the current date

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | `2008-02-29` |
| Submission date | `2026-03-01` |

**Steps**
| Step | Action |
|---|---|
| 1 | Call `validation.Validate` with the baseline record and the Test Data. |

**Expected Result**
- No reasons are returned.

### Scenario 7: Government-issued ID must not be expired

#### TC-001-049 - ID that expired the day before is rejected
**Type**: Negative <br>
**Priority**: High <br>
**Traces To**: AC-002-008 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| government-id-expiry-date | Submission date minus 1 day (`2026-10-31`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly one reason: `field` `customer.government-id-expiry-date`, `code` `GOVERNMENT_ID_EXPIRED`.

#### TC-001-050 - ID that expires on the submission date passes
**Type**: Positive <br>
**Priority**: High <br>
**Traces To**: AC-002-008 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| government-id-expiry-date | Submission date (`2026-11-01`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-051 - ID that expires the day after passes
**Type**: Positive <br>
**Priority**: Medium <br>
**Traces To**: AC-002-008 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| government-id-expiry-date | Submission date plus 1 day (`2026-11-02`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `PASSED` with no `reasons`.

#### TC-001-052 - Under-age customer with an expired ID gets both reasons
**Type**: Negative <br>
**Priority**: Medium <br>
**Traces To**: AC-002-006, AC-002-008 <br>

**Preconditions**
- The active partner holds the baseline customer record, with only the Test Data below changed

**Test Data**
| Field | Value |
|---|---|
| date-of-birth | Submission date minus 18 years, plus 1 day (`2008-11-02`) |
| government-id-expiry-date | Submission date minus 1 day (`2026-10-31`) |

**Steps**
| Step | Action |
|---|---|
| 1 | Submit the application with `POST /loan-applications/v1`, using the `sandbox-partner-key` API key. |
| 2 | Read `GET /loan-applications/v1/{application-reference}/validation-result` until `status` is no longer `PENDING`. |

**Expected Result**
- The lookup returns `status` `REJECTED` with exactly two reasons: `customer.date-of-birth` with `CUSTOMER_UNDER_AGE`, and `customer.government-id-expiry-date` with `GOVERNMENT_ID_EXPIRED`.
