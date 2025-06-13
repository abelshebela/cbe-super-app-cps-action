package bulkservices_inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	bulkservices_application "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/bulk_services"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/bulk_services"
)

type HttpStore struct {
	Application bulkservices_application.ApplicationAbstracts
}

func NewHttpBulkService(app bulkservices_application.ApplicationAbstracts) inbound.Inbound {
	return &HttpStore{
		Application: app,
	}
}

func InitServiceHandlerMaker(router chi.Router, handler inbound.Inbound) {

	routes := []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/enable_disable_service_maker",
			Handler: handler.EnableDisableServicesMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service",
			Handler: handler.FetchServices,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/enable_disable_service_maker",
			Handler: handler.EnableDisableServicesChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
