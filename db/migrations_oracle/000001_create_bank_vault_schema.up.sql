-- =============================
-- CPS Migration - Up
-- Only bank_vault_products table
-- =============================

-- ===========================================
-- 1️⃣ Drop any existing triggers
-- ===========================================
BEGIN
    FOR t IN (
        SELECT trigger_name
        FROM user_triggers
        WHERE trigger_name IN ('TRG_BANK_VAULT_PRODUCTS_ID', 'TRG_BANK_VAULT_PRODUCTS_UPDATED_AT')
    ) LOOP
        BEGIN
            EXECUTE IMMEDIATE 'DROP TRIGGER ' || t.trigger_name;
        EXCEPTION
            WHEN OTHERS THEN
                DBMS_OUTPUT.PUT_LINE('Skipping drop trigger ' || t.trigger_name || ': ' || SQLERRM);
        END;
    END LOOP;
END;
/

-- ===========================================
-- 2️⃣ Drop table if exists
-- ===========================================
BEGIN
    BEGIN
        EXECUTE IMMEDIATE 'DROP TABLE bank_vault_products CASCADE CONSTRAINTS';
    EXCEPTION
        WHEN OTHERS THEN
            IF SQLCODE != -942 THEN -- ORA-00942: table does not exist
                RAISE;
            END IF;
    END;
END;
/

-- ===========================================
-- 3️⃣ Drop sequence if exists
-- ===========================================
BEGIN
    FOR seq IN (
        SELECT sequence_name
        FROM user_sequences
        WHERE sequence_name = 'BANK_VAULT_PRODUCTS_SEQ'
    ) LOOP
        BEGIN
            EXECUTE IMMEDIATE 'DROP SEQUENCE ' || seq.sequence_name;
        EXCEPTION
            WHEN OTHERS THEN
                DBMS_OUTPUT.PUT_LINE('Skipping drop sequence ' || seq.sequence_name || ': ' || SQLERRM);
        END;
    END LOOP;
END;
/

-- ===========================================
-- 4️⃣ Create sequence
-- ===========================================
CREATE SEQUENCE bank_vault_products_seq
    START WITH 1
    INCREMENT BY 1;
/

-- ===========================================
-- 5️⃣ Create table with VARCHAR2(36) id
-- ===========================================
CREATE TABLE bank_vault_products (
    id                    VARCHAR2(36) NOT NULL,
    name                  VARCHAR2(500) NOT NULL,
    description           VARCHAR2(2000) DEFAULT '',
    currency              VARCHAR2(3) NOT NULL CHECK (LENGTH(currency)=3 AND currency=UPPER(currency)),
    rate_bps              NUMBER(19,4) NOT NULL CHECK (rate_bps >= 0),
    method                VARCHAR2(20) NOT NULL CHECK (method IN ('COMPOUND', 'SIMPLE')),
    frequency             VARCHAR2(20) NOT NULL CHECK (frequency IN ('DAILY', 'MONTHLY', 'QUARTERLY', 'ANNUALLY')),
    lock_period           NUMBER(19,0) NOT NULL CHECK (lock_period > 0),
    min_amount            NUMBER(19,4) NOT NULL CHECK (min_amount > 0),
    max_amount            NUMBER(19,4) NOT NULL CHECK (max_amount > 0),
    early_unlock_rate_bps NUMBER(1) DEFAULT 0 CHECK (early_unlock_rate_bps IN (0,1)),
    is_active             NUMBER(1) DEFAULT 0 CHECK (is_active IN (0,1)),
    created_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP WITH TIME ZONE,
    is_deleted            NUMBER(1) DEFAULT 0 NOT NULL,
    CONSTRAINT pk_bank_vault_products PRIMARY KEY (id),
    CONSTRAINT chk_amount_range CHECK (min_amount <= max_amount),
    CONSTRAINT chk_is_deleted CHECK (is_deleted IN (0, 1))
);
/

-- ===========================================
-- 7️⃣ Trigger to auto-update updated_at column
-- ===========================================
CREATE OR REPLACE TRIGGER trg_bank_vault_products_updated_at
BEFORE UPDATE ON bank_vault_products
FOR EACH ROW
BEGIN
    :NEW.updated_at := CURRENT_TIMESTAMP;
END;
/
