-- App expects column ENABLED (legacy DDL). Normalize IS_ENABLED -> ENABLED, or add ENABLED if missing.
DECLARE
  v_is_enabled NUMBER;
  v_enabled    NUMBER;
BEGIN
  SELECT COUNT(*) INTO v_is_enabled
  FROM user_tab_columns
  WHERE table_name = 'BUDGET_CATEGORIES' AND column_name = 'IS_ENABLED';

  SELECT COUNT(*) INTO v_enabled
  FROM user_tab_columns
  WHERE table_name = 'BUDGET_CATEGORIES' AND column_name = 'ENABLED';

  IF v_is_enabled > 0 AND v_enabled = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE BUDGET_CATEGORIES RENAME COLUMN IS_ENABLED TO ENABLED';
  ELSIF v_is_enabled = 0 AND v_enabled = 0 THEN
    EXECUTE IMMEDIATE q'[
      ALTER TABLE BUDGET_CATEGORIES ADD (ENABLED NUMBER(1) DEFAULT 1 NOT NULL)
    ]';
  END IF;
END;
/
