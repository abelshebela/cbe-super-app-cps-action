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
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

		response, err := service.DeviceLookup(ctx, deviceUUID, platform, appVersion)

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

		response, err := service.DeviceLookup(ctx, deviceUUID, platform, appVersion)

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

		response, err := prodService.DeviceLookup(ctx, deviceUUID, platform, appVersion)

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

		// Mock: No user found by device
		mockRepo.EXPECT().
			FindUserByDevice(ctx, deviceUUID).
			Return(nil, assert.AnError)

		// Mock: Get HQ data
		mockRepo.EXPECT().
			GetOneHQ(ctx, gomock.Any()).
			Return(nil, nil)

		response, err := service.DeviceLookup(ctx, deviceUUID, platform, appVersion)

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

		response, err := uatService.DeviceLookup(ctx, deviceUUID, platform, appVersion)

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
	repo := new(MockUserRepository) // You may need to define or import this mock
	logger := &TestLogger{}
	cfg := &config.VaultConfig{GoEnv: "dev"}
	service := &users.UserService{repository: repo, logger: logger, cfg: cfg}

	sessionID := "session123"
	phone := "0912345678"
	formattedPhone := utils.FormatPhoneNumber(phone)
	deviceUUID := "device-uuid-123"
	otp := "123456"
	encOtp, _, _ := utils.LocalEncryptPassword(otp, "otp", "", "", cfg)
	session := &users.PinResetSession{
		ID:          sessionID,
		PhoneNumber: formattedPhone,
		DeviceUUID:  deviceUUID,
		OTP:         encOtp,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		Status:      "pending",
		Attempts:    0,
		MaxAttempts: 3,
	}

	t.Run("valid OTP marks session as verified", func(t *testing.T) {
		repo.On("FindPinResetSession", mock.Anything, sessionID).Return(session, nil)
		repo.On("UpdatePinResetSession", mock.Anything, mock.MatchedBy(func(s *users.PinResetSession) bool {
			return s.Status == "verified"
		})).Return(nil)
		err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.NoError(t, err)
	})

	t.Run("expired session", func(t *testing.T) {
		s := *session
		s.ExpiresAt = time.Now().Add(-1 * time.Minute)
		repo.On("FindPinResetSession", mock.Anything, sessionID).Return(&s, nil)
		err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.Equal(t, users.ErrPinResetSessionExpired, err)
	})

	t.Run("invalid OTP increments attempts", func(t *testing.T) {
		s := *session
		s.Attempts = 0
		repo.On("FindPinResetSession", mock.Anything, sessionID).Return(&s, nil)
		repo.On("IncrementPinResetAttempts", mock.Anything, sessionID).Return(nil)
		err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, "000000")
		assert.Equal(t, users.ErrPinResetOTPInvalid, err)
	})

	t.Run("too many attempts", func(t *testing.T) {
		s := *session
		s.Attempts = 3
		repo.On("FindPinResetSession", mock.Anything, sessionID).Return(&s, nil)
		err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.Equal(t, users.ErrPinResetTooManyAttempts, err)
	})

	t.Run("phone/device mismatch", func(t *testing.T) {
		s := *session
		s.PhoneNumber = "other"
		repo.On("FindPinResetSession", mock.Anything, sessionID).Return(&s, nil)
		err := service.VerifyForgetPinOtp(context.Background(), sessionID, phone, deviceUUID, otp)
		assert.Equal(t, users.ErrPinResetSessionInvalid, err)
	})
}
