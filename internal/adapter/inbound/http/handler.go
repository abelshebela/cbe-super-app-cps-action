package http

import (
	"encoding/json"
	

	app "cbe-super-app-member-users/internal/application/users"
	domain "cbe-super-app-member-users/internal/domain/users"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTPHandler struct {
	appHandler app.ApplicationHandler
	logger     utils.Logger
}

func NewHTTPHandler(appHandler app.ApplicationHandler, logger utils.Logger) *HTTPHandler {
	return &HTTPHandler{
		appHandler: appHandler,
		logger:     logger,
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

func (h *HTTPHandler) FetchLinkedAccounts(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if _, err := primitive.ObjectIDFromHex(id); err != nil {
        h.logger.Errorf("Invalid ObjectID: %v", err)
        h.sendErrorResponse(w, http.StatusBadRequest, common.DefineError.General["INVALID_ID"].Code, common.DefineError.General["INVALID_ID"].Message)
        return
    }

    req := app.FetchLinkedAccountsRequest{UserID: id}
    response, err := h.appHandler.FetchLinkedAccounts(r.Context(), req)
    if err != nil {
        h.logger.Errorf("Failed to fetch linked accounts for user ID %s: %v", id, err)
        if serviceErr, ok := err.(*domain.ServiceError); ok {
            h.sendErrorResponse(w, http.StatusNotFound, serviceErr.Code, serviceErr.Message)
        } else {
            h.sendErrorResponse(w, http.StatusInternalServerError, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Code, common.DefineError.General["UNHANDLED_SERVER_ERROR"].Message)
        }
        return
    }

    resp := common.Response[*app.LinkedAccountResponseDto]{
        ResponseWriter: w,
        Status:         http.StatusOK,
        Data:           response,
    }
    resp.SendJSON()
}




func (h *HTTPHandler) GenerateEmailOTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req app.GenerateOTPRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	req.UserID = id
	
	response, err := h.appHandler.GenerateEmailOTP(r.Context(), req)
	
	if err != nil {
		h.handleOTPError(w, err)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) VerifyEmailOTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req app.VerifyOTPRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	req.UserID = id

	response, err := h.appHandler.VerifyEmailOTP(r.Context(), req)
	if err != nil {
		h.handleOTPError(w, err)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) handleOTPError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(*domain.ServiceError); ok {
		switch serviceErr.Code {
		case "INVALID_OTP", "EXPIRED_OTP", "EMAIL_IN_USE":
			h.sendErrorResponse(w, http.StatusBadRequest, serviceErr.Code, serviceErr.Message)
		case "INVALID_ID":
			h.sendErrorResponse(w, http.StatusBadRequest, serviceErr.Code, serviceErr.Message)
		case "NOT_FOUND":
			h.sendErrorResponse(w, http.StatusNotFound, serviceErr.Code, serviceErr.Message)
		default:
			h.sendErrorResponse(w, http.StatusInternalServerError, serviceErr.Code, serviceErr.Message)
		}
	} else {
		h.sendErrorResponse(w, http.StatusInternalServerError, 
			"INTERNAL_ERROR", "An unexpected error occurred")
	}
}

func (h *HTTPHandler) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := common.Response[interface{}]{
		ResponseWriter: w,
		Status:         status,
		Data:           data,
	}
	resp.SendJSON()
}