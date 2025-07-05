package service_details_inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/service_details"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

)

func InitServiceDetailsRoutes(router chi.Router, handler inbound.ServiceDetailsInbound, middleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service_details",
			Handler: handler.GetAllServiceDetails,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker, role.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_config/fetch_service_details/{id}",
			Handler: handler.GetServiceDetailsByID,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker, role.Checker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_service_details_maker",
			Handler: handler.UpdateServiceDetailsMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_service_details_checker",
			Handler: handler.UpdateServiceDetailsChecker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Checker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/daily_cap_maker",
			Handler: handler.ServiceDetailsDailyCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/single_cap_maker",
			Handler: handler.ServiceDetailsSingleCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/total_cap_maker",
			Handler: handler.TotalTransferCapMaker,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/update_min_cap",
			Handler: handler.UpdateCapMinAmountHandler,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/approve_service_details",
			Handler: handler.ApproveServiceDetailsHandler,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Checker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/service_fee_maker",
			Handler: handler.InitiateServiceFeeUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Checker, role.Checker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/service_fee_approve",
			Handler: handler.ApproveServiceFeeUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/total_cap_maker/service_fee_reject",
			Handler: handler.RejectServiceFeeUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				middleware.AuthenticateToken,
				middleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
