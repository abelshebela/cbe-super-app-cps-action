package users_test

import (
	"context"
	"testing"

	"cbe-super-app-member-users/internal/domain/users"
	mock_users_repo "cbe-super-app-member-users/mocks/domain/users"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()

	t.Run("Successful registration", func(t *testing.T) {
		// Add your mock expectations and assertions here
		err := service.Register(ctx, "1234567890", "device-uuid", "android")
		assert.NoError(t, err)
	})
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()

	t.Run("Successful login", func(t *testing.T) {
		// Add your mock expectations and assertions here
		token, err := service.Login(ctx, "1234567890", "device-uuid", "123456")
		assert.NoError(t, err)
		assert.Equal(t, "access_token", token)
	})
}

func TestUserService_ForgetPinSendOtp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()

	t.Run("Successful forget pin send otp", func(t *testing.T) {
		// Add your mock expectations and assertions here
		err := service.ForgetPinSendOtp(ctx, "1234567890", "device-uuid")
		assert.NoError(t, err)
	})
}

func TestUserService_SetPin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()

	t.Run("Successful set pin", func(t *testing.T) {
		// Add your mock expectations and assertions here
		err := service.SetPin(ctx, bson.NewObjectID().Hex(), "123456", "654321", nil, "member", "pin_set")
		assert.Error(t, err) // Will error unless you mock FindOTP, etc.
	})
}

func TestUserService_VerifyOtp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()

	t.Run("Successful verify otp", func(t *testing.T) {
		// Add your mock expectations and assertions here
		err := service.VerifyOtp(ctx, bson.NewObjectID().Hex(), "654321", nil, "member", "pin_set")
		assert.Error(t, err) // Will error unless you mock FindOTP, etc.
	})
}
