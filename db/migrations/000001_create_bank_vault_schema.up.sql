-- Drop sequence if exists
BEGIN
    BEGIN
        EXECUTE IMMEDIATE 'DROP SEQUENCE bank_vault_products_seq';
    EXCEPTION
        WHEN OTHERS THEN
            IF SQLCODE != -2289 THEN -- ORA-02289: sequence does not exist
                RAISE;
            END IF;
    END;
END;
/

-- Create sequence
CREATE SEQUENCE bank_vault_products_seq
    START WITH 1
    INCREMENT BY 1;
/

-- Create table
CREATE TABLE bank_vault_products (
    id                    VARCHAR2(36) NOT NULL,
    name                  VARCHAR2(500) NOT NULL,
    currency              VARCHAR2(3) NOT NULL CHECK (LENGTH(currency)=3 AND currency=UPPER(currency)),
    rate_bps              NUMBER(19,4) NOT NULL CHECK (rate_bps >= 0),
    method                VARCHAR2(20) NOT NULL CHECK (method IN ('COMPOUND', 'SIMPLE')),
    frequency             NUMBER(3) NOT NULL CHECK (frequency > 0 AND frequency <= 365),
    lock_period           NUMBER(19,0) NOT NULL CHECK (lock_period > 0),
    min_amount            NUMBER(19,4) NOT NULL CHECK (min_amount > 0),
    max_amount            NUMBER(19,4) NOT NULL CHECK (max_amount > 0),
    apply_interest_on_early_unlock NUMBER(1) DEFAULT 0 CHECK (apply_interest_on_early_unlock IN (0,1)),
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

-- Trigger for updated_at
CREATE OR REPLACE TRIGGER trg_bank_vault_products_updated_at
BEFORE UPDATE ON bank_vault_products
FOR EACH ROW
BEGIN
    :NEW.updated_at := CURRENT_TIMESTAMP;
END;
/
