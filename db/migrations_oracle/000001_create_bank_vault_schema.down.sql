-- =============================
-- CPS Migration - Down
-- Rollback for bank_vault_products
-- =============================

-- Drop triggers
BEGIN
    EXECUTE IMMEDIATE 'DROP TRIGGER trg_bank_vault_products_id';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/

BEGIN
    EXECUTE IMMEDIATE 'DROP TRIGGER trg_bank_vault_products_updated_at';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/

-- Drop table
BEGIN
    EXECUTE IMMEDIATE 'DROP TABLE bank_vault_products CASCADE CONSTRAINTS';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/

-- Drop sequence
BEGIN
    EXECUTE IMMEDIATE 'DROP SEQUENCE bank_vault_products_seq';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/
