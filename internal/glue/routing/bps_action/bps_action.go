package cpsaction

import (
	"net/http"

	bpsaction "cbe-super-app-cps-action/internal/constants/interfaces/bps_action"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler bpsaction.BPSActionAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPatch,
			Path:    "/bps_actions/{action_code}/approve",
			Handler: handler.ApproveBPSAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/bps_actions/{action_code}/reject",
			Handler: handler.RejectBPSAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/",
			Handler: handler.GetBPSActionsByDepartment,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/user/checked/actions",
			Handler: handler.GetUserCheckedActions,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/by-id/{action_id}",
			Handler: handler.GetBPSActionByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/by-action-code/{action_code}",
			Handler: handler.GetBPSActionByActionCode,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/counts",
			Handler: handler.GetActionCounts,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/approver/checker/actions",
			Handler: handler.GetUserApproverActions,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/bps_actions/auditor/checker/actions",
			Handler: handler.GetUserAuditorActions,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/bps_actions/{action_code}/auditor",
			Handler: handler.AuditorAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/bps_actions/{action_code}/auditor",
			Handler: handler.AuditorAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
