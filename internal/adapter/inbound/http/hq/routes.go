package hq

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"

	"github.com/go-chi/chi/v5"
)

func InitHQRoutes(r chi.Router, handler *HQHTTPHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/hq/{id}",
			Handler: handler.GetHQ,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/block-time/update_request",
			Handler: handler.UpdateBlockTimeRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/archive-time/update_request",
			Handler: handler.UpdateArchiveTimeRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/block-time/update_approve",
			Handler: handler.UpdateBlockTime,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/archive-time/update_approve",
			Handler: handler.UpdateArchiveTime,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-cheker"}),
			},
		},
	}

	route.RegisterRoutes(r, routes)
}
