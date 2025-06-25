package users

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"strings"

	"cbe-super-app-member-users/internal/application/dto"
	user_inbound "cbe-super-app-member-users/internal/port/inbound/users"
	"cbe-super-app-member-users/pkgs/common"
	"cbe-super-app-member-users/pkgs/utils"
	constant "cbe-super-app-member-users/pkgs/utils"
)

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

func (h UsersAdapter) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := common.Response[interface{}]{
		ResponseWriter: w,
		Status:         status,
		Data:           data,
	}
	resp.SendJSON()
}

