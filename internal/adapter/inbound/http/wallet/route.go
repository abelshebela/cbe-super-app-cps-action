package wallet

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	walletRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/wallet"
)

func InitWalletRoutes(router chi.Router, wallet walletRoutes.WalletAdapter, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/wallet", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: wallet.CreateWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/update/{id}",
				Handler: wallet.UpdateWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/delete/{id}",
				Handler: wallet.DeleteWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: wallet.GetWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/wallets",
				Handler: wallet.GetAllWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/{action_code}",
				Handler: wallet.Authorize,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
				},
			},

			{
				Method:  http.MethodPost,
				Path:    "/reject/{action_code}",
				Handler: wallet.Reject,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
				},
			},

			{
				Method:  http.MethodPost,
				Path:    "/enable/{id}",
				Handler: wallet.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/{id}",
				Handler: wallet.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
