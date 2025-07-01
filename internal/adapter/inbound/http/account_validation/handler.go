package accountvalidation_inbound

import (
	"encoding/json"
	"net/http"

	route "cbe-super-app-cps-action/internal/adapter/inbound/http"
	accountvalidation_app "cbe-super-app-cps-action/internal/application/account_validation"
	"cbe-super-app-cps-action/internal/application/middleware"
	inbound "cbe-super-app-cps-action/internal/port/inbound/account_validation"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HttpStore struct {
	Application accountvalidation_app.ApplicationAbstracts
	logger      utils.Logger
}

func NewHttpAccountValidation(app accountvalidation_app.ApplicationAbstracts, logger utils.Logger) inbound.Inbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *HttpStore) handleError(w http.ResponseWriter, err error) {
	h.logger.Errorf("Error occurred: %v", err)

	statusCode := http.StatusInternalServerError
	message := "Internal server error"

	if err.Error() == "validation rule ID cannot be empty" ||
		err.Error() == "maker ID cannot be empty" ||
		err.Error() == "checker ID cannot be empty" ||
		err.Error() == "action ID cannot be empty" {
		statusCode = http.StatusBadRequest
		message = err.Error()
	} else if err.Error() == "validation rule not found" {
		statusCode = http.StatusNotFound
		message = err.Error()
	} else if err.Error() == "validation rule already has a pending action" {
		statusCode = http.StatusConflict
		message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    statusCode,
		Message: message,
	})
}

func InitAccountValidationHandlerMaker(router chi.Router, handler inbound.Inbound, middleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_account_validation",
			Handler: handler.FetchAccountValidation,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_account_validation_maker",
			Handler: handler.UpdateAccountValidationMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_account_validation_checker",
			Handler: handler.UpdateAccountValidationChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
