package logistic_merchant_router

import (
	"cbe-super-app-cps-action/internal/constants"
	logistics_merchant_adaptor "cbe-super-app-cps-action/internal/constants/interfaces/logistics_merchant"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler logistics_merchant_adaptor.LogisticMerchantInboundAdaptor, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/logistics_merchants",
			Handler: handler.CreateLogisticMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/logistics_merchants/{id}",
			Handler: handler.UpdateLogisticMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/logistics_merchants/enable/{id}",
			Handler: handler.EnableLogisticMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/logistics_merchants/disable/{id}",
			Handler: handler.DisableLogisticMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/logistics_merchants/{id}",
			Handler: handler.DeleteLogisticMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/logistics_merchants/{id}",
			Handler: handler.GetLogisticMerchantByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker, constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/logistics_merchants",
			Handler: handler.GetLogisticMerchants,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker, constants.Maker, constants.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
