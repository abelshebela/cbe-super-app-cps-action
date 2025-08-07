package service

import (
	"context"
	"mime/multipart"

	"cbe-super-app-member-auth/internal/constants/dto"
)

type UserService interface {
	DeviceLookup(ctx context.Context, req dto.DeviceLookupRequest) (*dto.DeviceLookupResponse, error)
	PreLogin(ctx context.Context, phone string) (*dto.DeviceLookupResponse, error)
	ChangePin(ctx context.Context, req dto.ChangePinRequest) error
	VerifyOtp(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOtpResponse, error)
	ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	ResetPin(ctx context.Context, req dto.ResetPinRequest) (*dto.ResetPinResponse, error)
	SetPin(ctx context.Context, req dto.SetPinRequest) (*dto.SetPinResponse, error)
	VerifyForgetPinOtp(ctx context.Context, req dto.VerifyForgetPinOtpRequest) (*dto.VerifyOtpResponse, error)
	UpdateProfilePicture(ctx context.Context, userID string, req dto.UpdateProfilePicture, file multipart.File) error
	UpdateProfileTheme(ctx context.Context, id string, themeType string) error
}
