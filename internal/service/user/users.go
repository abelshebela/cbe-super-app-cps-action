package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/user/core"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/token"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
)

type UsersService struct {
	TokenService token.TokenService
	userRepo     storage.UserRepository
	otpRepo      storage.OTPRepository
	hqRepo       storage.HQRepository
	redisFactory storage.RedisRepository
	logger       utils.Logger
	Cfg          config.VaultConfig
}

func NewUserService(userRepo storage.UserRepository, otpRepo storage.OTPRepository, hqRepo storage.HQRepository, redisFactory storage.RedisRepository, tokenService token.TokenService, logger utils.Logger, config config.VaultConfig) service.UserService {
	return &UsersService{
		TokenService: tokenService,
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		hqRepo:       hqRepo,
		redisFactory: redisFactory,
		logger:       logger,
		Cfg:          config,
	}
}

func (us *UsersService) DeviceLookup(ctx context.Context, req dto.DeviceLookupRequest) (*dto.DeviceLookupResponse, error) {
	hqData, err := us.hqRepo.FindOne(ctx, nil)
	var response *dto.DeviceLookupResponse
	nextStep := constants.VerifyOtp

	if err != nil {
		return nil, err
	}
	if !core.IsAppVersionLatest(string(req.Platform), req.AppVersion, hqData) {

		return core.DeviceLookupResponseOldDevice(), nil
	}

	user, err := us.userRepo.FindByDeviceUUID(ctx, req.DeviceUUID)
	if err != nil {
		if err.Error() != errors.ErrUnexpected.Error() {
			return core.DeviceNotFoundResponse(), nil
		}
		return nil, err
	}

	if !user.IsVerified {
		nextStep = constants.Login
	}
	additional := core.DeviceLookupAdditionalBuilder(string(req.Platform), true, nextStep)
	token, err := us.TokenService.TempTokenMaker(user, additional, constants.DeviceLookUp)
	if err != nil {
		return nil, err
	}

	if err := core.ValidUserChecker(user, req.ApplicationInstallationDate.String()); err != nil {
		return nil, err
	}

	response = core.BuildDeviceLookupResponse(req.DeviceUUID, *user, true, token, nextStep)
	if !user.IsVerified {
		existingOtp, err := core.ExistinOTPCheck(ctx, us.otpRepo, *user)
		if err != nil {
			return nil, err
		}

		if existingOtp {
			return nil, errors.ErrOTPAlreadyExist
		}

		core.DeviceFoundButNotVerifiedPreparation(us.Cfg, *response)

		encOtpCode, _, err := us.TokenService.LocalEncryptPassword(response.OTPCode, constants.OTP, constants.OTP, constants.OTP)
		if err != nil {
			return nil, err
		}

		wait, err := strconv.Atoi(us.Cfg.OtpWaitingTime)
		if err != nil {
			wait = 10
		}

		expirationTime := time.Duration(wait) * time.Minute
		if err := s.otpCreator(ctx, response.UserID, encOtpCode, deviceUUID, expirationTime); err != nil {
			return nil, err
		}

	}

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
func (s *UserService) otpCreator(ctx context.Context, userId, encOtpCode, deviceUUID string, expirationTime time.Duration) error {
	now := time.Now().UTC()
	otpData, err := s.repository.GetOtpByUserID(ctx, userId)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			s.logger.Errorf("Failed to get otp data: %v", err)
			return err
		}
	} else if now.After(otpData.ExpiresAt) {
		s.repository.DeleteOtpHard(ctx, otpData.ID)
	} else {
		return fmt.Errorf("Time remaining:", otpData.ExpiresAt.Sub(now).Round(time.Second))
	}

	otpRecord := OTPRecord{
		UserCode:   userId,
		OTP:        encOtpCode,
		OTPFor:     string(enums.OTP_FOR_PIN_SET), // Using enum for "pin_set"
		UserRealm:  "member",                      // Using sourceApp as user_realm
		ExpiresAt:  time.Now().Add(expirationTime),
		CreatedAt:  time.Now(),
		DeviceUUID: deviceUUID,
	}

	// Store OTP in database
	if err := s.CreateOtp(ctx, otpRecord); err != nil {
		return errors.ErrFailedOtpCreation
	}
	return nil
}
