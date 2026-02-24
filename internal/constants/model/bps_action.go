package model

import "cbe-super-app-cps-action/internal/constants"

type BPSAction struct {
	ID              string `json:"id"`
	ActionCode      string `json:"action_code"`
	UserInformation struct {
		UserID         string   `json:"user_id"`
		UserCode       string   `json:"user_code"`
		FullName       string   `json:"full_name"`
		AccountNumbers []string `json:"account_numbers"`
		PhoneNumbers   string   `json:"phone_numbers"`
		BranchCode     string   `json:"branch_code"`
	} `json:"user_information"`
	BusinessInformation struct {
		BusinessID   string `json:"business_id"`
		TILLNumber   string `json:"till_number"`
		BusinessName string `json:"business_name"`
	} `json:"business_information"`
	MakerID            string   `json:"maker_id"`
	MakerName          string   `json:"maker_name"`
	MakerReason        string   `json:"maker_reason"`
	MakerPhoneNumber   string   `json:"maker_phone_number"`
	CheckersNeeded     int      `json:"checkers_needed"`
	CheckersApproved   int      `json:"checkers_approved"`
	CheckerID          []string `json:"checker_id"`
	CheckerName        string   `json:"checker_name"`
	CheckerPhoneNumber string   `json:"checker_phone_number"`
	ActionReason       struct {
		ActionType constants.ActionType `json:"action_type"`
		ActionNote string               `json:"action_note"`
		Identifier string               `json:"identifier"`
	} `json:"action_reason"`
	RequestAction      string `json:"request_action"`
	EntityIdentifier   string `json:"entity_identifier"`    // user_code, service_key
	HomeBranch         string `json:"home_branch"`          // maker, checker working branch
	AccountBranchCode  string `json:"account_branch_code"`  // applier/member account
	DistrictCode       string `json:"district_code"`        // home branch district code
	BranchCode         string `json:"branch_code"`          // maker, checker branch(IFB or CB)
	LinkedDistrictCode string `json:"linked_district_code"` // applier/member district code
	AccountNumber      string `json:"account_number"`
	AccountHolderName  string `json:"account_holder_name"`
	Auditors           struct {
		AuditorsRequired int      `json:"auditors_required"`
		AuditorID        []string `json:"auditor_id"`
		Audited          bool     `json:"audited"`
		AuditorApproval  bool     `json:"auditor_approval"`
		Reason           string   `json:"reason"`
	} `json:"auditors"`
	ServiceName    string                 `json:"service_name"`
	CurrentAction  map[string]interface{} `json:"current_action"`  // access list change
	PreviousAction map[string]interface{} `json:"previous_action"` // before the change
	VerifiedAt     string                 `json:"verified_at"`
	Status         string                 `json:"status"`
	CreatedAt      string                 `json:"created_at"`
	LastModifiedAt string                 `json:"last_modified_at"`
}
