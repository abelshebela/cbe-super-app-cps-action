package bps_actionrole_routing

import (
	actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/action_role"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler actionrole_inbound.BPSActionRoleHandler, auth middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/bps-action-roles",
			Handler: handler.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps-action-roles/action-list",
			Handler: handler.GetAllActionList,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps-action-roles/{code}",
			Handler: handler.GetByActionCode,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/bps-action-roles",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/bps-action-roles/{code}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/bps-action-roles/{code}/enable",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/bps-action-roles/{code}/disable",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
