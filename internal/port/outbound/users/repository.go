package users

import (
    "context"
    "cbe-super-app-member-users/internal/domain/users"
)

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*users.User, error)
    FindActiveLinkedAccounts(ctx context.Context, userID string) ([]users.LinkedAccountDetail, error)
    FindByEmail(ctx context.Context, email string) (*users.UserEmail, error)
    StoreOTP(ctx context.Context, otp *users.OTPRecord) error
    FindOTP(ctx context.Context, userID, email string) (*users.OTPRecord, error)
    UpdateUserEmail(ctx context.Context, userID, email string) error
}