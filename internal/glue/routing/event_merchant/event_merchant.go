package eventmerchant

import (
	event_merchant_port "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/event_merchant"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler event_merchant_port.EventMerchantInboundAdaptor, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/event_merchants",
			Handler: handler.CreateEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event_merchants/{id}",
			Handler: handler.UpdateEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event_merchants/enable",
			Handler: handler.EnableEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event_merchants/disable",
			Handler: handler.DisableEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/event_merchants/{id}",
			Handler: handler.DeleteEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/event_merchants/{id}",
			Handler: handler.GetEventMerchantByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/event_merchants",
			Handler: handler.GetEventMerchants,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/event-merchant/merchant-lookup/{merchant_id}",
			Handler: handler.EventMerchantLookup,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTokenOrMerchantIntegrationAPIKey,
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
