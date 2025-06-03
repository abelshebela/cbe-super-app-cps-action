package http

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	app "cbe-super-app-member-users/internal/application/users"
	domain "cbe-super-app-member-users/internal/domain/users"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func (h *HTTPHandler) RegisterRoutes(r chi.Router) chi.Router {
	r.Get("/{id}/linked-accounts", h.FetchLinkedAccounts)
	return r
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