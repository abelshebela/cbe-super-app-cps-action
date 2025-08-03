package service

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
)

type SpendingService interface {
	// Get(ctx context.Context) ([]dto.Spending, error)
	// GetOne(ctx context.Context, id string) (*dto.Spending, error)
	// Add(ctx context.Context, req dto.SpendingRequest) (*dto.Spending, error)
	// Modify(ctx context.Context, req dto.SpendingRequest) error
	// Remove(ctx context.Context, id string) error
}

type UserService interface {
	DeviceLookup(ctx context.Context, req dto.DeviceLookupResponse) (*dto.DeviceLookupResponse, error)
	PreLogin(ctx context.Context, header dto.Address, phone string) (*dto.DeviceLookupResponse, error)
	ChangePin(ctx context.Context, req dto.ChangePinRequest) error
	VerifyOtp(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOtpResponse, error)
	ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	ResetPin(ctx context.Context, req dto.ResetPinRequest) (*dto.ResetPinResponse, error)
	SetPin(ctx context.Context, req dto.SetPinRequest) (*dto.SetPinResponse, error)
	VerifyForgetPinOtp(ctx context.Context, req dto.VerifyForgetPinOtpRequest) (*dto.VerifyOtpResponse, error)
	UpdateProfilePicture(ctx context.Context, userID string, req dto.UpdateProfilePicture) (string, error)
	UpdateProfileTheme(ctx context.Context, id string, themeType string) (*dto.User, error)
}
