package ecommercemerchant

import (
	ecommerce_merchant "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/ecommerce_merchant"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler ecommerce_merchant.EcommerceMerchant, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/ecommerce-merchant",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant",
			Handler: handler.FindAllWithPagination,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.FindByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/enable",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/disable",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ecommerce-merchant/{id}",
			Handler: handler.Delete,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ecommerce-merchant/merchant-lookup/{merchant_id}",
			Handler: handler.MerchantLookup,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ecommerce-merchant/branch/delete/{id}",
			Handler: handler.DeleteBranch,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/branch/enable/{id}",
			Handler: handler.EnableBranch,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ecommerce-merchant/branch/disable/{id}",
			Handler: handler.DisableBranch,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
