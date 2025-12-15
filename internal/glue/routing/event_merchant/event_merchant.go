package eventmerchant

import (
	"cbe-super-app-cps-action/internal/constants"
	event_merchant_port "cbe-super-app-cps-action/internal/constants/interfaces/event_merchant"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
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
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event_merchants/{id}",
			Handler: handler.UpdateEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event_merchants/enable/{id}",
			Handler: handler.EnableEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/event_merchants/disable/{id}",
			Handler: handler.DisableEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/event_merchant/{id}",
			Handler: handler.DeleteEventMerchant,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/event_merchants/{id}",
			Handler: handler.GetEventMerchantByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker, constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/event_merchants",
			Handler: handler.GetEventMerchants,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker, constants.Maker, constants.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
>>>>>>>> 97d09d6944f92972be5047e3259c96888d8d2fce:internal/glue/routing/event_merchant/event_merchant.go
