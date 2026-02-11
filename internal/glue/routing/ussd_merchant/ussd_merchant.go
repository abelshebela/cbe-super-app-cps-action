package ussd_merchant

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/ussd_merchant"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Init(router chi.Router, handler ussd_merchant_interface.UssdMerchantInbound, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/ussd_merchant",
			Handler: handler.CreateUssdMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ussd_merchant/{id}",
			Handler: handler.UpdateUssdMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/ussd_merchant/{id}",
			Handler: handler.GetUssdMerchant,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ussd_merchant",
			Handler: handler.GetAllUssdMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ussd_merchant/disable/{id}",
			Handler: handler.DisableUssdMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/ussd_merchant/enable/{id}",
			Handler: handler.EnableUssdMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
