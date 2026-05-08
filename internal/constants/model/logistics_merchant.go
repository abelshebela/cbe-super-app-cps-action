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
	ID                    string    `json:"id" db:"id"`
	MerchantAccountNumber string    `json:"merchant_account_number" db:"merchant_account_number"`
	MerchantCode          string    `json:"merchant_code" db:"merchant_code"`
	MerchantName          string    `json:"merchant_name" db:"merchant_name"`
	SettlementMethod      string    `json:"settlement_method" db:"settlement_method"`
	MerchantType          string    `json:"merchant_type" db:"merchant_type"`
	ContactEmail          string    `json:"contact_email" db:"contact_email"`
	ContactPhone          string    `json:"contact_phone" db:"contact_phone"`
	IsEnabled             bool      `json:"is_enabled" db:"is_enabled"`
	IsDeleted             bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	LastModifiedAt        time.Time `json:"last_modified_at" db:"last_modified_at"`
	DeletedAt             time.Time `json:"deleted_at" db:"deleted_at"`
}
