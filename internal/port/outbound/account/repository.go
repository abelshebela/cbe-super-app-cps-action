package account

import (
	"context"
	
)

type AccountRepositoryPort interface {
	FindAccountUserByID(ctx context.Context, id string) (*AccountUser, error)
	UpdateUserCustomerNumber(ctx context.Context, userID, customerNumber string) error
	CreateLinkedAccount(ctx context.Context, account *LinkedAccount) error
}

type AccountUser struct {
	ID              string
	PhoneNumber     string
	KYCLevel        uint8
	BranchCode      string
	FullName        string
	RegistrationType string
	AndOrStatus     bool
}

type LinkedAccount struct {
	UserID            string
	CustomerNumber    string
	AccountNumber     string
	AccountHolderName string
	AccountType       string
	BranchCode        string
	RegistrationType  string
	AndOrStatus       bool
	CurrencyCode      string
	IsMain            bool
}