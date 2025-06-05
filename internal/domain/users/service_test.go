package users_test

import (
	"context"
	"testing"
	"time"

	"cbe-super-app-member-users/internal/domain/users"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"mocks "cbe-super-app-member-users/internal/domain/users/mocks"
)

func TestActiveLinkedAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger := utils.NewLogger()
	service := users.NewUserService(mockRepo, logger)

	user := MockUser()
	userID := user.ID
	linkedAccounts := MockLinkedAccounts()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(gomock.Any(), userID).Return(linkedAccounts, nil)

		response, err := service.ActiveLinkedAccounts(context.Background(), userID)
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
		_, err := service.ActiveLinkedAccounts(context.Background(), "invalid-id")
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["INVALID_ID"].Code, err.(*users.ServiceError).Code)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(nil, mongo.ErrNoDocuments)

		_, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["NOT_FOUND"].Code, err.(*users.ServiceError).Code)
	})

	t.Run("no linked accounts", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(gomock.Any(), userID).Return([]users.LinkedAccountDetail{}, nil)

		response, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.LinkedAccounts, 0)
	})

	t.Run("deleted user", func(t *testing.T) {
		deletedUser := &users.User{
			ID:        userID,
			FullName:  user.FullName,
			IsDeleted: true,
		}
		mockRepo.EXPECT().FindByID(gomock.Any(), userID).Return(deletedUser, nil)

		response, err := service.ActiveLinkedAccounts(context.Background(), userID)
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.LinkedAccounts, 0)
	})
}

func TestGenerateEmailOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger := utils.NewLogger()
	service := users.NewUserService(mockRepo, logger)

	user := MockUser()
	userID := user.ID
	email := "test@example.com"
	otpRequest := users.OTPRequest{
		UserID: userID,
		Email:  email,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(nil, mongo.ErrNoDocuments)
		mockRepo.EXPECT().StoreOTP(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, otp *users.OTPRecord) error {
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

	t.Run("invalid user ID", func(t *testing.T) {
		invalidRequest := users.OTPRequest{
			UserID: "invalid-id",
			Email:  email,
		}
		_, err := service.GenerateEmailOTP(context.Background(), invalidRequest)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["INVALID_ID"].Code, err.(*users.ServiceError).Code)
	})

	t.Run("email already in use", func(t *testing.T) {
		existingUser := &users.UserEmail{
			ID:    primitive.NewObjectID().Hex(),
			Email: email,
		}
		mockRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(existingUser, nil)

		_, err := service.GenerateEmailOTP(context.Background(), otpRequest)
		require.Error(t, err)
		assert.Equal(t, "EMAIL_IN_USE", err.(*users.ServiceError).Code)
	})

	t.Run("store OTP failure", func(t *testing.T) {
		mockRepo.EXPECT().FindByEmail(gomock.Any(), email).Return(nil, mongo.ErrNoDocuments)
		mockRepo.EXPECT().StoreOTP(gomock.Any(), gomock.Any()).Return(mongo.ErrClientDisconnected)

		_, err := service.GenerateEmailOTP(context.Background(), otpRequest)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code, err.(*users.ServiceError).Code)
	})
}

func TestVerifyEmailOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	logger := utils.NewLogger()
	service := users.NewUserService(mockRepo, logger)

	user := MockUser()
	userID := user.ID
	email := "test@example.com"
	otp := "123456"
	verification := users.OTPVerification{
		UserID: userID,
		Email:  email,
		OTP:    otp,
	}
	otpRecord := MockOTPRecord(userID, email, otp)

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(otpRecord, nil)
		mockRepo.EXPECT().UpdateUserEmail(gomock.Any(), userID, email).Return(nil)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.NoError(t, err)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		invalidVerification := users.OTPVerification{
			UserID: "invalid-id",
			Email:  email,
			OTP:    otp,
		}
		err := service.VerifyEmailOTP(context.Background(), invalidVerification)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["INVALID_ID"].Code, err.(*users.ServiceError).Code)
	})

	t.Run("invalid OTP", func(t *testing.T) {
		invalidOTPRecord := MockOTPRecord(userID, email, "654321")
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(invalidOTPRecord, nil)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, "INVALID_OTP", err.(*users.ServiceError).Code)
	})

	t.Run("expired OTP", func(t *testing.T) {
		expiredOTPRecord := MockOTPRecord(userID, email, otp)
		expiredOTPRecord.CreatedAt = time.Now().Add(-10 * time.Minute)
		expiredOTPRecord.ExpiresAt = time.Now().Add(-5 * time.Minute)
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(expiredOTPRecord, nil)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, "EXPIRED_OTP", err.(*users.ServiceError).Code)
	})

	t.Run("OTP not found", func(t *testing.T) {
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(nil, mongo.ErrNoDocuments)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, "INVALID_OTP", err.(*users.ServiceError).Code)
	})

	t.Run("update email failure", func(t *testing.T) {
		mockRepo.EXPECT().FindOTP(gomock.Any(), userID, email).Return(otpRecord, nil)
		mockRepo.EXPECT().UpdateUserEmail(gomock.Any(), userID, email).Return(mongo.ErrClientDisconnected)

		err := service.VerifyEmailOTP(context.Background(), verification)
		require.Error(t, err)
		assert.Equal(t, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code, err.(*users.ServiceError).Code)
	})
}