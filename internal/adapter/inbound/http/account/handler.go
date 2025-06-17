package account

import (
	"net/http"

	accountApp "cbe-super-app-member-users/internal/application/account"
	domainAccount "cbe-super-app-member-users/internal/domain/account"
	local "cbe-super-app-member-users/internal/shared"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HTTPHandler struct {
	appHandler *accountApp.ApplicationHandler
	logger     utils.Logger
}

func NewHTTPHandler(appHandler *accountApp.ApplicationHandler, logger utils.Logger) *HTTPHandler {
	return &HTTPHandler{
		appHandler: appHandler,
		logger:     logger,
	}
}

func (h *HTTPHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req := accountApp.CreateAccountRequest{UserID: id}
	
	response, err := h.appHandler.HandleCreateAccount(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"status":         "success",
		"message":        "Account created successfully",
		"customer_number": response.CustomerNumber,
		"account_number": response.AccountNumber,
	})
}

func (h *HTTPHandler) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := common.Response[interface{}]{
		ResponseWriter: w,
		Status:         status,
		Data:           data,
	}
	resp.SendJSON()
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(domainAccount.ServiceError); ok {
		status := http.StatusInternalServerError
		switch serviceErr.Code {
		case common.DefineError.User["USER_KYC_LEVEL_ZERO"].Code:
			status = http.StatusBadRequest
		case common.DefineError.User["USER_PHONE_EXISTS"].Code:
			status = http.StatusConflict
		case common.DefineError.General["NOT_FOUND"].Code:
			status = http.StatusNotFound
		case common.DefineError.General["INVALID_ID"].Code:
			status = http.StatusBadRequest
		case local.DefineError.Account["PHONE_LOOKUP_FAILED"].Code:
			status = http.StatusInternalServerError
		case local.DefineError.Account["API_REQUEST_FAILED"].Code:
			status = http.StatusInternalServerError
		case common.DefineError.General["GENERAL_SIF_GENERATION_FAILED"].Code:
			status = http.StatusInternalServerError
		case common.DefineError.General["GENERAL_DB_UPDATE_FAILED"].Code:
			status = http.StatusInternalServerError
		case common.DefineError.General["GENERAL_DB_INSERT_FAILED"].Code:
			status = http.StatusInternalServerError
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