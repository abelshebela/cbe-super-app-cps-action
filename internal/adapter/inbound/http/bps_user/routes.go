package bpsmakerhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bps_user"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func RegisterBPSUserMakerRoutes(router chi.Router, handler inbound.BPSUserHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/bps_users", func(r chi.Router) {
		routes := []sharedhttp.Route{
			// {
			// 	Method:  http.MethodGet,
			// 	Path:    "/pending_user_actions",
			// 	Handler: handler.GetPendingUserActions,
			// 	Middlewares: []func(next http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 		authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
			// 	},
			// },
			{
				Method:  http.MethodGet,
				Path:    "/{user_code}",
				Handler: handler.FetchUserByUserCode,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.GetAllBPSUsers,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/{user_code}",
				Handler: handler.DisableUser,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/enable/{user_code}",
				Handler: handler.EnableUser,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
		}

		sharedhttp.RegisterRoutes(r, routes)
	})
}
