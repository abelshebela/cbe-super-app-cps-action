package users_application

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	dto "cbe-super-app-member-users/internal/application/dto"
	domainUsers "cbe-super-app-member-users/internal/domain/users"
	userPort "cbe-super-app-member-users/internal/port/inbound/users"
	"cbe-super-app-member-users/pkgs/common"
	"cbe-super-app-member-users/pkgs/entities"
	localModel "cbe-super-app-member-users/pkgs/entities"
	"cbe-super-app-member-users/pkgs/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Constants for magic numbers and error messages
const (
	pinLength           = 6
	otpLength           = 6
	minPhoneLength      = 10
	minDeviceUUIDLength = 10
)

var (
	errMissingDeviceUUID  = fmt.Errorf("MISSING_DEVICE_UUID")
	errMissingPlatform    = fmt.Errorf("MISSING_PLATFORM")
	errInvalidPhoneNumber = fmt.Errorf("INVALID_PHONE_NUMBER")
	errInvalidDeviceUUID  = fmt.Errorf("INVALID_DEVICE_UUID")
	errInvalidPlatform    = fmt.Errorf("INVALID_PLATFORM")
	errInvalidPin         = fmt.Errorf("INVALID_PIN")
	errInvalidOTP         = fmt.Errorf("INVALID_OTP")
	errInvalidInputParams = fmt.Errorf("INVALID_INPUT_PARAMETERS")
)

// validatePin checks if the pin is valid according to business rules
func validatePin(pin string) error {
	if len(pin) != pinLength {
		return errInvalidPin
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return errInvalidPin
		}
	}
	return nil
}

// validateOTP checks if the OTP is valid according to business rules
func validateOTP(otp string) error {
	if len(otp) != otpLength {
		return errInvalidOTP
	}
	for _, c := range otp {
		if c < '0' || c > '9' {
			return errInvalidOTP
		}
	}
	return nil
}

// validatePhone checks if the phone number is valid
func validatePhone(phone string) error {
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < minPhoneLength {
		return errInvalidPhoneNumber
	}
	return nil
}

// validateDeviceUUID checks if the device UUID is valid
func validateDeviceUUID(deviceUUID string) error {
	if len(deviceUUID) < minDeviceUUIDLength {
		return errInvalidDeviceUUID
	}
	return nil
}

// Utility function to convert map[string]interface{} to bson.M
func mapToBsonM(m map[string]interface{}) bson.M {
	bsonMap := bson.M{}
	for k, v := range m {
		bsonMap[k] = v
	}
	return bsonMap
}

func (h UsersHandler) CheckDevice(ctx context.Context, header map[string]interface{}) (map[string]interface{}, error) {
	// Extract device information from headers
	deviceUUID, ok := header["device_uuid"].(string)
	if !ok || deviceUUID == "" {
		return nil, fmt.Errorf("MISSING_DEVICE_UUID")
	}

	platform, ok := header["platform"].(string)
	if !ok || platform == "" {
		return nil, fmt.Errorf("MISSING_PLATFORM")
	}

	appVersion, ok := header["app_version"].(string)
	if !ok {
		appVersion = ""
	}

	// Get HQ data for app version check
	hqData, err := h.userService.GetOneHQ(ctx, map[string]interface{}{})
	if err != nil {
		h.logger.Errorf("Failed to fetch HQ data: %v", err)
		return nil, fmt.Errorf("failed to fetch system configuration: %w", err)
	}

	// Look up user by device UUID
	userData, err := h.userService.GetOneUser(ctx, map[string]interface{}{"device.device_uuid": deviceUUID})
	if err != nil {
		h.logger.Infof("Device not found for UUID: %s", deviceUUID)
		return map[string]interface{}{
			"message":    "Device Not Found",
			"registered": false,
			"status":     "false",
			"code":       common.DefineSuccess.Auth["DEVICE_NOT_FOUND"].Code,
		}, nil
	}

	// Check if app version is latest
	isLatest := h.isAppVersionLatest(platform, appVersion, hqData)

	// Generate temporary token for device lookup
	tempToken, err := utils.TempTokenMaker(userData, []string{"access"}, "", nil, utils.DeviceLookup, h.vaultConfig)
	if err != nil {
		h.logger.Errorf("Failed to generate temp token: %v", err)
		return nil, fmt.Errorf("failed to generate temporary token: %w", err)
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"message":      common.DefineSuccess.Auth["DEVICE_FOUND"].Message,
		"registered":   true,
		"status":       "true",
		"user_id":      userData.ID.Hex(),
		"user_code":    userData.UserCode,
		"full_name":    userData.FullName,
		"phone_number": userData.PhoneNumber,
		"temp_token":   tempToken,
		"is_latest":    isLatest,
	}

	// Handle unverified users (generate OTP for PIN change)
	if !userData.IsVerified {
		if err := h.handleUnverifiedUser(ctx, userData, responseData); err != nil {
			return nil, err
		}
	}

	return responseData, nil
}

// isAppVersionLatest checks if the current app version is the latest
func (h UsersHandler) isAppVersionLatest(platform, appVersion string, hqData *domainUsers.HQ) bool {
	switch platform {
	case "android":
		return appVersion >= hqData.LatestAndroidVersion
	case "ios":
		return appVersion >= hqData.LatestiOSVersion
	default:
		return appVersion >= hqData.LatestiOSVersion // Default to iOS version
	}
}

// handleUnverifiedUser generates OTP for unverified users
func (h UsersHandler) handleUnverifiedUser(ctx context.Context, userData *localModel.User, responseData map[string]interface{}) error {
	responseData["otp_status"] = "unverified"

	// Generate OTP
	otpCode := utils.OTPGenerator(6)
	encOtpCode, _, err := utils.LocalEncryptPassword(otpCode, "otp", "", "", h.vaultConfig)
	if err != nil {
		h.logger.Errorf("Failed to encrypt OTP: %v", err)
		return fmt.Errorf("failed to encrypt OTP: %w", err)
	}

	// Add OTP code to response in development environments
	if h.vaultConfig.GoEnv == "dev" || h.vaultConfig.GoEnv == "uat" {
		responseData["otp_code"] = otpCode
	}

	// Calculate OTP expiration time
	wait, err := strconv.Atoi(h.vaultConfig.OtpWaitingTime)
	if err != nil {
		wait = 10 // Default to 10 minutes
	}
	expirationTime := time.Duration(wait) * time.Minute

	// Create OTP record
	otpRecord := domainUsers.OTPRecord{
		UserID:     userData.ID.Hex(),
		OTP:        encOtpCode,
		OTPFor:     "change_pin",
		UserRealm:  string(userData.Realm),
		ExpiresAt:  time.Now().Add(expirationTime),
		CreatedAt:  time.Now(),
		DeviceUUID: userData.Device.DeviceUUID,
	}

	// Store OTP in database
	if err := h.userService.CreateOtp(ctx, otpRecord); err != nil {
		h.logger.Errorf("Failed to create OTP: %v", err)
		return fmt.Errorf("failed to create OTP record: %w", err)
	}

	// Send OTP via SMS (async)
	go func() {
		message := common.DefineSuccess.Auth["AUTH_USER_OTP_SENT_USER"].Message + otpCode
		if err := utils.AxiosSendSms(ctx, userData.PhoneNumber, message); err != nil {
			h.logger.Errorf("Failed to send SMS: %v", err)
		}
	}()

	// Update response for unverified users
	responseData["status"] = "false"
	responseData["message"] = "Your pin is not changed, please check your phone for the new pin"

	return nil
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
	// Always fetch the latest HQ data, ignore the filter from the header
	filter := bson.M{"enabled": true} // or {} if you want the very latest regardless of status
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
func (h UsersHandler) UpdateProfileTheme(ctx context.Context, id string, themeType string) (*entities.User, error) {
	data, err := h.userService.UpdateProfileTheme(ctx, id, themeType)
	if err != nil {
		return nil, err
	}
	resp := &entities.User{
		FullName:    data.FullName,
		Avatar:      data.Avatar,
		PhoneNumber: data.PhoneNumber,
		Device: struct {
			DeviceUUID string "json:\"device_uuid\" bson:\"device_uuid\""
			AppVersion string "json:\"app_version\" bson:\"app_version\""
		}{
			DeviceUUID: data.Device.DeviceUUID,
		},
	}
	return resp, nil
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

func (h UsersHandler) DeviceLookup(ctx context.Context, header map[string]interface{}) (*dto.DeviceLookupResponse, error) {
	deviceUUID, ok := header["device_uuid"].(string)
	if !ok || deviceUUID == "" {
		return nil, fmt.Errorf("MISSING_DEVICE_UUID")
	}

	platform, ok := header["platform"].(string)
	if !ok || platform == "" {
		return nil, fmt.Errorf("MISSING_PLATFORM")
	}

	appVersion, ok := header["app_version"].(string)
	if !ok {
		appVersion = ""
	}

	sourceApp, ok := header["source_app"].(string)
	if !ok {
		sourceApp = ""
	}

	return h.userService.DeviceLookup(ctx, deviceUUID, platform, appVersion, sourceApp)
}

func (h UsersHandler) PreLogin(ctx context.Context, header map[string]interface{}, phone string) (*dto.DeviceLookupResponse, error) {
	deviceUUID, ok := header["device_uuid"].(string)
	if !ok || deviceUUID == "" {
		return nil, fmt.Errorf("MISSING_DEVICE_UUID")
	}

	platform, ok := header["platform"].(string)
	if !ok || platform == "" {
		return nil, fmt.Errorf("MISSING_PLATFORM")
	}

	appVersion, ok := header["app_version"].(string)
	if !ok {
		appVersion = ""
	}

	sourceApp, ok := header["source_app"].(string)
	if !ok {
		sourceApp = ""
	}

	return h.userService.PreLogin(ctx, deviceUUID, platform, appVersion, sourceApp, phone)
}

func (h UsersHandler) VerifyOtp(ctx context.Context, userID, phone_number, otp string, deviceUUID string, userRealm, otpFor, action string) (*dto.VerifyOtpResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	if err := h.validateOtpInputs(userID, otp, userRealm, otpFor); err != nil {
		return nil, err
	}
	defer func() { otp = "" }()
	encOtp, _, err := utils.LocalEncryptPassword(otp, "otp", "", "", h.vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt OTP: %w", err)
	}

	return h.userService.VerifyOtp(ctx, userID, phone_number, encOtp, deviceUUID, userRealm, otpFor, action)
}

func (h UsersHandler) SetPin(ctx context.Context, userID, newPin, deviceUUID string, userRealm string) (*dto.SetPinResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	if err := h.validateSetPinInputs(userID, newPin, userRealm); err != nil {
		return nil, err
	}

	return h.userService.SetPin(ctx, userID, newPin, deviceUUID, userRealm)
}
func (h UsersHandler) SavePinToHistory(ctx context.Context, userID string, newPin string) ([]string, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return []string{}, err
	}
	filter := bson.M{
		"_id":        id,
		"is_deleted": false,
	}
	user, err := h.userService.GetOneUser(ctx, filter)
	if err != nil {
		h.logger.Errorf("unable to get user data for Pin history: %v", err)
		return []string{}, err
	}

	passHistory := user.LoginPIN.PINHistory[:]
	// Append the new Pin to the history
	passHistory = append(passHistory, newPin)

	// Shorten the passHistory to a maximum length of 4
	if len(passHistory) > 4 {
		passHistory = passHistory[len(passHistory)-4:]
	}
	// Keep only the last 5 Pins (or adjust as needed)
	const maxHistory = 5
	if len(passHistory) > maxHistory {
		passHistory = passHistory[len(passHistory)-maxHistory:]
	}

	return passHistory, nil
}

// validateOtpInputs validates the input parameters for OTP verification
func (h UsersHandler) validateOtpInputs(userID, otp, userRealm, otpFor string) error {
	if userID == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: user ID is required")
	}
	if otp == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: OTP is required")
	}
	if userRealm == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: user realm is required")
	}
	if otpFor == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: OTP purpose is required")
	}
	return nil
}

func (h UsersHandler) Register(ctx context.Context, phone, full_name, email, deviceUUID, platform string) (*dto.RegisterResponse, error) {
	if err := h.validateRegistrationInputs(phone, deviceUUID, platform); err != nil {
		return nil, err
	}

	return h.userService.Register(ctx, phone, full_name, email, deviceUUID, platform)
}

// validateRegistrationInputs validates the registration input parameters
func (h UsersHandler) validateRegistrationInputs(phone, deviceUUID, platform string) error {
	if phone == "" {
		return fmt.Errorf("INVALID_INPUT")
	}
	if deviceUUID == "" {
		return fmt.Errorf("INVALID_INPUT")
	}
	if platform == "" {
		return fmt.Errorf("INVALID_INPUT")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER")
	}

	// Validate platform
	validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
	if !validPlatforms[platform] {
		return fmt.Errorf("INVALID_PLATFORM")
	}

	// Validate device UUID format (basic validation)
	if len(deviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID")
	}

	return nil
}

func (h UsersHandler) Login(ctx context.Context, phone, deviceUUID, pin string) (*dto.LoginResponse, error) {
	if err := h.validateLoginInputs(phone, deviceUUID, pin); err != nil {
		return nil, err
	}

	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" {
		return nil, fmt.Errorf("INVALID_PHONE_NUMBER")
	}

	loginResult, err := h.userService.Login(ctx, formattedPhone, deviceUUID, pin)
	if err != nil {
		return nil, err
	}

	return loginResult, nil
}

func (h UsersHandler) validateLoginInputs(phone, deviceUUID, pin string) error {
	if phone == "" {
		return fmt.Errorf("INVALID_INPUT: phone number is required")
	}
	if deviceUUID == "" {
		return fmt.Errorf("INVALID_INPUT: device UUID is required")
	}
	if pin == "" {
		return fmt.Errorf("INVALID_INPUT: PIN is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER")
	}

	// Validate PIN format (6 digits)
	if len(pin) != 6 {
		return fmt.Errorf("INVALID_PIN")
	}

	// Validate PIN contains only digits
	for _, char := range pin {
		if char < '0' || char > '9' {
			return fmt.Errorf("INVALID_PIN")
		}
	}

	// Validate device UUID format (basic validation)
	if len(deviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID")
	}

	return nil
}

func (h UsersHandler) ForgetPinSendOtp(ctx context.Context, phone, deviceUUID string) (*dto.ForgetPinSendOtpResponse, error) {
	return h.userService.ForgetPinSendOtp(ctx, phone, deviceUUID)
}

func (h UsersHandler) CompleteRegistration(ctx context.Context, registrationID, phone, deviceUUID, platform, fullName string) (*dto.CompleteRegistrationResponse, error) {
	// Call domain service for registration completion
	return h.userService.CompleteRegistration(ctx, registrationID, phone, deviceUUID, platform, fullName)
}

func (h UsersHandler) ResetPin(ctx context.Context, resetSessionID, phone, deviceUUID, otp, newPin string) (*dto.ResetPinResponse, error) {
	return h.userService.ResetPin(ctx, resetSessionID, phone, deviceUUID, otp, newPin)
}

// validateSetPinInputs validates the set PIN input parameters
func (h UsersHandler) validateSetPinInputs(userID, newPin, userRealm string) error {
	if userID == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: user ID is required")
	}
	if newPin == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: new PIN is required")
	}

	if userRealm == "" {
		return fmt.Errorf("INVALID_INPUT_PARAMETERS: user realm is required")
	}

	// Validate PIN format (6 digits)
	if len(newPin) != 6 {
		return fmt.Errorf("INVALID_PIN: PIN must be exactly 6 digits")
	}

	// Validate PIN contains only digits
	for _, char := range newPin {
		if char < '0' || char > '9' {
			return fmt.Errorf("INVALID_PIN: PIN must contain only digits")
		}
	}

	return nil
}

func (h UsersHandler) VerifyForgetPinOtp(ctx context.Context, resetSessionID, phone, deviceUUID, otp string) (*dto.VerifyOtpResponse, error) {
	return h.userService.VerifyForgetPinOtp(ctx, resetSessionID, phone, deviceUUID, otp)
}

func (h UsersHandler) ResetPinWithToken(ctx context.Context, resetSessionID, phone, deviceUUID, newPin string) (*dto.ResetPinResponse, error) {
	return h.userService.ResetPinWithToken(ctx, resetSessionID, phone, deviceUUID, newPin)
}
