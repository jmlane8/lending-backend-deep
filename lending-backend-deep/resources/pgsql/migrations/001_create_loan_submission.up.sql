CREATE TABLE partner (
    partner_id            UUID        NOT NULL,
    agreement_status      TEXT        NOT NULL CHECK (agreement_status IN ('ACTIVE', 'INACTIVE')),
    agreement_signed_date DATE        NOT NULL,
    notification_url      TEXT,
    notification_secret   TEXT,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted               BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at            TIMESTAMPTZ DEFAULT NULL,
    created_by            TEXT        NOT NULL DEFAULT 'system',
    updated_by            TEXT        NOT NULL DEFAULT 'system',
    deleted_by            TEXT        DEFAULT NULL,
    CONSTRAINT partner_pkey PRIMARY KEY (partner_id),
    CONSTRAINT partner_notification_secret_check
        CHECK (notification_url IS NULL OR notification_secret IS NOT NULL)
);

CREATE TABLE partner_credential (
    partner_credential_id UUID        NOT NULL,
    partner_id            UUID        NOT NULL,
    key_hash              BYTEA       NOT NULL,
    is_valid              BOOLEAN     NOT NULL,
    issued_at             TIMESTAMPTZ NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted               BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at            TIMESTAMPTZ DEFAULT NULL,
    created_by            TEXT        NOT NULL DEFAULT 'system',
    updated_by            TEXT        NOT NULL DEFAULT 'system',
    deleted_by            TEXT        DEFAULT NULL,
    CONSTRAINT partner_credential_pkey PRIMARY KEY (partner_credential_id),
    CONSTRAINT partner_credential_key_hash_key UNIQUE (key_hash),
    CONSTRAINT partner_credential_partner_fkey
        FOREIGN KEY (partner_id) REFERENCES partner (partner_id)
);

-- CHECK constraints can't read another table, so a trigger enforces this rule.
-- A credential issued later on the signing day passes.
CREATE FUNCTION partner_credential_check_issued_at() RETURNS trigger AS $$
BEGIN
    IF NEW.issued_at <= (SELECT agreement_signed_date
                         FROM partner
                         WHERE partner_id = NEW.partner_id) THEN
        RAISE EXCEPTION 'issued_at must be after agreement_signed_date';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER partner_credential_issued_at_trigger
    BEFORE INSERT OR UPDATE ON partner_credential
    FOR EACH ROW EXECUTE FUNCTION partner_credential_check_issued_at();

CREATE TABLE stored_photo (
    photo_id     UUID        NOT NULL,
    content_type TEXT        NOT NULL CHECK (content_type IN ('image/jpeg', 'image/png')),
    content      BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted      BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at   TIMESTAMPTZ DEFAULT NULL,
    created_by   TEXT        NOT NULL DEFAULT 'system',
    updated_by   TEXT        NOT NULL DEFAULT 'system',
    deleted_by   TEXT        DEFAULT NULL,
    CONSTRAINT stored_photo_pkey PRIMARY KEY (photo_id)
);

CREATE TABLE customer (
    customer_id               UUID           NOT NULL,
    national_id_number        TEXT           NOT NULL,
    full_legal_name           TEXT           NOT NULL,
    government_id_photo       UUID           NOT NULL,
    government_id_expiry_date DATE           NOT NULL,
    selfie_photo              UUID           NOT NULL,
    phone_number              TEXT           NOT NULL,
    gender                    TEXT           NOT NULL,
    date_of_birth             DATE           NOT NULL,
    place_of_birth            TEXT           NOT NULL,
    marital_status            TEXT           NOT NULL,
    occupation                TEXT           NOT NULL,
    work_duration_years       SMALLINT       NOT NULL,
    monthly_income            NUMERIC(20, 2) NOT NULL,
    number_of_dependents      SMALLINT       NOT NULL,
    email                     TEXT,
    domicile_street           TEXT           NOT NULL,
    domicile_district         TEXT           NOT NULL,
    domicile_city             TEXT           NOT NULL,
    domicile_province         TEXT           NOT NULL,
    created_at                TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted                   BOOLEAN        NOT NULL DEFAULT FALSE,
    deleted_at                TIMESTAMPTZ    DEFAULT NULL,
    created_by                TEXT           NOT NULL DEFAULT 'system',
    updated_by                TEXT           NOT NULL DEFAULT 'system',
    deleted_by                TEXT           DEFAULT NULL,
    CONSTRAINT customer_pkey PRIMARY KEY (customer_id),
    CONSTRAINT customer_national_id_number_key UNIQUE (national_id_number),
    CONSTRAINT customer_government_id_photo_fkey
        FOREIGN KEY (government_id_photo) REFERENCES stored_photo (photo_id),
    CONSTRAINT customer_selfie_photo_fkey
        FOREIGN KEY (selfie_photo) REFERENCES stored_photo (photo_id)
);

CREATE TABLE loan_application (
    loan_application_id   UUID        NOT NULL,
    application_reference TEXT        NOT NULL,
    partner_id            UUID        NOT NULL,
    customer_id           UUID,
    application_status    TEXT        NOT NULL DEFAULT 'PENDING'
        CHECK (application_status IN ('PENDING', 'PASSED', 'REJECTED')),
    rejection_reasons     JSONB,
    submitted_at          TIMESTAMPTZ NOT NULL,
    acknowledged_at       TIMESTAMPTZ NOT NULL,
    validated_at          TIMESTAMPTZ,
    rejected_at           TIMESTAMPTZ,
    rejection_stage       TEXT,
    validation_attempts   INTEGER     NOT NULL DEFAULT 0,
    next_validation_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted               BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at            TIMESTAMPTZ DEFAULT NULL,
    created_by            TEXT        NOT NULL DEFAULT 'system',
    updated_by            TEXT        NOT NULL DEFAULT 'system',
    deleted_by            TEXT        DEFAULT NULL,
    CONSTRAINT loan_application_pkey PRIMARY KEY (loan_application_id),
    CONSTRAINT loan_application_reference_key UNIQUE (application_reference),
    CONSTRAINT loan_application_partner_fkey
        FOREIGN KEY (partner_id) REFERENCES partner (partner_id),
    CONSTRAINT loan_application_customer_fkey
        FOREIGN KEY (customer_id) REFERENCES customer (customer_id)
);

CREATE INDEX loan_application_pending_idx
    ON loan_application (next_validation_at)
    WHERE application_status = 'PENDING';

CREATE TABLE application_payload (
    loan_application_id UUID        NOT NULL,
    body                BYTEA       NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted             BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at          TIMESTAMPTZ DEFAULT NULL,
    created_by          TEXT        NOT NULL DEFAULT 'system',
    updated_by          TEXT        NOT NULL DEFAULT 'system',
    deleted_by          TEXT        DEFAULT NULL,
    CONSTRAINT application_payload_pkey PRIMARY KEY (loan_application_id),
    CONSTRAINT application_payload_application_fkey
        FOREIGN KEY (loan_application_id) REFERENCES loan_application (loan_application_id)
        ON DELETE CASCADE
);

CREATE TABLE rejection_notification (
    event_id            UUID        NOT NULL,
    loan_application_id UUID        NOT NULL,
    status              TEXT        NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'DELIVERED', 'EXHAUSTED')),
    attempts            INTEGER     NOT NULL DEFAULT 0,
    next_attempt_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    delivered_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted             BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at          TIMESTAMPTZ DEFAULT NULL,
    created_by          TEXT        NOT NULL DEFAULT 'system',
    updated_by          TEXT        NOT NULL DEFAULT 'system',
    deleted_by          TEXT        DEFAULT NULL,
    CONSTRAINT rejection_notification_pkey PRIMARY KEY (event_id),
    CONSTRAINT rejection_notification_application_key UNIQUE (loan_application_id),
    CONSTRAINT rejection_notification_application_fkey
        FOREIGN KEY (loan_application_id) REFERENCES loan_application (loan_application_id)
);

CREATE INDEX rejection_notification_pending_idx
    ON rejection_notification (next_attempt_at)
    WHERE status = 'PENDING';

CREATE TABLE master_data_value (
    category   TEXT        NOT NULL CHECK (category IN ('GENDER', 'MARITAL_STATUS', 'OCCUPATION')),
    value      TEXT        NOT NULL,
    label      TEXT        NOT NULL,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted    BOOLEAN     NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_by TEXT        NOT NULL DEFAULT 'system',
    updated_by TEXT        NOT NULL DEFAULT 'system',
    deleted_by TEXT        DEFAULT NULL,
    CONSTRAINT master_data_value_pkey PRIMARY KEY (category, value)
);

INSERT INTO master_data_value (category, value, label) VALUES
    ('GENDER', 'MALE', 'Male'),
    ('GENDER', 'FEMALE', 'Female'),
    ('MARITAL_STATUS', 'SINGLE', 'Single'),
    ('MARITAL_STATUS', 'MARRIED', 'Married'),
    ('MARITAL_STATUS', 'DIVORCED', 'Divorced'),
    ('MARITAL_STATUS', 'WIDOWED', 'Widowed'),
    ('OCCUPATION', 'EMPLOYEE', 'Employee'),
    ('OCCUPATION', 'ENTREPRENEUR', 'Entrepreneur'),
    ('OCCUPATION', 'CIVIL_SERVANT', 'Civil Servant');
