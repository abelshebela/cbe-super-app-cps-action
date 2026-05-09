package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LogisticsMerchant struct {
	ID                bson.ObjectID `json:"id" bson:"_id"`
	MerchantID        string        `json:"merchant_id" bson:"merchant_id"`
	MerchantType      string        `json:"merchant_type" bson:"merchant_type"`
	SettlementMethod  string        `json:"settlement_method" bson:"settlement_method"`
	MerchantName      string        `json:"merchant_name" bson:"merchant_name"`
	BankAccountNumber string        `json:"bank_account_number" bson:"bank_account_number"`
	Enabled           bool          `json:"enabled" bson:"enabled"`
	IsDeleted         bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" bson:"updated_at"`
	DeletedAt         time.Time     `json:"deleted_at" bson:"deleted_at"`
}
type LogisticsMerchantOracle struct {
	ID                    string    `json:"id" bson:"id"`
	MerchantAccountNumber string    `json:"merchant_account_number" bson:"merchant_account_number"`
	MerchantCode          string    `json:"merchant_code" bson:"merchant_code"`
	MerchantName          string    `json:"merchant_name" bson:"merchant_name"`
	SettlementMethod      string    `json:"settlement_method" bson:"settlement_method"`
	MerchantType          string    `json:"merchant_type" bson:"merchant_type"`
	ContactEmail          string    `json:"contact_email" bson:"contact_email"`
	ContactPhone          string    `json:"contact_phone" bson:"contact_phone"`
	IsEnabled             bool      `json:"is_enabled" bson:"is_enabled"`
	IsDeleted             bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt             time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt        time.Time `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt             time.Time `json:"deleted_at" bson:"deleted_at"`
}
