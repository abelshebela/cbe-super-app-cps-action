package users

import (
	constants "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"

	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/response"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest/http/users/core"
	service "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	common "github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	pinLength = 6
)

type user struct {
	userService service.UserService
	logger      utils.Logger
}

func Init(userService service.UserService, logger utils.Logger) rest.Users {
	return &user{
		userService: userService,
		logger:      logger,
	}
}

func (u *user) Healthcheck(w http.ResponseWriter, r *http.Request) {
	response.SendSuccessResponse(w, http.StatusOK, "🚀 Member auth service is running", nil, nil)
}

func (u *user) Register(w http.ResponseWriter, r *http.Request) {
	header := core.ExtractHeader(r)
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Errorf("failed to decode register request: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}
	req.DeviceUUID = header.DeviceUUID
	req.Platform = header.Platform
	req.AppVersion = header.AppVersion
	req.SourceApp = header.SourceApp
	req.APPInstallationDate = header.ApplicationInstallationDate

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to register request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	// Call application service for registration
	res, err := u.userService.Register(r.Context(), req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "User Successfully registerd", res, nil)
}

func (u *user) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	header := core.ExtractHeader(r)
	userInfo, err := common.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	nextStep, err := common.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if nextStep != "verify_otp" {
		response.SendErrorResponse(w, errors.ErrUnauthorized)
		return
	}

	var req dto.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}
	req.UserID = userInfo.UserID
	req.PhoneNumber = userInfo.PhoneNumber
	req.DeviceUUID = header.DeviceUUID
	req.UserRealm = string(constants.MEMBER_REALM)
	req.Action = userInfo.UserID

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to verify otp request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	data, err := u.userService.VerifyOtp(r.Context(), req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}
	response.SendSuccessResponse(w, http.StatusOK, "OTP verified successfully", data, nil)
}

func (u *user) DeviceLookup(w http.ResponseWriter, r *http.Request) {
	header := core.ExtractHeader(r)

	if err := header.Validate(); err != nil {
		u.logger.Errorf("invalid device look up request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	res, err := u.userService.DeviceLookup(r.Context(), header)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK,
		"Successful Device Lookup", res, nil)
}

func (u *user) PreLogin(w http.ResponseWriter, r *http.Request) {
	header := core.ExtractHeader(r)

	if err := header.Validate(); err != nil {
		u.logger.Errorf("invalid input provided for pre login reqquest: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	var req dto.PhoneLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Errorf("failed to decode pre login request: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided for pre login request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	res, err := u.userService.PreLogin(r.Context(), req.Phone)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Successful prelogin", res, nil)
}

// func (u *user) ForgetPin(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	res, err := u.userService.FetchLinkedAccounts(r.Context(), userInfo.UserID)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "successful forget pin", res, nil)
// }

// func (u *user) ChangePin(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	var req dto.ChangePinRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}
// 	req.UserID = userInfo.UserID

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid input provided to change pin: %v", err)
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	if err := u.userService.ChangePin(r.Context(), req); err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "PIN changed successfully", nil, nil)
// }

// func (u *user) CheckPin(w http.ResponseWriter, r *http.Request) {
// 	var req dto.PinStrengthRequest

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		u.logger.Errorf("failed to decode pin strength request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	pin := strconv.Itoa(req.NewPin)
// 	err := validatePin(pin)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "PIN set Successfully", nil, nil)
// }

// func validatePin(pin string) error {
// 	if len(pin) != pinLength {
// 		return errors.ErrPINLength
// 	}
// 	for _, c := range pin {
// 		if c < '0' || c > '9' {
// 			return errors.ErrPINOnlyDigit
// 		}
// 	}

// 	// Check for weak patterns
// 	if isWeakPin(pin) {
// 		return errors.ErrWeakPIN
// 	}
// 	return nil
// }

// func isWeakPin(pin string) bool {
// 	// Check for repeated digits
// 	if strings.Count(pin, string(pin[0])) == len(pin) {
// 		return true
// 	}
// 	// Check for sequential patterns
// 	isAscending := true
// 	isDescending := true
// 	for i := 1; i < len(pin); i++ {
// 		if pin[i] != pin[i-1]+1 {
// 			isAscending = false
// 		}
// 		if pin[i] != pin[i-1]-1 {
// 			isDescending = false
// 		}
// 	}
// 	return isAscending || isDescending
// }

// func (u *user) CompleteRegistration(w http.ResponseWriter, r *http.Request) {
// 	var req dto.CompleteRegistrationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		u.logger.Errorf("failed to decode complete reqgistration request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid input provided to complete registration: %v", err)
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	_, err = u.userService.VerifyOtp(r.Context(), dto.VerifyOTPRequest{
// 		UserID:      req.RegistrationID,
// 		PhoneNumber: userInfo.PhoneNumber,
// 		OTP:         req.OTP,
// 		DeviceUUID:  req.DeviceUUID,
// 		UserRealm:   req.Platform,
// 		OtpFor:      "REGISTRATION",
// 		Action:      userInfo.UserID,
// 	})
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	// Complete the registration
// 	res, err := u.userService.CompleteRegistration(r.Context(), req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	// Send success response with token
// 	response.SendSuccessResponse(w, http.StatusOK, "registration complete", res, nil)
// }

// func (u *user) FetchLinkedAccounts(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	res, err := u.userService.FetchLinkedAccounts(r.Context(), userInfo.UserID)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "linked accounts succesfully fetched", res, nil)
// }

func (u *user) ForgetPinSendOtp(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgetPinSendOtpRequest
	header := core.ExtractHeader(r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Errorf("failed to decode forget pin request: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input for forget pin provided: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	result, err := u.userService.ForgetPinSendOtp(r.Context(), req.Phone, header.DeviceUUID)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "OTP sent for PIN reset", result, nil)
}

// func (u *user) GenerateEmailOTP(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	var req dto.OTPRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}
// 	req.UserID = userInfo.UserID

// 	res, err := u.userService.GenerateEmailOTP(r.Context(), req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "OTP generated and sent for email", res, nil)
// }

// func (u *user) Login(w http.ResponseWriter, r *http.Request) {
// 	nextStep, err := core.ExtractNextStep(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	if nextStep != "login" {
// 		response.SendErrorResponse(w, errors.ErrUnauthorized)
// 		return
// 	}

// 	header := core.ExtractHeader(r)

// 	if header.DeviceUUID == "" {
// 		u.logger.Errorf("device UUID is required")
// 		response.SendErrorResponse(w, errors.ErrBadRequest)
// 		return
// 	}

// 	var req dto.LoginRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		u.logger.Errorf("failed to decode login request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid input provided to login request: %v", err)
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	loginResult, err := u.userService.Login(r.Context(), req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "Login successful", loginResult, nil)
// }

// func (u *user) ResetPin(w http.ResponseWriter, r *http.Request) {
// 	nextStep, err := core.ExtractNextStep(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	if nextStep != "reset_pin" {
// 		response.SendErrorResponse(w, errors.ErrUnauthorized)
// 		return
// 	}

// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	header := core.ExtractHeader(r)

// 	var req dto.ResetPinRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.SendErrorResponse(w, errors.ErrBadRequest)
// 		return
// 	}
// 	req.Phone = userInfo.PhoneNumber
// 	req.DeviceUUID = header.DeviceUUID

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid input provided to reset pin request: %v", err)
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	res, err := u.userService.ResetPin(r.Context(), req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}
// 	response.SendSuccessResponse(w, http.StatusOK, "PIN reset successful", res, nil)
// }

// func (u *user) ResetPinWithToken(w http.ResponseWriter, r *http.Request) {
// 	header := core.ExtractHeader(r)
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}
// 	var req dto.ResetPinWithTokenRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		u.logger.Errorf("failed to decode reset pin with token request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	req.DeviceUUID = header.DeviceUUID
// 	req.Phone = userInfo.PhoneNumber

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid input provided to reset pin request: %v", err)
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	resp, err := u.userService.ResetPinWithToken(r.Context(), req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "PIN reset successful", resp, nil)
// }

func (u *user) SetPin(w http.ResponseWriter, r *http.Request) {
	nextStep, err := common.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}
	if nextStep != "set_pin" {
		response.SendErrorResponse(w, errors.ErrUnauthorized)
		return
	}

	header := core.ExtractHeader(r)
	user, err := common.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	var req dto.SetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Errorf("failed to decode set pin request: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}
	req.UserID = user.UserID
	req.DeviceUUID = header.DeviceUUID
	req.Realm = constants.MEMBER_REALM

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to set pin request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	res, err := u.userService.SetPin(r.Context(), req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "PIN set Successfully", res, nil)
}

// func (u *user) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	header := core.ExtractHeader(r)
// 	if header.DeviceUUID == "" {
// 		response.SendErrorResponse(w, errors.ErrDeviceIDRequired)
// 		return
// 	}

// 	if header.DeviceUUID == "" {
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	if err := u.userService.UnlinkDevice(r.Context(),
// 		userInfo.UserID, header.DeviceUUID); err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "Device Unlinked Successfully", nil, nil)
// }

// func (u *user) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
// 	var req dto.UpdateProfilePicture
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	if err := r.ParseMultipartForm(2 << 20); err != nil {
// 		u.logger.Errorf("failed to parse form data: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	file, fileHeader, err := r.FormFile("profile_picture")
// 	if err != nil {
// 		u.logger.Errorf("failed to get profile picture: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}
// 	defer file.Close()

// 	req.ProfilePicture = fileHeader

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	imageURL, err := u.userService.UpdateProfilePicture(r.Context(), userInfo.UserID, req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "image uploaded", imageURL, nil)
// }

// func (u *user) UpdateProfileTheme(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	var req dto.SetProfileThemeRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid request for update profile theme: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}

// 	_, err = u.userService.UpdateProfileTheme(r.Context(), userInfo.UserID, req.ThemeType)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "Profile theme set successfuly", nil, nil)
// }

// func (u *user) VerifyEmailOTP(w http.ResponseWriter, r *http.Request) {
// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	var req dto.OTPVerification
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		u.logger.Errorf("failed to decode otp verification request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrInvalidData)
// 		return
// 	}
// 	req.UserID = userInfo.UserID

// 	if err := u.userService.VerifyEmailOTP(r.Context(), req); err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "OTP verified successfully", nil, nil)
// }

// func (u *user) VerifyForgetPinOtp(w http.ResponseWriter, r *http.Request) {
// 	header := core.ExtractHeader(r)
// 	nextStep, err := core.ExtractNextStep(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	if nextStep != "forget_pin_verify_otp" {
// 		response.SendErrorResponse(w, errors.ErrUnauthorized)
// 		return
// 	}

// 	userInfo, err := core.ExtractUserInfo(r.Context(), u.logger)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	var req dto.VerifyForgetPinOtpRequest
// 	req.Phone = userInfo.PhoneNumber
// 	req.DeviceUUID = header.DeviceUUID

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		u.logger.Errorf("failed to decode verify foget pin request: %v", err)
// 		response.SendErrorResponse(w, errors.ErrBadRequest)
// 		return
// 	}

// 	if err := req.Validate(); err != nil {
// 		u.logger.Errorf("invalid input provided to verify forget pin request: %v", err)
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	res, err := u.userService.VerifyForgetPinOtp(r.Context(), req)
// 	if err != nil {
// 		response.SendErrorResponse(w, err)
// 		return
// 	}

// 	response.SendSuccessResponse(w, http.StatusOK, "OTP verified for PIN reset", res, nil)
// }
