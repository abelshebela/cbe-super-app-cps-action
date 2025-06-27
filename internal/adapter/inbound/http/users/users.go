package users

import (
	"encoding/json"
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
	// userID, ok := r.Context().Value("user_id").(string)
	var platform, appVersion, devideUUID, _, installationDate, _ = utils.HeaderRequirement(r, nil)

	var headerData = make(map[string]interface{})
	headerData["platform"] = platform
	headerData["app_version"] = appVersion
	headerData["device_uuid"] = devideUUID
	headerData["installation_date"] = installationDate

	HQ, hqerr := h.Application.CheckDevice(r.Context(), headerData)
	if hqerr != nil {
		utils.SendErrorResponse(w, hqerr.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.BaseResponseMaker(HQ, w, "Success", http.StatusOK)

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
	UserID     string  `json:"user_id"`
	Otp        string  `json:"otp"`
	DeviceUUID *string `json:"device_uuid,omitempty"`
	UserRealm  string  `json:"user_realm,omitempty"`
	OtpFor     string  `json:"otp_for,omitempty"`
}

func (h UsersAdapter) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	var req VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	if userID != "" {
		req.UserID = userID
	}
	err := h.Application.VerifyOtp(r.Context(), req.UserID, req.Otp, req.DeviceUUID, req.UserRealm, req.OtpFor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "EXPIRED_OTP" {
			status = http.StatusGone
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}
	h.sendSuccessResponse(w, http.StatusOK, "OTP verified successfully")
}

type SetPinRequest struct {
	UserID     string  `json:"user_id"`
	NewPin     string  `json:"new_pin"`
	Otp        string  `json:"otp"`
	DeviceUUID *string `json:"device_uuid,omitempty"`
	UserRealm  string  `json:"user_realm,omitempty"`
	OtpFor     string  `json:"otp_for,omitempty"`
}

func (h UsersAdapter) SetPin(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	var req SetPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	if userID != "" {
		req.UserID = userID
	}
	err := h.Application.SetPin(r.Context(), req.UserID, req.NewPin, req.Otp, req.DeviceUUID, req.UserRealm, req.OtpFor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "EXPIRED_OTP" {
			status = http.StatusGone
		}
		utils.SendErrorResponse(w, err.Error(), status, nil)
		return
	}
	h.sendSuccessResponse(w, http.StatusOK, "PIN set successfully")
}

type RegisterRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
	Platform   string `json:"platform"`
}

func (h UsersAdapter) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	err := h.Application.Register(r.Context(), req.Phone, req.DeviceUUID, req.Platform)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	h.sendSuccessResponse(w, http.StatusOK, "Registration started, OTP sent")
}

type LoginRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
	Pin        string `json:"pin"`
}

func (h UsersAdapter) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	token, err := h.Application.Login(r.Context(), req.Phone, req.DeviceUUID, req.Pin)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}
	h.sendSuccessResponse(w, http.StatusOK, map[string]string{"token": token})
}

type ForgetPinSendOtpRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
}

func (h UsersAdapter) ForgetPinSendOtp(w http.ResponseWriter, r *http.Request) {
	var req ForgetPinSendOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	err := h.Application.ForgetPinSendOtp(r.Context(), req.Phone, req.DeviceUUID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	h.sendSuccessResponse(w, http.StatusOK, "OTP sent for PIN reset")
}
