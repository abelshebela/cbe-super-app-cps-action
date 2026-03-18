CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS banks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    bank_name VARCHAR(32) NOT NULL,
    logo VARCHAR(255) NOT NULL,
    bic_code VARCHAR(16) NOT NULL,
    is_enabled INT DEFAULT 1,
    account_length INT NOT NULL,
    has_alpha_numeric INT DEFAULT 0,
    create_at TIMESTAMP DEFAULT now(),
    update_at TIMESTAMP DEFAULT now()
);