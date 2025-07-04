package department_handler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"

	"github.com/go-chi/chi/v5"
)

func InitDepartmentRoutes(router chi.Router, departmentHandler inbound.DepartmentPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/department", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: departmentHandler.CreateDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/request/{action_code}",
				Handler: departmentHandler.ApproveDepartmentRequest,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
