<<<<<<< HEAD
package amount_based_auth

import (
	amount_based "cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler amount_based.AmountBasedAuthAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPatch,
			Path:    "/amount_based_auth/update/{method}/{id}",
			Handler: handler.UpdateAmountBasedAuth,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/amount_based_auth",
			Handler: handler.GetAllAmountBasedAuth,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/amount_based_auth/reject/{id}",
			Handler: handler.RejectAmountBasedAuth,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/amount_based_auth/currency",
			Handler: handler.AddCurrency,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/amount_based_auth/reset/{currency}",
			Handler: handler.ResetConfig,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
=======
package amount_based_auth

import (
	amount_based "cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler amount_based.AmountBasedAuthAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPatch,
			Path:    "/amount_based_auth/update/{method}/{id}",
			Handler: handler.UpdateAmountBasedAuth,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/amount_based_auth",
			Handler: handler.GetAllAmountBasedAuth,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/amount_based_auth/reject/{id}",
			Handler: handler.RejectAmountBasedAuth,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
>>>>>>> 3640c5b5dde77249222ec7dd9f90a5770ab0dbc4
