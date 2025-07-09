package cpsmakerhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func RegisterCPSUserMakerRoutes(r chi.Router, handler inbound.CPSUserMakerHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []sharedhttp.Route{
		{
			Method:  http.MethodPost,
			Path:    "/cps_user_maker/create",
			Handler: handler.CreateUserRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPut,
			Path:    "/cps_user_maker/update",
			Handler: handler.UpdateUserRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/cps_user_maker/approve",
			Handler: handler.ApproveUserAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps_user_maker/pending_actions",
			Handler: handler.GetPendingUserActions,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
			},
		},
	}

	sharedhttp.RegisterRoutes(r, routes)
}
