package bulkservices_inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	bulkservices_application "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/bulk_services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/bulk_services"
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

	routes := []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/enable_disable_service_maker",
			Handler: handler.EnableDisableServicesMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service",
			Handler: handler.FetchServices,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/enable_disable_service_maker",
			Handler: handler.EnableDisableServicesChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/cif_search",
			Handler: handler.SearchAccountByCif,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/cif_remove_maker",
			Handler: handler.RemoveCifMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/cif_remove_checker",
			Handler: handler.RemoveCifChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
