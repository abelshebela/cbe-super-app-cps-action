package vaultgroupcategory

import (
	bankgroupcategory "cbe-super-app-cps-action/internal/constants/interfaces/vault_category"

	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler bankgroupcategory.VaultCategoryHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/vault/categories/create",
			Handler: handler.CreateVaultCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/categories",
			Handler: handler.FindAllVaultCategories,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/categories/{id}",
			Handler: handler.GetVaultCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault/categories/update/{id}",
			Handler: handler.UpdateVaultCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		// {
		// 	Method:  http.MethodDelete,
		// 	Path:    "/vault/categories/delete/{id}",
		// 	Handler: handler.DeleteVaultCategory,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		{
			Method:  http.MethodPatch,
			Path:    "/vault/categories/enable/{id}",
			Handler: handler.EnableVaultCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault/categories/disable/{id}",
			Handler: handler.DisableVaultCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},

		// VAULT TRANSACTION
		{
			Method:  http.MethodGet,
			Path:    "/vault/transactions",
			Handler: handler.GetVaultTransactions,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)

}
