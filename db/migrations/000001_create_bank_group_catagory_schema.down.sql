-- ===========================================
-- Oracle Migration: Drop Group Vault Categories Schema
-- ===========================================

-- Drop triggers
BEGIN
    EXECUTE IMMEDIATE 'DROP TRIGGER trg_group_vault_categories_updated_at';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -4080 THEN  -- ORA-04080: trigger does not exist
            RAISE;
        END IF;
END;
/

BEGIN
    EXECUTE IMMEDIATE 'DROP TRIGGER trg_group_vault_categories_insert';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -4080 THEN  -- ORA-04080: trigger does not exist
            RAISE;
        END IF;
END;
/

-- Drop indexes
BEGIN
    EXECUTE IMMEDIATE 'DROP INDEX idx_group_vault_categories_is_active';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE NOT IN (-1418, -942) THEN  -- ORA-01418 or ORA-00942
            RAISE;
        END IF;
END;
/

BEGIN
    EXECUTE IMMEDIATE 'DROP INDEX idx_group_vault_categories_display_order';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE NOT IN (-1418, -942) THEN  -- ORA-01418 or ORA-00942
            RAISE;
        END IF;
END;
/

-- Drop table
BEGIN
    EXECUTE IMMEDIATE 'DROP TABLE group_vault_categories CASCADE CONSTRAINTS';
EXCEPTION
    WHEN OTHERS THEN
        IF SQLCODE != -942 THEN  -- ORA-00942: table does not exist
            RAISE;
        END IF;
END;
/
