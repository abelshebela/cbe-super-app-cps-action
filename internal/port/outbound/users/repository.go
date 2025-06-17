package users

import (
	"context"
	"time"
)

type UserRepositoryPort interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindActiveLinkedAccounts(ctx context.Context, userID string) ([]LinkedAccountDetail, error)
	FindByEmail(ctx context.Context, email string) (*UserEmail, error)
	StoreOTP(ctx context.Context, otp *OTPRecord) error
	FindOTP(ctx context.Context, userID, email string) (*OTPRecord, error)
	UpdateUserEmail(ctx context.Context, userID, email string) error
}

type User struct {
	ID        string
	FullName  string
	IsDeleted bool
}

type UserEmail struct {
	ID    string
	Email string
}

type LinkedAccountDetail struct {
	AccountNumber     string
	AccountBranchCode string
	LinkedBranch      string
	IsAccountActive   bool
	LinkedStatus      bool
	CurrencyCode      string
}

type OTPRecord struct {
	UserID    string
	Email     string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}