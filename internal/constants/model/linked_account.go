package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LinkedAccount struct {
	ID                bson.ObjectID              `json:"id,omitempty" bson:"_id,omitempty"`
	UserID            bson.ObjectID              `json:"user_id" bson:"user_id"`
	CustomerNumber    string                     `json:"customer_number" bson:"customer_number"` // cif
	AccountNumber     string                     `json:"account_number" bson:"account_number"`
	AccountHolderName string                     `json:"account_holder_name" bson:"account_holder_name"`
	AccountType       string                     `json:"account_type" bson:"account_type"`
	BranchCode        string                     `json:"branch_code" bson:"branch_code"`
	LinkedStatus      bool                       `json:"linked_status" bson:"linked_status"`
	LastLinkedStatus  bool                       `json:"last_linked_status" bson:"last_linked_status"`
	LinkedAt          time.Time                  `json:"linked_at" bson:"linked_at"`
	LinkedBranch      string                     `json:"linker_branch" bson:"linker_branch"`
	RegistrationType  constants.RegistrationType `json:"registration_type" bson:"registration_type"`
	IsAccountActive   bool                       `json:"is_account_active" bson:"is_account_active"`
	AndOrStatus       bool                       `json:"and_or_status" bson:"and_or_status"`
	AccountBranchCode string                     `json:"account_branch_code" bson:"account_branch_code"`
	CurrencyCode      string                     `json:"currency" bson:"curreny"`
	IsMain            bool                       `json:"is_main" bson:"is_main"` // default: false, first account: true
	MakerAndChecker   types.MakerChecker
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"  bson:"updated_at"`
}
