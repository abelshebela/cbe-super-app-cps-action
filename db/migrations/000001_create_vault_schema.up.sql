-- ===========================================
-- Drop tables (safe)
-- ===========================================
BEGIN
  FOR t IN (
    SELECT table_name
    FROM user_tables
    WHERE table_name IN ('VAULT_TIERS', 'WITHDRAWALS', 'VAULT_CATEGORIES')
  ) LOOP
    EXECUTE IMMEDIATE 'DROP TABLE ' || t.table_name || ' CASCADE CONSTRAINTS';
  END LOOP;
END;
/

-- ===========================================
-- Drop indexes (safe)
-- ===========================================
BEGIN
  FOR i IN (
    SELECT index_name
    FROM user_indexes
    WHERE index_name IN ('IDX_VAULT_CATEGORIES_IS_ACTIVE', 'IDX_WITHDRAWALS_IS_ACTIVE', 'IDX_VAULT_TIERS_CATEGORY_ID')
  ) LOOP
    EXECUTE IMMEDIATE 'DROP INDEX ' || i.index_name;
  END LOOP;
END;
/

-- ===========================================
-- Create vault_categories
-- ===========================================
CREATE TABLE vault_categories (
    id                VARCHAR2(36) PRIMARY KEY,
    name              VARCHAR2(255) NOT NULL,
    cover_image_url   VARCHAR2(255),
    interest_type     VARCHAR2(50) NOT NULL CHECK (interest_type IN ('FLAT', 'DYNAMIC')),
    category_interest VARCHAR2(50),
    deadlock          NUMBER(1) DEFAULT 0 NOT NULL,
    is_active         NUMBER(1) DEFAULT 0 NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT uq_vault_categories_name UNIQUE (name),
    CONSTRAINT chk_vc_deadlock CHECK (deadlock IN (0,1)),
    CONSTRAINT chk_vc_is_active CHECK (is_active IN (0,1))
);

-- ===========================================
-- Create vault_tiers
-- ===========================================
CREATE TABLE vault_tiers (
    id             VARCHAR2(36) PRIMARY KEY,
    category_id    VARCHAR2(36) NOT NULL,
    name           VARCHAR2(100) NOT NULL,
    tier_interest  VARCHAR2(50),
    min_amount     VARCHAR2(50),
    max_amount     VARCHAR2(50),
    CONSTRAINT fk_vt_category FOREIGN KEY (category_id) REFERENCES vault_categories(id) ON DELETE CASCADE
);

-- ===========================================
-- Create withdrawals
-- ===========================================
CREATE TABLE withdrawals (
    id                       VARCHAR2(36) PRIMARY KEY,
    locked_vault_id          VARCHAR2(36) NOT NULL,
    amount                   VARCHAR2(50) NOT NULL,
    withdrawer_name          VARCHAR2(255),
    withdrawer_phone_number  VARCHAR2(50),
    status                   VARCHAR2(50),
    is_active                NUMBER(1) DEFAULT 0 NOT NULL,
    created_at               TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    updated_at               TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
    CONSTRAINT chk_w_is_active CHECK (is_active IN (0,1))
);

-- ===========================================
-- Indexes
-- ===========================================
CREATE INDEX idx_vault_categories_is_active ON vault_categories(is_active);
CREATE INDEX idx_withdrawals_is_active ON withdrawals(is_active);
CREATE INDEX idx_vault_tiers_category_id ON vault_tiers(category_id);

-- ===========================================
-- Triggers
-- ===========================================
CREATE OR REPLACE TRIGGER trg_vault_categories_biu
BEFORE INSERT OR UPDATE ON vault_categories
FOR EACH ROW
BEGIN
  IF INSERTING THEN
    IF :NEW.created_at IS NULL THEN
      :NEW.created_at := SYSTIMESTAMP;
    END IF;
  END IF;
  :NEW.updated_at := SYSTIMESTAMP;
END;
/

CREATE OR REPLACE TRIGGER trg_withdrawals_biu
BEFORE INSERT OR UPDATE ON withdrawals
FOR EACH ROW
BEGIN
  IF INSERTING THEN
    IF :NEW.created_at IS NULL THEN
      :NEW.created_at := SYSTIMESTAMP;
    END IF;
  END IF;
  :NEW.updated_at := SYSTIMESTAMP;
END;
/