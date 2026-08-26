-- Phase 1 Foundation: Minimal schema to prove migration execution.
-- Avoid adding users, organizations, etc. unless explicitly required.
-- Later phases will introduce core identity and operational tables.

CREATE TABLE IF NOT EXISTS schema_migrations_history (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Note: In a typical setup, golang-migrate manages its own schema_migrations table.
-- We are just creating a simple table here to prove standard table creation works 
-- and can be rolled back.
CREATE TABLE IF NOT EXISTS system_settings (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
