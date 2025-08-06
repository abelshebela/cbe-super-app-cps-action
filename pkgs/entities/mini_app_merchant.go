package entities

import (
	"time"

	"cbe-super-app-member-users/pkgs/entities/enums"
	"cbe-super-app-member-users/pkgs/entities/type_definition"
)

type MiniAppMerchant struct {
	ID             string
	MiniAppID      string
	Code           string
	Name           string
	Email          string
	PhoneNumber    string
	AccountNumber  string
	Type           string
	KYCStatus      enums.KYCStatus
	KYCInformation type_definition.KYCInformation
	Branch         []type_definition.BranchInformation
	Enabled        bool
	IsDeleted      bool
	CreatedAt      time.Time
	LastUpdatedAt  time.Time
	DeletedAt      time.Time
}
