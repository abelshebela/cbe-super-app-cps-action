package account

import "context"


type Repository interface {
	FindAccountUserByID(ctx context.Context, id string) (*AccountUser, error)
	UpdateUserCustomerNumber(ctx context.Context, userID, customerNumber string) error
	CreateLinkedAccount(ctx context.Context, account *LinkedAccount) error
}