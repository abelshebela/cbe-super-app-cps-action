package account_test

import (
	"cbe-super-app-member-users/internal/domain/account"
	port_account "cbe-super-app-member-users/internal/port/outbound/account"
	mock_account_api "cbe-super-app-member-users/internal/port/outbound/account/mocks"
	mock_account_repo "cbe-super-app-member-users/internal/port/outbound/account/mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.uber.org/zap"
)

func TestAccountService_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_account_repo.NewMockAccountRepositoryPort(ctrl)
	mockAPIClient := mock_account_api.NewMockAccountAPIPort(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}

	service := account.NewAccountService(mockRepo, mockAPIClient, zapLogger, cfg)

	ctx := context.Background()
	userID := "test-user-id"

	t.Run("Successful account creation", func(t *testing.T) {
		user := &port_account.AccountUser{
			ID:          userID,
			KYCLevel:    1,
			PhoneNumber: "1234567890",
		}
		mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)
		mockAPIClient.EXPECT().LookupAccountByPhone(ctx, user.PhoneNumber, "").Return(false, nil)
		mockRepo.EXPECT().UpdateUserCustomerNumber(ctx, userID, gomock.Any()).Return(nil)
		mockRepo.EXPECT().CreateLinkedAccount(ctx, gomock.Any()).Return(nil)

		result, err := service.CreateAccount(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CustomerNumber)
		assert.NotEmpty(t, result.AccountNumber)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(nil, account.ErrNotFound)

		result, err := service.CreateAccount(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("User KYC level is 0", func(t *testing.T) {
		user := &port_account.AccountUser{
			ID:       userID,
			KYCLevel: 0,
		}
		mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)

		result, err := service.CreateAccount(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "USER_KYC_LEVEL_ZERO", err.Error())
	})

	t.Run("Phone number already exists", func(t *testing.T) {
		user := &port_account.AccountUser{
			ID:          userID,
			KYCLevel:    1,
			PhoneNumber: "1234567890",
		}
		mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)
		mockAPIClient.EXPECT().LookupAccountByPhone(ctx, user.PhoneNumber, "").Return(true, nil)

		result, err := service.CreateAccount(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "USER_PHONE_EXISTS", err.Error())
	})

	t.Run("Phone lookup fails", func(t *testing.T) {
		user := &port_account.AccountUser{
			ID:          userID,
			KYCLevel:    1,
			PhoneNumber: "1234567890",
		}
		mockRepo.EXPECT().FindAccountUserByID(ctx, userID).Return(user, nil)
		mockAPIClient.EXPECT().LookupAccountByPhone(ctx, user.PhoneNumber, "").Return(false, fmt.Errorf("some api error"))

		result, err := service.CreateAccount(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "API_REQUEST_FAILED", err.Error())
	})

}
