package cpsaction

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	avatar "cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler avatar.AvatarInbound, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/avatar",
			Handler: handler.CreateAvatar,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
				authMiddleware.RequireFormContentType(),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/avatar/{id}",
			Handler: handler.DeleteAvatar,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},

		{
			Method:  http.MethodPatch,
			Path:    "/avatar/disable/{id}",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/avatar/enable/{id}",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/avatar",
			Handler: handler.FetchAvatars,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/avatar/{id}",
			Handler: handler.FetchAvatar,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},

		{
			Method:  http.MethodPatch,
			Path:    "/avatar/{id}",
			Handler: handler.UpdateAvatar,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
