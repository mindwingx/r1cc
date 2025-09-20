CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- Create enum types only if they don't exist
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_status') THEN
             CREATE TYPE message_status AS ENUM ('queued','sending','sent','delivered','failed');
    END IF;
END$$;

-- +migrate Up
CREATE TABLE IF NOT EXISTS messages (
    id           SERIAL PRIMARY KEY,
    uuid UUID    DEFAULT uuid_generate_v4() NOT NULL,
    tenant_id    INTEGER NOT NULL,
    mobile       VARCHAR(255) NOT NULL,
    message_text TEXT NOT NULL,
    message_hash BYTEA NOT NULL UNIQUE,
    status       message_status NOT NULL DEFAULT 'queued',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_msg_hash ON messages (message_hash);

-- +migrate Down