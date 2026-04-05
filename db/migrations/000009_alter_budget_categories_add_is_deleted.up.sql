-- Add IS_DELETED when BUDGET_CATEGORIES exists without it (manual DDL / old schema).
DECLARE
  v_count NUMBER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_tab_columns
  WHERE table_name = 'BUDGET_CATEGORIES'
    AND column_name = 'IS_DELETED';

  IF v_count = 0 THEN
    EXECUTE IMMEDIATE q'[
      ALTER TABLE BUDGET_CATEGORIES ADD (IS_DELETED NUMBER(1) DEFAULT 0 NOT NULL)
    ]';
  END IF;
END;
/
