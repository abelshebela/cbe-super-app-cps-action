package users_test

import (
	"context"
	"testing"
	"time"

	"cbe-super-app-member-users/internal/domain/users"
	"cbe-super-app-member-users/internal/shared"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	userPortMocks "cbe-super-app-member-users/internal/port/outbound/users/mocks"
	userPort "cbe-super-app-member-users/internal/port/outbound/users"
)

func MockPortUser() *userPort.User {
	return &userPort.User{
		ID:        "6644c8e37f41d2c9c1c293fa",
		FullName:  "John Doe",
		IsDeleted: false,
	}
}

func MockPortOTPRecord(userID, email, otp string) *userPort.OTPRecord {
	return &userPort.OTPRecord{
		UserID:    userID,
		Email:     email,
		OTP:       otp,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func MockPortLinkedAccounts() []userPort.LinkedAccountDetail {
	return []userPort.LinkedAccountDetail{
		{
			AccountNumber: "9876543210",
			CurrencyCode:  "EUR",
		},
		{
			AccountNumber: "1122334455",
			CurrencyCode:  "GBP",
		},
	}
}

func TestVerifyEmailOTP(t *testing.T) {
	user := MockPortUser()
	userID := user.ID
	email := "test@example.com"
	otp := "123456"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		verification := users.OTPVerification{
			UserID: userID,
			Email:  email,
			OTP:    otp,
		}
		otpRecord := MockPortOTPRecord(userID, email, otp)

		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(otpRecord, nil)
		mockRepo.EXPECT().UpdateUserEmail(gomock.Any(), userID, email).Return(nil)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.NoError(t, err)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		invalidVerification := users.OTPVerification{
			UserID: "invalid-id",
			Email:  email,
			OTP:    otp,
		}
		mockRepo.EXPECT().FindOTP(gomock.Any(), "invalid-id", email).Return(nil, shared.ErrNotFound)

		err := service.VerifyEmailOTP(context.Background(), invalidVerification)
		require.Error(t, err)
		assert.Equal(t, shared.DefineError.OTP["INVALID_OTP"].Code, err.(users.ServiceError).Code)
	})

	t.Run("invalid OTP", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		verification := users.OTPVerification{
			UserID: userID,
			Email:  email,
			OTP:    otp,
		}
		invalidOTPRecord := MockPortOTPRecord(userID, email, "654321")
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(invalidOTPRecord, nil)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, shared.DefineError.OTP["INVALID_OTP"].Code, err.(users.ServiceError).Code)
	})

	t.Run("expired OTP", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		verification := users.OTPVerification{
			UserID: userID,
			Email:  email,
			OTP:    otp,
		}
		expiredOTPRecord := MockPortOTPRecord(userID, email, otp)
		expiredOTPRecord.CreatedAt = time.Now().Add(-10 * time.Minute)
		expiredOTPRecord.ExpiresAt = time.Now().Add(-5 * time.Minute)
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(expiredOTPRecord, nil)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, shared.DefineError.OTP["EXPIRED_OTP"].Code, err.(users.ServiceError).Code)
	})

	t.Run("OTP not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		verification := users.OTPVerification{
			UserID: userID,
			Email:  email,
			OTP:    otp,
		}
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(nil, shared.ErrNotFound)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, shared.DefineError.OTP["INVALID_OTP"].Code, err.(users.ServiceError).Code)
	})

	t.Run("update email failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		verification := users.OTPVerification{
			UserID: userID,
			Email:  email,
			OTP:    otp,
		}
		otpRecord := MockPortOTPRecord(userID, email, otp)
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(otpRecord, nil)
		mockRepo.EXPECT().UpdateUserEmail(gomock.Any(), userID, email).Return(assert.AnError)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code, err.(users.ServiceError).Code)
	})
}

func TestFetchLinkedAccounts(t *testing.T) {
	user := MockPortUser()
	userID := user.ID
	linkedAccounts := MockPortLinkedAccounts()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(gomock.Any(), userID).Return(linkedAccounts, nil)

		response, err := service.FetchLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, userID, response.UserID)
		assert.Equal(t, user.FullName, response.FullName)
		assert.Len(t, response.LinkedAccounts, 2)
		assert.Equal(t, "9876543210", response.LinkedAccounts[0].AccountNumber)
		assert.Equal(t, "EUR", response.LinkedAccounts[0].CurrencyCode)
		assert.Equal(t, "1122334455", response.LinkedAccounts[1].AccountNumber)
		assert.Equal(t, "GBP", response.LinkedAccounts[1].CurrencyCode)
	})

	t.Run("invalid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		mockRepo.EXPECT().FindByID(gomock.Any(), "invalid-id").Return(nil, shared.ErrNotFound)

		_, err := service.FetchLinkedAccounts(context.Background(), "invalid-id")
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["NOT_FOUND"].Code, err.(users.ServiceError).Code)
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, assert.AnError)

		_, err := service.FetchLinkedAccounts(context.Background(), userID)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code, err.(users.ServiceError).Code)
	})

	t.Run("no linked accounts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(gomock.Any(), userID).Return([]userPort.LinkedAccountDetail{}, nil)

		response, err := service.FetchLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.LinkedAccounts, 0)
	})

	t.Run("deleted user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		deletedUser := &userPort.User{
			ID:        userID,
			FullName:  user.FullName,
			IsDeleted: true,
		}
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(deletedUser, nil)

		response, err := service.FetchLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.LinkedAccounts, 0)
	})
}

func TestGenerateEmailOTP(t *testing.T) {
	user := MockPortUser()
	userID := user.ID
	email := "test@example.com"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		otpRequest := users.OTPRequest{
			UserID: userID,
			Email:  email,
		}

		mockRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(nil, shared.ErrNotFound)
		mockRepo.EXPECT().StoreOTP(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, otp *userPort.OTPRecord) error {
				assert.Equal(t, userID, otp.UserID)
				assert.Equal(t, email, otp.Email)
				assert.Len(t, otp.OTP, 6)
				return nil
			},
		)

		otp, err := service.GenerateEmailOTP(context.Background(), otpRequest)
		require.NoError(t, err)
		assert.NotEmpty(t, otp)
		assert.Len(t, otp, 6)
	})

	t.Run("email already in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		otpRequest := users.OTPRequest{
			UserID: userID,
			Email:  email,
		}

		existingUser := &userPort.UserEmail{
			ID:    "different-user-id",
			Email: email,
		}
		mockRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(existingUser, nil)

		_, err := service.GenerateEmailOTP(context.Background(), otpRequest)
		require.Error(t, err)
		assert.Equal(t, shared.DefineError.OTP["EMAIL_IN_USE"].Code, err.(users.ServiceError).Code)
	})

	t.Run("store OTP failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := userPortMocks.NewMockUserRepositoryPort(ctrl)
		logger := utils.NewLogger()
		service := users.NewUserService(mockRepo, logger)

		otpRequest := users.OTPRequest{
			UserID: userID,
			Email:  email,
		}

		mockRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(nil, shared.ErrNotFound)
		mockRepo.EXPECT().StoreOTP(gomock.Any(), gomock.Any()).Return(assert.AnError)

		_, err := service.GenerateEmailOTP(context.Background(), otpRequest)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code, err.(users.ServiceError).Code)
	})
}