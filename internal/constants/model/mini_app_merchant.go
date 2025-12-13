package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniAppMerchant struct {
	ID                bson.ObjectID `bson:"_id" json:"id"`
	MerchantName      string        `bson:"merchant_name" json:"merchant_name"`
	MerchantCode      string        `bson:"merchant_code" json:"merchant_code" `
	SettlementMethod  string        `bson:"settlement_method" json:"settlement_method"`
	BankAccountNumber string        `bson:"bank_account_number" json:"bank_account_number"`
	Enabled           bool          `bson:"enabled" json:"enabled"`
	IsDeleted         bool          `bson:"is_deleted" json:"is_deleted"`
	KYC               types.KYC     `bson:"kyc" json:"kyc"`
	CreatedAt         time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time    `bson:"deleted_at" json:"deleted_at"`
}
