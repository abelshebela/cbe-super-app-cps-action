package vaultamounttier

import (
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	role "cbe-super-app-cps-action/internal/constants"
	vaultAmountTier "cbe-super-app-cps-action/internal/constants/interfaces/vault_amount_tier"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler vaultAmountTier.VaultAmountTierHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/vault-amount-tier/create",
			Handler: handler.CreateAmountTier,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault-amount-tier/find-all",
			Handler: handler.FindAllAmountTiers,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/vault-amount-tier/{id}",
			Handler: handler.GetAmountTier,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},

		{
			Method:  http.MethodPatch,
			Path:    "/vault-amount-tier/{id}/update",
			Handler: handler.UpdateAmountTier,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/vault-amount-tier/{id}/delete",
			Handler: handler.DeleteAmountTier,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault-amount-tier/{id}/disable",
			Handler: handler.DisableAmountTier,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/vault-amount-tier/{id}/enable",
			Handler: handler.EnableAmountTier,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}
	glue.RegisterRoutes(router, routes)

}
