package users

import (
	"encoding/json"
	"fmt"

	// "fmt"
	"net/http"
	"strings"

	"cbe-super-app-member-users/internal/application/dto"
	user_inbound "cbe-super-app-member-users/internal/port/inbound/users"

	// "cbe-super-app-member-users/pkgs/common"
	"cbe-super-app-member-users/pkgs/utils"
	constant "cbe-super-app-member-users/pkgs/utils"
)

func (h UsersAdapter) DeviceLookup(w http.ResponseWriter, r *http.Request) {
	// Extract required headers
	platform, appVersion, deviceUUID, sourceApp, installationDate, additionalHeaders := utils.HeaderRequirement(r, nil)

	// Validate required headers
	if deviceUUID == "" || platform == "" {
		utils.SendErrorResponse(w, "MISSING_REQUIRED_HEADERS", http.StatusBadRequest, nil)
		return
	}

	// Prepare header data for device lookup
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

	// Call application service for device lookup
	response, err := h.Application.DeviceLookup(r.Context(), headerData)
	if err != nil {
		utils.SendErrorResponse(w, "DEVICE_LOOKUP_FAILED", http.StatusInternalServerError, nil)
		return
	}

	// Send success response with token
	h.sendSuccessResponse(w, http.StatusOK, response)
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
		utils.SendErrorResponse(w, "NO_FILE", http.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	const MaxFileSize = 2 << 20
	if fileHeader.Size > MaxFileSize {
		utils.SendErrorResponse(w, "FILE_TOO_LARGE", http.StatusBadRequest, nil)
		return
	}

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

// func (h UsersAdapter) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
// 	resp := common.Response[interface{}]{
// 		ResponseWriter: w,
// 		Status:         status,
// 		Data:           data,
// 	}
// 	resp.SendJSON()
// }

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
	// Extract user ID from context
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok || userID == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	// Extract phone number from context
	phoneNumber := r.Context().Value("phone_number")

	// Extract device information from headers
	_, _, deviceUUID, sourceApp, _, _ := utils.HeaderRequirement(r, nil)

	// Parse and validate request body
	var req VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}

	// Validate required fields
	if req.Otp == "" {
		utils.SendErrorResponse(w, "MISSING_OTP", http.StatusBadRequest, nil)
		return
	}

	// Set default OTP purpose if not provided
	if req.OtpFor == "" {
		req.OtpFor = "pin_set"
	}

	// Call application service for OTP verification
	token, err := h.Application.VerifyOtp(r.Context(), userID, req.Otp, &deviceUUID, sourceApp, req.OtpFor)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "EXPIRED_OTP":
			status = http.StatusGone
		case "INVALID_OTP":
			status = http.StatusBadRequest
		case "INVALID_INPUT_PARAMETERS":
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		utils.BaseResponseMaker(map[string]interface{}{}, w, "OTP verification failed", status)
		return
	}

	// Prepare success response
	response := map[string]interface{}{
		"token":        token,
		"user_id":      userID,
		"phone_number": phoneNumber,
	}

	utils.BaseResponseMaker(response, w, "OTP verified successfully", http.StatusOK)
}

type SetPinRequest struct {
	NewPin     string  `json:"new_pin" validate:"required,min=6,max=6"`
	Otp        string  `json:"otp" validate:"required,min=6,max=6"`
	DeviceUUID *string `json:"device_uuid,omitempty"`
	UserRealm  string  `json:"user_realm,omitempty"`
	OtpFor     string  `json:"otp_for,omitempty"`
}

func (h UsersAdapter) SetPin(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value(constant.ContextKey("user_id")).(string)
	if !ok || userID == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	// Parse and validate request body
	var req SetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}

	// Validate required fields
	if req.NewPin == "" || req.Otp == "" {
		utils.SendErrorResponse(w, "MISSING_REQUIRED_FIELDS", http.StatusBadRequest, nil)
		return
	}

	// Set default values if not provided
	if req.OtpFor == "" {
		req.OtpFor = "pin_set"
	}
	if req.UserRealm == "" {
		req.UserRealm = "member"
	}

	// Call application service
	response, err := h.Application.SetPin(r.Context(), userID, req.NewPin, req.Otp, req.DeviceUUID, req.UserRealm, req.OtpFor)
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

	// Send success response with token
	h.sendSuccessResponse(w, http.StatusOK, response)
}

type RegisterRequest struct {
	Phone      string `json:"phone" validate:"required"`
	DeviceUUID string `json:"device_uuid" validate:"required"`
	Platform   string `json:"platform" validate:"required,oneof=android ios web"`
	FullName   string `json:"full_name,omitempty"`
	Email      string `json:"email,omitempty"`
}

func (h UsersAdapter) Register(w http.ResponseWriter, r *http.Request) {
	// Extract device information from headers for additional validation
	platform, _, deviceUUID, _, _, _ := utils.HeaderRequirement(r, nil)

	// Parse and validate request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := h.validateRegisterRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate device UUID consistency between header and body
	if deviceUUID != "" && deviceUUID != req.DeviceUUID {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Device UUID mismatch", http.StatusBadRequest)
		return
	}

	// Use platform from header if not provided in body
	if req.Platform == "" && platform != "" {
		req.Platform = platform
	}

	// Call application service for registration
	response, err := h.Application.Register(r.Context(), req.Phone, req.DeviceUUID, req.Platform)
	if err != nil {
		status := http.StatusBadRequest
		switch err.Error() {
		case "PHONE_ALREADY_EXISTS":
			status = http.StatusConflict
		case "DEVICE_ALREADY_REGISTERED":
			status = http.StatusConflict
		case "INVALID_PHONE_NUMBER":
			status = http.StatusBadRequest
		case "INVALID_PLATFORM":
			status = http.StatusBadRequest
		case "REGISTRATION_RATE_LIMITED":
			status = http.StatusTooManyRequests
		default:
			status = http.StatusInternalServerError
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}

	// Send success response with token
	h.sendSuccessResponse(w, http.StatusOK, response)
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

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(req.Phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER: please provide a valid phone number")
	}

	// Validate platform
	validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
	if !validPlatforms[req.Platform] {
		return fmt.Errorf("INVALID_PLATFORM: platform must be android, ios, or web")
	}

	// Validate device UUID format (basic validation)
	if len(req.DeviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID: device UUID appears to be invalid")
	}

	return nil
}

type LoginRequest struct {
	Phone      string `json:"phone" validate:"required"`
	DeviceUUID string `json:"device_uuid" validate:"required"`
	Pin        string `json:"pin" validate:"required,min=6,max=6"`
}

func (h UsersAdapter) Login(w http.ResponseWriter, r *http.Request) {
	// Extract device information from headers for additional validation
	platform, appVersion, deviceUUID, _, _, _ := utils.HeaderRequirement(r, nil)

	// Parse and validate request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := h.validateLoginRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate device UUID consistency between header and body
	if deviceUUID != "" && deviceUUID != req.DeviceUUID {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Device UUID mismatch", http.StatusBadRequest)
		return
	}

	// Call application service for login
	loginResult, err := h.Application.Login(r.Context(), req.Phone, req.DeviceUUID, req.Pin)
	if err != nil {
		status := http.StatusUnauthorized
		switch err.Error() {
		case "USER_NOT_FOUND":
			status = http.StatusNotFound
		case "INVALID_PIN":
			status = http.StatusUnauthorized
		case "ACCOUNT_BLOCKED":
			status = http.StatusForbidden
		case "ACCOUNT_DELETED":
			status = http.StatusGone
		case "DEVICE_NOT_LINKED":
			status = http.StatusForbidden
		case "TOO_MANY_LOGIN_ATTEMPTS":
			status = http.StatusTooManyRequests
		case "INVALID_PHONE_NUMBER":
			status = http.StatusBadRequest
		case "INVALID_DEVICE_UUID":
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), status)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"token":           loginResult.Token,
		"user_id":         loginResult.UserID,
		"user_code":       loginResult.UserCode,
		"full_name":       loginResult.FullName,
		"phone_number":    loginResult.PhoneNumber,
		"kyc_level":       loginResult.KYCLevel,
		"is_verified":     loginResult.IsVerified,
		"device_uuid":     req.DeviceUUID,
		"platform":        platform,
		"app_version":     appVersion,
		"login_time":      loginResult.LoginTime,
		"session_expires": loginResult.SessionExpires,
	}

	utils.BaseResponseMaker(response, w, "Login successful", http.StatusOK)
}

// validateLoginRequest validates the login request
func (h UsersAdapter) validateLoginRequest(req LoginRequest) error {
	if req.Phone == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: phone number is required")
	}
	if req.DeviceUUID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: device UUID is required")
	}
	if req.Pin == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: PIN is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(req.Phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER: please provide a valid phone number")
	}

	// Validate PIN format (6 digits)
	if len(req.Pin) != 6 {
		return fmt.Errorf("INVALID_PIN: PIN must be exactly 6 digits")
	}

	// Validate PIN contains only digits
	for _, char := range req.Pin {
		if char < '0' || char > '9' {
			return fmt.Errorf("INVALID_PIN: PIN must contain only digits")
		}
	}

	// Validate device UUID format (basic validation)
	if len(req.DeviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID: device UUID appears to be invalid")
	}

	return nil
}

type ForgetPinSendOtpRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
}

func (h UsersAdapter) ForgetPinSendOtp(w http.ResponseWriter, r *http.Request) {
	var req ForgetPinSendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := h.validateForgetPinSendOtpRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call application service
	result, err := h.Application.ForgetPinSendOtp(r.Context(), req.Phone, req.DeviceUUID)
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

	// Send success response
	response := map[string]interface{}{
		"phone_number":       result.PhoneNumber,
		"device_uuid":        result.DeviceUUID,
		"otp_sent":           result.OTPSent,
		"otp_expiry_minutes": result.OTPExpiryMinutes,
		"reset_session_id":   result.ResetSessionID,
		"next_step":          result.NextStep,
	}

	utils.BaseResponseMaker(response, w, "OTP sent for PIN reset", http.StatusOK)
}

// validateForgetPinSendOtpRequest validates the forget PIN send OTP request
func (h UsersAdapter) validateForgetPinSendOtpRequest(req ForgetPinSendOtpRequest) error {
	if req.Phone == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: phone number is required")
	}
	if req.DeviceUUID == "" {
		return fmt.Errorf("MISSING_REQUIRED_FIELDS: device UUID is required")
	}

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(req.Phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER: please provide a valid phone number")
	}

	// Validate device UUID format (basic validation)
	if len(req.DeviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID: device UUID appears to be invalid")
	}

	return nil
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
		"next_step":         result.NextStep,
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

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(req.Phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER: please provide a valid phone number")
	}

	// Validate OTP format (6 digits)
	if len(req.OTP) != 6 {
		return fmt.Errorf("INVALID_OTP: OTP must be exactly 6 digits")
	}

	// Validate OTP contains only digits
	for _, char := range req.OTP {
		if char < '0' || char > '9' {
			return fmt.Errorf("INVALID_OTP: OTP must contain only digits")
		}
	}

	// Validate new PIN format (6 digits)
	if len(req.NewPin) != 6 {
		return fmt.Errorf("INVALID_PIN: PIN must be exactly 6 digits")
	}

	// Validate new PIN contains only digits
	for _, char := range req.NewPin {
		if char < '0' || char > '9' {
			return fmt.Errorf("INVALID_PIN: PIN must contain only digits")
		}
	}

	// Validate device UUID format (basic validation)
	if len(req.DeviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID: device UUID appears to be invalid")
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

	// Validate required fields
	if err := h.validateCompleteRegistrationRequest(req); err != nil {
		utils.BaseResponseMaker(map[string]interface{}{}, w, err.Error(), http.StatusBadRequest)
		return
	}

	// First verify the OTP
	_, err := h.Application.VerifyOtp(r.Context(), req.RegistrationID, req.OTP, &req.DeviceUUID, req.Platform, "registration")
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

	// Validate phone number format
	formattedPhone := utils.FormatPhoneNumber(req.Phone)
	if formattedPhone == "" || len(formattedPhone) < 10 {
		return fmt.Errorf("INVALID_PHONE_NUMBER: please provide a valid phone number")
	}

	// Validate platform
	validPlatforms := map[string]bool{"android": true, "ios": true, "web": true}
	if !validPlatforms[req.Platform] {
		return fmt.Errorf("INVALID_PLATFORM: platform must be android, ios, or web")
	}

	// Validate device UUID format (basic validation)
	if len(req.DeviceUUID) < 10 {
		return fmt.Errorf("INVALID_DEVICE_UUID: device UUID appears to be invalid")
	}

	return nil
}
