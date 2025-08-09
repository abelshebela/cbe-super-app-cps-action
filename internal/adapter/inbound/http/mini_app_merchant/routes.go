package miniappmerchant

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/mini_app_merchant"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func InitMiniAppMerchantHandlerMaker(router chi.Router, handler Inbound.MiniAppMerchantInbound, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/mini-app-merchants", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: handler.CreateMiniAppMerchant,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateMiniAppMerchant)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: handler.UpdateMiniAppMerchant,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateMiniAppMerchant)),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: handler.DeleteMiniAppMerchant,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDeleteMiniAppMerchant)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.GetAllMiniAppMerchant,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.GetMiniAppMerchant,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/enable/{id}",
				Handler: handler.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestEnableMiniAppMerchant)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/disable/{id}",
				Handler: handler.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDisableMiniAppMerchant)),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
