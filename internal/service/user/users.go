package user

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
)

type UsersService struct {
	userRepo     storage.UserRepository
	otpRepo      storage.OTPRepository
	hqRepo       storage.HQRepository
	redisFactory storage.RedisRepository
	logger       utils.Logger
	cfg          config.VaultConfig
}

func NewUserService(userRepo storage.UserRepository, otpRepo storage.OTPRepository, hqRepo storage.HQRepository, redisFactory storage.RedisRepository, logger utils.Logger, config config.VaultConfig) service.UserService {
	return &UsersService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		hqRepo:       hqRepo,
		redisFactory: redisFactory,
		logger:       logger,
		cfg:          config,
	}
}

func (us *UsersService) DeviceLookup(ctx context.Context, req dto.DeviceLookupResponse) (*dto.DeviceLookupResponse, error) {

}
func (us *UsersService) PreLogin(ctx context.Context, header dto.Address, phone string) (*dto.DeviceLookupResponse, error) {
}
func (us *UsersService) ChangePin(ctx context.Context, req dto.ChangePinRequest) error {}
func (us *UsersService) VerifyOtp(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOtpResponse, error) {
}
func (us *UsersService) ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error) {
}
func (us *UsersService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
}
func (us *UsersService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
}
func (us *UsersService) ResetPin(ctx context.Context, req dto.ResetPinRequest) (*dto.ResetPinResponse, error) {
}
func (us *UsersService) SetPin(ctx context.Context, req dto.SetPinRequest) (*dto.SetPinResponse, error) {
}
func (us *UsersService) VerifyForgetPinOtp(ctx context.Context, req dto.VerifyForgetPinOtpRequest) (*dto.VerifyOtpResponse, error) {
}
func (us *UsersService) UpdateProfilePicture(ctx context.Context, userID string, req dto.UpdateProfilePicture) (string, error) {
}
func (us *UsersService) UpdateProfileTheme(ctx context.Context, id string, themeType string) (*dto.User, error) {
}
