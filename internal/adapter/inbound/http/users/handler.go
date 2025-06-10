package users

import (
	"encoding/json"
	"net/http"

	appUsers "cbe-super-app-member-users/internal/application/users"
	domainUsers "cbe-super-app-member-users/internal/domain/users"
	local "cbe-super-app-member-users/internal/shared"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HTTPHandler struct {
	appHandler *appUsers.ApplicationHandler
	logger     utils.Logger
}

func NewHTTPHandler(appHandler *appUsers.ApplicationHandler, logger utils.Logger) *HTTPHandler {
	return &HTTPHandler{
		appHandler: appHandler,
		logger:     logger,
	}
}

func (h *HTTPHandler) FetchLinkedAccounts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := appUsers.FetchLinkedAccountsRequest{UserID: id}
	
	response, err := h.appHandler.HandleFetchLinkedAccounts(r.Context(), req)
	if err != nil {
		h.logger.Errorf("Failed to fetch linked accounts for user ID %s: %s", id, err.Error())
		h.handleError(w, err)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) GenerateEmailOTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req appUsers.GenerateOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	req.UserID = id

	response, err := h.appHandler.HandleGenerateEmailOTP(r.Context(), req)
	if err != nil {
		h.handleOTPError(w, err)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) VerifyEmailOTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req appUsers.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	req.UserID = id

	response, err := h.appHandler.HandleVerifyEmailOTP(r.Context(), req)
	if err != nil {
		h.handleOTPError(w, err)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(domainUsers.ServiceError); ok {
		status := http.StatusInternalServerError
		switch serviceErr.Code {
		case common.DefineError.General["NOT_FOUND"].Code:
			status = http.StatusNotFound
		case common.DefineError.General["INVALID_ID"].Code:
			status = http.StatusBadRequest
		}
		h.sendErrorResponse(w, status, serviceErr.Code, serviceErr.Message)
	} else {
		h.sendErrorResponse(w, http.StatusInternalServerError,
			common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code,
			common.DefineError.General["UNHANDLED_SERVER_ERROR"].Message)
	}
}

func (h *HTTPHandler) handleOTPError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(domainUsers.ServiceError); ok {
		status := http.StatusInternalServerError
		switch serviceErr.Code {
		case local.DefineError.OTP["INVALID_OTP"].Code, local.DefineError.OTP["EXPIRED_OTP"].Code, local.DefineError.OTP["EMAIL_IN_USE"].Code:
			status = http.StatusBadRequest
		case common.DefineError.General["INVALID_ID"].Code:
			status = http.StatusBadRequest
		case common.DefineError.General["NOT_FOUND"].Code:
			status = http.StatusNotFound
		}
		h.sendErrorResponse(w, status, serviceErr.Code, serviceErr.Message)
	} else {
		h.sendErrorResponse(w, http.StatusInternalServerError,
			common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code,
			common.DefineError.General["UNHANDLED_SERVER_ERROR"].Message)
	}
}

func (h *HTTPHandler) sendErrorResponse(w http.ResponseWriter, status int, code, message string) {
	resp := common.Response[struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}]{
		ResponseWriter: w,
		Status:         status,
		Data: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}
	resp.SendJSON()
}

func (h *HTTPHandler) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := common.Response[interface{}]{
		ResponseWriter: w,
		Status:         status,
		Data:           data,
	}
	resp.SendJSON()
}