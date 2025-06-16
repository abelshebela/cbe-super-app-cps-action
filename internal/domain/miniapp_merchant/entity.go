package miniappmerchant

import "time"

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "PENDING"
	KYCStatusApproved KYCStatus = "APPROVED"
	KYCStatusRejected KYCStatus = "REJECTED"
)

type KYCInformation struct {
	Name        string
	Email       string
	PhoneNumber string
}

type BranchInformation struct {
	ID            string
	Code          string
	Name          string
	Address       string
	Owner         string
	AccountNumber string
}

type MiniAppMerchant struct {
	ID             string
	MiniAppID      string
	Code           string
	Name           string
	Email          string
	PhoneNumber    string
	AccountNumber  string
	Type           string
	KYCStatus      KYCStatus
	KYCInformation KYCInformation
	Branch         []BranchInformation
	Enabled        bool
	IsDeleted      bool
	CreatedAt      time.Time
	LastUpdatedAt  time.Time
	DeletedAt      time.Time
}
