package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type KYCInformation struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Phone string `bson:"phone"`
}

type KYC struct {
	Status         string         `bson:"status"`
	Representative KYCInformation `bson:"representative"`
}

type BranchInformation struct {
	BranchCode          string `json:"branch_code"`
	BranchName          string `json:"branch_name"`
	BranchAddress       string `json:"branch_address"`
	BranchOwner         string `json:"branch_owner"`
	BranchAccountNumber string `json:"branch_account_number"`
}

type MiniAppMerchant struct {
	ID                bson.ObjectID       `bson:"_id,omitempty"`
	Code              string              `bson:"merchant_code"`
	MerchantName      string              `bson:"merchant_name"`
	MerchantType      string              `bson:"merchant_type"`
	KYC               KYC                 `bson:"kyc"`
	BankAccountNumber string              `bson:"bank_account_number"`
	Branches          []BranchInformation `bson:"branches"`
	Email             string              `bson:"email"`
	PhoneNumber       string              `bson:"phone_number"`
	MiniApps          []MiniApps          `bson:"mini_apps"`
	Enabled           bool                `bson:"enabled"`
	IsDeleted         bool                `bson:"is_deleted"`
	CreatedAt         time.Time           `bson:"created_at"`
	LastUpdatedAt     time.Time           `bson:"last_updated_at"`
	DeletedAt         time.Time           `bson:"deleted_at,omitempty"`
}

type MiniApps struct {
	ID        bson.ObjectID `bson:"id"`
	Enabled   bool          `bson:"enabled"`
	IsDeleted bool          `bson:"is_deleted"`
}
