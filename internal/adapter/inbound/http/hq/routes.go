package hq

import (
	"net/http"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"

	"github.com/go-chi/chi/v5"
)

func InitHQRoutes(r chi.Router, handler *HQHTTPHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/cbesuperapp/cps_action/hq/{id}",
			Handler: handler.GetHQ,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_action/hq/block-time/request",
			Handler: handler.UpdateBlockTimeRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_action/hq/archive-time/request",
			Handler: handler.UpdateArchiveTimeRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_action/hq/block-time/update",
			Handler: handler.UpdateBlockTime,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_action/hq/archive-time/update",
			Handler: handler.UpdateArchiveTime,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
			},
		},
	}

	route.RegisterRoutes(r, routes)
}
