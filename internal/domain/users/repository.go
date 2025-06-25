package users

import (
	"context"
	"errors"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindActiveLinkedAccounts(ctx context.Context, userID string) ([]LinkedAccountDetail, error)
	FindByEmail(ctx context.Context, email string) (*UserEmail, error)
	StoreOTP(ctx context.Context, otp *OTPRecord) error
	FindOTP(ctx context.Context, userID, email string) (*OTPRecord, error)
	UpdateUserEmail(ctx context.Context, userID, email string) error
	UpdateProfileImageURL(ctx context.Context, id string, imageURL string) error
	DeleteOtp(ctx context.Context, id string) (error )
	UnlinkDevice(ctx context.Context, userID string, deviceID string)error

}



var (
	ErrNotFound = errors.New("not found")
)
