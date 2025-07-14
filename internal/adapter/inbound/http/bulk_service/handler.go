package bulkservices_inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	bulkservices_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bulk_services"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bulk_services"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type HttpStore struct {
	Application bulkservices_application.ApplicationAbstracts
	Logger      utils.Logger
}

func NewHttpBulkService(app bulkservices_application.ApplicationAbstracts, logger utils.Logger) inbound.Inbound {
	return &HttpStore{
		Application: app,
		Logger:      logger,
	}
}

func InitServiceHandlerMaker(router chi.Router, handler inbound.Inbound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/bulk-services", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/actions/enable-disable/request",
				Handler: handler.EnableDisableServicesMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/services",
				Handler: handler.FetchServices,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{
						role.Maker, role.IFBMaker, role.Checker, role.IFBChecker,
					}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/actions/enable-disable/approve",
				Handler: handler.EnableDisableServicesChecker,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/accounts/search",
				Handler: handler.SearchAccountByCif,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{
						role.Maker, role.IFBMaker, role.Checker, role.IFBChecker,
					}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/accounts/remove/request",
				Handler: handler.RemoveCifMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/accounts/remove/approve",
				Handler: handler.RemoveCifChecker,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
