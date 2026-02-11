package job_role

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/job_role"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler inbound.RolesInbound, auth middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/job_roles",
			Handler: handler.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/job_roles/{id}",
			Handler: handler.GetByID,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/job_roles",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/job_roles/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
