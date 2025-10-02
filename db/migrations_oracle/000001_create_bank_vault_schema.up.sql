-- =============================
-- CPS Migration - Up
-- Only bank_vault_products table
-- =============================

-- Drop existing sequence (ignore error if it doesn't exist)
BEGIN
    EXECUTE IMMEDIATE 'DROP SEQUENCE bank_vault_products_seq';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/

-- Drop existing table (ignore error if it doesn't exist)
BEGIN
    EXECUTE IMMEDIATE 'DROP TABLE bank_vault_products CASCADE CONSTRAINTS';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/

-- Create sequence for IDs (if you ever need numeric PKs)
CREATE SEQUENCE bank_vault_products_seq START WITH 1 INCREMENT BY 1;

-- Create parent table: bank_vault_products
CREATE TABLE bank_vault_products (
    id                   VARCHAR2(36) NOT NULL,
    name                 VARCHAR2(500)   NOT NULL,
    description          VARCHAR2(2000)  DEFAULT '',
    currency             VARCHAR2(3)     NOT NULL CHECK (LENGTH(currency) = 3 AND currency = UPPER(currency)),
    rate_bps             NUMBER(19,4)    NOT NULL CHECK (rate_bps >= 0),
    method               VARCHAR2(20)    NOT NULL CHECK (method IN ('COMPOUND', 'SIMPLE')),
    frequency            VARCHAR2(20)    NOT NULL CHECK (frequency IN ('DAILY', 'MONTHLY', 'QUARTERLY', 'ANNUALLY')),
    lock_period          NUMBER(19,0)    NOT NULL CHECK (lock_period > 0),
    min_amount           NUMBER(19,4)    NOT NULL CHECK (min_amount > 0),
    max_amount           NUMBER(19,4)    NOT NULL CHECK (max_amount > 0),
    early_unlock_fee_bps NUMBER(19,4)    DEFAULT 0 CHECK (early_unlock_fee_bps >= 0),
    is_active            NUMBER(1)       DEFAULT 0 CHECK (is_active IN (0,1)),
    is_deleted           NUMBER(1) DEFAULT 0 NOT NULL,
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at           TIMESTAMP WITH TIME ZONE,
    CONSTRAINT pk_bank_vault_products PRIMARY KEY (id),
    CONSTRAINT chk_amount_range CHECK (min_amount <= max_amount)
);

-- Trigger to auto-generate GUIDs
CREATE OR REPLACE TRIGGER trg_bank_vault_products_id
    BEFORE INSERT ON bank_vault_products
    FOR EACH ROW
    WHEN (NEW.id IS NULL)
BEGIN
    :NEW.id := SYS_GUID();
END;
/

-- Trigger to auto-update updated_at column
CREATE OR REPLACE TRIGGER trg_bank_vault_products_updated_at
    BEFORE UPDATE ON bank_vault_products
    FOR EACH ROW
BEGIN
    :NEW.updated_at := CURRENT_TIMESTAMP;
END;
/
