package bps_action

import (
	"time"

	customer_dto "cbe-super-app-cps-action/internal/constants/dto/customer"
	bps_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

type BPSActionCountResponse struct {
	Pending    int `json:"total_pending" bson:"pending"`
	Approved   int `json:"total_approved" bson:"approved"`
	Rejected   int `json:"total_rejected" bson:"rejected"`
	Inprogress int `json:"total_inprogress_audit" bson:"inprogress"`
	Completed  int `json:"total_completed_audit" bson:"completed"`
}

type BPSAuditorActionCountResponse struct {
	AllAction  int `json:"all_action"`
	UnAudited  int `json:"un_audited"`
	Inprogress int `json:"inprogress"`
	Audited    int `json:"audited"`
}

type BPSCheckerActionCountResponse struct {
	AllAction int `json:"all_action"`
	Pending   int `json:"pending"`
	Approved  int `json:"approved"`
	Rejected  int `json:"rejected"`
}

// BPSActionUserInfo is the trimmed, unified shape used to surface checker /
// auditor / maker user records on a BPSAction detail response. The same shape
// is used regardless of whether the underlying record was resolved against the
// bps_user collection or fell back to cps_user; the Source field tells the FE
// which collection the record came from.
type BPSActionUserInfo struct {
	ID          string   `json:"id"`
	UserCode    string   `json:"user_code,omitempty"`
	FullName    string   `json:"full_name,omitempty"`
	UserName    string   `json:"username,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	JobTitle    string   `json:"job_title,omitempty"`
	Role        string   `json:"role,omitempty"`
	Department  string   `json:"department,omitempty"` // CPS users only
	BranchCode  []string `json:"branch_code,omitempty"`
	BranchName  string   `json:"branch_name,omitempty"`
	Source      string   `json:"source,omitempty"` // "bps_user" or "cps_user"
}

// CheckerLevelInfo describes the approval status of a single checker slot.
type CheckerLevelInfo struct {
	Level      int        `json:"level"`
	CheckerID  string     `json:"checker_id"`
	Status     string     `json:"status"`               // APPROVED | REJECTED | PENDING
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
}

// BPSActionDetailResponse is the enriched payload returned by the
// /actions/by-action-code/{action_code} endpoint. It bundles the action document
// with the customer's member detail, currently linked accounts, previously
// linked (unlinked / archived) accounts, and the resolved checker / auditor
// user records.
//
// All enrichment fields are best-effort: empty slices / nil when the underlying
// records cannot be resolved.
//
// LinkedAccounts mirrors the shape returned by the customer module's Oracle
// detail flow (customer_dto.LinkedAccount) so the FE can reuse the same renderer.
type BPSActionDetailResponse struct {
	Action                 *bps_model.BPSAction                 `json:"action"`
	MemberDetail           *customer_dto.CustomerDetailResponse `json:"member_detail,omitempty"`
	LinkedAccounts         []customer_dto.LinkedAccount         `json:"linked_accounts"`
	UnlinkedAccounts       []model.ArchivedLinkedAccount        `json:"unlinked_accounts"`
	// Resolved user records
	Maker                  *BPSActionUserInfo                   `json:"maker,omitempty"`
	Checkers               []BPSActionUserInfo                  `json:"checkers"`
	Auditors               []BPSActionUserInfo                  `json:"auditors"`
	// Per-checker-level approval status derived from CheckerID + CheckerTime
	CheckerLevels          []CheckerLevelInfo                   `json:"checker_levels"`
	// Surfaced rejection reasons
	RejectionReason        string                               `json:"rejection_reason,omitempty"`
	AuditorRejectionReason string                               `json:"auditor_rejection_reason,omitempty"`
}
