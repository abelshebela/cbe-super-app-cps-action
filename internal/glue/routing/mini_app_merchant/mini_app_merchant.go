package miniappmerchant

import (
	miniappmerchat "cbe-super-app-cps-action/internal/constants/interfaces/mini_app_merchant"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler miniappmerchat.MiniAppMerchant, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/mini_app_merchant/create",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},

		{
			Method:  http.MethodGet,
			Path:    "/mini_app_merchant/get_all",
			Handler: handler.FindAllWithPagination,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/mini_app_merchant/{id}",
			Handler: handler.FindByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/mini_app_merchant/enable/{id}",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/mini_app_merchant/disable/{id}",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/mini_app_merchant/update/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/mini_app_merchant/delete/{id}",
			Handler: handler.Delete,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
