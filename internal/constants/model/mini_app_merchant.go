package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusComplete KYCStatus = "COMPLETE"
	KYCStatusRejected KYCStatus = "REJECTED"
)

type KYCInformation struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type KYC struct {
	Status         KYCStatus      `json:"status" bson:"status"`
	Representative KYCInformation `json:"representative" bson:"representative"`
}

type BranchInformation struct {
	BranchCode          string `json:"branch_code"`
	BranchName          string `json:"branch_name"`
	BranchAddress       string `json:"branch_address"`
	BranchOwner         string `json:"branch_owner"`
	BranchAccountNumber string `json:"branch_account_number"`
}

type MiniApps struct {
	ID        string `json:"id" bson:"id"`
	Enabled   bool   `json:"enabled" bson:"enabled"`
	IsDeleted bool   `json:"is_deleted" bson:"is_deleted"`
}

type MiniAppMerchant struct {
	ID                bson.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	Code              string              `json:"merchant_code" bson:"merchant_code"`
	MerchantName      string              `json:"merchant_name" bson:"merchant_name"`
	MerchantType      string              `json:"merchant_type" bson:"merchant_type"`
	KYC               KYC                 `json:"kyc" bson:"kyc"`
	BankAccountNumber string              `json:"bank_account_number" bson:"bank_account_number"`
	Branches          []BranchInformation `json:"branches" bson:"branches"`
	Email             string              `json:"email" bson:"email"`
	PhoneNumber       string              `json:"phone_number" bson:"phone_number"`
	Enabled           bool                `json:"enabled" bson:"enabled"`
	IsDeleted         bool                `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time           `json:"created_at" bson:"created_at"`
	LastModifiedAt    time.Time           `json:"last_modified" bson:"last_updated_at"`
	DeletedAt         time.Time           `json:"deleted_at" bson:"deleted_at"`
}

type CheckMiniAppMerchant struct {
	BankAccountNumber string `json:"bank_account_number"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
}
type MiniAppMerchantExistOptions struct {
	ExcludeID string
}
