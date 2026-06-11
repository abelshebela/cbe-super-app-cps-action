package account_product_routing

import (
	ap_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_product"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, h ap_interface.AccountProductHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:      http.MethodGet,
			Path:        "/account_product",
			Handler:     h.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodGet,
			Path:        "/account_product/{id}",
			Handler:     h.GetByID,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPost,
			Path:        "/account_product",
			Handler:     h.Create,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/account_product/{id}",
			Handler:     h.Update,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/account_product/{id}/enable",
			Handler:     h.Enable,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/account_product/{id}/disable",
			Handler:     h.Disable,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/account_product/{id}",
			Handler:     h.Delete,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
	}

	glue.RegisterRoutes(router, routes)
}
