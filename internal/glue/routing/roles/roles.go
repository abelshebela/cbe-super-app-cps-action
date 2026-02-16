package job_role

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/roles"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler inbound.RolesInbound, auth middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/roles",
			Handler: handler.FindAllWithPagination,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/roles/all",
			Handler: handler.FindAll,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/roles/{id}",
			Handler: handler.FindById,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/roles",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/roles/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/roles/{id}/enable",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/roles/{id}/disable",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/roles/{id}",
			Handler: handler.Delete,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
