package model

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BPSAction struct {
	ID                bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ActionCode        string        `json:"action_code" bson:"action_code"`
	IsAuditorApproved bool          `json:"is_auditor_approved" bson:"is_auditor_approved"`
	UserInformation   struct {
		UserID         string   `json:"user_id,omitempty" bson:"user_id,omitempty"`
		UserCode       string   `json:"user_code" bson:"user_code"`
		FullName       string   `json:"full_name" bson:"full_name"`
		AccountNumbers []string `json:"account_numbers" bson:"account_numbers"`
		PhoneNumbers   string   `json:"phone_numbers" bson:"phone_numbers"`
		BranchCode     string   `json:"branch_code" bson:"branch_code"`
	} `json:"user_information" bson:"user_information"`
	BusinessInformation struct {
		BusinessID   bson.ObjectID `json:"business_id" bson:"business_id"`
		TILLNumber   string        `json:"till_number" bson:"till_number"`
		BusinessName string        `json:"business_name" bson:"business_name"`
	} `json:"business" bson:"business"`
	CheckersNeeded     int      `json:"checkers_needed" bson:"checkers_needed"`
	CheckersApproved   int      `json:"checkers_approved" bson:"checkers_approved"`
	CheckerID          []string `json:"checker_id" bson:"checker_id"`
	MakerID            string   `json:"maker_user" bson:"maker_user"`
	MakerName          string   `json:"maker_name" bson:"maker_name"`
	MakerReason        string   `json:"maker_reason" bson:"maker_reason"`
	MakerPhoneNumber   string   `json:"maker_phone_number" bson:"maker_phone_number"`
	CheckerName        string   `json:"checker_name" bson:"checker_name"`
	CheckerPhoneNumber string   `json:"checker_phone_number" bson:"checker_phone_number"`
	ActionReason       struct {
		ActionType constants.ActionType `json:"action_type" bson:"action_type"` //  reject, enable, disable
		ActionNote string               `json:"action_note" bson:"action_note"`
		Identifier string               `json:"identifier" bson:"identifier"`
	} `json:"action_reason" bson:"action_reason"`
	MakerMID        string   `json:"maker_mid" bson:"maker_mid"`
	CheckerMID      []string `json:"checker_mid" bson:"checker_mid"`
	AuditorMID      []string `json:"auditor_mid" bson:"auditor_mid"`
	CheckerNameList []string `json:"checker_name_list" bson:"checker_name_list"`
	AuditorNameList []string `json:"auditor_name_list" bson:"auditor_name_list"`
	Auditors        struct {
		AuditorName        string   `json:"auditor_name" bson:"auditor_name"`
		AuditorPhoneNumber string   `json:"auditor_phone_number" bson:"auditor_phone_number"`
		AuditorsRequired   int      `json:"auditors_required" bson:"auditors_required"`
		AuditorID          []string `json:"auditor_id" bson:"auditor_id"`
		Audited            bool     `json:"audited" bson:"audited"`
		AuditorApproval    bool     `json:"auditor_approval" bson:"auditor_approval"`
		Reason             string   `json:"reason" bson:"reason"`
	} `json:"auditors" bson:"auditors"`
	CheckerTime        []time.Time             `json:"checker_time" bson:"checker_time"`
	AuditorTime        []time.Time             `json:"auditor_time" bson:"auditor_time"`
	RequestAction      constants.RequestAction `json:"request_action" bson:"request_action"`
	EntityIdentifyer   string                  `json:"value" bson:"value"`                               // user_code, service_key
	HomeBranch         string                  `json:"home_branch" bson:"home_branch"`                   // maker, checker working branch
	AccountBranchCode  string                  `json:"account_branch_code" bson:"account_branch_code"`   // applier/member account
	DistrictCode       string                  `json:"district_code" bson:"district_code"`               // home branch district code
	BranchCode         string                  `json:"branch_code" bson:"branch_code"`                   // maker, checker branch(IFB or CB)
	LinkedDistrictCode string                  `json:"linked_district_code" bson:"linked_district_code"` // applier/member district code
	AccountNumber      string                  `json:"account_number" bson:"account_number"`
	AccountHolderName  string                  `json:"account_holder_name" bson:"account_holder_name"`
	ServiceName        string                  `json:"service_name" bson:"service_name"`
	CurrentAction      interface{}             `json:"current_action" bson:"current_action"`   // access list chang
	PreviousAction     interface{}             `json:"previous_action" bson:"previous_action"` // befor the change
	VerifiedAt         *time.Time              `json:"time,omitempty" bson:"time,omitempty"`
	Status             string                  `json:"status" bson:"status"`
	CustomerBarred     bool                    `json:"customer_barred" bson:"customer_barred"`
	CreatedAt          time.Time               `json:"created_at" bson:"created_at"`
	LastModifiedAt     time.Time               `json:"last_modified_at" bson:"last_modified_at"`
}
