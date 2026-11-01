-- Local and sandbox only. Never run this in production.
-- Key hashes are the SHA-256 of the key shown in each comment.

INSERT INTO partner (partner_id, agreement_status, agreement_signed_date) VALUES
    ('01900000-0000-7000-8000-000000000001', 'ACTIVE', DATE '2026-01-01'),
    ('01900000-0000-7000-8000-000000000002', 'ACTIVE', DATE '2026-01-01'),
    ('01900000-0000-7000-8000-000000000003', 'INACTIVE', DATE '2026-01-01')
ON CONFLICT (partner_id) DO NOTHING;

INSERT INTO partner_credential (partner_credential_id, partner_id, key_hash, is_valid, issued_at) VALUES
    -- sandbox-partner-key
    ('01900000-0000-7000-8000-000000000011', '01900000-0000-7000-8000-000000000001',
     sha256('sandbox-partner-key'::bytea), TRUE, TIMESTAMPTZ '2026-01-02 00:00:00+00'),
    -- sandbox-partner-2-key
    ('01900000-0000-7000-8000-000000000012', '01900000-0000-7000-8000-000000000002',
     sha256('sandbox-partner-2-key'::bytea), TRUE, TIMESTAMPTZ '2026-01-02 00:00:00+00'),
    -- inactive-partner-key
    ('01900000-0000-7000-8000-000000000013', '01900000-0000-7000-8000-000000000003',
     sha256('inactive-partner-key'::bytea), TRUE, TIMESTAMPTZ '2026-01-02 00:00:00+00')
ON CONFLICT (partner_credential_id) DO NOTHING;
