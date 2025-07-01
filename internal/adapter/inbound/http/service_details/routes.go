package service_details_inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "cbe-super-app-cps-action/internal/adapter/inbound/http"
	"cbe-super-app-cps-action/internal/application/middleware"
	inbound "cbe-super-app-cps-action/internal/port/inbound/service_details"
)

func InitServiceDetailsRoutes(router chi.Router, handler inbound.ServiceDetailsInbound, middleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service_details",
			Handler: handler.GetAllServiceDetails,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "checker"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service_details/{id}",
			Handler: handler.GetServiceDetailsByID,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_service_details_maker",
			Handler: handler.UpdateServiceDetailsMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_service_details_checker",
			Handler: handler.UpdateServiceDetailsChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/daily_cap_maker",
			Handler: handler.ServiceDetailsDailyCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/single_cap_maker",
			Handler: handler.ServiceDetailsSingleCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/total_cap_maker",
			Handler: handler.TotalTransferCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker"}),
			},
		},
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/api/v1/cbesuperapp/cps_config/service_fee_initiate",
		// 	Handler: handler.ServiceFeeMaker,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		middleware.AuthenticateToken,
		// 		middleware.AccessControl([]string{"MAKER"}),
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/api/v1/cbesuperapp/cps_config/service_fee_Approve",
		// 	Handler: handler.ServiceFeeApprove,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		middleware.AuthenticateToken,
		// 		middleware.AccessControl([]string{"CHEKER"}),
		// 	},
		// },
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/api/v1/cbesuperapp/cps_config/service_fee_Reject",
		// 	Handler: handler.ServiceFeeReject,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		middleware.AuthenticateToken,
		// 		middleware.AccessControl([]string{"CHEKER"}),
		// 	},
		// },
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_min_cap",
			Handler: handler.UpdateCapMinAmountHandler,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/approve_service_details",
			Handler: handler.ApproveServiceDetailsHandler,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/service_fee_maker",
			Handler: handler.InitiateServiceFeeUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/service_fee_approve",
			Handler: handler.ApproveServiceFeeUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/total_cap_maker/service_fee_reject",
			Handler: handler.RejectServiceFeeUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
