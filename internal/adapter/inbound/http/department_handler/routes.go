package department_handler

import (
	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitDepartmentRoutes(router chi.Router, departmentHandler inbound.DepartmentPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/departments", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: departmentHandler.CreateDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/",
				Handler: departmentHandler.UpdateDepartmentRequest,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{action_code}/approve",
				Handler: departmentHandler.ApproveDepartmentRequest,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/{action_code}/reject",
				Handler: departmentHandler.RejectDepartmentRequest,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
