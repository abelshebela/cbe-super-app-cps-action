package core

import (
	"context"
	"crypto/subtle"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/types"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/token"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	local_util "github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsAppVersionLatest(platform, appVersion string, hqData *model.HQ) bool {
	deviceVer, err := version.NewVersion(appVersion)
	if err != nil {
		return false
	}

	var latestVersion *version.Version

	switch platform {
	case string(constants.Android):
		latestVersion, err = version.NewVersion(hqData.LatestAndroidVersion)
	case string(constants.Ios):
		latestVersion, err = version.NewVersion(hqData.LatestiOSVersion)
	default:
		return false
	}

	if err != nil {
		return false
	}

	return deviceVer.GreaterThanOrEqual(latestVersion)
}

func DeviceLookupResponseOldDevice() *dto.DeviceLookupResponse {
	return &dto.DeviceLookupResponse{
		NextStep: constants.UpdateApp,
		IsLatest: false,
		Message:  errors.ErrOldDevice.Error(),
		Status:   400,
	}
}

func DeviceNotFoundResponse() *dto.DeviceLookupResponse {
	return &dto.DeviceLookupResponse{
		UserFound: false,
		NextStep:  constants.Prelogin,
		Message:   errors.ErrDeviceNotFound.Error(),
		Status:    400,
	}
}

func ValidUserChecker(userData *model.User, installationData string) error {
	deviceAppDate, err := time.Parse(time.RFC3339, installationData)
	if err != nil {
		deviceAppDate, err = time.Parse("2006-01-02", installationData)
		if err != nil {
			return err
		}
	}
	// deviceAppDate, err := time.Parse(time.RFC3339, installationData)
	duration := userData.APPInstallationDate.Sub(deviceAppDate)
	dParsed, err := time.ParseDuration(duration.String())
	if err != nil {
		return err
	}

	if dParsed != 0 {
		return errors.ErrDeviceDiffInstallationDate
	}

	if userData.IsAccountBlocked {
		return errors.ErrAccBlocked
	}

	if userData.BlockedOnCPS {
		return errors.ErrUserBlockedByCps
	}

	if userData.IsAccountBlocked {
		return errors.ErrUserAccountBlockedByCps
	}

	return nil
}

func BuildDeviceLookupResponse(deviceUUID string, user model.User, isLatest bool, token, nextStep string) *dto.DeviceLookupResponse {
	return &dto.DeviceLookupResponse{
		DeviceUUID:  deviceUUID,
		UserID:      user.ID.Hex(),
		UserCode:    user.UserCode,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsLatest:    isLatest,
		Token:       token,
		TokenExpiry: time.Now().Add(5 * time.Minute),
		NextStep:    nextStep,
		Message:     constants.DeviceFound,
		Status:      200,
	}
}

func DeviceLookupAdditionalBuilder(platform string, userFound bool, nextStep string) map[string]interface{} {
	return map[string]interface{}{
		"platform":   platform,
		"token_type": "device_lookup",
		"pin_login":  userFound,
		"register":   !userFound,
		"next_step":  nextStep,
	}
}

func DeviceFoundButNotVerifiedPreparation(cfg config.VaultConfig, response dto.DeviceLookupResponse) {

	otpCode := local_util.OTPGenerator(constants.OTPLength)
	if cfg.GoEnv == constants.DEV || cfg.GoEnv == constants.UAT {
		response.OTPCode = otpCode
		response.OTPFor = string(constants.OTPForEnable)
	}

}

func ExistinOTPCheck(ctx context.Context, otpRepo storage.OTPRepository, user model.User) (bool, error) {
	filter := bson.M{
		"phone_number": user.PhoneNumber,
	}
	data, err := otpRepo.Find(ctx, filter)
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil
	}
	return true, nil
}

func BuildOTPRecord(user model.User, encOtpCode, deviceUUID string, expirationTime time.Duration, otpFor string) model.OTP {
	return model.OTP{
		UserCode:    user.ID.Hex(),
		OTPCode:     encOtpCode,
		PhoneNumber: user.PhoneNumber,
		OTPFor:      constants.OTPFor(otpFor),
		UserRealm:   constants.MEMBER_REALM,
		ExpiresAt:   time.Now().Add(expirationTime),
		CreatedAt:   time.Now(),
		DeviceUUID:  &deviceUUID,
	}
}

func PhoneNotFoundResponse() *dto.DeviceLookupResponse {
	return &dto.DeviceLookupResponse{
		DeviceUUID: constants.Empty,
		NextStep:   constants.Register,
		IsLatest:   true,
		Message:    errors.ErrPhoneNotFound.Error(),
		Status:     404,
	}
}

func BuildPhoneLookupResponse(deviceUUID string, user model.User, isLatest bool, token, nextStep string) *dto.DeviceLookupResponse {
	return &dto.DeviceLookupResponse{
		DeviceUUID:  deviceUUID,
		UserID:      user.ID.Hex(),
		UserCode:    user.UserCode,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsLatest:    isLatest,
		Token:       token,
		TokenExpiry: time.Now().Add(5 * time.Minute),
		NextStep:    nextStep,
		Message:     constants.PhoneFound,
		Status:      200,
	}
}

func PhoneFoundButNotVerifiedPreparation(cfg config.VaultConfig, response dto.DeviceLookupResponse) {

	otpCode := local_util.OTPGenerator(constants.OTPLength)
	if cfg.GoEnv == constants.DEV || cfg.GoEnv == constants.UAT {
		response.OTPCode = otpCode
		response.OTPFor = string(constants.OTPForPINSet)
	}

}

func PhoneLookupAdditionalBuilder(platform string, userFound bool, nextStep string) map[string]interface{} {
	return map[string]interface{}{
		"platform":   platform,
		"token_type": "pre_login",
		"next_step":  nextStep,
	}
}

func VerifyOtpAdditionalBuilder(otpFor string, userFound bool, nextStep string) map[string]interface{} {
	return map[string]interface{}{
		"otp_for":    otpFor,
		"token_type": "verify_otp",
		"next_step":  nextStep,
	}
}

func OtpValidator(ctx context.Context, otpRepo storage.OTPRepository, req dto.VerifyOTPRequest) (string, error) {

	if req.OtpFor == constants.OTPForRegistration && req.DeviceUUID != "" {
		filter := bson.M{
			"devide_uuid":  req.DeviceUUID,
			"phone_number": req.PhoneNumber,
		}
		regisration, err := otpRepo.Find(ctx, filter)
		if err != nil {
			{
				return constants.Empty, err
			}
		}

		if time.Now().After(regisration.ExpiresAt) {
			if err := otpRepo.Delete(ctx, regisration.ID.Hex()); err != nil {
				return constants.Empty, err
			}

			return constants.Empty, errors.ErrOtpExpired
		}

		if subtle.ConstantTimeCompare([]byte(regisration.OTPCode), []byte(req.OTP)) == 0 {
			return constants.Empty, errors.ErrInvalidOTP
		}

		if err := otpRepo.Delete(ctx, regisration.ID.Hex()); err != nil {
			return constants.Empty, errors.ErrOTPNotFound
		}

		return regisration.FullName, nil
	}

	filter := bson.M{
		"user_code": req.UserID,
		"otp_for":   req.OtpFor,
	}
	otpRecord, err := otpRepo.Find(ctx, filter)
	if err != nil {
		return constants.Empty, err
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		if err := otpRepo.Delete(ctx, otpRecord.ID.Hex()); err != nil {
			return constants.Empty, err
		}

		return constants.Empty, errors.ErrOtpExpired
	}

	if subtle.ConstantTimeCompare([]byte(otpRecord.OTPCode), []byte(req.OTP)) == 0 {
		return constants.Empty, errors.ErrInvalidOTP
	}

	if err := otpRepo.Delete(ctx, otpRecord.ID.Hex()); err != nil {
		return constants.Empty, errors.ErrOTPNotFound
	}

	return otpRecord.FullName, nil
}

func BuildUserData(fullName, phoneNumber, deviceUUID string) *model.User {
	return &model.User{
		Realm:       constants.MEMBER_REALM,
		UserCode:    utils.GenerateRandom(20),
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		KYCLevel:    0,
		DeviceUUID:  deviceUUID,
		IsVerified:  true,
		IsBlocked:   false,
		Enabled:     true,
	}
}

func BuildVerifyOtpResponse(userID, phone, token, nextStep string) *dto.VerifyOtpResponse {
	return &dto.VerifyOtpResponse{
		UserID:      userID,
		PhoneNumber: phone,
		OTPVerified: true,
		Token:       token,
		TokenType:   "verify_otp",
		TokenExpiry: time.Now().Add(10 * time.Minute),
		NextStep:    nextStep,
	}
}

func PinValidator(userPin, inputPin string) error {
	if subtle.ConstantTimeCompare([]byte(userPin), []byte(inputPin)) == 0 {
		return errors.ErrInvalidPIN
	}
	return nil
}

func BuildLoginResponse(user *model.User, token, deviceUUID string, sessionExpires time.Time) *dto.LoginResponse {
	return &dto.LoginResponse{
		Token:          token,
		UserID:         user.ID.Hex(),
		UserCode:       user.UserCode,
		FullName:       user.FullName,
		PhoneNumber:    user.PhoneNumber,
		KYCLevel:       user.KYCLevel,
		IsVerified:     user.IsVerified,
		DeviceUUID:     deviceUUID,
		LoginTime:      time.Now(),
		SessionExpires: sessionExpires,
		LastLogin:      user.LastLogin,
		LoginAttempts:  int(user.LoginAttemptCount),
	}
}

func CheckPinInHistory(pinHistory []string, newPin string) bool {
	for _, pin := range pinHistory {
		if pin == newPin {
			return true
		}
	}
	return false
}

func SetPinHistory(user *model.User, newPin string) types.LoginPIN {
	var newHistory [4]string
	copy(newHistory[1:], user.LoginPIN.PINHistory[:3])
	newHistory[0] = user.LoginPIN.PIN

	return types.LoginPIN{
		PIN:              newPin,
		PINHistory:       newHistory,
		LastPINCreatedAt: time.Now(),
	}
}

func BuildSetPinResponse(userID, token string) *dto.SetPinResponse {
	return &dto.SetPinResponse{
		UserID:      userID,
		PinSet:      true,
		Token:       token,
		TokenType:   "set_pin",
		TokenExpiry: time.Now().Add(24 * time.Hour),
	}
}

func OtpProvider(tokenService token.TokenService) (*string, *string, error) {
	otpCode := utils.GenerateRandom(6)

	encOtpCode, _, err := tokenService.LocalEncryptPassword(otpCode, constants.OTP, constants.OTP, constants.OTP)
	if err != nil {
		return nil, nil, err
	}

	return &otpCode, &encOtpCode, nil
}

func BuildRegistrationRecord(req dto.RegisterRequest, encOtpCode string) types.RegistrationRecord {
	registrationID := uuid.New().String()
	expirationTime := 10 * time.Minute

	return types.RegistrationRecord{
		ID:          registrationID,
		PhoneNumber: req.Phone,
		DeviceUUID:  req.DeviceUUID,
		Platform:    string(req.Platform),
		FullName:    req.FullName,
		OTP:         encOtpCode,
		OTPFor:      constants.OTPForRegistration,
		Status:      string(constants.Pending),
		ExpiresAt:   time.Now().Add(expirationTime),
		CreatedAt:   time.Now(),
		Attempts:    0,
		MaxAttempts: 3,
	}
}

func BuildOTPFromRegistration(registration types.RegistrationRecord) model.OTP {
	return model.OTP{
		PhoneNumber: registration.PhoneNumber,
		DeviceUUID:  &registration.DeviceUUID,
		OTPCode:     registration.OTP,
		FullName:    registration.FullName,
		OTPFor:      constants.OTPFor(registration.OTPFor),
		Status:      constants.OTPStatus(registration.Status),
		ExpiresAt:   registration.ExpiresAt,
		CreatedAt:   registration.CreatedAt,
	}
}

func BuildRegisterResponse(registrationID string, req dto.RegisterRequest, env string, otp string, wait int, token string) *dto.RegisterResponse {
	otpCode := constants.Empty
	if env == constants.DEV || env == constants.UAT {
		otpCode = otp
	}
	return &dto.RegisterResponse{
		RegistrationID:     registrationID,
		PhoneNumber:        req.Phone,
		DeviceUUID:         req.DeviceUUID,
		Platform:           string(req.Platform),
		Otp:                otpCode,
		OTPSent:            true,
		OTPExpiryMinutes:   wait,
		Token:              token,
		TokenType:          constants.Permanent,
		TokenExpiry:        time.Now().Add(24 * time.Hour),
		NextStep:           constants.VerifyOtp,
		RegistrationStatus: constants.Incomplete,
	}
}

func BuildPinResetSession(user model.User, encOtpCode string, expirationTime time.Duration) model.PinResetSession {

	return model.PinResetSession{
		ID:               bson.NewObjectID(),
		UserID:           user.ID.Hex(),
		PhoneNumber:      user.PhoneNumber,
		DeviceUUID:       user.DeviceUUID,
		OTP:              encOtpCode,
		OTPFor:           string(constants.OTPForPINReset),
		Status:           string(constants.Pending),
		ExpiresAt:        time.Now().Add(expirationTime + 5*time.Minute),
		CreatedAt:        time.Now(),
		Attempts:         0,
		Enabled:          false,
		MaxAttempts:      3,
		AccessRestricted: true,
	}
}

func ResetPinAdditionalBuilder(sessionID string) map[string]interface{} {
	return map[string]interface{}{
		"reset_session_id": sessionID,
		"token_type":       constants.OTPForForgetPin,
		"next_step":        constants.ForgetPinVerifyOtp,
	}
}

func BuildForgetPinSendOtpResponse(formattedPhone, deviceUUID, sessionID, token, otp, env string, wait int) *dto.ForgetPinSendOtpResponse {
	otpCode := ""
	if env == constants.DEV || env == constants.UAT {
		otpCode = otp
	}
	return &dto.ForgetPinSendOtpResponse{
		PhoneNumber:      formattedPhone,
		DeviceUUID:       deviceUUID,
		OTPSent:          true,
		OTPExpiryMinutes: wait,
		ResetSessionID:   sessionID,
		Token:            token,
		TokenType:        string(constants.OTPForForgetPin),
		TokenExpiry:      time.Now().Add(10 * time.Minute),
		NextStep:         constants.ForgetPinVerifyOtp,
		OTP:              otpCode,
	}
}

func ValidatePinResetSession(session *model.PinResetSession, req dto.VerifyForgetPinOtpRequest) error {

	if session.DeviceUUID != req.DeviceUUID || session.PhoneNumber != req.Phone {
		return errors.ErrPinResetDeviceMismatch
	}

	if !session.Enabled {
		return errors.ErrPinNeedActivation
	}

	if session.Status != string(constants.Pending) {
		return errors.ErrPinResetSessionInvalid
	}

	if session.Attempts >= session.MaxAttempts {
		return errors.ErrTooManyResetAttempts
	}

	return nil
}

func BuildResetPinVerifyOtpResponse(id, phone, token string) *dto.VerifyOtpResponse {
	return &dto.VerifyOtpResponse{
		UserID:      id,
		PhoneNumber: phone,
		OTPVerified: true,
		Token:       token,
		TokenType:   constants.ResetPin,
		TokenExpiry: time.Now().Add(24 * time.Hour),
		NextStep:    constants.ResetPin,
	}
}

func ValidatePinResetSessionToken(session *model.PinResetSession, req dto.ResetPinRequest) error {

	if time.Now().After(session.ExpiresAt) {
		return errors.ErrPinResetSessionExpired
	}
	if session.DeviceUUID != req.DeviceUUID || session.PhoneNumber != req.Phone {
		return errors.ErrPinResetDeviceMismatch
	}

	if !session.Enabled {
		return errors.ErrPinNeedActivation
	}

	if session.Status != string(constants.Pending) {
		return errors.ErrPinResetSessionInvalid
	}

	if session.Attempts >= session.MaxAttempts {
		return errors.ErrTooManyResetAttempts
	}

	return nil
}

func BuildResetPinResponse(user *model.User, session *model.PinResetSession, token string) *dto.ResetPinResponse {
	return &dto.ResetPinResponse{
		UserID:           user.ID.Hex(),
		UserCode:         user.UserCode,
		FullName:         user.FullName,
		PhoneNumber:      user.PhoneNumber,
		PinReset:         true,
		ResetTime:        time.Now(),
		AccessRestricted: session.AccessRestricted,
		Restrictions:     session.Restrictions,
		Token:            token,
		TokenType:        constants.Permanent,
		TokenExpiry:      time.Now().Add(24 * time.Hour),
		NextStep:         "login_with_new_pin",
	}
}

func FileBucketUploader(ctx context.Context, minioServer config.MinioClientInterface, bucketName, objectName, filePath string) (string, error) {
	exists, err := minioServer.BucketExist(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		if _, err = minioServer.MakeBucket(ctx, bucketName); err != nil {
			return "", errors.ErrFailedToUpload
		}
	}

	bucketResp, err := minioServer.SaveObject(ctx, config.SaveObjectBody{
		BucketName:  bucketName,
		ObjectName:  objectName,
		File:        filePath,
		ContentType: "jpeg",
	})
	if err != nil {
		return "", errors.ErrFailedToUpload
	}

	return bucketResp.Key, nil
}
