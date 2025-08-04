package user

import (
	"context"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/user/core"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/external_call"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/token"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
)

type UsersService struct {
	smsService   external_call.SMSPersistence
	TokenService token.TokenService
	userRepo     storage.UserRepository
	otpRepo      storage.OTPRepository
	hqRepo       storage.HQRepository
	logger       utils.Logger
	Cfg          config.VaultConfig
}

func NewUserService(userRepo storage.UserRepository, smsService external_call.SMSPersistence, otpRepo storage.OTPRepository, hqRepo storage.HQRepository, tokenService token.TokenService, logger utils.Logger, config config.VaultConfig) service.UserService {
	return &UsersService{
		smsService:   smsService,
		TokenService: tokenService,
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		hqRepo:       hqRepo,
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

	if user.IsVerified {
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
		core.DeviceFoundButNotVerifiedPreparation(us.Cfg, *response)
		if err := us.NotVerifiedUser(ctx, user, response.OTPCode, response.OTPFor); err != nil {
			return nil, err
		}
	}

	return response, nil
}
func (us *UsersService) PreLogin(ctx context.Context, phone string) (*dto.DeviceLookupResponse, error) {
	var response *dto.DeviceLookupResponse
	nextStep := constants.VerifyOtp

	user, err := us.userRepo.FindByPhoneNumber(ctx, phone)
	if err != nil {
		if err.Error() != errors.ErrUnexpected.Error() {
			return core.PhoneNotFoundResponse(), nil
		}
		return nil, err
	}

	additional := core.PhoneLookupAdditionalBuilder(string(user.Platform), true, nextStep)
	token, err := us.TokenService.TempTokenMaker(user, additional, constants.DeviceLookUp)
	if err != nil {
		return nil, err
	}

	if err := core.ValidUserChecker(user, user.APPInstallationDate.GoString()); err != nil {
		return nil, err
	}

	response = core.BuildPhoneLookupResponse(user.DeviceUUID, *user, true, token, nextStep)

	core.PhoneFoundButNotVerifiedPreparation(us.Cfg, *response)
	if err := us.NotVerifiedUser(ctx, user, response.OTPCode, response.OTPFor); err != nil {
		return nil, err
	}

	return response, nil
}

// func (us *UsersService) ChangePin(ctx context.Context, req dto.ChangePinRequest) error {

// 	return nil
// }

func (us *UsersService) VerifyOtp(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOtpResponse, error) {

	return nil, nil
}

// func (us *UsersService) ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error) {
// }
// func (us *UsersService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
// }
// func (us *UsersService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
// }
// func (us *UsersService) ResetPin(ctx context.Context, req dto.ResetPinRequest) (*dto.ResetPinResponse, error) {
// }
// func (us *UsersService) SetPin(ctx context.Context, req dto.SetPinRequest) (*dto.SetPinResponse, error) {
// }
// func (us *UsersService) VerifyForgetPinOtp(ctx context.Context, req dto.VerifyForgetPinOtpRequest) (*dto.VerifyOtpResponse, error) {
// }
// func (us *UsersService) UpdateProfilePicture(ctx context.Context, userID string, req dto.UpdateProfilePicture) (string, error) {
// }
// func (us *UsersService) UpdateProfileTheme(ctx context.Context, id string, themeType string) (*dto.User, error) {
// }
func (us *UsersService) NotVerifiedUser(ctx context.Context, user *model.User, otpCode string, otpFor string) error {
	existingOtp, err := core.ExistinOTPCheck(ctx, us.otpRepo, *user)
	if err != nil {
		return err
	}

	if existingOtp {
		return errors.ErrOTPAlreadyExist
	}

	encOtpCode, _, err := us.TokenService.LocalEncryptPassword(otpCode, constants.OTP, constants.OTP, constants.OTP)
	if err != nil {
		return err
	}

	wait, err := strconv.Atoi(us.Cfg.OtpWaitingTime)
	if err != nil {
		wait = 10
	}

	expirationTime := time.Duration(wait) * time.Minute
	otpRecord := core.BuildOTPRecord(*user, encOtpCode, user.DeviceUUID, expirationTime, otpFor)

	if err := us.otpRepo.Save(ctx, &otpRecord); err != nil {
		return err
	}

	go func() {
		if err := us.smsService.SendOTP(ctx, user.PhoneNumber, otpCode, int(expirationTime)); err != nil {
			us.logger.Errorf(err.Error())
		}

	}()
	return nil
}
