package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniAppMerchant struct {
	ID                bson.ObjectID             `json:"id,omitempty" bson:"_id,omitempty"`
	Code              string                    `json:"merchant_code" bson:"merchant_code"`
	MerchantName      string                    `json:"merchant_name" bson:"merchant_name"`
	MerchantType      string                    `json:"merchant_type" bson:"merchant_type"`
	KYC               types.KYC                 `json:"kyc" bson:"kyc"`
	BankAccountNumber string                    `json:"bank_account_number" bson:"bank_account_number"`
	Branches          []types.BranchInformation `json:"branches" bson:"branches"`
	Email             string                    `json:"email" bson:"email"`
	PhoneNumber       string                    `json:"phone_number" bson:"phone_number"`
	Enabled           bool                      `json:"enabled" bson:"enabled"`
	IsDeleted         bool                      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time                 `json:"created_at" bson:"created_at"`
	LastModifiedAt    time.Time                 `json:"last_modified" bson:"last_updated_at"`
	DeletedAt         time.Time                 `json:"deleted_at" bson:"deleted_at"`
}
