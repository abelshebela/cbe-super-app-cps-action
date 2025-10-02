CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Enums
CREATE TYPE accrual_method AS ENUM ('COMPOUND', 'SIMPLE');
CREATE TYPE accrual_frequency AS ENUM ('DAILY', 'MONTHLY', 'QUARTERLY', 'ANNUALLY');
CREATE TYPE locked_vault_status AS ENUM ('ACTIVE','MATURED','UNLOCKED_EARLY','PAID_OUT');
CREATE TYPE transaction_type AS ENUM ('PAYOUT','UNLOCK');
CREATE TYPE receipt_kind AS ENUM ('LOCK','UNLOCK','PAYOUT');

-- Parent table: bank_vault_products
CREATE TABLE bank_vault_products (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                 TEXT            NOT NULL,
    description          TEXT            NOT NULL DEFAULT '',
    currency             TEXT            NOT NULL CHECK (char_length(currency) = 3 AND currency = upper(currency)),
    rate_bps             NUMERIC(19,4)   NOT NULL CHECK (rate_bps >= 0),
    method               accrual_method  NOT NULL,
    frequency            accrual_frequency NOT NULL,
    lock_period          BIGINT          NOT NULL CHECK (lock_period > 0),
    min_amount           NUMERIC(19,4)   NOT NULL CHECK (min_amount > 0),
    max_amount           NUMERIC(19,4)   NOT NULL CHECK (max_amount > 0),
    early_unlock_fee_bps NUMERIC(19,4)   NOT NULL DEFAULT 0 CHECK (early_unlock_fee_bps >= 0),
    is_active            BOOLEAN         NOT NULL DEFAULT FALSE,
    is_deleted           BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ,
    CONSTRAINT chk_amount_range CHECK (min_amount <= max_amount)
);