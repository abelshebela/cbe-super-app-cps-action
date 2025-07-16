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
	router.Route("/service-details", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.GetAllServiceDetails,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.GetServiceDetailsByID,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/maker/update",
				Handler: handler.UpdateServiceDetailsMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/checker/approve/{action_code}",
				Handler: handler.UpdateServiceDetailsChecker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Checker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/checker/reject/{action_code}",
				Handler: handler.UpdateServiceDetailsChecker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Checker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cap/daily/maker",
				Handler: handler.ServiceDetailsDailyCapMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cap/single/maker",
				Handler: handler.ServiceDetailsSingleCapMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cap/total/maker",
				Handler: handler.TotalTransferCapMaker,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cap/min/update",
				Handler: handler.UpdateCapMinAmountHandler,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve",
				Handler: handler.ApproveServiceDetailsHandler,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Checker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/fee/maker",
				Handler: handler.InitiateServiceFeeUpdate,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Checker, role.Checker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/fee/approve",
				Handler: handler.ApproveServiceFeeUpdate,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cap/total/maker/fee/reject",
				Handler: handler.RejectServiceFeeUpdate,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
