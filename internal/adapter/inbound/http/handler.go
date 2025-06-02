package http

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
   app "cbe-super-app-member-users/internal/application/users"
   domain "cbe-super-app-member-users/internal/domain/users"
    // "cbe-super-app-member-users/internal/port/inbound"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HTTPHandler struct {
    appService app.ApplicationService
    logger     utils.Logger
}

func NewHTTPHandler(appService app.ApplicationService, logger utils.Logger) *HTTPHandler {
    return &HTTPHandler{
        appService: appService,
        logger:     logger,
    }
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) chi.Router {
    r.Get("/{id}/linked-accounts", h.FetchLinkedAccounts)
    return r
}


func (h *HTTPHandler) FetchLinkedAccounts(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if _, err := primitive.ObjectIDFromHex(id); err != nil {
        h.logger.Errorf("Invalid ObjectID: %v", err)
        resp := common.Response[struct {
            Code    string `json:"code"`
            Message string `json:"message"`
        }]{
            ResponseWriter: w,
            Status:         http.StatusBadRequest,
            Data: struct {
                Code    string `json:"code"`
                Message string `json:"message"`
            }{
                Code:    common.DefineError.General["INVALID_ID"].Code,
                Message: common.DefineError.General["INVALID_ID"].Message,
            },
        }
        resp.SendJSON()
        return
    }

    response, err := h.appService.FetchLinkedAccounts(r.Context(), id)
    if err != nil {
        h.logger.Errorf("Failed to fetch linked accounts for user ID %s: %v", id, err)
        var resp common.Response[struct {
            Code    string `json:"code"`
            Message string `json:"message"`
        }]
        if serviceErr, ok := err.(*domain.ServiceError); ok && serviceErr.Code == "NOT_FOUND" {
            resp = common.Response[struct {
                Code    string `json:"code"`
                Message string `json:"message"`
            }]{
                ResponseWriter: w,
                Status:         http.StatusNotFound,
                Data: struct {
                    Code    string `json:"code"`
                    Message string `json:"message"`
                }{
                    Code:    common.DefineError.General["NOT_FOUND"].Code,
                    Message: serviceErr.Message,
                },
            }
        } else {
            resp = common.Response[struct {
                Code    string `json:"code"`
                Message string `json:"message"`
            }]{
                ResponseWriter: w,
                Status:         http.StatusInternalServerError,
                Data: struct {
                    Code    string `json:"code"`
                    Message string `json:"message"`
                }{
                    Code:    "INTERNAL_SERVER_ERROR",
                    Message: "An unexpected error occurred.",
                },
            }
        }
        resp.SendJSON()
        return
    }

    resp := common.Response[*domain.LinkedAccountResponse]{
        ResponseWriter: w,
        Status:         http.StatusOK,
        Data:           response,
    }
    resp.SendJSON()
}