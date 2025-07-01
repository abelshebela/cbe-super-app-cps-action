package avatar

import (
	"net/http"

	route "cbe-super-app-cps-action/internal/adapter/inbound/http"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/internal/port/inbound/avatar"

	"github.com/go-chi/chi/v5"
)

func InitAvatarRoutes(router chi.Router, handler avatar.AvatarInbound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/avatar", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: handler.CreateAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/delete/{id}",
				Handler: handler.DeleteAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/{action_code}",
				Handler: handler.Authorize,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/reject/{action_code}",
				Handler: handler.Reject,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/{id}",
				Handler: handler.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/enable/{id}",
				Handler: handler.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.GetAllAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "CHECKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.GetAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "CHECKER"}),
				},
			},

			{
				Method:  http.MethodPatch,
				Path:    "/update/{id}",
				Handler: handler.UpdateAvatar,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})
}
