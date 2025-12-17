package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniAppMerchant struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Code         string        `json:"merchant_code" bson:"merchant_code"`
	MerchantName string        `json:"merchant_name" bson:"merchant_name"`
	// KYC               types.KYC                 `json:"kyc" bson:"kyc"`
	SettlementMethod  string                    `json:"settlement_method" bson:"settlement_method"`
	Branches          []types.BranchInformation `json:"branches" bson:"branches"`
	Email             string                    `json:"email" bson:"email"`
	PhoneNumber       string                    `json:"phone_number" bson:"phone_number"`
	Enabled           bool                      `json:"enabled" bson:"enabled"`
	IsDeleted         bool                      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time                 `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at" bson:"updated_at"`
	DeletedAt         *time.Time                `json:"deleted_at" bson:"deleted_at"`
	BankAccountNumber string                    `json:"bank_account_number,omitempty" bson:"bank_account_number,omitempty"`
	// MerchantType      string                    `json:"merchant_type" bson:"merchant_type"`
}
