package account

import "context"

type AccountPort interface {
	CreateAccount(ctx context.Context, userID string) (*AccountCreationResult, error)
}

type AccountCreationResult struct {
	CustomerNumber string
	AccountNumber  string
}