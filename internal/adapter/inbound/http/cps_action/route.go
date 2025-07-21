package cpsaction

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	actions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/cps_actions"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
)

func InitCPSActionsRoutes(router chi.Router, actions actions.CPSActionAdapter, authMiddleware middleware.AuthMiddleware) {
	router.Route("/actions", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/{action_code}/approve",
				Handler: actions.ApproveCPSAction,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{action_code}/reject",
				Handler: actions.RejectCPSAction,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: actions.GetCPSActionsByDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/by-id/{id}",
				Handler: actions.GetCPSActionByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/by-action-code/{action_code}",
				Handler: actions.GetCPSActionByActionCode,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})
}
