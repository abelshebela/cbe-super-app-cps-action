package users_application

import (
	"context"
	
	"mime/multipart"

	"cbe-super-app-member-users/internal/application/dto"
	domainUsers "cbe-super-app-member-users/internal/domain/users"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"
)

func (h UsersHandler) FetchLinkedAccounts(ctx context.Context, userID string) (*userPort.LinkedAccountResponse, error) {
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



func (h UsersHandler) GenerateEmailOTP(ctx context.Context, req userPort.OTPRequest) (string, error) {
	domainReq := domainUsers.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	}
	return h.userService.GenerateEmailOTP(ctx, domainReq)
}

func (h UsersHandler) VerifyEmailOTP(ctx context.Context, verification userPort.OTPVerification) error {
	domainVerification := domainUsers.OTPVerification{
		UserID: verification.UserID,
		Email:  verification.Email,
		OTP:    verification.OTP,
	}
	return h.userService.VerifyEmailOTP(ctx, domainVerification)
}

func (h UsersHandler) HandleFetchLinkedAccounts(ctx context.Context, req dto.FetchLinkedAccountsRequest) (*userPort.LinkedAccountResponse, error) {
	return h.FetchLinkedAccounts(ctx, req.UserID)
}

func (h UsersHandler) HandleGenerateEmailOTP(ctx context.Context, req dto.GenerateOTPRequest) (*dto.GenerateOTPResponse, error) {
	portReq := userPort.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	}
	

	otp, err := h.GenerateEmailOTP(ctx, portReq)
	if err != nil {
		return nil, err
	}

	return &dto.GenerateOTPResponse{OTP: otp}, nil
}

func (h UsersHandler) HandleVerifyEmailOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	portVerification := userPort.OTPVerification{
		UserID: req.UserID,
		Email:  req.Email,
		OTP:    req.OTP,
	}

	err := h.VerifyEmailOTP(ctx, portVerification)
	if err != nil {
		return nil, err
	}

	return &dto.VerifyOTPResponse{Success: true}, nil
}

func (h UsersHandler) UpdateProfilePicture(ctx context.Context, id string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	imageURL, err := h.userService.UpdateProfilePicture(ctx, id, file, fileHeader)
	if err != nil {
		return "", err
	}
	return imageURL, nil
}

func (h UsersHandler) UnlinkDevice(ctx context.Context, userID string, deviceID string) error {
	err := h.userService.UnlinkDevice(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	return nil


}