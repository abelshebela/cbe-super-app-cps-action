package ecommercemerchant

import (
	"cbe-super-app-cps-action/internal/constants"
	miniappmerchat "cbe-super-app-cps-action/internal/constants/interfaces/ecommerce_merchant"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler miniappmerchat.EcommerceMerchant, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/ecommerce-merchant",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},

		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant",
			Handler: handler.FindAllWithPagination,
			// Middlewares: []func(next http.Handler) http.Handler{
			// 	authMiddleware.AuthenticateToken,
			// 	authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			// },
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.FindByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/enable/{id}",
			Handler: handler.Enable,
			// Middlewares: []func(next http.Handler) http.Handler{
			// 	authMiddleware.AuthenticateToken,
			// 	authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			// },
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/disable/{id}",
			Handler: handler.Disable,
			// Middlewares: []func(next http.Handler) http.Handler{
			// 	authMiddleware.AuthenticateToken,
			// 	authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			// },
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.Delete,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant/merchant-lookup/{merchant_id}",
			Handler: handler.MerchantLookup,
			// Middlewares: []func(next http.Handler) http.Handler{
			// 	authMiddleware.AuthenticateToken,
			// 	authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			// },
		},
	}

	glue.RegisterRoutes(router, routes)

}
