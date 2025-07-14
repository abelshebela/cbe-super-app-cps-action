package accountvalidation_inbound

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	accountvalidation_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_validation"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
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
	router.Route("/account_validation", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.FetchAccountValidation,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/update/{id}",
				Handler: handler.UpdateAccountValidationMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve_reject",
				Handler: handler.UpdateAccountValidationChecker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Checker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)

	})
}
