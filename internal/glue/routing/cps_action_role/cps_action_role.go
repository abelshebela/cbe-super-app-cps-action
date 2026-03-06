package cps_actionrole_routing

import (
	actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action_role"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler actionrole_inbound.CPSActionRoleHandler, auth middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/cps-action-roles",
			Handler: handler.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps-action-list",
			Handler: handler.GetAllActionList,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps-action-roles/{code}",
			Handler: handler.GetByActionCode,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/cps-action-roles",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cps-action-roles/{code}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cps-action-roles/{code}/enable",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cps-action-roles/{code}/disable",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps-action-roles/{code}/versions",
			Handler: handler.GetVersions,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/cps-action-roles/{code}/roles",
			Handler: handler.GetConfiguredRoles,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cps-action-roles/{code}/versions/{version}/role",
			Handler: handler.UpdateVersionRoleCode,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
