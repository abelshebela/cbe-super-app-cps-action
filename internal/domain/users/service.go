package users

import (
	"context"
	"time"
	userRepoPort "cbe-super-app-member-users/internal/port/outbound/users"
	"cbe-super-app-member-users/internal/shared"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UserService struct {
	repository userRepoPort.UserRepositoryPort
	logger     utils.Logger
}

func NewUserService(repository userRepoPort.UserRepositoryPort, logger utils.Logger) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
	}
}

func (s *UserService) FetchLinkedAccounts(ctx context.Context, id string) (*LinkedAccountResponse, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == shared.ErrNotFound {
			s.logger.Warnf("User with ID %s not found", id)
			return nil, NewServiceError(common.DefineError.General["NOT_FOUND"])
		}
		s.logger.Errorf("Failed to fetch user with ID %s: %v", id, err)
		return nil, NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	if user.IsDeleted {
		s.logger.Infof("User with ID %s is deleted", id)
		return &LinkedAccountResponse{
			UserID:         user.ID,
			FullName:       user.FullName,
			LinkedAccounts: []LinkedAccountDetail{},
		}, nil
	}

	linkedAccounts, err := s.repository.FindActiveLinkedAccounts(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch linked accounts for user ID %s: %v", id, err)
		return nil, NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}
var linkedAccountsResponse []LinkedAccountDetail
for _, account := range linkedAccounts {
	linkedAccountsResponse = append(linkedAccountsResponse, LinkedAccountDetail{
		AccountNumber:     account.AccountNumber,
		AccountBranchCode: account.AccountBranchCode,
		LinkedBranch:      account.LinkedBranch,
		IsAccountActive:   account.IsAccountActive,
		LinkedStatus:      account.LinkedStatus,
		CurrencyCode:      account.CurrencyCode,
	})
}

response := &LinkedAccountResponse{
	UserID:         user.ID,
	FullName:       user.FullName,
	LinkedAccounts: linkedAccountsResponse,
}


	s.logger.Infof("Successfully fetched active linked accounts for user ID %s", id)
	return response, nil
}

func (s *UserService) GenerateEmailOTP(ctx context.Context, req OTPRequest) (string, error) {
	existingUser, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil && err != shared.ErrNotFound {
		s.logger.Errorf("Email check failed: %v", err)
		return "", NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}
	if existingUser != nil && existingUser.ID != req.UserID {
		return "", NewInternalServiceError(shared.DefineError.OTP["EMAIL_IN_USE"])
	}

	otp := utils.OTPGenerator(6)
	expiresAt := time.Now().Add(5 * time.Minute)

	record := &OTPRecord{
		UserID:    req.UserID,
		Email:     req.Email,
		OTP:       otp,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
	recordRepo := &userRepoPort.OTPRecord{
    UserID:            record.UserID,
	Email:             record.Email,
	OTP:               record.OTP,
	CreatedAt:         record.CreatedAt,
	ExpiresAt:         record.ExpiresAt,
}

	if err := s.repository.StoreOTP(ctx, recordRepo); err != nil {
		s.logger.Errorf("OTP storage failed: %v", err)
		return "", NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	return otp, nil
}

func (s *UserService) VerifyEmailOTP(ctx context.Context, verification OTPVerification) error {
	record, err := s.repository.FindOTP(ctx, verification.UserID, verification.Email)
	if err != nil {
		if err == shared.ErrNotFound {
			s.logger.Errorf("OTP lookup failed for user %s: no documents found", verification.UserID)
			return NewInternalServiceError(shared.DefineError.OTP["INVALID_OTP"])
		}
		s.logger.Errorf("OTP lookup failed for user %s: %v", verification.UserID, err)
		return NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	if record.OTP != verification.OTP {
		s.logger.Warnf("Invalid OTP provided for user %s", verification.UserID)
		return NewInternalServiceError(shared.DefineError.OTP["INVALID_OTP"])
	}

	if time.Now().After(record.ExpiresAt) {
		s.logger.Warnf("Expired OTP for user %s", verification.UserID)
		return NewInternalServiceError(shared.DefineError.OTP["EXPIRED_OTP"])
	}

	if err := s.repository.UpdateUserEmail(ctx, verification.UserID, verification.Email); err != nil {
		s.logger.Errorf("Email update failed for user %s: %v", verification.UserID, err)
		return NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	s.logger.Infof("Email OTP verified successfully for user %s", verification.UserID)
	return nil
}