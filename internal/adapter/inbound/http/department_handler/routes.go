package department_handler

import (
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound/department"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitDepartmentRoutes(router chi.Router, departmentHandler inbound.DepartmentPortHandler) {
	router.Route("/api/v1/cbesuperapp/cps_action/department", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: departmentHandler.CreateDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/request/{action_code}",
				Handler: departmentHandler.ApproveDepartmentRequest,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{"CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
