package service_details_inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/service_details"
)

func InitServiceDetailsRoutes(router chi.Router, handler inbound.ServiceDetailsInbound, middleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service_details",
			Handler: handler.GetAllServiceDetails,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER", "CHECKER"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service_details/{id}",
			Handler: handler.GetServiceDetailsByID,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER", "CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_service_details_maker",
			Handler: handler.UpdateServiceDetailsMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_service_details_checker",
			Handler: handler.UpdateServiceDetailsChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/daily_cap_maker",
			Handler: handler.ServiceDetailsDailyCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/single_cap_maker",
			Handler: handler.ServiceDetailsSingleCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/total_cap_maker",
			Handler: handler.TotalTransferCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"MAKER"}),
			},
		},
		 {
            Method:  http.MethodPost,
            Path:    "/api/v1/cbesuperapp/cps_config/update_min_cap",
            Handler: handler.UpdateCapMinAmountHandler,
            Middlewares: []func(next http.Handler) http.Handler{
                middleware.AuthenticateToken,
                middleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/api/v1/cbesuperapp/cps_config/approve_service_details",
            Handler: handler.ApproveServiceDetailsHandler,
            Middlewares: []func(next http.Handler) http.Handler{
                middleware.AuthenticateToken,
                middleware.AccessControl([]string{"CHECKER"}),
            },
        },
	}

	route.RegisterRoutes(router, routes)
}
