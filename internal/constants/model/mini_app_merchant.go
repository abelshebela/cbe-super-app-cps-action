package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniAppMerchant struct {
	ID                bson.ObjectID             `bson:"_id,omitempty"`
	Code              string                    `bson:"merchant_code"`
	MerchantName      string                    `bson:"merchant_name"`
	MerchantType      string                    `bson:"merchant_type"`
	KYC               types.KYC                 `bson:"kyc"`
	BankAccountNumber string                    `bson:"bank_account_number"`
	Branches          []types.BranchInformation `bson:"branches"`
	Email             string                    `bson:"email"`
	PhoneNumber       string                    `bson:"phone_number"`
	MiniApps          []types.MiniApps          `bson:"mini_apps"`
	Enabled           bool                      `bson:"enabled"`
	IsDeleted         bool                      `bson:"is_deleted"`
	CreatedAt         time.Time                 `bson:"created_at"`
	LastUpdatedAt     time.Time                 `bson:"last_updated_at"`
	DeletedAt         time.Time                 `bson:"deleted_at,omitempty"`
}
