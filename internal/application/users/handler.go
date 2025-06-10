package users

import (
	"context"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"
	domainUsers "cbe-super-app-member-users/internal/domain/users"
)

type ApplicationHandler struct {
	userService *domainUsers.UserService
}

func NewApplicationHandler(userService *domainUsers.UserService) *ApplicationHandler {
	return &ApplicationHandler{
		userService: userService,
	}
}

type FetchLinkedAccountsRequest struct {
	UserID string `json:"user_id"`
}

type GenerateOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	OTP    string `json:"otp"`
}

type GenerateOTPResponse struct {
	OTP string `json:"otp"`
}

type VerifyOTPResponse struct {
	Success bool `json:"success"`
}

func (h *ApplicationHandler) FetchLinkedAccounts(ctx context.Context, userID string) (*userPort.LinkedAccountResponse, error) {
	domainResponse, err := h.userService.FetchLinkedAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	portAccounts := make([]userPort.LinkedAccountDetail, len(domainResponse.LinkedAccounts))
	for i, acc := range domainResponse.LinkedAccounts {
		portAccounts[i] = userPort.LinkedAccountDetail{
			AccountNumber:     acc.AccountNumber,
			AccountBranchCode: acc.AccountBranchCode,
			LinkedBranch:      acc.LinkedBranch,
			IsAccountActive:   acc.IsAccountActive,
			LinkedStatus:      acc.LinkedStatus,
			CurrencyCode:      acc.CurrencyCode,
		}
	}
	
	return &userPort.LinkedAccountResponse{
		UserID:         domainResponse.UserID,
		FullName:       domainResponse.FullName,
		LinkedAccounts: portAccounts,
	}, nil
}

func (h *ApplicationHandler) GenerateEmailOTP(ctx context.Context, req userPort.OTPRequest) (string, error) {
	domainReq := domainUsers.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	}
	return h.userService.GenerateEmailOTP(ctx, domainReq)
}

func (h *ApplicationHandler) VerifyEmailOTP(ctx context.Context, verification userPort.OTPVerification) error {
	domainVerification := domainUsers.OTPVerification{
		UserID: verification.UserID,
		Email:  verification.Email,
		OTP:    verification.OTP,
	}
	return h.userService.VerifyEmailOTP(ctx, domainVerification)
}

func (h *ApplicationHandler) HandleFetchLinkedAccounts(ctx context.Context, req FetchLinkedAccountsRequest) (*userPort.LinkedAccountResponse, error) {
	return h.FetchLinkedAccounts(ctx, req.UserID)
}

func (h *ApplicationHandler) HandleGenerateEmailOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	portReq := userPort.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	}
	
	otp, err := h.GenerateEmailOTP(ctx, portReq)
	if err != nil {
		return nil, err
	}
	
	return &GenerateOTPResponse{OTP: otp}, nil
}

func (h *ApplicationHandler) HandleVerifyEmailOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	portVerification := userPort.OTPVerification{
		UserID: req.UserID,
		Email:  req.Email,
		OTP:    req.OTP,
	}
	
	err := h.VerifyEmailOTP(ctx, portVerification)
	if err != nil {
		return nil, err
	}
	
	return &VerifyOTPResponse{Success: true}, nil
}