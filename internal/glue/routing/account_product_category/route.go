package account_product_category_routing

import (
	apc_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_product_category"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, h apc_interface.AccountProductCategoryHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:      http.MethodGet,
			Path:        "/apc",
			Handler:     h.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodGet,
			Path:        "/apc/{id}",
			Handler:     h.GetByID,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPost,
			Path:        "/apc",
			Handler:     h.Create,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/apc/{id}",
			Handler:     h.Update,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/apc/{id}/enable",
			Handler:     h.Enable,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/apc/{id}/disable",
			Handler:     h.Disable,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/apc/{id}",
			Handler:     h.Delete,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
	}

	glue.RegisterRoutes(router, routes)
}
