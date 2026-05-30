Hi DevOps / Release Management Team,

We are requesting a change to the **Staging** environment for the **CBE Super App — CPS Action** project. Below is the summary of the proposed changes, the impact assessment, and the implementation plan for your review and approval.

---

## Overview

- **Requested By:** CBE Super App Backend Team
- **Priority:** High
- **Target Deployment:** June 2, 2026, at 10:00 AM EAT
- **Associated Ticket(s):** dev → staging promotion (79 commits)

---

## Description & Justification

This release promotes accumulated feature work and bug fixes from the `dev` branch to `staging`. The changes deliver a formal KYC reviewer assignment workflow, duplicate action submission prevention, extended advanced filtering for checkers and auditors, expanded BPS action type coverage, logistics merchant migration to Oracle, and several data correctness fixes identified during QA and internal testing.

---

## Impact & Components Affected

- **Code Repository:** cbe-super-app-cps-action — dev → staging
- **Database Schema:** Yes — a new `kyc_reviews` collection is introduced in MongoDB to store KYC reviewer session records (reviewer identity, start time, expiry, pick status). No existing collections are structurally altered.
- **Environment Variables:** Yes — a new Jenkins credential `GITHUB_EMAIL` must be added to the `kr-jenkins-slave-1` build agent. This is required by the new ArgoCD pipeline stage that automatically commits the updated image tag to the GitOps manifest repository after a successful build.
- **Compliance Risk:** Medium — changes touch the KYC review flow, which is a regulated process. All changes have passed internal review. No financial ledger calculations are modified.

### Business Changes

**1. Duplicate Action Submission Prevention**
A maker can no longer submit the same CREATE action twice while an identical request is already waiting for checker approval. The system computes a checksum from the business payload (excluding timestamps and identities) and blocks the submission with an explicit error if a matching PENDING action already exists. This prevents duplicate work from entering the checker queue unnoticed.

**2. KYC Review Workflow — Start Review**
A new "Start Review" step has been introduced. A reviewer can initiate a session on a PENDING KYC application, which records the reviewer's identity, start time, and an expiry window. If an active (non-expired) session already exists for that KYC, no second session can be opened. Only KYC records in PENDING status can enter this step. The KYC detail response now also returns the review start time, expiry time, and assigned reviewer information when the record is in IN_REVIEW status.

**3. KYC Review Workflow — Pick Review**
A reviewer can now formally "pick" (claim) a KYC review that has been started. If the review session is still active, only the originally assigned reviewer can pick it. If the session has expired, any eligible reviewer can pick it and extend ownership. The pick action is recorded as a maker action, making the assignment auditable through the normal CPS checker approval flow.

**4. KYC Enable/Disable Guard**
KYC records in PENDING status can no longer be enabled or disabled, preventing accidental status changes on applications that are actively awaiting review.

**5. Advanced Filtering for Checkers and Auditors (CPS & BPS)**
The action list views for checkers and auditors now support multi-value filters for approval level, service, and action status. Checkers can filter their queue by level (e.g., Level 1 only), specific services, and status (PENDING, APPROVED, REJECTED). The system resolves matching records through the action audit log so that history-based filters work correctly for actions that have already moved past the PENDING stage. The PENDING-only view is handled separately to ensure no pending items are missed.

**6. Auditor Status Tracking**
The action audit log now accurately reflects auditor progress. When an auditor claims an action the log is updated to AUDITORINPROGRESS. When an auditor marks an action, the system determines whether all required auditor groups have completed their review: if yes, the log updates to AUDITORCHECKED; if not, it remains AUDITORINPROGRESS. Audit reports and dashboards will show accurate auditor workload status in real time.

**7. Checker Level Fix on Approve / Reject**
The checker's approval level (1, 2, 3…) was being stored as blank in the action log on approve and reject operations due to a type mismatch in the request context. This is now read and formatted correctly. Checker-level breakdowns in reports and audit trails are now fully populated.

**8. Customer Segmentation — Permanent Delete**
Deleting a customer segmentation record is now a permanent (hard) delete. Previously, deleted segments remained in the database with a soft-delete flag. They are now fully removed once the delete action passes checker approval.

**9. USSD Merchant — Service Name in Response**
The USSD merchant detail endpoint now returns the service name alongside the service ID. Consumers of this endpoint no longer need a second request to resolve the service name.

**10. Services — P&L Account Validation**
When creating or updating a service, the system now validates that a GL account provided as a Product GL Account is correctly flagged as a P&L account. If it is not flagged accordingly, the request is rejected, preventing misconfiguration that could affect financial reporting.

**11. CPS User Creation — Department and Job Title Validation**
When creating a new CPS user, the system now validates that the specified department and job title both exist in the system registries. The job title and role are resolved from the job role registry (not taken directly from the request), ensuring users are always assigned canonical roles. Requests with invalid or unregistered departments or job titles are rejected.

**12. BPS Action Type Coverage Expansion**
The BPS approval workflow now recognises approximately 50 additional action types that were previously only handled in CPS. These include: manage BPS users (enable/disable), manage banks and wallets (create/enable/disable), manage advertisements (create/update/enable/disable/delete), manage service fees and daily limits (create/update/delete), manage validation rules and password rules, manage events and event categories, manage mini-app merchants, configure access control, archive expiry, password expiry, and upgrade KYC level. These operations now go through the full BPS checker and auditor approval flow.

**13. Logistics Merchant — Oracle-Backed Data**
Logistics merchant data is now read from and written to Oracle (the bank's core system) instead of MongoDB. All queries include a merchant type filter to prevent merchants of different types from appearing in each other's lists. Date fields are handled as nullable to correctly represent merchants with no expiry or creation date recorded in the core system.

**14. Image Upload Security — Structure Validation**
File uploads that claim to be images are now validated not only by MIME type but also by their internal binary structure. A file renamed to `.jpg` or `.png` that is not a valid image will be rejected at upload time.

**15. CPS User — First-Time Login Indicator**
The CPS user profile response now includes an `is_first_time_login` field, enabling the portal frontend to prompt first-time users to change their password on initial login.

**16. Donation Date Handling Fix**
Donation records with no end date no longer return a zero-value timestamp in the response. Open-ended campaigns now return an empty value instead of `0001-01-01`.

---

## Deployment & Rollback Plan

**Deployment Steps:**
1. Merge `dev` → `staging` via pull request on GitHub.
2. Add the `GITHUB_EMAIL` Jenkins credential to the `kr-jenkins-slave-1` build agent before triggering the pipeline.
3. Trigger the Jenkins pipeline for the `staging` branch. The new ArgoCD stage will automatically commit the updated image tag and ArgoCD will sync the deployment to the Kubernetes staging namespace.
4. Confirm all staging pods are healthy after ArgoCD sync completes.

**Rollback Plan:** If any issues arise, revert the `staging` branch to the previous HEAD (`e2ec31bb` — last known-good staging commit) and re-trigger the Jenkins pipeline. ArgoCD will automatically roll back to the previous image tag once the pipeline writes it to the manifest repository.

---

## QA / Verification Steps

Test the following scenarios on staging after deployment:

1. **Duplicate prevention** — Submit the same CREATE action twice without approving the first. The second submission must be rejected with a duplicate pending error.
2. **KYC Start Review** — Start a review on a PENDING KYC (should succeed and move it to IN_REVIEW). Attempt to start a second review on the same record (should be rejected). Attempt to start a review on an APPROVED KYC (should be rejected).
3. **KYC Pick Review** — As the assigned reviewer, pick the review (should succeed). As a different reviewer while the session is still active, attempt to pick (should be rejected). Let the session expire and pick as a different reviewer (should succeed).
4. **KYC Enable/Disable Guard** — Attempt to enable or disable a PENDING KYC record (should be blocked).
5. **Advanced Filtering** — In the checker queue, apply filters for Level 1, a specific service, and status APPROVED. Verify only matching records are returned. Apply a PENDING-only filter and confirm pending items are not missing.
6. **Auditor Status** — Claim an action as auditor and verify the log shows AUDITORINPROGRESS. Mark the final auditor group and verify the log updates to AUDITORCHECKED.
7. **Checker Level in Log** — Approve an action as a Level 2 checker. Verify the action log records level "2" and not blank.
8. **Segmentation Delete** — Delete a customer segment through the approval flow. After approval confirm the record is no longer retrievable.
9. **USSD Merchant** — Fetch a USSD merchant by ID and verify the `service_name` field is populated.
10. **P&L Account** — Create a service with a GL account not flagged as P&L (should be rejected).
11. **CPS User Creation** — Attempt to create a CPS user with an invalid department (should be rejected). Attempt with an invalid job title (should be rejected).
12. **Image Upload** — Upload a text file with a `.jpg` extension (should be rejected).
13. **Donation Date** — Retrieve an open-ended donation and verify `end_date` is empty, not `0001-01-01`.
14. **First-Time Login** — Log in as a newly created CPS user and verify `is_first_time_login: true` in the profile response.

Please review the attached/linked details and reply with your approval or feedback by **June 1, 2026 at 17:00 EAT**.

Best regards,

**[Requested By]** Dawit — Backend Engineer — future-world3000@proton.me

**[Checked By]** ______________________

**[Reviewed By]** ______________________

**[Approved By]** ______________________
