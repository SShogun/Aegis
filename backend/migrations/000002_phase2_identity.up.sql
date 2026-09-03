-- Phase 2 Identity: Internal user records.
-- Authentication providers establish external identity;
-- Aegis maintains its own internal user record.

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(320) NOT NULL,
    name VARCHAR(255),

    identity_provider VARCHAR(100) NOT NULL,
    provider_subject VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_provider_subject_unique
        UNIQUE (identity_provider, provider_subject)
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique
    ON users (LOWER(email));
