package users

import (
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"
	"strings"

	"cbe-super-app-member-users/internal/application/dto"
	user_inbound "cbe-super-app-member-users/internal/port/inbound/users"

	"cbe-super-app-member-users/pkgs/utils"
	constant "cbe-super-app-member-users/pkgs/utils"
)

// Constants for magic numbers and error messages
const (
	pinLength           = 6
	otpLength           = 6
	minPhoneLength      = 10
	minDeviceUUIDLength = 10
)

var (
	errMissingDeviceUUID  = "MISSING_DEVICE_UUID"
	errMissingPlatform    = "MISSING_PLATFORM"
	errInvalidPhoneNumber = "INVALID_PHONE_NUMBER"
	errInvalidDeviceUUID  = "INVALID_DEVICE_UUID"
	errInvalidPlatform    = "INVALID_PLATFORM"
	errInvalidPin         = "INVALID_PIN"
	errInvalidOTP         = "INVALID_OTP"
	errInvalidInputParams = "INVALID_INPUT_PARAMETERS"
)

// validatePin checks if the pin is valid according to business rules
func validatePin(pin string) error {
	if len(pin) != pinLength {
		return fmt.Errorf(errInvalidPin+": PIN must be exactly %d digits", pinLength)
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return fmt.Errorf(errInvalidPin + ": PIN must contain only digits")
		}
	}

	// Check for weak patterns
	if isWeakPin(pin) {
		return fmt.Errorf("PIN contains weak patterns")
		return fmt.Errorf(errInvalidPin + ": PIN contains weak patterns")
	}
	return nil
}

func isWeakPin(pin string) bool {
	// Check for repeated digits
	if strings.Count(pin, string(pin[0])) == len(pin) {
		return true
	}
	// Check for sequential patterns
	isAscending := true
	isDescending := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			isAscending = false
		}
		if pin[i] != pin[i-1]-1 {
			isDescending = false
		}
	}
	return isAscending || isDescending
}

// validateOTP checks if the OTP is valid according to business rules
func validateOTP(otp string) error {
	if len(otp) != otpLength {
		return fmt.Errorf(errInvalidOTP+": OTP must be exactly %d digits", otpLength)
	}

	for _, c := range otp {
		if c < '0' || c > '9' {
			return fmt.Errorf(errInvalidOTP + ": OTP must contain only digits")
		}
	}
	return nil
}

// validatePhone checks if the phone number is valid
func validatePhone(phone string) error {
	formattedPhone := utils.FormatPhoneNumber(phone)
	if formattedPhone == "" || len(formattedPhone) < minPhoneLength {
		return fmt.Errorf(errInvalidPhoneNumber + ": please provide a valid phone number")
	}
	return nil
}

// validateDeviceUUID checks if the device UUID is valid
func validateDeviceUUID(deviceUUID string) error {
	if len(deviceUUID) < minDeviceUUIDLength {
		return fmt.Errorf(errInvalidDeviceUUID + ": device UUID appears to be invalid")
	}
	return nil
}

func (h UsersAdapter) DeviceLookup(w http.ResponseWriter, r *http.Request) {
	platform, appVersion, deviceUUID, sourceApp, installationDate, additionalHeaders := utils.HeaderRequirement(r, nil)

	if deviceUUID == "" || platform == "" || appVersion == "" || installationDate == "" {
		utils.BaseResponseMaker(nil, w, "Missing required headers", http.StatusBadRequest)
		return
	}

	headerData := map[string]interface{}{
		"platform":          platform,
		"app_version":       appVersion,
		"device_uuid":       deviceUUID,
		"source_app":        sourceApp,
		"installation_date": installationDate,
	}

	// Add any additional headers
	for key, value := range additionalHeaders {
		headerData[key] = value
	}

	response, err := h.Application.DeviceLookup(r.Context(), headerData)
	if err != nil {
		utils.BaseResponseMaker(nil, w, fmt.Sprintf("Device lookup failed %v", err.Error()), http.StatusNoContent)
		return
	}

	data, _ := utils.StructToMap(response)
	utils.BaseResponseMaker(data, w, response.Message, response.Status)
}

func (h UsersAdapter) PreLogin(w http.ResponseWriter, r *http.Request) {
	platform, appVersion, deviceUUID, sourceApp, installationDate, additionalHeaders := utils.HeaderRequirement(r, nil)

	if deviceUUID == "" || platform == "" {
		utils.BaseResponseMaker(nil, w, "Missing required headers", http.StatusBadRequest)
		return
	}
	var req dto.PhoneLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_PAYLOAD", 403, nil)
		return
	}

	headerData := map[string]interface{}{
		"platform":          platform,
		"app_version":       appVersion,
		"device_uuid":       deviceUUID,
		"source_app":        sourceApp,
		"installation_date": installationDate,
	}

	for key, value := range additionalHeaders {
		headerData[key] = value
	}

	response, err := h.Application.PreLogin(r.Context(), headerData, req.Phone, installationDate)
	if err != nil {

		if err.Error() == "BLOCK_BY_MULTIPLE_TRIES" {
			utils.SendErrorResponse(w, "user is blocked by multiple tries contact the nearest branch", 400, nil)
			return
		} else if err.Error() == "USER_DISABLED_BLOCKED" {
			utils.SendErrorResponse(w, "user is blocked or disabled contact nearest branch", 400, nil)
			return
		}
		utils.SendErrorResponse(w, err.Error(), 400, nil)
		return
	}

	data, _ := utils.StructToMap(response)
	utils.BaseResponseMaker(data, w, response.Message, response.Status)
}

func (h UsersAdapter) ForgetPin(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	// h.logger.Infof("User ID from context: %s", userID, "here we are")
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	req := dto.FetchLinkedAccountsRequest{UserID: userID}

	response, err := h.Application.FetchLinkedAccounts(r.Context(), req.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		// h.logger.Infof("here we are ")
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h UsersAdapter) FetchLinkedAccounts(w http.ResponseWriter, r *http.Request) {
	// userID, ok := r.Context().Value("user_id").(string)

	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	// h.logger.Infof("User ID from context: %s", userID, "here we are")
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	req := dto.FetchLinkedAccountsRequest{UserID: userID}

	response, err := h.Application.FetchLinkedAccounts(r.Context(), req.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		// h.logger.Infof("here we are ")
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h UsersAdapter) GenerateEmailOTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	var req user_inbound.OTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	req.UserID = userID

	response, err := h.Application.GenerateEmailOTP(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h UsersAdapter) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		utils.SendErrorResponse(w, "INVALID_FORM", http.StatusBadRequest, nil)
		return
	}

	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		fmt.Println("error", err)
		utils.SendErrorResponse(w, "NO_FILE", http.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		utils.SendErrorResponse(w, "INVALID_FILE_TYPE", http.StatusBadRequest, nil)
		return
	}

	imageURL, err := h.Application.UpdateProfilePicture(r.Context(), userID, file, fileHeader)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, struct {
		Message string `json:"message"`
		URL     string `json:"url"`
	}{
		Message: "Image uploaded!",
		URL:     imageURL,
	})
}

func (h UsersAdapter) UpdateProfileTheme(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	var req dto.SetProfileThemeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_PAYLOAD", 403, nil)
		return
	}

	_, err := h.Application.UpdateProfileTheme(r.Context(), userID, req.ThemeType)
	if err != nil {
		status := http.StatusInternalServerError

		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	utils.BaseResponseMaker(map[string]interface{}{}, w, "Profile theme set successfuly", 200)
}
func (h UsersAdapter) VerifyEmailOTP(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	var req user_inbound.OTPVerification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	req.UserID = userID

	err := h.Application.VerifyEmailOTP(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "INVALID_OTP" || err.Error() == "EXPIRED_OTP" {
			status = http.StatusBadRequest
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "OTP verified successfully")
}

func (h UsersAdapter) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	deviceUuid := r.Header.Get("device-uuid")

	if deviceUuid == "" {
		// fmt.Println("Device UUID is required")
		utils.SendErrorResponse(w, "DEVICE_ID_REQUIRED", http.StatusBadRequest, nil)
		return
	}
	err := h.Application.UnlinkDevice(r.Context(), userID, deviceUuid)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}
	h.sendSuccessResponse(w, http.StatusOK, "Device Unlinked Successfully")

}

func (h UsersAdapter) ChangePin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	var req user_inbound.ChangePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	req.UserID = userID

	err := h.Application.ChangePin(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "INVALID_PIN" || err.Error() == "WEAK_PIN" {
			status = http.StatusBadRequest
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "PIN changed successfully")
}

func (h UsersAdapter) Healthcheck(w http.ResponseWriter, r *http.Request) {
	returndata := make(map[string]interface{})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	returndata["message"] = "LDAP AUTH IS ACTIVE"
	json.NewEncoder(w).Encode(returndata)
	return
}
func (h UsersAdapter) CheckPin(w http.ResponseWriter, r *http.Request) {

	var payload dto.PinStrengthRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// this is util based reusable function make to pass the error key and will send to user with status code
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}

	pin := strconv.Itoa(payload.NewPin)
	err := validatePin(pin)
	if err != nil {
		// this is util based reusable function make to pass the error key and will send to user with status code
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, map[string]interface{}{"errors": []string{err.Error()}})
		return
	}

	utils.BaseResponseMaker(nil, w, "PIN set Successfully", 200)
}

func (h UsersAdapter) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: status,
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

type VerifyOtpRequest struct {
	Otp    string `json:"otp_code" validate:"required,min=6,max=6"`
	OtpFor string `json:"otp_for,omitempty"`
}

func (h UsersAdapter) VerifyOtp(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok || userID == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	nextStep, ok := r.Context().Value(constant.ContextKey("next_step")).(string)
	if !ok || nextStep != "verify_otp" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	action, ok := r.Context().Value(constant.ContextKey("action")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	phoneNumber := r.Context().Value(constant.ContextKey("phone_number")).(string)

	_, _, deviceUUID, sourceApp, _, _ := utils.HeaderRequirement(r, nil)

	var req VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}

	if req.Otp == "" {
		utils.SendErrorResponse(w, "MISSING_OTP", http.StatusBadRequest, nil)
		return
	}

	allowedOtpTypes := map[string]bool{"PIN_SET": true, "REGISTRATION": true, "PIN_RESET": true, "LOGIN": true, "ENABLE": true}
	if !allowedOtpTypes[req.OtpFor] {
		utils.SendErrorResponse(w, "INVALID_OTP_TYPE", http.StatusBadRequest, nil)
		return
	}
	fmt.Println("===========checkpoint 1==========")
	data, err := h.Application.VerifyOtp(r.Context(), userID, phoneNumber, req.Otp, deviceUUID, sourceApp, req.OtpFor, action)

	fmt.Println("***************************")
	fmt.Println(err)
	fmt.Println("***************************")
	fmt.Println(data)
	if err != nil {

		utils.BaseResponseMaker(map[string]interface{}{}, w, "OTP verification failed", 500)
		return
	}

	dataResponse, err := utils.StructToMap(data)
	if err != nil {

		utils.BaseResponseMaker(map[string]interface{}{}, w, "OTP verification failed", 500)
		return
	}

	utils.BaseResponseMaker(dataResponse, w, "OTP verified successfully", http.StatusOK)
}

type SetPinRequest struct {
	NewPin string `json:"new_pin" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) SetPin(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok || userID == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	nextStep, ok := r.Context().Value(constant.ContextKey("next_step")).(string)
	if !ok || nextStep != "set_pin" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	_, _, deviceUUID, _, _, _ := utils.HeaderRequirement(r, nil)

	headerData := map[string]interface{}{
		"user_id":     userID,
		"device_uuid": deviceUUID,
		"user_realm":  "member",
	}

	var req SetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}

	if err := h.validateSetPinRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.Application.SetPin(r.Context(), userID, req.NewPin, deviceUUID, headerData["user_realm"].(string))
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "EXPIRED_OTP":
			status = http.StatusGone
		case "INVALID_OTP":
			status = http.StatusBadRequest
		case "INVALID_PIN":
			status = http.StatusBadRequest
		case "NOT_FOUND":
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	responseMap, _ := utils.StructToMap(response)
	utils.BaseResponseMaker(responseMap, w, "PIN set successfully", http.StatusOK)
}

type RegisterRequest struct {
	Phone      string `json:"phone" validate:"required"`
	DeviceUUID string `json:"device_uuid" validate:"required"`
	Platform   string `json:"platform" validate:"required,oneof=android ios web"`
	FullName   string `json:"full_name,omitempty"`
	Email      string `json:"email,omitempty"`
}

func (h UsersAdapter) Register(w http.ResponseWriter, r *http.Request) {
	platform, _, deviceUUID, _, _, _ := utils.HeaderRequirement(r, nil)

	// Parse and validate request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.validateRegisterRequest(req); err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	if deviceUUID != "" && deviceUUID != req.DeviceUUID {
		utils.SendErrorResponse(w, "DEVICE_UUID_MISMATCH", 0, nil)
		return
	}

	// Use platform from header if not provided in body
	if req.Platform == "" && platform != "" {
		req.Platform = platform
	}

	// Call application service for registration
	response, err := h.Application.Register(r.Context(), req.Phone, req.FullName, req.Email, req.DeviceUUID, req.Platform)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 400, nil)
		return
	}

	data, _ := constant.StructToMap(response)
	constant.BaseResponseMaker(data, w, "User Successfully registerd", http.StatusOK)
}

// validateRegisterRequest validates the registration request
func (h UsersAdapter) validateRegisterRequest(req RegisterRequest) error {
	if req.Phone == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: phone number is required")
	}
	if req.DeviceUUID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: device UUID is required")
	}
	if req.Platform == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: platform is required")
	}

	if err := validatePhone(req.Phone); err != nil {
		return err
	}

	validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
	if !validPlatforms[req.Platform] {
		return fmt.Errorf(errInvalidPlatform + ": platform must be android, ios, or web")
	}

	if err := validateDeviceUUID(req.DeviceUUID); err != nil {
		return err
	}

	return nil
}

type LoginRequest struct {
	Pin string `json:"pin" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) Login(w http.ResponseWriter, r *http.Request) {
	_, _, deviceUUID, _, _, _ := utils.HeaderRequirement(r, []string{})
	nextStep, ok := r.Context().Value(constant.ContextKey("next_step")).(string)

	if !ok || nextStep != "login" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.validateLoginRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	if deviceUUID == "" {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Device UUID mismatch", http.StatusBadRequest)
		return
	}
	phone_number := r.Context().Value(utils.ContextKey("phone_number")).(string)
	phone := utils.FormatPhoneNumber(phone_number)
	loginResult, err := h.Application.Login(r.Context(), phone, deviceUUID, req.Pin)

	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]interface{}{
		"access_token":    loginResult.Token,
		"user_id":         loginResult.UserID,
		"user_code":       loginResult.UserCode,
		"full_name":       loginResult.FullName,
		"phone_number":    loginResult.PhoneNumber,
		"kyc_level":       loginResult.KYCLevel,
		"is_verified":     loginResult.IsVerified,
		"login_time":      loginResult.LoginTime,
		"session_expires": loginResult.SessionExpires,
	}

	utils.BaseResponseMaker(response, w, "Login successful", http.StatusOK)
}

func (h UsersAdapter) validateSetPinRequest(req SetPinRequest) error {

	if req.NewPin == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: PIN is required")
	}

	if err := validatePin(req.NewPin); err != nil {
		return err
	}
	return nil
}

func (h UsersAdapter) validateLoginRequest(req LoginRequest) error {

	if req.Pin == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: PIN is required")
	}

	return nil
}

type ForgetPinSendOtpRequest struct {
	Phone string `json:"phone"`
}

func (h UsersAdapter) ForgetPinSendOtp(w http.ResponseWriter, r *http.Request) {
	var req ForgetPinSendOtpRequest
	_, _, deviceUUID, _, _, _ := utils.HeaderRequirement(r, nil)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.validateForgetPinSendOtpRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Application.ForgetPinSendOtp(r.Context(), req.Phone, deviceUUID)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "PIN_RESET_USER_NOT_FOUND":
			status = http.StatusNotFound
		case "ACCOUNT_DELETED":
			status = http.StatusGone
		case "ACCOUNT_BLOCKED":
			status = http.StatusForbidden
		case "PIN_RESET_DEVICE_MISMATCH":
			status = http.StatusBadRequest
		case "PIN_RESET_ALREADY_IN_PROGRESS":
			status = http.StatusConflict
		case "PIN_RESET_FAILED":
			status = http.StatusInternalServerError
		default:
			status = http.StatusInternalServerError
		}
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), status)
		return
	}

	response := map[string]interface{}{
		"phone_number":       result.PhoneNumber,
		"device_uuid":        result.DeviceUUID,
		"otp_sent":           result.OTPSent,
		"otp":                result.OTP,
		"temp_token":         result.Token,
		"otp_expiry_minutes": result.OTPExpiryMinutes,
		"reset_session_id":   result.ResetSessionID,
		"next_step":          result.NextStep,
	}

	utils.BaseResponseMaker(response, w, "OTP sent for PIN reset visit the nearest branch", http.StatusOK)
}

func (h UsersAdapter) validateForgetPinSendOtpRequest(req ForgetPinSendOtpRequest) error {
	if req.Phone == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: phone number is required")
	}

	if err := validatePhone(req.Phone); err != nil {
		return err
	}

	return nil
}

type VerifyForgetPinOtpRequest struct {
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) VerifyForgetPinOtp(w http.ResponseWriter, r *http.Request) {
	nextStep, ok := r.Context().Value(constant.ContextKey("next_step")).(string)
	if !ok || nextStep != "forget_pin_verify_otp" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	var req VerifyForgetPinOtpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(nil, w, "Invalid JSON payload"+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation (reuse logic from validateResetPinRequest, but without new_pin)
	if req.ResetSessionID == "" {
		utils.BaseResponseMaker(nil, w, "MISSING_REQUIRED_FIELDS: reset session ID is required", http.StatusBadRequest)
		return
	}
	if req.Phone == "" {
		utils.BaseResponseMaker(nil, w, "MISSING_REQUIRED_FIELDS: phone number is required", http.StatusBadRequest)
		return
	}
	if req.DeviceUUID == "" {
		utils.BaseResponseMaker(nil, w, "MISSING_REQUIRED_FIELDS: device UUID is required", http.StatusBadRequest)
		return
	}
	if req.OTP == "" {
		utils.BaseResponseMaker(nil, w, "MISSING_REQUIRED_FIELDS: OTP is required", http.StatusBadRequest)
		return
	}
	if req.Phone == "" || len(req.Phone) < minPhoneLength {
		utils.BaseResponseMaker(nil, w, "INVALID_PHONE_NUMBER: please provide a valid phone number", http.StatusBadRequest)
		return
	}

	formattedPhone := utils.FormatPhoneNumber(req.Phone)

	if formattedPhone == "" || len(formattedPhone) < minPhoneLength {
		utils.BaseResponseMaker(nil, w, "INVALID_PHONE_NUMBER: please provide a valid phone number", http.StatusBadRequest)
		return
	}

	if err := validateOTP(req.OTP); err != nil {
		utils.BaseResponseMaker(nil, w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.DeviceUUID) < 10 {
		utils.BaseResponseMaker(nil, w, "INVALID_DEVICE_UUID: device UUID appears to be invalid", http.StatusBadRequest)
		return
	}

	if err := validateDeviceUUID(req.DeviceUUID); err != nil {
		utils.BaseResponseMaker(nil, w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Application.VerifyForgetPinOtp(r.Context(), req.ResetSessionID, formattedPhone, req.DeviceUUID, req.OTP)
	if err != nil {
		status := map[string]int{
			"PIN_RESET_SESSION_NOT_FOUND": http.StatusNotFound,
			"PIN_RESET_SESSION_EXPIRED":   http.StatusGone,
			"PIN_RESET_SESSION_INVALID":   http.StatusBadRequest,
			"PIN_RESET_TOO_MANY_ATTEMPTS": http.StatusTooManyRequests,
			"PIN_RESET_OTP_INVALID":       http.StatusBadRequest,
			"PIN_RESET_FAILED":            http.StatusInternalServerError,
		}[err.Error()]
		if status == 0 {
			status = http.StatusInternalServerError
		}
		utils.BaseResponseMaker(nil, w, err.Error(), status)
		return
	}

	responseMap, _ := utils.StructToMap(resp)
	utils.BaseResponseMaker(responseMap, w, "OTP verified for PIN reset", http.StatusOK)
}

type ResetPinRequest struct {
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
	NewPin         string `json:"new_pin" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) ResetPin(w http.ResponseWriter, r *http.Request) {
	var req ResetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(nil, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	nextStep, ok := r.Context().Value(constant.ContextKey("next_step")).(string)
	if !ok || nextStep != "reset_pin" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	if err := h.validateResetPinRequest(req); err != nil {
		utils.BaseResponseMaker(nil, w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Application.ResetPin(r.Context(), req.ResetSessionID, req.Phone, req.DeviceUUID, req.OTP, req.NewPin)
	if err != nil {
		status := map[string]int{
			"PIN_RESET_SESSION_NOT_FOUND": http.StatusNotFound,
			"PIN_RESET_SESSION_EXPIRED":   http.StatusGone,
			"PIN_RESET_SESSION_INVALID":   http.StatusBadRequest,
			"PIN_RESET_TOO_MANY_ATTEMPTS": http.StatusTooManyRequests,
			"PIN_RESET_OTP_INVALID":       http.StatusBadRequest,
			"PIN_RESET_USER_NOT_FOUND":    http.StatusNotFound,
			"PIN_RESET_FAILED":            http.StatusInternalServerError,
		}[err.Error()]

		if status == 0 {
			status = http.StatusInternalServerError
		}

		utils.BaseResponseMaker(nil, w, err.Error(), status)
		return
	}

	response := map[string]interface{}{
		"user_id":           result.UserID,
		"user_code":         result.UserCode,
		"full_name":         result.FullName,
		"phone_number":      result.PhoneNumber,
		"pin_reset":         result.PinReset,
		"reset_time":        result.ResetTime,
		"access_restricted": result.AccessRestricted,
		"restrictions":      result.Restrictions,
	}

	utils.BaseResponseMaker(response, w, "PIN reset completed successfully", http.StatusOK)
}

// validateResetPinRequest validates the reset PIN request
func (h UsersAdapter) validateResetPinRequest(req ResetPinRequest) error {
	if req.ResetSessionID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: reset session ID is required")
	}
	if req.Phone == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: phone number is required")
	}
	if req.DeviceUUID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: device UUID is required")
	}
	if req.OTP == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: OTP is required")
	}
	if req.NewPin == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: new PIN is required")
	}

	if err := validatePhone(req.Phone); err != nil {
		return err
	}
	if err := validateOTP(req.OTP); err != nil {
		return err
	}
	if err := validatePin(req.NewPin); err != nil {
		return err
	}
	if err := validateDeviceUUID(req.DeviceUUID); err != nil {
		return err
	}

	return nil
}

type CompleteRegistrationRequest struct {
	RegistrationID string `json:"registration_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	Platform       string `json:"platform" validate:"required,oneof=android ios web"`
	FullName       string `json:"full_name" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) CompleteRegistration(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request body
	var req CompleteRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	action, ok := r.Context().Value(constant.ContextKey("action")).(string)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	// Validate required fields
	if err := h.validateCompleteRegistrationRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}
	phoneNumber := r.Context().Value("phone_number").(string)

	_, err := h.Application.VerifyOtp(r.Context(), req.RegistrationID, phoneNumber, req.OTP, req.DeviceUUID, req.Platform, "registration", action)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "EXPIRED_OTP":
			status = http.StatusGone
		case "INVALID_OTP":
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), status)
		return
	}

	// Complete the registration
	response, err := h.Application.CompleteRegistration(r.Context(), req.RegistrationID, req.Phone, req.DeviceUUID, req.Platform, req.FullName)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "REGISTRATION_NOT_FOUND":
			status = http.StatusNotFound
		case "REGISTRATION_EXPIRED":
			status = http.StatusGone
		case "INVALID_REGISTRATION_STATUS":
			status = http.StatusBadRequest
		case "USER_CREATION_FAILED":
			status = http.StatusInternalServerError
		default:
			status = http.StatusInternalServerError
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	// Send success response with token
	h.sendSuccessResponse(w, http.StatusOK, response)
}

// validateCompleteRegistrationRequest validates the complete registration request
func (h UsersAdapter) validateCompleteRegistrationRequest(req CompleteRegistrationRequest) error {
	if req.RegistrationID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: registration ID is required")
	}
	if req.Phone == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: phone number is required")
	}
	if req.DeviceUUID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: device UUID is required")
	}
	if req.Platform == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: platform is required")
	}
	if req.FullName == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: full name is required")
	}
	if req.OTP == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: OTP is required")
	}

	if err := validatePhone(req.Phone); err != nil {
		return err
	}
	validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
	if !validPlatforms[req.Platform] {
		return fmt.Errorf(errInvalidPlatform + ": platform must be android, ios, or web")
	}
	if err := validateDeviceUUID(req.DeviceUUID); err != nil {
		return err
	}

	return nil
}

type ResetPinWithTokenRequest struct {
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	NewPin         string `json:"new_pin" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) ResetPinWithToken(w http.ResponseWriter, r *http.Request) {
	var req ResetPinWithTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(nil, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if req.ResetSessionID == "" || req.Phone == "" || req.DeviceUUID == "" || req.NewPin == "" {
		utils.BaseResponseMaker(nil, w, "MISSING_REQUIRED_FIELDS: all fields are required", http.StatusBadRequest)
		return
	}

	if err := validatePin(req.NewPin); err != nil {
		utils.BaseResponseMaker(nil, w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Application.ResetPinWithToken(r.Context(), req.ResetSessionID, req.Phone, req.DeviceUUID, req.NewPin)
	if err != nil {
		utils.BaseResponseMaker(nil, w, err.Error(), http.StatusBadRequest)
		return
	}

	responseMap, _ := utils.StructToMap(resp)
	utils.BaseResponseMaker(responseMap, w, "PIN reset successfully", http.StatusOK)
}
