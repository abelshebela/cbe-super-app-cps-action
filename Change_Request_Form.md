Hi DevOps / Release Management Team,

We are requesting a change to the **Staging** environment for the **CBE Super App — CPS Action** project. Below is the summary of the proposed changes, the impact assessment, and the implementation plan for your review and approval.

---

## 📋 Overview

- **Requested By:** CBE Super App Backend Team — Dawit (Backend Engineer)
- **Priority:** 🔴 High
- **Target Deployment:** June 2, 2026, at 10:00 AM EAT
- **Associated Ticket(s):** dev → staging promotion (79 commits) — [cbe-super-app-cps-action, branch: dev]

---

## 🔍 Description & Justification

This release promotes accumulated feature work and bug fixes from the `dev` branch to `staging`. The changes deliver a formal KYC reviewer assignment workflow, duplicate action submission prevention, extended advanced filtering for checkers and auditors, expanded BPS action type coverage, logistics merchant migration to Oracle, and several data correctness fixes identified during QA and internal testing. These changes are required to meet the agreed sprint deliverables and unblock QA validation on staging.

---

## ⚡ Impact & Components Affected

- **Code Repository:** `cbe-super-app-cps-action` — PR: dev → staging (81 files changed, +3,898 / −797 lines)
- **Database Schema:** Yes
  - **MongoDB (new collection):** `started_kyc_reviews` — stores KYC reviewer session records (reviewer identity, start time, expiry window, pick status, pick count). No existing MongoDB collections are structurally altered.
  - **MongoDB (new fields on `cps_actions`):** `checksum` (SHA-256 of business payload) and `unique_tokens` (extracted identifying values), both used for duplicate CREATE action detection.
  - **MongoDB (new field on `user_action_logs`):** `action_auditor_status` — tracks the overall auditor workflow state (`AUDITORINPROGRESS` / `AUDITORCHECKED`) per action code, written in bulk across all related log entries when an auditor claims or marks an action.
  - **Oracle:** Customer segmentation delete is now a cascading hard delete across `CUSTOMER_SUB_SEGMENTS`, `CUSTOMER_SEGMENTATIONS`, and `CUSTOMER_GROUPS`. Previously rows were soft-deleted with `IS_DELETED = 1`; they are now permanently removed once checker approval is granted.
- **Environment Variables:** Yes — a new Jenkins credential `GITHUB_EMAIL` must be added to the `kr-jenkins-slave-1` build agent before the pipeline is triggered. This is required by the new ArgoCD stage that automatically commits the updated image tag to the GitOps manifest repository (`cbe-superapp-deployment`) after a successful build.
- **Compliance Risk:** 🟡 Medium — changes touch the KYC review flow, which is a regulated process. The hard delete of customer segmentation records is also irreversible once approved. All changes have passed internal code review. No financial ledger calculations are modified.

### Business Changes

**1. Duplicate Action Submission Prevention**
A maker can no longer submit the same CREATE action twice while an identical request is already pending checker approval. The system computes a SHA-256 checksum from the business payload and also extracts unique field tokens. Both are checked against existing PENDING actions before saving. If a duplicate is found the submission is blocked with an explicit `409 Conflict` error, preventing duplicate work from entering the checker queue.

**2. KYC Review Workflow — Start Review**
A reviewer can now initiate a formal review session on a PENDING KYC application. The session records the reviewer's identity, start time, and a configurable expiry window. Only one active (non-expired) session can exist per KYC record. Starting a second session while one is active, or starting on a non-PENDING record, is rejected. The KYC status moves to `IN_REVIEW`. The KYC detail response now returns `started_at`, `expires_at`, and reviewer information when the record is in `IN_REVIEW` status.

**3. KYC Review Workflow — Pick Review**
A reviewer can formally claim (pick) a KYC review session. While the session is still active only the originally assigned reviewer can pick it. After the session expires any eligible reviewer can pick it and extend ownership. The pick is submitted as a maker CPS action (`PICK_KYC_REVIEW`) and goes through the full checker approval flow, making assignment fully auditable.

**4. KYC Enable/Disable Guard**
KYC records in `PENDING` status can no longer be enabled or disabled, preventing accidental status changes on applications that are actively awaiting review.

**5. Advanced Filtering for Checkers and Auditors (CPS & BPS)**
The action list views for checkers and auditors now support multi-value filters for `approval_level`, `service`, and `action_status`. For non-PENDING statuses the system resolves matching `action_codes` through the `user_action_log` before querying `cps_actions`, ensuring history-based filters work correctly. The PENDING-only view bypasses the log lookup and is handled directly by the repository for correctness.

**6. Auditor Status Tracking Fix**
The `user_action_log` now accurately reflects auditor progress. When an auditor claims an action the log is bulk-updated to `AUDITORINPROGRESS`. When an auditor marks, the system checks whether all required auditor groups have completed: if yes the log updates to `AUDITORCHECKED`, otherwise it remains `AUDITORINPROGRESS`.

**7. Checker Level Fix on Approve / Reject**
The checker's approval level (1, 2, 3…) was being stored blank in the action log due to a Go type assertion mismatch (`int` vs `float64` from JSON-decoded context). This is now read as `float64` and formatted correctly. Checker-level breakdowns in audit trails and reports are now fully populated.

**8. Customer Segmentation — Permanent Delete**
Deleting a customer segmentation record now performs a cascading hard delete in Oracle. The operation removes sub-segments, then parent segmentations with no remaining children, then parent groups with no remaining segmentations — all within a single transaction. Records are permanently removed once the delete action passes checker approval.

**9. USSD Merchant — Additional Fields in Response**
The USSD merchant detail response now includes `merchant_code` and `credential` fields in addition to all existing fields, removing the need for a secondary lookup.

**10. Services — P&L Account Validation**
When creating or updating a service, the system now validates that a GL account provided as the Product GL Account is correctly flagged as a P&L account in the core system. Requests with a non-P&L GL account are rejected, preventing financial misconfiguration.

**11. CPS User Creation — Department and Job Title Validation**
When creating a new CPS user, the system now validates that the specified department and job title both exist in the system registries. The role is resolved from the job role registry rather than taken directly from the request, ensuring canonical role assignment. Requests with invalid department or job title are rejected.

**12. BPS Action Type Coverage Expansion**
The BPS approval workflow now covers approximately 50 additional action types previously only handled in CPS, including: manage BPS users, banks, wallets, advertisements, service fees, daily limits, validation rules, password rules, events, event categories, mini-app merchants, access control, archive expiry, password expiry, and KYC level upgrades. All now go through the full BPS checker and auditor approval flow.

**13. Logistics Merchant — Oracle-Backed Data**
Logistics merchant read and write operations are now backed by Oracle instead of MongoDB. All queries include a `merchant_type` filter to prevent cross-type contamination. Date fields are nullable to correctly represent merchants with no expiry or creation date in the core system.

**14. Image Upload Security — Binary Structure Validation**
File uploads are now validated against their internal binary structure (magic bytes) in addition to MIME type. A file renamed to `.jpg` or `.png` that is not a valid image is rejected at upload time.

**15. CPS User — First-Time Login Indicator**
The CPS user profile response now includes an `is_first_time_login` field, enabling the portal frontend to prompt first-time users to change their password on initial login.

**16. Donation Date Handling Fix**
Donation records with no end date no longer return the Go zero-value timestamp (`0001-01-01`). Open-ended campaigns return an empty value instead.

**17. Password Rule Validation Fix**
The `PasswordRuleUpdate` validator has been rewritten with explicit checks to correctly enforce `MinLength < MaxLength`. The previous framework-based validator failed to catch the `>=` boundary case.

**18. Customer Segmentation — Duplicate Entry Validator Removed**
The `validateUniqueCustomerEntries` DTO-level validator has been removed. Duplicate-entry detection at the payload level is no longer enforced at this layer.

---

## 🛠 Deployment & Rollback Plan

**Deployment Steps:**

1. Add the `GITHUB_EMAIL` Jenkins credential to the `kr-jenkins-slave-1` build agent.
2. Merge `dev` → `staging` via pull request on GitHub.
3. Trigger the Jenkins pipeline for the `staging` branch. The new ArgoCD stage will clone `cbe-superapp-deployment`, update the image tag in `k8s-manifests/overlays/staging/superapp/cps/kustomization.yml`, commit, and push. ArgoCD will sync the deployment to the Kubernetes staging namespace automatically.
4. Confirm all staging pods are healthy after ArgoCD sync completes (`kubectl get pods -n staging`).

**Rollback Plan:** If any issues arise, revert the `staging` branch to the previous HEAD (`e2ec31bb` — last known-good staging commit) and re-trigger the Jenkins pipeline. ArgoCD will automatically apply the prior image tag once the pipeline writes it back to the manifest repository. No database migration scripts are required for rollback; the new MongoDB fields and collection are additive and the application degrades gracefully without them. The Oracle hard-delete is irreversible — segmentation records deleted after deployment cannot be recovered from the database. Ensure the staging dataset is acceptable before approving delete actions post-deployment.

---

## 🧪 QA / Verification Steps

Test the following scenarios on staging after deployment:

1. **Duplicate prevention** — Submit the same CREATE action twice without approving the first. The second submission must be rejected with a `409 ERROR_DUPLICATE_PENDING_CREATE_ACTION` error.
2. **KYC Start Review** — Start a review on a PENDING KYC (should succeed, status moves to IN_REVIEW). Attempt to start a second review on the same record (should be rejected). Attempt to start on an APPROVED KYC (should be rejected).
3. **KYC Pick Review** — As the assigned reviewer, pick the review (should succeed). As a different reviewer while the session is active, attempt to pick (should be rejected). Let the session expire and pick as a different reviewer (should succeed).
4. **KYC Enable/Disable Guard** — Attempt to enable or disable a PENDING KYC record (should be blocked with an error).
5. **Advanced Filtering** — In the checker queue, apply combined filters for Level 1 + a specific service + status APPROVED. Verify only matching records are returned. Apply a PENDING-only filter and confirm pending items are present and complete.
6. **Auditor Status** — Claim an action as auditor and verify the log shows `AUDITORINPROGRESS`. Mark the final auditor group and verify the log updates to `AUDITORCHECKED`.
7. **Checker Level in Log** — Approve an action as a Level 2 checker. Verify the action log records level `"2"` and not blank.
8. **Segmentation Delete** — Delete a customer segment through the full approval flow. After checker approval confirm the record is permanently gone and cannot be retrieved.
9. **USSD Merchant** — Fetch a USSD merchant by ID and verify `merchant_code` and `credential` fields are populated in the response.
10. **P&L Account** — Create a service providing a GL account not flagged as P&L in the core system (should be rejected).
11. **CPS User Creation** — Attempt to create a CPS user with an invalid department (should be rejected). Attempt with an invalid job title (should be rejected). Create with valid data and confirm the role is resolved from the job role registry.
12. **Image Upload** — Upload a text file with a `.jpg` extension (should be rejected at upload time).
13. **Donation Date** — Retrieve an open-ended donation campaign and verify `end_date` is empty (not `0001-01-01`).
14. **First-Time Login** — Log in as a newly created CPS user and verify `is_first_time_login: true` in the profile response.
15. **Logistics Merchant** — Fetch the logistics merchant list and verify no cross-type merchants appear. Verify records with no expiry date display correctly (not a zero date).

Please review the attached/linked details and reply with your approval or feedback by **June 1, 2026 at 17:00 EAT**.

Best regards,

---

**[Requested By]** Dawit — Backend Engineer — future-world3000@proton.me

**[Checked By]** ______________________ — ______________________ — ______________________

**[Reviewed By]** ______________________ — ______________________ — ______________________

**[Approved By]** ______________________ — ______________________ — ______________________
