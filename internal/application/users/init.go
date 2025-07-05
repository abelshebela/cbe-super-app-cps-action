package users_application

import (
	"context"
	"mime/multipart"

	dto "cbe-super-app-member-users/internal/application/dto"
	domainUsers "cbe-super-app-member-users/internal/domain/users"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"

	// "cbe-super-app-member-users/pkgs/config"

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
	// CheckDevice(ctx context.Context, header map[string]interface{}) (map[string]interface{}, error)
	GetOneHQ(ctx context.Context, req map[string]interface{}) (*domainUsers.HQ, error)
	DeviceLookup(ctx context.Context, header map[string]interface{}) (*dto.DeviceLookupResponse, error)
	VerifyOtp(ctx context.Context, userID, otp string, deviceUUID string, userRealm, otpFor string) (*dto.VerifyOtpResponse, error)
	SetPin(ctx context.Context, userID, newPin, deviceUUID string, userRealm string) (*dto.SetPinResponse, error)
	Register(ctx context.Context, phone, deviceUUID, platform string) (*dto.RegisterResponse, error)
	CompleteRegistration(ctx context.Context, registrationID, phone, deviceUUID, platform, fullName string) (*dto.CompleteRegistrationResponse, error)
	Login(ctx context.Context, phone, deviceUUID, pin string) (*dto.LoginResponse, error)
	ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error)
	ResetPin(ctx context.Context, resetSessionID, phone, deviceUUID, otp, newPin string) (*dto.ResetPinResponse, error)
	CheckDevice(ctx context.Context, header map[string]interface{}) (map[string]interface{}, error)
	VerifyForgetPinOtp(ctx context.Context, resetSessionID, phone, deviceUUID, otp string) (*dto.VerifyOtpResponse, error)
	ResetPinWithToken(ctx context.Context, resetSessionID, phone, deviceUUID, newPin string) (*dto.ResetPinResponse, error)
}

type UsersHandler struct {
	userService *domainUsers.UserService
	logger      utils.Logger
	vaultConfig *config.VaultConfig
}

func InitUsersHandler(userService *domainUsers.UserService, logger utils.Logger, minIO config.MinioClientInterface, cfg *config.VaultConfig) ApplicationService {
	return UsersHandler{
		userService: userService,
		logger:      logger,
		vaultConfig: cfg,
	}
}
