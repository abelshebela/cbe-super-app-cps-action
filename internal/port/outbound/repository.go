package outbound

import (
	"context"
	"cbe-super-app-member-users/internal/domain/users"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*users.User, error)
	FindActiveLinkedAccounts(ctx context.Context, userID string) ([]users.LinkedAccountDetail, error)
}