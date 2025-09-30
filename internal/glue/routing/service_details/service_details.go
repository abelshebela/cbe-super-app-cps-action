package service

import (
	"net/http"
	"cbe-super-app-cps-action/internal/constants"
	serviceHandler "cbe-super-app-cps-action/internal/constants/interfaces/service_details"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, serviceHandler serviceHandler.ServiceAdapter, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
			{
				Method:  http.MethodGet,
				Path:    "/service",
				Handler: serviceHandler.GetAllService,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service/minimum",
				Handler: serviceHandler.GetAllMinimumTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service/maximum",
				Handler: serviceHandler.GetAllMaximumTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service/service_fee",
				Handler: serviceHandler.GetAllServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service/total/transfer_cap",
				Handler: serviceHandler.GetAllTotalTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/service/service_fee/detail/{id}",
				Handler: serviceHandler.GetServiceFeeDetail,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/service/service_fee/update/{id}",
				Handler: serviceHandler.UpdateServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
					
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/service/single_transfer_max/update/{id}",
				Handler: serviceHandler.UpdateSingleMaxTransfer,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/service/total_transfer_max/update",
				Handler: serviceHandler.UpdateTotalMaxTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/service/minimum_transfer/update/{id}",
				Handler: serviceHandler.UpdateMinimumTransferCap,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/service/service_fee/delete/{id}",
				Handler: serviceHandler.DeleteServiceFeeTire,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				
				},
			},
		}

	
	glue.RegisterRoutes(router, routes)
}
