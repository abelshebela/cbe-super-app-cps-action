package users

import (
	constants "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	local_util "github.com/CBE-Super-App/cbe-super-app-member-auth/pkgs/utils"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/response"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest/http/users/core"
	service "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"encoding/json"
	"net/http"
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
	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)

	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	nextStep, err := local_util.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if nextStep != constants.VerifyOtp {
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
	req.Action = userInfo.Action

	fmt.Println("999999999///999999999")

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

func (u *user) ForgetPin(w http.ResponseWriter, r *http.Request) {
	header := core.ExtractHeader(r)

	if err := header.Validate(); err != nil {
		u.logger.Errorf("invalid input provided for pre login reqquest: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	res, err := u.userService.ForgetPinSendOtp(r.Context(), userInfo.PhoneNumber, header.DeviceUUID)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "successful forget pin", res, nil)
}

func (u *user) ChangePin(w http.ResponseWriter, r *http.Request) {
	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)

	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	var req dto.ChangePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}
	req.UserID = userInfo.UserID

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to change pin: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	if err := u.userService.ChangePin(r.Context(), req); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "PIN changed successfully", nil, nil)
}

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

func (u *user) Login(w http.ResponseWriter, r *http.Request) {
	nextStep, err := local_util.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if nextStep != constants.Login {
		response.SendErrorResponse(w, errors.ErrUnauthorized)
		return
	}

	header := core.ExtractHeader(r)

	if header.DeviceUUID == "" {
		u.logger.Errorf("device UUID is required")
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Errorf("failed to decode login request: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to login request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	loginResult, err := u.userService.Login(r.Context(), req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Login successful", loginResult, nil)
}

func (u *user) ResetPin(w http.ResponseWriter, r *http.Request) {
	nextStep, err := local_util.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if nextStep != constants.ResetPin {
		response.SendErrorResponse(w, errors.ErrUnauthorized)
		return
	}

	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	header := core.ExtractHeader(r)

	var req dto.ResetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}
	req.Phone = userInfo.PhoneNumber
	req.DeviceUUID = header.DeviceUUID

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to reset pin request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	res, err := u.userService.ResetPin(r.Context(), req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}
	response.SendSuccessResponse(w, http.StatusOK, "PIN reset successful", res, nil)
}
func (u *user) SetPin(w http.ResponseWriter, r *http.Request) {
	nextStep, err := local_util.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}
	if nextStep != constants.SetPin {
		response.SendErrorResponse(w, errors.ErrUnauthorized)
		return
	}

	header := core.ExtractHeader(r)
	user, err := local_util.ExtractUserInfo(r.Context(), u.logger)
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

func (u *user) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateProfilePicture
	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if err := r.ParseMultipartForm(2 << 20); err != nil {
		u.logger.Errorf("failed to parse form data: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		u.logger.Errorf("failed to get profile picture: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}
	defer file.Close()

	req.ProfilePicture = fileHeader

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid request: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	err = u.userService.UpdateProfilePicture(r.Context(), userInfo.UserID, req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, constants.ImageUploadSuccess, constants.Empty, nil)
}

func (u *user) UpdateProfileTheme(w http.ResponseWriter, r *http.Request) {
	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	var req dto.SetProfileThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid request for update profile theme: %v", err)
		response.SendErrorResponse(w, errors.ErrInvalidData)
		return
	}

	err = u.userService.UpdateProfileTheme(r.Context(), userInfo.UserID, req.ThemeType)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, constants.ProfileUploadSucess, nil, nil)
}

func (u *user) VerifyForgetPinOtp(w http.ResponseWriter, r *http.Request) {
	header := core.ExtractHeader(r)
	nextStep, err := local_util.ExtractNextStep(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if nextStep != constants.ForgetPinVerifyOtp {
		response.SendErrorResponse(w, errors.ErrUnauthorized)
		return
	}

	userInfo, err := local_util.ExtractUserInfo(r.Context(), u.logger)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	var req dto.VerifyForgetPinOtpRequest
	req.Phone = userInfo.PhoneNumber
	req.DeviceUUID = header.DeviceUUID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logger.Errorf("failed to decode verify foget pin request: %v", err)
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		u.logger.Errorf("invalid input provided to verify forget pin request: %v", err)
		response.SendErrorResponse(w, err)
		return
	}

	res, err := u.userService.VerifyForgetPinOtp(r.Context(), req)
	if err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, constants.PINResetSuccess, res, nil)
}
