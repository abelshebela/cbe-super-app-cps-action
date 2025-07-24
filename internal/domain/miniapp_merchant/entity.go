package miniappmerchant

import "time"

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusComplete KYCStatus = "COMPLETE"
	KYCStatusRejected KYCStatus = "REJECTED"
)

type KYCInformation struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type KYC struct {
	Status         KYCStatus      `json:"status"`
	Representative KYCInformation `json:"representative"`
}

type BranchInformation struct {
	BranchCode          string `json:"branch_code"`
	BranchName          string `json:"branch_name"`
	BranchAddress       string `json:"branch_address"`
	BranchOwner         string `json:"branch_owner"`
	BranchAccountNumber string `json:"branch_account_number"`
}

type MiniAppMerchant struct {
	ID                string              `json:"id,omitempty"`
	Code              string              `json:"merchant_code"`
	MerchantName      string              `json:"merchant_name"`
	MerchantType      string              `json:"merchant_type"`
	KYC               KYC                 `json:"kyc"`
	BankAccountNumber string              `json:"bank_account_number"`
	Branches          []BranchInformation `json:"branches"`
	Email             string              `json:"email"`
	PhoneNumber       string              `json:"phone_number"`
	MiniAppIDs        []string            `json:"mini_apps"`
	Enabled           bool                `json:"enabled"`
	IsDeleted         bool                `json:"is_deleted"`
	CreatedAt         time.Time           `json:"created_at"`
	LastModifiedAt    time.Time           `json:"last_modified"`
	DeletedAt         time.Time           `json:"deleted_at"`
}

type CheckMiniAppMerchant struct {
	BankAccountNumber string `json:"bank_account_number"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
}
