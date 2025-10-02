-- Oracle Migration: Drop Group Vault Categories Schema

-- Drop triggers
BEGIN
  EXECUTE IMMEDIATE 'DROP TRIGGER trg_group_vault_categories_updated_at';
EXCEPTION WHEN OTHERS THEN IF SQLCODE != -4080 THEN RAISE; END IF; -- trigger does not exist
END;
/

BEGIN
  EXECUTE IMMEDIATE 'DROP TRIGGER trg_group_vault_categories_insert';
EXCEPTION WHEN OTHERS THEN IF SQLCODE != -4080 THEN RAISE; END IF; -- trigger does not exist
END;
/

-- Drop indexes (safe-guarded)
BEGIN
  EXECUTE IMMEDIATE 'DROP INDEX idx_group_vault_categories_is_active';
EXCEPTION WHEN OTHERS THEN IF SQLCODE != -1418 AND SQLCODE != -942 THEN RAISE; END IF;
END;
/

BEGIN
  EXECUTE IMMEDIATE 'DROP INDEX idx_group_vault_categories_display_order';
EXCEPTION WHEN OTHERS THEN IF SQLCODE != -1418 AND SQLCODE != -942 THEN RAISE; END IF;
END;
/

-- Drop table
BEGIN
   EXECUTE IMMEDIATE 'DROP TABLE group_vault_categories CASCADE CONSTRAINTS';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -942 THEN RAISE; END IF;
END;
/
