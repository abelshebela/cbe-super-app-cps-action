package bankvault

import (
	bankvault "cbe-super-app-cps-action/internal/constants/interfaces/bankvault"

	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler bankvault.BankVaultHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/vault/products/create",
			Handler: handler.CreateBankVault,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/products",
			Handler: handler.FindAllBankVaults,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/products/{id}",
			Handler: handler.GetBankVault,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault/products/update/{id}",
			Handler: handler.UpdateBankVault,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/vault/products/delete/{id}",
			Handler: handler.DeleteBankVault,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault/products/disable/{id}",
			Handler: handler.DisableBankVault,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault/products/enable/{id}",
			Handler: handler.EnableBankVault,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault/bank-vaults",
			Handler: handler.GetAllLockedBankVaults,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		// {
		// 	Method:  http.MethodGet,
		// 	Path:    "/vault/bank-vaults/{transaction_reference}",
		// 	Handler: handler.GetTransaction,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
		{
			Method:  http.MethodGet,
			Path:    "/vault/group-vaults",
			Handler: handler.GetAllGroupVaults,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		// {
		// 	Method:  http.MethodGet,
		// 	Path:    "/vault/group-vaults/{id}",
		// 	Handler: handler.GetGroupVault,
		// 	Middlewares: []func(next http.Handler) http.Handler{
		// 		authMiddleware.AuthenticateToken,
		// 	},
		// },
	}
	glue.RegisterRoutes(router, routes)

}
