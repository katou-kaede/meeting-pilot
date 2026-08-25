CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,

    password_hash VARCHAR(255) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_users_name_not_blank
        CHECK (BTRIM(name) <> ''),

    CONSTRAINT chk_users_email_not_blank
        CHECK (BTRIM(email) <> '')
);

CREATE UNIQUE INDEX uq_users_email_lower
ON users (LOWER(email));