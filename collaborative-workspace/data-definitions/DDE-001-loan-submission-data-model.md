# DDE-001 — Loan Submission

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| `0.0.1` | 2026-10-29 | - Initial version derived from UC-001.<br>- Added ERD-001. |

## Scope

Covers the data model for partner loan submission through the Loan Submission API. It supports partner access control, customer and spouse data capture, car collateral data, loan terms, master data validation (vehicles, loan amounts, tenures, coverage areas), and the resubmission rules that depend on rejected applications and active loans.

## Source Artifacts

| Source ID | Type | Description |
|-----------|------|-------------|
| [UC-001](../use-cases/UC-001-submit-loan-application.md) | Use Case | Submit Loan Application |

## Entity Summary

| Entity | Description | Source |
|--------|-------------|--------|
| [Partner](#partner) | A contracted company that submits loan applications through the API (the loan channel) | UC-001 |
| [Partner Credential](#partner-credential) | API credentials a partner uses to authenticate submissions | UC-001 |
| [Loan Application](#loan-application) | A car-collateral loan application submitted by a partner on behalf of a customer | UC-001 |
| [Customer](#customer) | The individual applying for a loan through a partner | UC-001 |
| [Spouse](#spouse) | The customer's spouse, recorded when the customer is married | UC-001 |
| [Collateral](#collateral) | The car pledged as collateral for a loan application | UC-001 |
| [Loan](#loan) | An approved loan resulting from a loan application | UC-001 |
| [Vehicle Brand](#vehicle-brand) | Master data: a vehicle brand accepted as collateral | UC-001 |
| [Vehicle Model](#vehicle-model) | Master data: a vehicle model accepted as collateral | UC-001 |
| [Manufacturing Year](#manufacturing-year) | Master data: a vehicle manufacturing year accepted as collateral | UC-001 |
| [Loan Amount Option](#loan-amount-option) | Master data: a loan amount partners can select | UC-001 |
| [Loan Tenure](#loan-tenure) | Master data: a loan tenure partners can select | UC-001 |
| [Coverage Area](#coverage-area) | Master data: a geographic area where the company accepts applications | UC-001 |

## Entity Definitions

### Partner

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| partner_id | Unique identifier of the partner | UUID | Yes | PK |
| agreement_status | Whether the partner's partnership agreement is active | Text | Yes | Allowed values: ACTIVE, INACTIVE. Only ACTIVE partners may submit |
| agreement_signed_date | Date the legal partnership agreement was signed | Date | Yes |  |

### Partner Credential

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| partner_credential_id | Unique identifier of the credential | UUID | Yes | PK |
| partner_id | Partner that owns the credential | UUID | Yes | FK → Partner |
| is_valid | Whether the credential can currently be used to authenticate | Boolean | Yes |  |
| issued_at | When the credential was provisioned | Date | Yes | Must be after Partner.agreement_signed_date |

### Loan Application

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| loan_application_id | Unique identifier of the loan application | UUID | Yes | PK |
| application_reference | Unique reference returned to the partner in the acknowledgement | Text | Yes | Unique. Based on a UC-001 assumption |
| partner_id | Partner that submitted the application | UUID | Yes | FK → Partner |
| customer_id | Customer the application is for | UUID | Yes | FK → Customer |
| loan_amount_option_id | Requested loan amount, selected from master data | UUID | Yes | FK → Loan Amount Option |
| loan_tenure_id | Requested loan tenure, selected from master data | UUID | Yes | FK → Loan Tenure |
| application_status | Current status of the application | Text | Yes | Allowed values not defined in UC-001; must at least distinguish rejected and approved applications |
| submitted_at | When the application was recorded | Date | Yes | Acknowledged within 2 seconds of submission |
| rejected_at | When the application was rejected, if it was | Date | No | Used for the 30-day resubmission rule |

### Customer

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| customer_id | Unique identifier of the customer | UUID | Yes | PK |
| national_id_number | Customer's national ID number | Text | Yes | Unique. Identifies "the same application" for the 30-day resubmission rule (UC-001 assumption) |
| full_legal_name | Customer's full legal name | Text | Yes |  |
| government_id_photo | Photo of the customer's government-issued ID | Text | Yes | Reference to the stored image |
| government_id_expiry_date | Expiry date of the government-issued ID | Date | Yes | Must not be in the past at submission. How it is captured is not specified in UC-001 |
| selfie_photo | Selfie photo of the customer | Text | Yes | Reference to the stored image |
| phone_number | Customer's phone number | Text | Yes |  |
| gender | Customer's gender | Text | Yes | Allowed values not specified in UC-001 |
| date_of_birth | Customer's date of birth | Date | Yes | Customer must be at least 18 on the submission date |
| place_of_birth | Customer's place of birth | Text | Yes |  |
| marital_status | Customer's marital status | Text | Yes | Allowed values: SINGLE, MARRIED, DIVORCED, WIDOWED |
| occupation | Customer's occupation | Text | Yes | Allowed values: EMPLOYEE, ENTREPRENEUR, CIVIL_SERVANT |
| work_duration_years | How long the customer has worked, in years | Number | Yes |  |
| monthly_income | Customer's monthly income | Number | Yes | Currency not specified in UC-001 |
| number_of_dependents | Number of people who depend on the customer | Number | Yes |  |
| email | Customer's email address | Text | No |  |
| domicile_street | Street of the customer's domicile address | Text | Yes |  |
| domicile_district | District of the customer's domicile address | Text | Yes | Checked against Coverage Area |
| domicile_city | City of the customer's domicile address | Text | Yes | Checked against Coverage Area |
| domicile_province | Province of the customer's domicile address | Text | Yes | Checked against Coverage Area |

### Spouse

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| spouse_id | Unique identifier of the spouse record | UUID | Yes | PK |
| loan_application_id | Application the spouse data was submitted with | UUID | Yes | FK → Loan Application. Unique |
| national_id_number | Spouse's national ID number | Text | Yes |  |
| full_legal_name | Spouse's full legal name | Text | Yes |  |
| date_of_birth | Spouse's date of birth | Date | Yes | Spouse must be at least 18 on the submission date |
| phone_number | Spouse's phone number | Text | Yes |  |
| email | Spouse's email address | Text | No |  |
| government_id_photo | Photo of the spouse's government-issued ID | Text | No | Reference to the stored image |

### Collateral

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| collateral_id | Unique identifier of the collateral record | UUID | Yes | PK |
| loan_application_id | Application the collateral secures | UUID | Yes | FK → Loan Application. Unique |
| vehicle_model_id | Vehicle model, selected from master data | UUID | Yes | FK → Vehicle Model (which determines the brand) |
| manufacturing_year_id | Manufacturing year, selected from master data | UUID | Yes | FK → Manufacturing Year |
| license_plate_number | Vehicle's license plate number | Text | Yes | Must not be tied to an Active Loan |
| ownership_holder | Who legally owns the vehicle | Text | Yes | Allowed values: CUSTOMER, SPOUSE |

### Loan

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| loan_id | Unique identifier of the loan | UUID | Yes | PK |
| loan_application_id | Application the loan was approved from | UUID | Yes | FK → Loan Application. Unique |
| settlement_status | Whether the loan is still outstanding | Text | Yes | Allowed values: ACTIVE, SETTLED. An ACTIVE loan blocks new applications for the same license plate |

### Vehicle Brand

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| vehicle_brand_id | Unique identifier of the brand | UUID | Yes | PK |
| brand_name | Name of the brand | Text | Yes | Unique |

### Vehicle Model

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| vehicle_model_id | Unique identifier of the model | UUID | Yes | PK |
| vehicle_brand_id | Brand the model belongs to | UUID | Yes | FK → Vehicle Brand |
| model_name | Name of the model | Text | Yes |  |

### Manufacturing Year

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| manufacturing_year_id | Unique identifier of the manufacturing year | UUID | Yes | PK |
| year | Manufacturing year | Number | Yes | Unique |

### Loan Amount Option

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| loan_amount_option_id | Unique identifier of the amount option | UUID | Yes | PK |
| amount | Loan amount | Number | Yes | Between 500 and 200,000 inclusive. Currency not specified in UC-001 |

### Loan Tenure

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| loan_tenure_id | Unique identifier of the tenure | UUID | Yes | PK |
| tenure | Length of the loan | Number | Yes | Unit (e.g., months) not specified in UC-001 |

### Coverage Area

| Field | Description | Data Type | Mandatory | Notes |
|-------|-------------|-----------|-----------|-------|
| coverage_area_id | Unique identifier of the coverage area | UUID | Yes | PK |
| province | Province covered | Text | Yes |  |
| city | City covered | Text | No | Granularity of coverage areas is not specified in UC-001 |
| district | District covered | Text | No | Granularity of coverage areas is not specified in UC-001 |

## Relationships

| From Entity | Cardinality | To Entity | Description |
|-------------|-------------|-----------|-------------|
| [Partner](#partner) | one-to-many | [Partner Credential](#partner-credential) | A partner can hold many credentials; each credential belongs to one partner |
| [Partner](#partner) | one-to-many | [Loan Application](#loan-application) | A partner can submit many applications; each application comes from one partner |
| [Customer](#customer) | one-to-many | [Loan Application](#loan-application) | A customer can have many applications over time; each application is for one customer |
| [Loan Application](#loan-application) | one-to-zero-or-one | [Spouse](#spouse) | An application includes spouse data only when the customer is married |
| [Loan Application](#loan-application) | one-to-one | [Collateral](#collateral) | Each application is secured by exactly one car |
| [Loan Application](#loan-application) | one-to-zero-or-one | [Loan](#loan) | An application results in a loan only if it is approved |
| [Vehicle Brand](#vehicle-brand) | one-to-many | [Vehicle Model](#vehicle-model) | A brand has many models; each model belongs to one brand |
| [Vehicle Model](#vehicle-model) | one-to-many | [Collateral](#collateral) | A model can classify many collateral records |
| [Manufacturing Year](#manufacturing-year) | one-to-many | [Collateral](#collateral) | A manufacturing year can apply to many collateral records |
| [Loan Amount Option](#loan-amount-option) | one-to-many | [Loan Application](#loan-application) | An amount option can be selected by many applications |
| [Loan Tenure](#loan-tenure) | one-to-many | [Loan Application](#loan-application) | A tenure can be selected by many applications |

## Entity Relationship Diagrams

| ID |
|----|
| [ERD-001](../entity-relationship-diagrams/ERD-001-loan-submission.md) |
