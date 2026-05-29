-- Ensure IS_ENABLED exists alongside ENABLED so NVL(IS_ENABLED, ENABLED) works in the app.
-- Copies from ENABLED when the column is new.
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

  IF v_is_enabled = 0 AND v_enabled > 0 THEN
    EXECUTE IMMEDIATE q'[
      ALTER TABLE BUDGET_CATEGORIES ADD (IS_ENABLED NUMBER(1) DEFAULT 1 NOT NULL)
    ]';
    EXECUTE IMMEDIATE q'[
      UPDATE BUDGET_CATEGORIES SET IS_ENABLED = ENABLED
    ]';
  END IF;
END;
/
