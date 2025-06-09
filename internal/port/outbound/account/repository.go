package account

import (
    "context"
    // "cbe-super-app-member-users/internal/domain/account"
)

type AccountRepository interface {
    FindAccountUserByID(ctx context.Context, id string) (AccountUser, error)
    UpdateUserCustomerNumber(ctx context.Context, userID, customerNumber string) error
    CreateLinkedAccount(ctx context.Context, account *LinkedAccount) error
}