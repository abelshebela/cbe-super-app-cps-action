package cpsuser

import (
	role "cbe-super-app-cps-action/internal/constants"
	cps_user "cbe-super-app-cps-action/internal/constants/interfaces/cps_user"

	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler cps_user.CPSUserHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/cps_users/create",
			Handler: handler.CreateUserRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cps_users/update/{user_code}",
			Handler: handler.UpdateUserRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/cps_users/delete/{user_code}",
			Handler: handler.DeleteUserRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps_users/{user_code}",
			Handler: handler.FetchUserByUserCode,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps_users",
			Handler: handler.GetAllCPSUsers,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/cps_users/disable/{user_code}",
			Handler: handler.DisableUser,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/cps_users/enable/{user_code}",
			Handler: handler.EnableUser,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
