package logistics_merchant_dto

import "time"

type LogisticsMerchantResponse struct {
	ID                string `json:"id" bson:"_id"`
	MerchantID        string `json:"merchant_id" bson:"merchant_id"`
	MerchantType      string `json:"merchant_type" bson:"merchant_type"`
	SettlementMethod  string `json:"settlement_method" bson:"settlement_method"`
	MerchantName      string `json:"merchant_name" bson:"merchant_name"`
	BankAccountNumber string `json:"bank_account_number" bson:"bank_account_number"`
	// Email             string    `json:"email" bson:"email"`
	// PhoneNumber       string    `json:"phone_number" bson:"phone_number"`
	Enabled   bool      `json:"enabled" bson:"enabled"`
	IsDeleted bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at"`
}

type CreateLogisticsMerchantRequest struct {
	MerchantID          string `json:"merchant_id" bson:"merchant_id"`
	MerchantType        string `json:"merchant_type" bson:"merchant_type"`
	SettlementMethod    string `json:"settlement_method" bson:"settlement_method"`
	MerchantName        string `json:"merchant_name" bson:"merchant_name"`
	BankAccountNumber   string `json:"bank_account_number" bson:"bank_account_number"`
	IsLogisticsMerchant *bool  `json:"is_logistics_merchant" bson:"is_logistics_merchant"`
	// Email             string `json:"email" bson:"email"`
	// PhoneNumber       string `json:"phone_number" bson:"phone_number"`
}

type UpdateLogisticsMerchantRequest struct {
	MerchantID          string `json:"merchant_id" bson:"merchant_id"`
	MerchantType        string `json:"merchant_type" bson:"merchant_type"`
	SettlementMethod    string `json:"settlement_method" bson:"settlement_method"`
	MerchantName        string `json:"merchant_name" bson:"merchant_name"`
	BankAccountNumber   string `json:"bank_account_number" bson:"bank_account_number"`
	IsLogisticsMerchant *bool  `json:"is_logistics_merchant" bson:"is_logistics_merchant"`
	// Email             string `json:"email" bson:"email"`
	// PhoneNumber       string `json:"phone_number" bson:"phone_number"`
}
