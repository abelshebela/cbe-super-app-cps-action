-- =============================
-- CPS Migration - Down
-- Rollback for bank_vault_products
-- =============================

-- Drop triggers safely
BEGIN
    EXECUTE IMMEDIATE 'DROP TRIGGER trg_bank_vault_products_id';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -4080 THEN -- ORA-04080: trigger does not exist
            RAISE;
        END IF;
END;
/

BEGIN
    EXECUTE IMMEDIATE 'DROP TRIGGER trg_bank_vault_products_updated_at';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -4080 THEN
            RAISE;
        END IF;
END;
/

-- Drop table safely
BEGIN
    EXECUTE IMMEDIATE 'DROP TABLE bank_vault_products CASCADE CONSTRAINTS';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -942 THEN -- ORA-00942: table or view does not exist
            RAISE;
        END IF;
END;
/

-- Drop sequence safely
BEGIN
    EXECUTE IMMEDIATE 'DROP SEQUENCE bank_vault_products_seq';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -2289 THEN -- ORA-02289: sequence does not exist
            RAISE;
        END IF;
END;
/
