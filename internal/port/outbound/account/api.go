package account

import "context"

type AccountAPIClient interface {
    LookupAccountByPhone(ctx context.Context, phoneNumber string) (bool, error)
}