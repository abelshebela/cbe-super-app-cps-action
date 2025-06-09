package users

import (
	"context"
	"cbe-super-app-member-users/internal/domain/users"
)

type ApplicationHandler interface {
    FetchLinkedAccounts(ctx context.Context, req FetchLinkedAccountsRequest) (*LinkedAccountResponseDto, error)
    GenerateEmailOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error)
    VerifyEmailOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error)
}

type applicationHandler struct {
	domainService *users.UserService
}

func NewApplicationHandler(domainService *users.UserService) ApplicationHandler {
	return &applicationHandler{domainService: domainService}
}

func (h *applicationHandler) FetchLinkedAccounts(ctx context.Context, req FetchLinkedAccountsRequest) (*LinkedAccountResponseDto, error) {
	response, err := h.domainService.ActiveLinkedAccounts(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return ToLinkedAccountResponseDto(response), nil
}

func (h *applicationHandler) GenerateEmailOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	otp, err := h.domainService.GenerateEmailOTP(ctx, users.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &GenerateOTPResponse{OTP: otp}, nil
}

func (h *applicationHandler) VerifyEmailOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	if err := h.domainService.VerifyEmailOTP(ctx, users.OTPVerification{
		UserID: req.UserID,
		Email:  req.Email,
		OTP:    req.OTP,
	}); err != nil {
		return nil, err
	}
	return &VerifyOTPResponse{
		Success: true,
		Email:   req.Email,
	}, nil
}