-- Add IS_ENABLED when BUDGET_CATEGORIES was created without it (manual DDL / old schema).
DECLARE
  v_count NUMBER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_tab_columns
  WHERE table_name = 'BUDGET_CATEGORIES'
    AND column_name = 'IS_ENABLED';

  IF v_count = 0 THEN
    EXECUTE IMMEDIATE q'[
      ALTER TABLE BUDGET_CATEGORIES ADD (IS_ENABLED NUMBER(1) DEFAULT 1 NOT NULL)
    ]';
  END IF;
END;
/
