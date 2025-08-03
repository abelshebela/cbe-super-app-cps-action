package core

import (
	"context"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	local_util "github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsAppVersionLatest(platform, appVersion string, hqData *model.HQ) bool {
	deviceApp, err := strconv.ParseFloat(appVersion, 32)
	if err != nil {
		return false
	}

	hqAndroidVersion, androidErr := strconv.ParseFloat(hqData.LatestAndroidVersion, 32)
	hqIosVersion, iosErr := strconv.ParseFloat(hqData.LatestiOSVersion, 32)
	if androidErr != nil || iosErr != nil {
		return false
	}

	switch platform {
	case string(constants.Android):
		return deviceApp >= hqAndroidVersion
	case string(constants.Ios):
		return deviceApp >= hqIosVersion
	default:
		return false
	}
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
