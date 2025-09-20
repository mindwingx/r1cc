CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate Up
CREATE TABLE IF NOT EXISTS tenants (
    id          SERIAL PRIMARY KEY,
    uuid        UUID DEFAULT uuid_generate_v4() NOT NULL,
    username    VARCHAR(255) NOT NULL UNIQUE,
    tenant_name VARCHAR(255) NOT NULL,
    active      BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_tenant_username ON tenants (username);

-- +migrate Down
