# Data Model: Submit Loan Application with Customer Data

Based on `collaborative-workspace/data-definitions/DDE-001-loan-submission-data-model.md`. Differences from DDE-001 are listed at the end.

## Tables

All timestamps are `TIMESTAMPTZ`. Audit columns follow the template (`created_at`, `updated_at`); the template's `deleted` columns are not used on these tables.

### partner (minimum for this feature; US-001 owns the rest)

| Column | Type | Rules |
|---|---|---|
| partner_id | UUID | PK |
| agreement_status | TEXT | `ACTIVE` or `INACTIVE` |
| agreement_signed_date | DATE | NOT NULL |
| notification_url | TEXT | Nullable. HTTPS required outside `local` |
| notification_secret | TEXT | Nullable. Required when `notification_url` is set |

### partner_credential

| Column | Type | Rules |
|---|---|---|
| partner_credential_id | UUID | PK |
| partner_id | UUID | FK to partner |
| key_hash | BYTEA | UNIQUE. SHA-256 of the API key |
| is_valid | BOOLEAN | NOT NULL |
| issued_at | TIMESTAMPTZ | NOT NULL. Must be after `partner.agreement_signed_date` |

**Rule enforcement**: PostgreSQL `CHECK` constraints can't read another table, so a `BEFORE INSERT OR UPDATE` trigger on this table enforces the `issued_at` rule. A credential issued later on the signing day passes.

### loan_application

| Column | Type | Rules |
|---|---|---|
| loan_application_id | UUID | PK (UUIDv7) |
| application_reference | TEXT | UNIQUE. `LA-<uuid>` |
| partner_id | UUID | FK to partner, NOT NULL |
| customer_id | UUID | FK to customer. NULL until passed |
| application_status | TEXT | `PENDING`, `PASSED`, or `REJECTED`. Later stories add decision statuses |
| rejection_reasons | JSONB | Set only when `REJECTED` |
| submitted_at | TIMESTAMPTZ | NOT NULL. When the submission was received; fixes the submission date |
| acknowledged_at | TIMESTAMPTZ | NOT NULL |
| validated_at | TIMESTAMPTZ | Set on pass or rejection |
| rejected_at | TIMESTAMPTZ | Set on rejection |
| rejection_stage | TEXT | `CUSTOMER_DATA` when rejected here |
| validation_attempts | INTEGER | Default 0. Counts failed attempts |
| next_validation_at | TIMESTAMPTZ | NOT NULL. Delays a retry after an internal failure |

**Index**: partial index on `(next_validation_at)` where `application_status = 'PENDING'`.

**Minimal Rejected record (FR-011)**: on a rejected row, only the reference, partner, dates, and reasons hold values. `customer_id` stays NULL, and no `application_payload` row remains.

**30-day rule (FR-016)**: `rejection_stage = 'CUSTOMER_DATA'` marks rows that US-007 must ignore. These rows hold no national ID anyway.

### application_payload

| Column | Type | Rules |
|---|---|---|
| loan_application_id | UUID | PK, FK to loan_application `ON DELETE CASCADE` |
| body | BYTEA | Exact request body |

**Lifecycle**: created with the application. Deleted on rejection. On pass, `body` is rewritten to hold only the `collateral` and `loan` sections, for US-004 and US-005.

### customer

Columns follow DDE-001: `customer_id` (PK), `national_id_number` (UNIQUE), `full_legal_name`, `government_id_photo` (UUID, FK to stored_photo), `government_id_expiry_date`, `selfie_photo` (UUID, FK to stored_photo), `phone_number`, `gender`, `date_of_birth`, `place_of_birth`, `marital_status`, `occupation`, `work_duration_years`, `monthly_income` (`NUMERIC(20,2)`), `number_of_dependents`, `email` (nullable), and `domicile_street`, `domicile_district`, `domicile_city`, `domicile_province`.

Written only when an application passes, using an upsert on `national_id_number`.

### stored_photo

| Column | Type | Rules |
|---|---|---|
| photo_id | UUID | PK |
| content_type | TEXT | `image/jpeg` or `image/png` |
| content | BYTEA | Decoded image |

### rejection_notification

| Column | Type | Rules |
|---|---|---|
| event_id | UUID | PK. Sent to the partner for de-duplication |
| loan_application_id | UUID | FK to loan_application, UNIQUE |
| status | TEXT | `PENDING`, `DELIVERED`, or `EXHAUSTED` |
| attempts | INTEGER | Default 0 |
| next_attempt_at | TIMESTAMPTZ | NOT NULL |
| delivered_at | TIMESTAMPTZ | Set on delivery |

No row is created if the partner has no `notification_url`; the lookup still works, and a warning is logged.

**Index**: partial index on `(next_attempt_at)` where `status = 'PENDING'`.

### master_data_value

| Column | Type | Rules |
|---|---|---|
| category | TEXT | `GENDER`, `MARITAL_STATUS`, or `OCCUPATION` |
| value | TEXT | The API value |
| label | TEXT | Display name |
| is_active | BOOLEAN | NOT NULL |

**PK**: `(category, value)`.

**Seeded by migration**: `MALE`, `FEMALE`; `SINGLE`, `MARRIED`, `DIVORCED`, `WIDOWED`; `EMPLOYEE`, `ENTREPRENEUR`, `CIVIL_SERVANT`. Values match DDE-001, and matching is case-sensitive. Only `GENDER` is meant to change without a code change.

## State Transitions

```text
                 submission recorded
                         │
                         ▼
                      PENDING ──internal failure──► PENDING (retry with backoff)
                     │       │
       all rules pass│       │any rule fails
                     ▼       ▼
                  PASSED   REJECTED
```

`PASSED` and `REJECTED` are final for this feature. Each transition happens in one transaction:

- **To PASSED**: upsert customer, insert photos, set `customer_id`, rewrite the payload, set dates.
- **To REJECTED**: delete the payload, set reasons and dates, insert the notification.

## Validation Rules

Rules run in `internal/validation`, with no I/O except the master data values passed in. Every rule runs, so all problems appear in one rejection. A field with the wrong JSON type is `INVALID`.

| Field (`customer.`) | Rule | Reason code | Requirement |
|---|---|---|---|
| Any mandatory field | Absent, null, or empty or whitespace-only text | `MISSING` | FR-002, FR-004 |
| `marital-status` | Exact accepted value | `INVALID` | FR-005 |
| `occupation` | Exact accepted value | `INVALID` | FR-006 |
| `date-of-birth` | Real date in `YYYY-MM-DD`, not after the submission date | `INVALID` | FR-015 |
| `date-of-birth` | Customer is 18 on the submission date. 29 February births reach 18 on 1 March in non-leap years | `CUSTOMER_UNDER_AGE` | FR-007 |
| `government-id-expiry-date` | Real date in `YYYY-MM-DD` | `INVALID` | FR-015 |
| `government-id-expiry-date` | Not before the submission date | `GOVERNMENT_ID_EXPIRED` | FR-008 |
| `national-id-number` | Letters and digits only | `INVALID` | FR-015 |
| `phone-number` | Optional `+`, then 8 to 15 digits | `INVALID` | FR-015 |
| `gender` | Active `GENDER` master data value | `INVALID` | FR-015 |
| `work-duration-years`, `number-of-dependents` | Whole number, 0 or more, at most 32,767 | `INVALID` | FR-015 |
| `monthly-income` | Greater than 0, plain decimal, up to 18 whole digits and 2 decimals | `INVALID` | FR-015 |
| `government-id-photo`, `selfie-photo` | Base64 JPEG or PNG that fully decodes, at most 5 MiB and 50 megapixels | `INVALID` | FR-015 |
| Any text field | No NUL characters, which PostgreSQL text can't store | `INVALID` | FR-015 |

`email` is never a rejection reason (FR-003). A non-text `email` is treated as absent.

The numeric limits and the NUL rule go slightly beyond FR-015. They stop a value that can't be stored from making a row fail forever.

## Differences from DDE-001

| DDE-001 | This model | Reason |
|---|---|---|
All differences were reviewed and resolved on 2026-10-31. Each one needs a matching change to DDE-001 or ERD-001.

| # | DDE-001 / ERD-001 | This model | Reason |
|---|---|---|---|
| 1 | `customer_id` mandatory on Loan Application | Nullable | Rejected records hold no customer data (FR-011) |
| 2 | `loan_amount_option_id`, `loan_tenure_id` mandatory | Left out; added as nullable by US-005 | A submission missing them must still be recorded (Principle III) |
| 3 | Application has exactly one Collateral | No row here; added by US-004 as one-to-zero-or-one | Same reason as 2 |
| 4 | `rejected_at` feeds the 30-day rule for every rejection | Customer data rejections are exempt | FR-016 |
| 5 | `application_status` values undefined | `PENDING`, `PASSED`, `REJECTED` | FR-014 |
| 6 | `submitted_at` only | Adds `acknowledged_at` | SC-001 |
| 7 | `submitted_at`, `rejected_at`, `issued_at` typed `Date` | `TIMESTAMPTZ` | A date can't measure 2 seconds |
| 8 | Photos are Text references, no photo table | Same field names as UUID references to `stored_photo` | Images are stored in PostgreSQL |
| 9 | No validation tracking fields | Adds `rejection_reasons`, `rejection_stage`, `validated_at`, `validation_attempts`, `next_validation_at` | FR-014, FR-016, SC-006, Principle III |
| 10 | No raw submission storage | Adds `application_payload` | Principle III |
| 11 | No notification data | Adds `rejection_notification`, plus `notification_url` and `notification_secret` on `partner` | FR-013 |
| 12 | No stored key field | Adds `key_hash` on `partner_credential` | Principle IV |
| 13 | Gender values undefined, no enumeration storage | Adds `master_data_value`, with `MALE` and `FEMALE` pending the business stakeholder's confirmation | FR-015, Principle V |
| 15 | Spouse, Collateral, Loan, vehicle master data, Loan Amount Option, Loan Tenure, Coverage Area defined | Not created | Deferred to US-003 to US-007 and later stories |

Item 14 was a plan gap rather than a difference: the `issued_at` rule is now enforced above.
