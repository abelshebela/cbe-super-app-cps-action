package accountvalidation_inbound

import (
	
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


func InitAccountValidationHandlerMaker(router chi.Router, handler inbound.Inbound, middleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/fetch_account_validation",
			Handler: handler.FetchAccountValidation,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/update_account_validation_maker",
			Handler: handler.UpdateAccountValidationMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/update_account_validation_checker",
			Handler: handler.UpdateAccountValidationChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"checker"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
