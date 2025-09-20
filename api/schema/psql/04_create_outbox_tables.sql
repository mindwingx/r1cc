CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- Create enum types only if they don't exist
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'outbox_status') THEN
             CREATE TYPE outbox_status AS ENUM ('pending', 'publishing', 'published', 'failed');
    END IF;
END$$;


-- +migrate Up
CREATE TABLE IF NOT EXISTS outboxes (
    id            SERIAL PRIMARY KEY,
    uuid          UUID DEFAULT uuid_generate_v4() NOT NULL,
    event_type    VARCHAR(255) NOT NULL,
    message_id    INTEGER NOT NULL,
    payload       JSONB NOT NULL,
    status        outbox_status NOT NULL DEFAULT 'pending',
    retries       INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    retry_at      TIMESTAMP NULL,
    deleted_at    TIMESTAMP NULL,
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_msg_hash ON messages (message_hash);