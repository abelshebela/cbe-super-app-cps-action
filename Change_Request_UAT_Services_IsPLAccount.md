Hi DevOps / Release Management Team,

We are requesting a change to the **UAT** environment for the **CBE Super App — CPS Action** project. Below is the summary of the proposed changes, the impact assessment, and the implementation plan for your review and approval.

---

## 📋 Overview

- **Requested By:** DAWIT GIRMA — CPS Team Lead
- **Priority:** 🟡 Medium
- **Target Deployment:** Upon approval
- **Associated Ticket(s):** SERVICES table schema update — add IS_PL_ACCOUNT column (UAT)

---

## 🔍 Description & Justification

The CPS Action service now validates that a GL account configured as a Product GL Account on a service record is correctly flagged as a P&L (Profit & Loss) account in the core banking system before allowing a create or update operation to proceed. To support this validation and persist the P&L flag alongside the service record, the `SERVICES` table in Oracle requires a new column `IS_PL_ACCOUNT NUMBER(1) DEFAULT 0 NOT NULL` protected by a CHECK constraint.

---

## ⚡ Impact & Components Affected

- **Code Repository:** cbe-super-app-cps-action — already merged on `dev` (promoted to `staging`)
- **Database Schema:** Yes — one column added and one CHECK constraint added to the existing `SERVICES` table. No existing columns, indexes, or foreign keys are altered or removed.
- **Environment Variables:** No changes required.
- **Compliance Risk:** 🟢 Low — the column defaults to `0` (false) for all existing rows, so no existing service records are affected. No financial ledger calculations are modified. The change is fully backward-compatible.

### What Is Changing

**Column: `IS_PL_ACCOUNT NUMBER(1) DEFAULT 0 NOT NULL`**
A boolean-style flag stored as `NUMBER(1)`. When `1`, the Product GL Account linked to this service is confirmed to be a P&L account. When `0` (default), no P&L flag is set. The CPS Action service writes this flag at create/update time after querying the core banking account lookup.

**Constraint: `CHK_SERVICES_IS_PL_ACCOUNT CHECK (IS_PL_ACCOUNT IN (0, 1))`**
Ensures only the values `0` or `1` can be stored in the column, preventing invalid data at the database layer.

### Full Target Table Definition (for reference)

```sql
CREATE TABLE SERVICES (
    ID                           RAW(16)        DEFAULT SYS_GUID() PRIMARY KEY,
    ACCESS_LIST_ID               RAW(16)        UNIQUE NOT NULL,
    SERVICE_CODE                 VARCHAR2(64)   NOT NULL,
    MINIMUM_FRAUD_AMOUNT         NUMBER(19, 4)  NOT NULL,
    IS_PL_ACCOUNT                NUMBER(1)      DEFAULT 0 NOT NULL,
    PRODUCT_GL_ACCOUNT_NUMBER    VARCHAR2(16),
    PRODUCT_GL_ACCOUNT_CURRENCY  VARCHAR2(3),
    CREATED_AT                   TIMESTAMP      DEFAULT SYSTIMESTAMP,
    LAST_MODIFIED_AT             TIMESTAMP      DEFAULT SYSTIMESTAMP,
    IS_DELETED                   NUMBER(1)      DEFAULT 0,
    DELETED_AT                   TIMESTAMP,
    CONSTRAINT CHK_SERVICES_IS_DELETED
        CHECK (IS_DELETED IN (0, 1)),
    CONSTRAINT CHK_SERVICES_IS_PL_ACCOUNT
        CHECK (IS_PL_ACCOUNT IN (0, 1)),
    CONSTRAINT FK_SERVICES_ACCESS_LIST_ID
        FOREIGN KEY (ACCESS_LIST_ID) REFERENCES ACCESS_LISTS (ID) ON DELETE CASCADE,
    CONSTRAINT FK_SERVICES_PRODUCT_GL_ACCOUNT_NUMBER
        FOREIGN KEY (PRODUCT_GL_ACCOUNT_NUMBER) REFERENCES ACCOUNTS (ACCOUNT_NUMBER) ON DELETE CASCADE,
    CONSTRAINT CHK_SERVICES_PRODUCT_GL_ACCOUNT_CURRENCY
        CHECK (PRODUCT_GL_ACCOUNT_CURRENCY IN ('ETB', 'USD', 'GBP', 'EUR'))
);
CREATE INDEX IDX_SERVICES_ACCESS_LIST_ID ON SERVICES(ACCESS_LIST_ID);
```

---

## 🛠 Deployment & Rollback Plan

**Deployment Steps:**

1. Connect to the UAT Oracle database instance.
2. Execute the migration script below against the UAT schema:

```sql
-- Step 1: Add the IS_PL_ACCOUNT column
ALTER TABLE SERVICES
    ADD (IS_PL_ACCOUNT NUMBER(1) DEFAULT 0 NOT NULL);

-- Step 2: Add the CHECK constraint
ALTER TABLE SERVICES
    ADD CONSTRAINT CHK_SERVICES_IS_PL_ACCOUNT
    CHECK (IS_PL_ACCOUNT IN (0, 1));
```

3. Verify the column and constraint were created (see QA step 1 below).
4. Deploy the updated CPS Action application image to the UAT Kubernetes namespace.

**Rollback Plan:**

If any issues arise, execute the rollback script below to remove the column and constraint:

```sql
-- Rollback Step 1: Drop the CHECK constraint
ALTER TABLE SERVICES
    DROP CONSTRAINT CHK_SERVICES_IS_PL_ACCOUNT;

-- Rollback Step 2: Drop the column
ALTER TABLE SERVICES
    DROP COLUMN IS_PL_ACCOUNT;
```

No data loss occurs on rollback — the column contains no business-critical data (it is a derived flag populated by the application at write time).

---

## 🧪 QA / Verification Steps

Test the following on UAT after the migration and deployment:

1. **Schema verification** — Run the query below and confirm `IS_PL_ACCOUNT` appears with `DATA_TYPE = NUMBER`, `NULLABLE = N`, and `DATA_DEFAULT = 0`:
   ```sql
   SELECT COLUMN_NAME, DATA_TYPE, NULLABLE, DATA_DEFAULT
   FROM   ALL_TAB_COLUMNS
   WHERE  TABLE_NAME  = 'SERVICES'
   AND    COLUMN_NAME = 'IS_PL_ACCOUNT';
   ```

2. **Constraint verification** — Run the query below and confirm `CHK_SERVICES_IS_PL_ACCOUNT` is listed:
   ```sql
   SELECT CONSTRAINT_NAME, CONSTRAINT_TYPE, SEARCH_CONDITION
   FROM   ALL_CONSTRAINTS
   WHERE  TABLE_NAME       = 'SERVICES'
   AND    CONSTRAINT_NAME  = 'CHK_SERVICES_IS_PL_ACCOUNT';
   ```

3. **Existing row default** — Confirm all existing rows have `IS_PL_ACCOUNT = 0`:
   ```sql
   SELECT COUNT(*) FROM SERVICES WHERE IS_PL_ACCOUNT != 0;
   -- Expected: 0
   ```

4. **P&L validation — accepted** — Via the CPS Action API, create or update a service using a GL account that IS flagged as a P&L account in the core system. The request should succeed and the service record should be saved with `IS_PL_ACCOUNT = 1`.

5. **P&L validation — rejected** — Via the CPS Action API, create or update a service using a GL account that is NOT flagged as a P&L account. The request should be rejected with an appropriate validation error before any database write occurs.

6. **Constraint guard** — Attempt a direct INSERT with an invalid value to confirm the constraint is enforced at the DB layer:
   ```sql
   -- Should raise ORA-02290: check constraint violated
   INSERT INTO SERVICES (ACCESS_LIST_ID, SERVICE_CODE, MINIMUM_FRAUD_AMOUNT, IS_PL_ACCOUNT)
   VALUES (SYS_GUID(), 'TEST', 0, 5);
   ```

Please review the attached/linked details and reply with your approval or feedback at your earliest convenience.

Best regards,

---

**[Requested By]** DAWIT GIRMA — CPS Team Lead — future-world3000@proton.me

**[Checked By]** ______________________ — ______________________ — ______________________

**[Reviewed By]** ______________________ — ______________________ — ______________________

**[Approved By]** ______________________ — ______________________ — ______________________
