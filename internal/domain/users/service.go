package users

import (
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UserService struct {
	repository UserRepository
	logger     utils.Logger
}

func NewUserService(repository UserRepository, logger utils.Logger) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
	}
}

func (s *UserService) ActiveLinkedAccounts(ctx context.Context, id string) (*LinkedAccountResponse, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
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

	response := &LinkedAccountResponse{
		UserID:         user.ID,
		FullName:       user.FullName,
		LinkedAccounts: linkedAccounts,
	}

	s.logger.Infof("Successfully fetched active linked accounts for user ID %s", id)
	return response, nil
}

func (s *UserService) GenerateEmailOTP(ctx context.Context, req OTPRequest) (string, error) {
	existingUser, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil && err != ErrNotFound {
		s.logger.Errorf("Email check failed: %v", err)
		// return "", NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}
	if existingUser != nil && existingUser.ID != req.UserID {
		return "", NewServiceError(common.ErrorDefinition{
			Code:    "EMAIL_IN_USE",
			Message: "Email address is already registered",
		})
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

	if err := s.repository.StoreOTP(ctx, record); err != nil {
		s.logger.Errorf("OTP storage failed: %v", err)
		return "", NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	return otp, nil
}

func (s *UserService) VerifyEmailOTP(ctx context.Context, verification OTPVerification) error {
	// s.logger.Errorf("OTP lookup : %v",verification )
	record, err := s.repository.FindOTP(ctx, verification.UserID, verification.Email)
	// s.logger.Errorf("OTP lookup : %v", record)
	if err != nil {
		if err == ErrNotFound {
			return NewInvalidOTPError()
		}
		s.logger.Errorf("OTP lookup failed: %v", err)
		return NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	if record.OTP != verification.OTP {
		return NewInvalidOTPError()
	}

	if time.Now().After(record.ExpiresAt) {
		// s.logger.Errorf("vvvv", time.Now().After(record.ExpiresAt))
		return NewExpiredOTPError()
	}

	if err := s.repository.UpdateUserEmail(ctx, verification.UserID, verification.Email); err != nil {
		s.logger.Errorf("Email update failed: %v", err)
		return NewServiceError(common.DefineError.General["UNHANDLED_SERVER_ERROR"])
	}

	return nil
}