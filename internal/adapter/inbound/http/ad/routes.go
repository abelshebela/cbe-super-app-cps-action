package ad

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	adRoutes "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/go-chi/chi/v5"
)

func InitADRoutes(router chi.Router, ad adRoutes.ADAdapter, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/adverts", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: ad.CreateAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					middleware.RequireFormContentType(),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateAdvert)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: ad.UpdateAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateAdvert)),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: ad.DeleteAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDeleteAdvert)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: ad.FetchAdvertByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: ad.FetchAdverts,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/enable/{id}",
				Handler: ad.EnableAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestEnableAdvert)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/disable/{id}",
				Handler: ad.DisableAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDisableAdvert)),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
