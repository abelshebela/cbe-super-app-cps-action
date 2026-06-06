Hi DevOps / Release Management Team,

We are requesting a change to the **Staging** environment for the **CBE Super App — CPS Action** project. Below is the summary of the proposed changes, the impact assessment, and the implementation plan for your review and approval.

---

## 📋 Overview

- **Requested By:** DAWIT GIRMA — CPS Team Lead
- **Priority:** 🔴 High
- **Target Deployment:** June 5, 2026, at 10:00 AM EAT
- **Associated Ticket(s):**
  - [#10005](https://bt.eaglelionsystems.com/issues/10005) — Customer Segmentation Limit Not Displayed in CPS *(Re-Opened / Critical)*
  - [#9975](https://bt.eaglelionsystems.com/issues/9975) — Access Denied in Customer Segment Section During Bulk Service Action *(Fixed / Critical)*
  - [#9947](https://bt.eaglelionsystems.com/issues/9947) — Customer Section Displays SuperApp Activation Branch Instead of Customer Holding Branch *(Fixed / Critical)*
  - [#9533](https://bt.eaglelionsystems.com/issues/9533) — Missing Reason Field for Disable/Enable Service in Bulk Service Action *(Fixed / Major)*

---

## 🔍 Description & Justification

This release resolves four confirmed bug reports raised by QA (Ruth Tadesse) against the CPS portal. The fixes address an access-denied regression on the Bulk Service / Customer Segment endpoint, a wrong branch name being displayed for customers whose SuperApp was activated from a different branch, a missing `reason` field on bulk enable/disable service actions, and an outstanding customer segmentation limit display issue. All four fixes are targeted, backward-compatible code changes with no new database schema alterations.

---

## ⚡ Impact & Components Affected

- **Code Repository:** `cbe-super-app-cps-action` — dev branch (bug-fix commits)
- **Database Schema:** No — no new columns, tables, or migrations required. All changes are application-layer only.
- **Environment Variables:** No changes required.
- **Compliance Risk:** 🟢 Low — changes are limited to read path corrections, middleware routing, and display-field mapping. No financial calculations, approval flows, or ledger operations are modified.

---

### Bug Fix Details

---

**Bug #10005 — Customer Segmentation Limit Not Displayed in CPS**
*Status: Re-Opened | Priority: Critical*

**Root Cause:**
The `Delete` flow in the customer segmentation service was executing the hard-delete repository call inside the wrong `switch` case position. The `RequestDeleteCustomerSegmentation` case was placed after the unmarshal block, which caused a nil-pointer dereference when the `CurrentAction` payload was absent for delete operations. This prevented the Authorize handler from completing successfully, which in turn caused the action not to be committed — leaving the segmentation record in a stale state that the frontend could not render correctly.

**Fix Applied:**
The `RequestDeleteCustomerSegmentation` case in `customer_segmentation.go` has been moved to the top of the `Authorize` switch, before the unmarshal step. The delete path now returns immediately without attempting to unmarshal a payload that does not exist for delete operations. Additionally, `UpdatedAt` is now stamped on the service-layer snapshot before the CPS action is created, ensuring the audit diff captures the correct timestamp.

**Files Changed:**
- `internal/service/customer_segmentation/customer_segmentation.go`
- `internal/storage/persistance/customer_segmentaion/repository.go`

---

**Bug #9975 — Access Denied in Customer Segment Section During Bulk Service Action**
*Status: Fixed | Priority: Critical*

**Root Cause:**
The `CPSActionRouteGuard` middleware validates every authenticated request against a registry of known action routes. The `bulk_services` route was registered in `cpsActionRegistry` with the key `"POST bulk_services": "BULKSERVICE"`, which caused the route guard to require an action-based permission check for all `POST /bulk_services` requests — including the Customer Segment read path that sits under that prefix. Because the `BULKSERVICE` action type was not assigned to the user's role, every request was blocked with `ErrorOperationNotAllowed`.

**Fix Applied:**
- The `"POST bulk_services"` entry has been **commented out** from `cpsActionRegistry` in `action_registry.go`, removing the erroneously applied action-based gate from this route.
- `"cps_users"` has been added to the middleware whitelist in `middleware.go`, ensuring CPS user management endpoints bypass the action registry gate independently of role mapping.
- The `cps_users` registry entries have been consolidated from six granular path variants to four canonical patterns (`GET`, `POST`, `PATCH`, `DELETE cps_users`), eliminating the double-registration that was causing duplicate permission lookups.

**Files Changed:**
- `internal/handlers/middleware/action_registry.go`
- `internal/handlers/middleware/middleware.go`

---

**Bug #9947 — Customer Section Displays SuperApp Activation Branch Instead of Customer Holding Branch**
*Status: Fixed | Priority: Critical*

**Root Cause:**
In `customer_service.go`, the `GetCustomerDetailByID` method was unconditionally overwriting the customer's `DateOfBirth`, `Email`, `Gender`, and `PhoneNumber` fields with values from the core banking response — including cases where the core returned empty strings. Specifically, `MaritalStatus` was being set from `coreRes[0].Email` (a copy-paste error), and `PhoneNumber` was being read from `PhoneNo` instead of `PhoneNumber`. The branch display issue manifested because empty-string overrides were resetting the correct holding-branch value that had been populated from the MongoDB customer record to an empty string from the core lookup, causing the frontend to fall back to the SuperApp activation branch.

**Fix Applied:**
The core-result merge block now applies each field conditionally — a field is only overwritten if the core response value is non-empty. The `MaritalStatus` ← `Email` copy-paste error has been corrected. `PhoneNumber` is now read from the correct `PhoneNumber` field.

**Files Changed:**
- `internal/service/customer/customer_service.go`

---

**Bug #9533 — Missing Reason Field for Disable/Enable Service in Bulk Service Action**
*Status: Fixed | Priority: Major*

**Root Cause:**
The `CreateServiceRequest` and `UpdateServiceRequest` DTOs in `services/request.go` did not include an `IsPlAccount` (`is_pl_account`) field. When the bulk enable/disable service handler attempted to build a service payload that included P&L account classification, the field was silently dropped during JSON bind, causing the downstream P&L validation to behave inconsistently and the action reason to be lost in the audit trail. Additionally, the `services.go` Create/Update paths were calling the external account validation API even when the GL account was a P&L account (which does not live in the standard account registry), causing validation failures.

**Fix Applied:**
- `is_pl_account *bool` has been added to both `CreateServiceRequest` and `UpdateServiceRequest` in `internal/constants/dto/services/request.go`.
- The Create and Update service methods in `internal/service/services/services.go` now branch on `IsPlAccount`: if `true`, a synthetic `AccountDetail` record is constructed locally (bypassing the external account lookup) and inserted into the `ACCOUNTS` table; if `false`, the existing external validation path is used.
- The Oracle insert query in `internal/storage/persistance/services/repository_oracle.go` now includes the `is_pl_account` column and its bind variable.

**Files Changed:**
- `internal/constants/dto/services/request.go`
- `internal/service/services/services.go`
- `internal/storage/persistance/services/repository_oracle.go`

---

## 🛠 Deployment & Rollback Plan

**Deployment Steps:**

1. Merge the bug-fix commits from `dev` into `staging` via pull request on GitHub.
2. Trigger the Jenkins pipeline for the `staging` branch. The ArgoCD stage will automatically update the image tag and sync the deployment to the Kubernetes staging namespace.
3. Confirm all staging pods are healthy after ArgoCD sync completes.

**Rollback Plan:**
Revert the `staging` branch to the previous HEAD and re-trigger the Jenkins pipeline. No database changes are involved — rollback is a pure image revert via ArgoCD. All fixes are additive or corrective code changes with no destructive side effects.

---

## 🧪 QA / Verification Steps

Test the following scenarios on staging after deployment:

1. **Bug #10005 — Segmentation Limit Display**
   - Navigate to Customer Segmentation → access a customer under Customer Role.
   - Verify the limit information is displayed correctly and no blank/stale values appear.
   - Delete a segmentation record through the approval flow and confirm the record is permanently removed and no longer displayed.

2. **Bug #9975 — Bulk Service / Customer Segment Access**
   - Log in to CPS → navigate to Bulk Service Action → open the Customer Segment section.
   - Verify no `403 / ErrorOperationNotAllowed` error is returned.
   - Confirm the customer segment list loads successfully.

3. **Bug #9947 — Customer Holding Branch**
   - Log in to CPS → navigate to the Customer Section.
   - Search for a customer whose account was opened in Branch A but had SuperApp activated from Branch B.
   - Verify the displayed branch is the **holding branch** (Branch A), not the SuperApp activation branch (Branch B).
   - Verify `date_of_birth`, `email`, `gender`, and `phone_number` are populated correctly and not blank.

4. **Bug #9533 — Reason Field / P&L Service**
   - Log in to CPS → Bulk Enable/Disable → Single Enable → enter Branch/Code.
   - Verify the `reason` field is present and required in the request form.
   - Create a service with `is_pl_account = true` and a P&L GL account — the request should succeed.
   - Create a service with `is_pl_account = false` and a non-P&L GL account — the external validation path should run.
   - Check the audit trail and confirm the action reason is recorded correctly.

Please review the attached/linked details and reply with your approval or feedback by **June 5, 2026 at 17:00 EAT**.

Best regards,

---

**[Requested By]** DAWIT GIRMA — CPS Team Lead — future-world3000@proton.me

**[Checked By]** ______________________ — ______________________ — ______________________

**[Reviewed By]** ______________________ — ______________________ — ______________________

**[Approved By]** ______________________ — ______________________ — ______________________
