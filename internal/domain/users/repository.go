package users

import (
	"context"
	"time"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindActiveLinkedAccounts(ctx context.Context, id string) ([]LinkedAccountDetail, error)
		StoreOTP(ctx context.Context, otp *OTPRecord) error
	FindOTP(ctx context.Context, userID, email string) (*OTPRecord, error)
	UpdateUserEmail(ctx context.Context, userID, email string) error
	FindByEmail(ctx context.Context, email string) (*UserEmail, error)
}
type OTPRecord struct {
	UserID    string
	Email     string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}