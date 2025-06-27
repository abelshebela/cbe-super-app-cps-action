package users_application

import (
	"context"
	"mime/multipart"
	"time"

	dto "cbe-super-app-member-users/internal/application/dto"
	domainUsers "cbe-super-app-member-users/internal/domain/users"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"
	"cbe-super-app-member-users/pkgs/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Utility function to convert map[string]interface{} to bson.M
func mapToBsonM(m map[string]interface{}) bson.M {
	bsonMap := bson.M{}
	for k, v := range m {
		bsonMap[k] = v
	}
	return bsonMap
}

func (h UsersHandler) CheckDevice(ctx context.Context, header map[string]interface{}) (map[string]interface{}, error) {
	var returndata = make(map[string]interface{})
	hqdata, err := h.GetOneHQ(ctx, header)
	if err != nil {
		h.logger.Errorf("Failed to fetch HQ data: %v", err)
		return nil, err
	}

	userdata, err := h.userService.GetOneUser(ctx, header)
	if err != nil {
		h.logger.Errorf("Device Not Found: %v", err)
		return nil, err
	}

	if userdata.PinStatus == "" {
		h.logger.Infof("Your pin already changed")
		returndata["pin_status"] = "pin_changed"
	} else {
		var OTPCode = utils.OTPGenerator(6)
		encCode, _, _ := utils.LocalEncryptPassword(OTPCode, "otp", "", "", nil)

		utils.AxiosSendSms(ctx, userdata.PhoneNumber, "Your new pin is for CBE Super App is "+encCode)
		returndata["pin_status"] = "pin_not_changed"
		returndata["otp_code"] = OTPCode
		h.userService.CreateOtp(ctx, domainUsers.OTPRecord{
			UserID:     userdata.ID.Hex(),
			OTP:        encCode,
			ExpiresAt:  time.Now().Add(time.Minute * 5),
			CreatedAt:  time.Now(),
			UserRealm:  "member",
			DeviceUUID: header["device_uuid"].(*string),
			OTPFor:     "pin_set",
		})
	}

	appVersionStr, _ := header["appVersion"].(string)
	var isLatest bool
	if platform, ok := header["platform"].(string); ok && platform == "android" {
		isLatest = appVersionStr >= hqdata.LatestAndroidVersion
	} else {
		isLatest = appVersionStr >= hqdata.LatestiOSVersion
	}
	returndata["is_latest"] = isLatest
	returndata["message"] = "Your pin is not changed, please check your phone for the new pin"

	return returndata, nil
}

func (h UsersHandler) FetchLinkedAccounts(ctx context.Context, userID string) (*userPort.LinkedAccountResponse, error) {
	domainResponse, err := h.userService.FetchLinkedAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}

	portAccounts := make([]userPort.LinkedAccountDetail, len(domainResponse.LinkedAccounts))
	for i, acc := range domainResponse.LinkedAccounts {
		portAccounts[i] = userPort.LinkedAccountDetail{
			AccountNumber:     acc.AccountNumber,
			AccountBranchCode: acc.AccountBranchCode,
			LinkedBranch:      acc.LinkedBranch,
			IsAccountActive:   acc.IsAccountActive,
			LinkedStatus:      acc.LinkedStatus,
			CurrencyCode:      acc.CurrencyCode,
		}
	}

	return &userPort.LinkedAccountResponse{
		UserID:         domainResponse.UserID,
		FullName:       domainResponse.FullName,
		LinkedAccounts: portAccounts,
	}, nil
}

func (h UsersHandler) GetOneHQ(ctx context.Context, req map[string]interface{}) (*domainUsers.HQ, error) {
	filter := mapToBsonM(req)
	hq, err := h.userService.GetOneHQ(ctx, filter)
	if err != nil {
		return nil, err
	}
	return hq, nil
}

func (h UsersHandler) GenerateEmailOTP(ctx context.Context, req userPort.OTPRequest) (string, error) {
	domainReq := domainUsers.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	}
	return h.userService.GenerateEmailOTP(ctx, domainReq)
}

func (h UsersHandler) VerifyEmailOTP(ctx context.Context, verification userPort.OTPVerification) error {
	domainVerification := domainUsers.OTPVerification{
		UserID: verification.UserID,
		Email:  verification.Email,
		OTP:    verification.OTP,
	}
	return h.userService.VerifyEmailOTP(ctx, domainVerification)
}

func (h UsersHandler) HandleFetchLinkedAccounts(ctx context.Context, req dto.FetchLinkedAccountsRequest) (*userPort.LinkedAccountResponse, error) {
	return h.FetchLinkedAccounts(ctx, req.UserID)
}

func (h UsersHandler) HandleGenerateEmailOTP(ctx context.Context, req dto.GenerateOTPRequest) (*dto.GenerateOTPResponse, error) {
	portReq := userPort.OTPRequest{
		UserID: req.UserID,
		Email:  req.Email,
	}

	otp, err := h.GenerateEmailOTP(ctx, portReq)
	if err != nil {
		return nil, err
	}

	return &dto.GenerateOTPResponse{OTP: otp}, nil
}

func (h UsersHandler) HandleVerifyEmailOTP(ctx context.Context, req dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	portVerification := userPort.OTPVerification{
		UserID: req.UserID,
		Email:  req.Email,
		OTP:    req.OTP,
	}

	err := h.VerifyEmailOTP(ctx, portVerification)
	if err != nil {
		return nil, err
	}

	return &dto.VerifyOTPResponse{Success: true}, nil
}

func (h UsersHandler) UpdateProfilePicture(ctx context.Context, id string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	imageURL, err := h.userService.UpdateProfilePicture(ctx, id, file, fileHeader)
	if err != nil {
		return "", err
	}
	return imageURL, nil
}

func (h UsersHandler) UnlinkDevice(ctx context.Context, userID string, deviceID string) error {
	err := h.userService.UnlinkDevice(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	return nil
}

func (h UsersHandler) ChangePin(ctx context.Context, req userPort.ChangePinRequest) error {
	domainChangePinReq := domainUsers.ChangePinRequest{
		UserID: req.UserID,
		OldPin: req.OldPin,
		NewPin: req.NewPin,
	}
	err := h.userService.ChangePin(ctx, domainChangePinReq)
	if err != nil {
		return err
	}
	return nil
}

func (h UsersHandler) VerifyOtp(ctx context.Context, userID, otp string, deviceUUID *string, userRealm, otpFor string) error {
	return h.userService.VerifyOtp(ctx, userID, otp, deviceUUID, userRealm, otpFor)
}

func (h UsersHandler) Register(ctx context.Context, phone, deviceUUID, platform string) error {
	return h.userService.Register(ctx, phone, deviceUUID, platform)
}

func (h UsersHandler) Login(ctx context.Context, phone, deviceUUID, pin string) (string, error) {
	return h.userService.Login(ctx, phone, deviceUUID, pin)
}

func (h UsersHandler) ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) error {
	return h.userService.ForgetPinSendOtp(ctx, phone, deviceUUID)
}
