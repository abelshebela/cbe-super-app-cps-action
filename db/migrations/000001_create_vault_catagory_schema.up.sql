-- ===========================================
-- Drop ANY existing object named vault_categories
-- ===========================================

DECLARE
  v_sql VARCHAR2(4000);
BEGIN
  FOR obj IN (
    SELECT object_name, object_type
    FROM user_objects
    WHERE object_name = 'vault_categories'
  ) LOOP
    v_sql := 'DROP ' || obj.object_type || ' ' || obj.object_name;

    -- tables require CASCADE
    IF obj.object_type = 'TABLE' THEN
      v_sql := v_sql || ' CASCADE CONSTRAINTS';
    END IF;

    EXECUTE IMMEDIATE v_sql;
  END LOOP;
END;
/
 
-- ===========================================
-- Drop ANY existing object named idx_vault_categories_is_active
-- ===========================================

DECLARE
  v_sql VARCHAR2(4000);
BEGIN
  FOR obj IN (
    SELECT object_name, object_type
    FROM user_objects
    WHERE object_name = 'IDX_vault_categories_IS_ACTIVE'
  ) LOOP
    EXECUTE IMMEDIATE 'DROP INDEX ' || obj.object_name;
  END LOOP;
END;
/

-- ===========================================
-- Create Table
-- ===========================================

CREATE TABLE vault_categories (
    id              VARCHAR2(36) NOT NULL PRIMARY KEY,
    name            VARCHAR2(255) NOT NULL,
    interest        NUMBER(19,4) NOT NULL CHECK (interest >= 0),
    category_type   VARCHAR2(20) NOT NULL CHECK (method IN ('PERSONAL', 'GROUP')),
    cover_image     VARCHAR2(255),
    is_active       NUMBER(1) DEFAULT 0 NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    deleted_at      TIMESTAMP WITH TIME ZONE,
    is_deleted      NUMBER(1) DEFAULT 0 NOT NULL,
    CONSTRAINT uq_vault_categories_name UNIQUE (name),
    CONSTRAINT chk_vault_categories_is_active CHECK (is_active IN (0,1)),
    CONSTRAINT chk_vault_categories_is_deleted CHECK (is_deleted IN (0,1))
);
/

COMMENT ON COLUMN vault_categories.deleted_at IS 'Used for soft deletes';
/

-- ===========================================
-- Recreate Index
-- ===========================================

CREATE INDEX idx_vault_categories_is_active
    ON vault_categories(is_active);
/

-- ===========================================
-- Triggers
-- ===========================================

CREATE OR REPLACE TRIGGER trg_vault_categories_updated_at
  BEFORE UPDATE ON vault_categories
  FOR EACH ROW
BEGIN
  :NEW.updated_at := SYSTIMESTAMP;
END;
/

CREATE OR REPLACE TRIGGER trg_vault_categories_insert
  BEFORE INSERT ON vault_categories
  FOR EACH ROW
BEGIN
  IF :NEW.created_at IS NULL THEN
    :NEW.created_at := SYSTIMESTAMP;
  END IF;
  :NEW.updated_at := SYSTIMESTAMP;
END;
/
