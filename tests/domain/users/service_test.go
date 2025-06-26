package users_test

import (
	"cbe-super-app-member-users/internal/domain/users"
	mock_users_repo "cbe-super-app-member-users/mocks/domain/users"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.uber.org/zap"
)

func TestUserService_FetchLinkedAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}

	// We are not testing MinIO functionality here, so we can pass a nil client.
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()
	userID := "test-user-id"

	t.Run("Successful fetch of linked accounts", func(t *testing.T) {
		user := &users.User{
			ID:       userID,
			FullName: "Test User",
		}
		linkedAccounts := []users.LinkedAccountDetail{
			{AccountNumber: "123", AccountBranchCode: "001"},
			{AccountNumber: "456", AccountBranchCode: "002"},
		}

		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(ctx, userID).Return(linkedAccounts, nil)

		result, err := service.FetchLinkedAccounts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, "Test User", result.FullName)
		assert.Len(t, result.LinkedAccounts, 2)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, userID).Return(nil, users.ErrNotFound)

		result, err := service.FetchLinkedAccounts(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("User is deleted", func(t *testing.T) {
		user := &users.User{
			ID:        userID,
			FullName:  "Test User",
			IsDeleted: true,
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		result, err := service.FetchLinkedAccounts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.LinkedAccounts, 0)
	})

	t.Run("No active linked accounts found", func(t *testing.T) {
		user := &users.User{
			ID:       userID,
			FullName: "Test User",
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().FindActiveLinkedAccounts(ctx, userID).Return([]users.LinkedAccountDetail{}, nil)

		result, err := service.FetchLinkedAccounts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.LinkedAccounts, 0)
	})
}

func TestUserService_UnlinkDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()
	userID := "test-user-id"
	deviceID := "test-device-id"

	t.Run("Successful unlink", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			Device: &users.DeviceInfo{
				DeviceUUID: deviceID,
				AppVersion: "1.0.0",
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().UnlinkDevice(ctx, userID, deviceID).Return(nil)

		err := service.UnlinkDevice(ctx, userID, deviceID)
		assert.NoError(t, err)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, userID).Return(nil, users.ErrNotFound)

		err := service.UnlinkDevice(ctx, userID, deviceID)
		assert.Error(t, err)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("No linked device", func(t *testing.T) {
		user := &users.User{
			ID:     userID,
			Device: nil,
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		err := service.UnlinkDevice(ctx, userID, deviceID)
		assert.Error(t, err)
		assert.Equal(t, "NO_LINKED_DEVICES", err.Error())
	})

	t.Run("Device ID does not match", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			Device: &users.DeviceInfo{
				DeviceUUID: "other-device-id",
				AppVersion: "1.0.0",
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		err := service.UnlinkDevice(ctx, userID, deviceID)
		assert.Error(t, err)
		assert.Equal(t, "NO_LINKED_DEVICES", err.Error())
	})

	t.Run("Repository unlink error", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			Device: &users.DeviceInfo{
				DeviceUUID: deviceID,
				AppVersion: "1.0.0",
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().UnlinkDevice(ctx, userID, deviceID).Return(fmt.Errorf("db error"))

		err := service.UnlinkDevice(ctx, userID, deviceID)
		assert.Error(t, err)
		assert.Equal(t, "COULD_NOT_UNLINK_DEVICE", err.Error())
	})
}

func TestUserService_ChangePin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()
	userID := "test-user-id"
	oldPin := "135792"
	newPin := "246801"

	t.Run("Successful pin change", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			LoginPIN: users.LoginPIN{
				PIN:        oldPin,
				PINHistory: [4]string{"111111", "222222", "333333", "444444"},
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().ChangePin(ctx, userID, gomock.Any()).Return(nil)

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: newPin,
		})
		assert.NoError(t, err)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(ctx, userID).Return(nil, users.ErrNotFound)

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: newPin,
		})
		assert.Error(t, err)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("Old PIN mismatch", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			LoginPIN: users.LoginPIN{
				PIN:        "999999", // Different from oldPin
				PINHistory: [4]string{"111111", "222222", "333333", "444444"},
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: newPin,
		})
		assert.Error(t, err)
		assert.Equal(t, "OLD_PIN_MISMATCH", err.Error())
	})

	t.Run("Same PIN", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			LoginPIN: users.LoginPIN{
				PIN:        oldPin,
				PINHistory: [4]string{"111111", "222222", "333333", "444444"},
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: oldPin, // Same as old PIN
		})
		assert.Error(t, err)
		assert.Equal(t, "SAME_PIN", err.Error())
	})

	t.Run("PIN in history", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			LoginPIN: users.LoginPIN{
				PIN:        oldPin,
				PINHistory: [4]string{"246801", "222222", "333333", "444444"}, // newPin is "246801"
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		// Do NOT expect ChangePin to be called here

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: newPin, // "246801"
		})
		assert.Error(t, err)
		assert.Equal(t, "PIN_IN_HISTORY", err.Error())
	})

	t.Run("Invalid PIN format", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			LoginPIN: users.LoginPIN{
				PIN:        oldPin,
				PINHistory: [4]string{"111111", "222222", "333333", "444444"},
			},
		}
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: "12345", // Too short
		})
		assert.Error(t, err)
		assert.Equal(t, "PIN_LIMIT", err.Error())
	})

	t.Run("Repository error", func(t *testing.T) {
		user := &users.User{
			ID: userID,
			LoginPIN: users.LoginPIN{
				PIN:        oldPin,
				PINHistory: [4]string{"111111", "222222", "333333", "444444"},
			},
		}
		repoErr := fmt.Errorf("db error")
		mockRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().ChangePin(ctx, userID, gomock.Any()).Return(repoErr)

		err := service.ChangePin(ctx, users.ChangePinRequest{
			UserID: userID,
			OldPin: oldPin,
			NewPin: newPin,
		})
		assert.Error(t, err)
		assert.Equal(t, "ERROR_CHANGING_PIN", err.Error())
	})
}
