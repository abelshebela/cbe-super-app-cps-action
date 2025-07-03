package users_test

import (
	"context"
	"testing"
	"time"

	"cbe-super-app-member-users/internal/domain/users"
	mock_users_repo "cbe-super-app-member-users/mocks/domain/users"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

func TestUserService_DeviceLookup_WithOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{
		GoEnv:          "dev", // Set to dev environment
		OtpWaitingTime: "10",
		Key:            "234567890-=1234567890-=1234567890-=1234567890-=",
		IV:             "1234567890-=12",
	}
	var minIOClient config.MinioClientInterface

	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)
	ctx := context.Background()

	t.Run("Device lookup with user found, not verified in dev environment - should include OTP", func(t *testing.T) {
		deviceUUID := "test-device-uuid-12345"
		platform := "android"
		appVersion := "1.0.0"
		sourceApp := "cbe-super-app"
		// Mock user found with IsVerified = false
		user := &users.User{
			ID:          bson.NewObjectID(),
			PhoneNumber: "1234567890",
			IsVerified:  false, // User is not verified
			Device: struct {
				DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
				AppVersion string `json:"app_version" bson:"app_version"`
			}{
				DeviceUUID: deviceUUID,
				AppVersion: appVersion,
			},
		}

		// Mock: User found by device
		mockRepo.EXPECT().
			FindUserByDevice(ctx, deviceUUID).
			Return(user, nil)

		// Mock: Get HQ data
		mockRepo.EXPECT().
			GetOneHQ(ctx, gomock.Any()).
			Return(nil, nil)

		// Mock: Create OTP record
		mockRepo.EXPECT().
			CreateOtp(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, otpRecord interface{}) error {
				// Verify the OTP purpose is "pin_set"
				if otp, ok := otpRecord.(users.OTPRecord); ok {
					assert.Equal(t, "pin_set", otp.OTPFor)
				}
				return nil
			})

		response, err := service.DeviceLookup(ctx, deviceUUID, platform, appVersion, sourceApp)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, deviceUUID, response.DeviceUUID)
		assert.Equal(t, platform, response.Platform)
		assert.Equal(t, appVersion, response.AppVersion)
		assert.True(t, response.UserFound)
		assert.Equal(t, "login", response.NextStep)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "device_lookup", response.TokenType)
		assert.NotEmpty(t, response.OTPCode) // OTP should be included when user is not verified in dev environment
		assert.Len(t, response.OTPCode, 6)   // OTP should be 6 digits
	})

	t.Run("Device lookup with user found, verified in dev environment - should NOT include OTP", func(t *testing.T) {
		deviceUUID := "test-device-uuid-12345"
		platform := "android"
		appVersion := "1.0.0"
		sourceApp := "cbe-super-app"
		// Mock user found with IsVerified = true
		user := &users.User{
			ID:          bson.NewObjectID(),
			PhoneNumber: "1234567890",
			IsVerified:  true, // User is verified
			Device: struct {
				DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
				AppVersion string `json:"app_version" bson:"app_version"`
			}{
				DeviceUUID: deviceUUID,
				AppVersion: appVersion,
			},
		}

		// Mock: User found by device
		mockRepo.EXPECT().
			FindUserByDevice(ctx, deviceUUID).
			Return(user, nil)

		// Mock: Get HQ data
		mockRepo.EXPECT().
			GetOneHQ(ctx, gomock.Any()).
			Return(nil, nil)

		response, err := service.DeviceLookup(ctx, deviceUUID, platform, appVersion, sourceApp)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, deviceUUID, response.DeviceUUID)
		assert.Equal(t, platform, response.Platform)
		assert.Equal(t, appVersion, response.AppVersion)
		assert.True(t, response.UserFound)
		assert.Equal(t, "login", response.NextStep)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "device_lookup", response.TokenType)
		assert.Empty(t, response.OTPCode) // OTP should NOT be included when user is verified, even in dev environment
	})

	t.Run("Device lookup with user found, not verified in production environment - should NOT include OTP", func(t *testing.T) {
		// Change config to production environment
		prodCfg := &config.VaultConfig{
			GoEnv:          "production", // Set to production environment
			OtpWaitingTime: "10",
			Key:            "234567890-=1234567890-=1234567890-=1234567890-=",
			IV:             "1234567890-=12",
		}
		prodService := users.NewUserService(mockRepo, zapLogger, minIOClient, prodCfg)

		deviceUUID := "test-device-uuid-12345"
		platform := "android"
		appVersion := "1.0.0"
		sourceApp := "cbe-super-app"
		// Mock user found with IsVerified = false
		user := &users.User{
			ID:          bson.NewObjectID(),
			PhoneNumber: "1234567890",
			IsVerified:  false, // User is not verified
			Device: struct {
				DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
				AppVersion string `json:"app_version" bson:"app_version"`
			}{
				DeviceUUID: deviceUUID,
				AppVersion: appVersion,
			},
		}

		// Mock: User found by device
		mockRepo.EXPECT().
			FindUserByDevice(ctx, deviceUUID).
			Return(user, nil)

		// Mock: Get HQ data
		mockRepo.EXPECT().
			GetOneHQ(ctx, gomock.Any()).
			Return(nil, nil)

		response, err := prodService.DeviceLookup(ctx, deviceUUID, platform, appVersion, sourceApp)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, deviceUUID, response.DeviceUUID)
		assert.Equal(t, platform, response.Platform)
		assert.Equal(t, appVersion, response.AppVersion)
		assert.True(t, response.UserFound)
		assert.Equal(t, "login", response.NextStep)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "device_lookup", response.TokenType)
		assert.Empty(t, response.OTPCode) // OTP should NOT be included in production environment, even if user is not verified
	})

	t.Run("Device lookup with no user found - should not include OTP", func(t *testing.T) {
		deviceUUID := "test-device-uuid-12345"
		platform := "android"
		appVersion := "1.0.0"
		sourceApp := "cbe-super-app"
		// Mock: No user found by device
		mockRepo.EXPECT().
			FindUserByDevice(ctx, deviceUUID).
			Return(nil, assert.AnError)

		// Mock: Get HQ data
		mockRepo.EXPECT().
			GetOneHQ(ctx, gomock.Any()).
			Return(nil, nil)

		response, err := service.DeviceLookup(ctx, deviceUUID, platform, appVersion, sourceApp)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, deviceUUID, response.DeviceUUID)
		assert.Equal(t, platform, response.Platform)
		assert.Equal(t, appVersion, response.AppVersion)
		assert.False(t, response.UserFound)
		assert.Equal(t, "register", response.NextStep)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "device_lookup", response.TokenType)
		assert.Empty(t, response.OTPCode) // OTP should not be included when no user found
	})

	t.Run("Device lookup with user found, not verified in uat environment - should include OTP", func(t *testing.T) {
		// Change config to uat environment
		uatCfg := &config.VaultConfig{
			GoEnv:          "uat", // Set to uat environment
			OtpWaitingTime: "10",
			Key:            "234567890-=1234567890-=1234567890-=1234567890-=",
			IV:             "1234567890-=12",
		}
		uatService := users.NewUserService(mockRepo, zapLogger, minIOClient, uatCfg)

		deviceUUID := "test-device-uuid-12345"
		platform := "android"
		appVersion := "1.0.0"
		sourceApp := "cbe-super-app"
		// Mock user found with IsVerified = false
		user := &users.User{
			ID:          bson.NewObjectID(),
			PhoneNumber: "1234567890",
			IsVerified:  false, // User is not verified
			Device: struct {
				DeviceUUID string `json:"device_uuid" bson:"device_uuid"`
				AppVersion string `json:"app_version" bson:"app_version"`
			}{
				DeviceUUID: deviceUUID,
				AppVersion: appVersion,
			},
		}

		// Mock: User found by device
		mockRepo.EXPECT().
			FindUserByDevice(ctx, deviceUUID).
			Return(user, nil)

		// Mock: Get HQ data
		mockRepo.EXPECT().
			GetOneHQ(ctx, gomock.Any()).
			Return(nil, nil)

		// Mock: Create OTP record
		mockRepo.EXPECT().
			CreateOtp(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, otpRecord interface{}) error {
				// Verify the OTP purpose is "pin_set"
				if otp, ok := otpRecord.(users.OTPRecord); ok {
					assert.Equal(t, "pin_set", otp.OTPFor)
				}
				return nil
			})

		response, err := uatService.DeviceLookup(ctx, deviceUUID, platform, appVersion, sourceApp)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, deviceUUID, response.DeviceUUID)
		assert.Equal(t, platform, response.Platform)
		assert.Equal(t, appVersion, response.AppVersion)
		assert.True(t, response.UserFound)
		assert.Equal(t, "login", response.NextStep)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "device_lookup", response.TokenType)
		assert.NotEmpty(t, response.OTPCode) // OTP should be included when user is not verified in uat environment
		assert.Len(t, response.OTPCode, 6)   // OTP should be 6 digits
	})
}

func TestVerifyForgetPinOtp(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_users_repo.NewMockUserRepository(ctrl)
	logger, _ := zap.NewDevelopment()
	zapLogger := logger.Sugar()
	cfg := &config.VaultConfig{
		GoEnv:          "dev", // Set to dev environment
		OtpWaitingTime: "10",
		Key:            "234567890-=1234567890-=1234567890-=1234567890-=",
		IV:             "1234567890-=12",
	}
	var minIOClient config.MinioClientInterface
	service := users.NewUserService(mockRepo, zapLogger, minIOClient, cfg)

	sessionID := "session123"
	phone := "0912345678"
	deviceUUID := "device-uuid-123"
	otp := "123456"
	session := &users.PinResetSession{
		ID:          sessionID,
		PhoneNumber: phone,
		DeviceUUID:  deviceUUID,
		OTP:         otp,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		Status:      "pending",
		Attempts:    0,
		MaxAttempts: 3,
	}

	t.Run("valid OTP marks session as verified", func(t *testing.T) {
		mockRepo.EXPECT().
			FindPinResetSession(gomock.Any(), sessionID).
			Return(session, nil)
		mockRepo.EXPECT().
			UpdatePinResetSession(gomock.Any(), mock.MatchedBy(func(s *users.PinResetSession) bool {
				return s.Status == "verified"
			})).Return(nil)
		_, err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.NoError(t, err)
	})

	t.Run("expired session", func(t *testing.T) {
		s := *session
		s.ExpiresAt = time.Now().Add(-1 * time.Minute)
		mockRepo.EXPECT().
			FindPinResetSession(gomock.Any(), sessionID).
			Return(&s, nil)
		_, err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.Equal(t, users.ErrPinResetSessionExpired, err)
	})

	t.Run("invalid OTP increments attempts", func(t *testing.T) {
		s := *session
		s.Attempts = 0
		mockRepo.EXPECT().
			FindPinResetSession(gomock.Any(), sessionID).
			Return(&s, nil)
		mockRepo.EXPECT().
			IncrementPinResetAttempts(gomock.Any(), sessionID).
			Return(nil)
		_, err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, "000000")
		assert.Equal(t, users.ErrPinResetOTPInvalid, err)
	})

	t.Run("too many attempts", func(t *testing.T) {
		s := *session
		s.Attempts = 3
		mockRepo.EXPECT().
			FindPinResetSession(gomock.Any(), sessionID).
			Return(&s, nil)
		_, err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.Equal(t, users.ErrPinResetTooManyAttempts, err)
	})

	t.Run("phone/device mismatch", func(t *testing.T) {
		s := *session
		s.PhoneNumber = "other"
		mockRepo.EXPECT().
			FindPinResetSession(gomock.Any(), sessionID).
			Return(&s, nil)
		_, err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.Equal(t, users.ErrPinResetSessionInvalid, err)
	})
}
