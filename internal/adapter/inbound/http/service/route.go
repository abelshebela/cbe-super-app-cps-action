package service

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	"github.com/go-chi/chi/v5"
)

func InteServiceRoute(router chi.Router, serviceHandler inbound.Service, authMiddleware middleware.AuthMiddleware) {
	router.Route("/service", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: serviceHandler.GetAllService,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/minimum",
				Handler: serviceHandler.GetAllMinimumTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/maximum",
				Handler: serviceHandler.GetAllMaximumTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service_fee",
				Handler: serviceHandler.GetAllServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/total/transfer_cap",
				Handler: serviceHandler.GetAllTotalTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service_fee/detail/{id}",
				Handler: serviceHandler.GetServiceFeeDetail,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/service_fee/update/{id}",
				Handler: serviceHandler.UpdateServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/single_transfer_max/update/{id}",
				Handler: serviceHandler.UpdateSingleMaxTransfer,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/total_transfer_max/update/{id}",
				Handler: serviceHandler.UpdateTotalMaxTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/minimum_transfer/update/{id}",
				Handler: serviceHandler.UpdateMinimumTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/service_fee/delete/{id}",
				Handler: serviceHandler.DeleteServiceFeeTire,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
