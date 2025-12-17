-- ===========================================
-- Drop ANY existing object named VAULT_CATEGORIES (table, view, index, etc.)
-- ===========================================
DECLARE
  v_sql VARCHAR2(4000);
BEGIN
  FOR obj IN (
    SELECT object_name, object_type
    FROM user_objects
    WHERE object_name = 'VAULT_CATEGORIES'
  ) LOOP
    v_sql := 'DROP ' || obj.object_type || ' ' || obj.object_name;
    IF obj.object_type = 'TABLE' THEN
      v_sql := v_sql || ' CASCADE CONSTRAINTS';
    END IF;
    EXECUTE IMMEDIATE v_sql;
  END LOOP;
END;
/
 
-- ===========================================
-- Drop ANY existing index IDX_VAULT_CATEGORIES_IS_ACTIVE
-- ===========================================
DECLARE
  v_sql VARCHAR2(4000);
BEGIN
  FOR obj IN (
    SELECT index_name AS object_name
    FROM user_indexes
    WHERE index_name = 'IDX_VAULT_CATEGORIES_IS_ACTIVE'
  ) LOOP
    EXECUTE IMMEDIATE 'DROP INDEX ' || obj.object_name;
  END LOOP;
END;
/

-- ===========================================
-- Drop ANY existing object named AMOUNT_TIER
-- ===========================================
DECLARE
  v_sql VARCHAR2(4000);
BEGIN
  FOR obj IN (
    SELECT object_name, object_type
    FROM user_objects
    WHERE object_name = 'AMOUNT_TIER'
  ) LOOP
    v_sql := 'DROP ' || obj.object_type || ' ' || obj.object_name;
    IF obj.object_type = 'TABLE' THEN
      v_sql := v_sql || ' CASCADE CONSTRAINTS';
    END IF;
    EXECUTE IMMEDIATE v_sql;
  END LOOP;
END;
/

-- ===========================================
-- Recreate vault_categories
-- ===========================================
CREATE TABLE vault_categories (
    id              VARCHAR2(36) NOT NULL PRIMARY KEY,
    name            VARCHAR2(255) NOT NULL,
    category_type   VARCHAR2(20) NOT NULL CHECK (category_type IN ('PERSONAL', 'GROUP')),
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

-- ===========================================
-- Recreate amount_tier
-- ===========================================
CREATE TABLE amount_tier (
    id                  VARCHAR2(36) NOT NULL PRIMARY KEY,
    vault_category_id   VARCHAR2(36) NOT NULL,
    min_amount          NUMBER(19,4) NOT NULL,
    max_amount          NUMBER(19,4) NOT NULL,
    interest            NUMBER(5,2)  NOT NULL CHECK (interest >= 0),
    is_active           NUMBER(1) DEFAULT 0 NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    deleted_at          TIMESTAMP WITH TIME ZONE,
    is_deleted          NUMBER(1) DEFAULT 0 NOT NULL,
    CONSTRAINT fk_amount_tier_category FOREIGN KEY (vault_category_id) REFERENCES vault_categories (id) ON DELETE CASCADE,
    CONSTRAINT chk_amount_minmax CHECK (max_amount > min_amount)
);
/

COMMENT ON COLUMN vault_categories.deleted_at IS 'Used for soft deletes';
/
COMMENT ON COLUMN amount_tier.deleted_at IS 'Used for soft deletes';
/

-- ===========================================
-- Index
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
