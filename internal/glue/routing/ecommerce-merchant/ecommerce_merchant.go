package ecommercemerchant

import (
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
			},
		},

		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant",
			Handler: handler.FindAllWithPagination,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.FindByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/enable/{id}",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/disable/{id}",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.Delete,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant/merchant-lookup/{merchant_id}",
			Handler: handler.MerchantLookup,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
