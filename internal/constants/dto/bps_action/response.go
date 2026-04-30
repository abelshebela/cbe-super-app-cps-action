package bps_action

import (
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

// BPSActionDetailResponse is the enriched payload returned by the
// /actions/by-action-code/{action_code} endpoint. It bundles the action document
// with the customer's member detail, currently linked accounts, and previously
// linked (unlinked / archived) accounts.
//
// MemberDetail / LinkedAccounts / UnlinkedAccounts are nil-or-empty when the
// underlying user is not resolvable (e.g. the action does not target a customer).
//
// LinkedAccounts mirrors the shape returned by the customer module's Oracle
// detail flow (customer_dto.LinkedAccount) so the FE can reuse the same renderer.
type BPSActionDetailResponse struct {
	Action           *bps_model.BPSAction                 `json:"action"`
	MemberDetail     *customer_dto.CustomerDetailResponse `json:"member_detail,omitempty"`
	LinkedAccounts   []customer_dto.LinkedAccount         `json:"linked_accounts"`
	UnlinkedAccounts []model.ArchivedLinkedAccount        `json:"unlinked_accounts"`
}
