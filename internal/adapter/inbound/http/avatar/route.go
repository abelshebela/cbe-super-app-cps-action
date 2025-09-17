package avatar

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/avatar"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
)

func InitAvatarRoutes(router chi.Router, handler avatar.AvatarInbound, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/avatar", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: handler.CreateAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					middleware.RequireFormContentType(),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateAvatar)),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: handler.DeleteAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDeleteAvatar)),
				},
			},

			{
				Method:  http.MethodPatch,
				Path:    "/disable/{id}",
				Handler: handler.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDisableAvatar)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/enable/{id}",
				Handler: handler.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestEnableAvatar)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.FetchAvatars,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.FetchAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},

			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: handler.UpdateAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateAvatar)),
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})
}
