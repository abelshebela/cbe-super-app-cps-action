-- ===========================================
-- Drop triggers
-- ===========================================

BEGIN
  FOR t IN (
    SELECT trigger_name
    FROM user_triggers
    WHERE trigger_name IN (
      'TRG_VAULT_CATEGORIES_BIU',
      'TRG_WITHDRAWALS_BIU'
    )
  ) LOOP
    EXECUTE IMMEDIATE 'DROP TRIGGER ' || t.trigger_name;
  END LOOP;
END;
/
------------------------------------------------


-- ===========================================
-- Drop indexes
-- ===========================================

BEGIN
  FOR i IN (
    SELECT index_name
    FROM user_indexes
    WHERE index_name IN (
      'IDX_VAULT_CATEGORIES_IS_ACTIVE',
      'IDX_WITHDRAWALS_IS_ACTIVE',
      'IDX_VAULT_TIERS_CATEGORY_ID'
    )
  ) LOOP
    EXECUTE IMMEDIATE 'DROP INDEX ' || i.index_name;
  END LOOP;
END;
/
------------------------------------------------


-- ===========================================
-- Drop tables
-- ===========================================

BEGIN
  FOR t IN (
    SELECT table_name
    FROM user_tables
    WHERE table_name IN (
      'WITHDRAWALS',
      'VAULT_TIERS',
      'VAULT_CATEGORIES'
    )
  ) LOOP
    EXECUTE IMMEDIATE 'DROP TABLE ' || t.table_name || ' CASCADE CONSTRAINTS';
  END LOOP;
END;
/
------------------------------------------------
