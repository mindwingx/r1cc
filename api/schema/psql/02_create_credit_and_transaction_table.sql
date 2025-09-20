CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate Up
CREATE TABLE IF NOT EXISTS credits (
    id          SERIAL PRIMARY KEY,
    uuid        UUID DEFAULT uuid_generate_v4() NOT NULL,
    tenant_id   INTEGER NOT NULL UNIQUE,
    balance     NUMERIC(20, 4) NOT NULL CHECK (balance >= 0),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_credits_tenant ON credits(tenant_id);


CREATE TABLE IF NOT EXISTS credit_transactions (
    id              BYTEA PRIMARY KEY,
    credit_id       INTEGER NOT NULL,
    amount          NUMERIC(20, 4) NOT NULL,
    message_hash_id BYTEA NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (credit_id) REFERENCES credits(id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_credit_transactions_tenant ON credit_transactions(credit_id);
CREATE INDEX IF NOT EXISTS idx_credit_transactions_created ON credit_transactions(created_at);


-- +migrate Down