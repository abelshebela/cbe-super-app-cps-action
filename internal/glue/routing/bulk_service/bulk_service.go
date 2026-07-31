package bulk_service

import (
	bulk_service "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler bulk_service.BulkServiceHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/bulk_services",
			Handler: handler.GetAllBulkServices,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/bulk_services/disable",
			Handler: handler.DisableBulkService,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/bulk_services/enable",
			Handler: handler.EnableBulkService,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
