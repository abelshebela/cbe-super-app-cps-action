package users_application

import (
	"context"
	"mime/multipart"

	domainUsers "cbe-super-app-member-users/internal/domain/users"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	FetchLinkedAccounts(ctx context.Context, userID string) (*userPort.LinkedAccountResponse, error)
	GenerateEmailOTP(ctx context.Context, req userPort.OTPRequest) (string, error)
	VerifyEmailOTP(ctx context.Context, req userPort.OTPVerification) error
	UpdateProfilePicture(ctx context.Context, id string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UnlinkDevice(ctx context.Context, userID string, deviceID string) error
	ChangePin(ctx context.Context, req userPort.ChangePinRequest) error
	CheckDevice(ctx context.Context, header map[string]interface{}) (map[string]interface{}, error)
	GetOneHQ(ctx context.Context, req map[string]interface{}) (*domainUsers.HQ, error)
	VerifyOtp(ctx context.Context, userID, otp string, deviceUUID *string, userRealm, otpFor string) error
	SetPin(ctx context.Context, userID, newPin, otp string, deviceUUID *string, userRealm, otpFor string) error
	Register(ctx context.Context, phone, deviceUUID, platform string) error
	Login(ctx context.Context, phone, deviceUUID, pin string) (string, error)
	ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) error
	DeviceLookup(ctx context.Context, deviceUUID, platform, appVersion string) (map[string]interface{}, error)
}

type UsersHandler struct {
	userService *domainUsers.UserService
	logger      utils.Logger
}

func InitUsersHandler(userService *domainUsers.UserService, logger utils.Logger, minIO config.MinioClientInterface) ApplicationService {
	return UsersHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h UsersHandler) SetPin(ctx context.Context, userID, newPin, otp string, deviceUUID *string, userRealm, otpFor string) error {
	return h.userService.SetPin(ctx, userID, newPin, otp, deviceUUID, userRealm, otpFor)
}

func (h UsersHandler) DeviceLookup(ctx context.Context, deviceUUID, platform, appVersion string) (map[string]interface{}, error) {
	return h.userService.DeviceLookup(ctx, deviceUUID, platform, appVersion)
}
