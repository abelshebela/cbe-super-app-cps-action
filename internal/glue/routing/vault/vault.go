package vault

import (
	bankgroupcategory "cbe-super-app-cps-action/internal/constants/interfaces/vault"

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
		{
			Method:  http.MethodGet,
			Path:    "/vault/transactions/{transaction_id}",
			Handler: handler.GetVaultTransaction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},

		// VAULT WITHDRAWAL REQUEST
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/vault/withdrawals/create",
		// 	Handler: handler.CreateWithdrawalRequest,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		// {
		// 	Method:  http.MethodPatch,
		// 	Path:    "/vault/withdrawals/{id}/cancel",
		// 	Handler: handler.UpdateWithDrawalRequest,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		{
			Method:  http.MethodPatch,
			Path:    "/vault/withdrawals/{id}/approve",
			Handler: handler.UpdateWithDrawalRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault/withdrawals/{id}/reject",
			Handler: handler.UpdateWithDrawalRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/withdrawals",
			Handler: handler.GetAllWithdrawalRequests,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/withdrawals/{id}",
			Handler: handler.GetWithdrawalRequestById,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)

}
