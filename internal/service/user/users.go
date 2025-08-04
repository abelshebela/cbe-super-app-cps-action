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
	local_util "github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

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

func (us *UsersService) VerifyOtp(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOtpResponse, error) {

	var phone, fullName string
	nextStep := constants.SetPin

	fullName, err := core.OtpValidator(ctx, us.otpRepo, req)
	if err != nil {
		return nil, err
	}

	if req.OtpFor == constants.OTPForRegistration {
		userEntity := core.BuildUserData(fullName, req.PhoneNumber, req.DeviceUUID)

		if err := us.userRepo.Save(ctx, userEntity); err != nil {
			return nil, errors.ErrRegistrationExpired
		}
		nextStep = constants.SetPin
	}

	user, err := us.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	phone = req.PhoneNumber
	fullName = user.FullName

	userEntity := core.BuildUserData(fullName, phone, req.DeviceUUID)

	if req.Action == constants.Prelogin {
		nextStep = constants.Login
		if err := us.userRepo.Update(ctx, req.UserID, userEntity); err != nil {
			return nil, err
		}

	}

	additional := core.PhoneLookupAdditionalBuilder(string(user.Platform), true, nextStep)

	token, err := us.TokenService.TempTokenMaker(user, additional, req.Action)
	if err != nil {
		return nil, err
	}

	response := core.BuildVerifyOtpResponse(req.UserID, req.PhoneNumber, token, nextStep)

	return response, nil
}

func (us *UsersService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := us.userRepo.FindByPhoneNumber(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	if err := core.ValidUserChecker(user, user.APPInstallationDate.String()); err != nil {
		return nil, err
	}

	encryptedPIn, _, err := us.TokenService.LocalEncryptPassword(req.Pin, constants.Pin, constants.Pin, constants.Pin)
	if err != nil {
		return nil, err
	}

	if err := local_util.CheckLoginThrottle(user.LoginAttemptCount, user.LastLoginAttempt); err != nil {
		return nil, err
	}

	if err := core.PinValidator(user.LoginPIN.PIN, encryptedPIn); err != nil {
		return nil, errors.ErrInvalidPIN
	}

	if err := us.userRepo.Update(ctx, user.ID.Hex(), &model.User{LoginAttemptCount: 0, LastLoginAttempt: time.Now()}); err != nil {
		return nil, err
	}

	userEntity := core.BuildUserData(user.FullName, user.PhoneNumber, user.PhoneNumber)

	token, err := us.TokenService.TokenMaker(userEntity, &us.Cfg, constants.Permanent)
	if err != nil {
		return nil, err
	}

	sessionExpires := time.Now().Add(24 * time.Hour)

	response := core.BuildLoginResponse(userEntity, token, user.DeviceUUID, sessionExpires)

	return response, nil

}

func (us *UsersService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existing, err := us.userRepo.FindByPhoneNumber(ctx, req.Phone)
	expirationTime := 10 * time.Minute
	wait := int(expirationTime.Minutes())

	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrPhoneNumberAlreadyExists
	}

	pendingRegistration, err := us.otpRepo.Find(ctx, bson.M{"phone_number": req.Phone, "device_uuid": req.DeviceUUID, "otp_for": constants.OTPForRegistration})
	if err != nil {
		return nil, err
	}

	if pendingRegistration != nil {
		return nil, errors.ErrPinResetAlreadyInProgress
	}

	otp, encOtp, err := core.OtpProvider(us.TokenService)
	if err != nil {
		return nil, err
	}

	registration := core.BuildRegistrationRecord(req, *encOtp)

	otpRecord := core.BuildOTPFromRegistration(registration)

	if err = us.otpRepo.Save(ctx, &otpRecord); err != nil {
		return nil, err
	}

	go func() {
		if err := us.smsService.SendOTP(ctx, req.Phone, *otp, wait); err != nil {
			us.logger.Errorf("Sending sms is failed %v", err)
		}
	}()

	userEntity := core.BuildUserData(req.FullName, req.Phone, req.DeviceUUID)

	additional := core.VerifyOtpAdditionalBuilder(constants.OTPForRegistration, true, constants.VerifyOtp)

	token, err := us.TokenService.TempTokenMaker(userEntity, additional, constants.Register)
	if err != nil {
		return nil, err
	}

	response := core.BuildRegisterResponse(registration.ID, req, us.Cfg.GoEnv, *otp, wait, token)

	return response, nil
}

func (us *UsersService) SetPin(ctx context.Context, req dto.SetPinRequest) (*dto.SetPinResponse, error) {
	user, err := us.userRepo.FindById(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if core.CheckPinInHistory(user.LoginPIN.PINHistory[:], req.NewPin) {
		return nil, errors.ErrPinInHistory
	}

	encryptedPin, _, err := us.TokenService.LocalEncryptPassword(req.NewPin, constants.Pin, constants.Pin, constants.Pin)
	if err != nil {
		return nil, err
	}

	loginHistory := core.SetPinHistory(user, encryptedPin)
	if err := us.userRepo.Update(ctx, req.UserID, &model.User{LoginPIN: loginHistory}); err != nil {
		return nil, err
	}

	userEnntity := core.BuildUserData(user.FullName, user.PhoneNumber, user.DeviceUUID)

	token, err := us.TokenService.TokenMaker(userEnntity, &us.Cfg, constants.Permanent)
	if err != nil {
		return nil, err
	}

	response := core.BuildSetPinResponse(req.UserID, token)

	return response, nil
}
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

// func (us *UsersService) ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error) {
// }
// func (us *UsersService) VerifyForgetPinOtp(ctx context.Context, req dto.VerifyForgetPinOtpRequest) (*dto.VerifyOtpResponse, error) {
// }
// func (us *UsersService) ResetPin(ctx context.Context, req dto.ResetPinRequest) (*dto.ResetPinResponse, error) {
// }

// func (us *UsersService) ChangePin(ctx context.Context, req dto.ChangePinRequest) error {
// 	return nil
// }

// func (us *UsersService) UpdateProfilePicture(ctx context.Context, userID string, req dto.UpdateProfilePicture) (string, error) {
// }
// func (us *UsersService) UpdateProfileTheme(ctx context.Context, id string, themeType string) (*dto.User, error) {
// }
