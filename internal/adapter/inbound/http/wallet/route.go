package wallet

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	walletRoutes "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/wallet"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	"github.com/go-chi/chi/v5"
)

func InitWalletRoutes(router chi.Router, wallet walletRoutes.WalletAdapter, authMiddleware middleware.AuthMiddleware,cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/wallets", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: wallet.CreateWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					middleware.RequireFormContentType(),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateWallet)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: wallet.UpdateWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateWallet)),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: wallet.DeleteWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDeleteWallet)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: wallet.GetWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: wallet.GetAllWallet,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}/enable",
				Handler: wallet.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestEnableWallet)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}/disable",
				Handler: wallet.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestEnableWallet)),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
