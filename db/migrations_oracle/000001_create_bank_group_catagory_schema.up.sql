-- Drop table if exists
BEGIN
   EXECUTE IMMEDIATE 'DROP TABLE group_vault_categories CASCADE CONSTRAINTS';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -942 THEN RAISE; END IF; -- ORA-00942: table or view does not exist
END;
/

-- Drop triggers if they exist
BEGIN
   EXECUTE IMMEDIATE 'DROP TRIGGER trg_group_vault_categories_updated_at';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -4080 THEN RAISE; END IF; -- ORA-04080: trigger does not exist
END;
/

BEGIN
   EXECUTE IMMEDIATE 'DROP TRIGGER trg_group_vault_categories_insert';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -4080 THEN RAISE; END IF;
END;
/

-- Drop index if exists
BEGIN
   EXECUTE IMMEDIATE 'DROP INDEX idx_group_vault_categories_is_active';
EXCEPTION
   WHEN OTHERS THEN
      IF SQLCODE != -1418 THEN RAISE; END IF; -- ORA-01418: specified index does not exist
END;
/

-- Create group_vault_categories reference table
CREATE TABLE group_vault_categories (
    id RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
    name VARCHAR2(255) NOT NULL,
    description VARCHAR2(1000),
    is_active NUMBER(1) DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    is_deleted NUMBER(1) DEFAULT 0 NOT NULL,
    CONSTRAINT uq_group_vault_categories_name UNIQUE (name),
    CONSTRAINT chk_group_vault_categories_is_active CHECK (is_active IN (0,1)),
    CONSTRAINT chk_group_vault_categories_is_deleted CHECK (is_deleted IN (0,1))
);
/

-- Optional documentation
COMMENT ON COLUMN group_vault_categories.deleted_at IS 'Used for soft deletes';
/

-- Index
CREATE INDEX idx_group_vault_categories_is_active 
    ON group_vault_categories(is_active);
/

-- Trigger to maintain updated_at
CREATE OR REPLACE TRIGGER trg_group_vault_categories_updated_at
  BEFORE UPDATE ON group_vault_categories
  FOR EACH ROW
BEGIN
  :NEW.updated_at := SYSTIMESTAMP;
END;
/

-- Trigger to set timestamps on insert
CREATE OR REPLACE TRIGGER trg_group_vault_categories_insert
  BEFORE INSERT ON group_vault_categories
  FOR EACH ROW
BEGIN
  IF :NEW.created_at IS NULL THEN
    :NEW.created_at := SYSTIMESTAMP;
  END IF;
  :NEW.updated_at := SYSTIMESTAMP;
END;
/
