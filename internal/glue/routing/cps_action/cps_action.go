package cpsaction

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	cpsaction "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler cpsaction.CPSActionAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPatch,
			Path:    "/actions/{action_code}/approve",
			Handler: handler.ApproveCPSAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/actions/{action_code}/reject",
			Handler: handler.RejectCPSAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/actions/",
			Handler: handler.GetCPSActionsByDepartment,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/actions/by-id/{action_id}",
			Handler: handler.GetCPSActionByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/actions/by-action-code/{action_code}",
			Handler: handler.GetCPSActionByActionCode,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/actions/counts",
			Handler: handler.GetActionCounts,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
