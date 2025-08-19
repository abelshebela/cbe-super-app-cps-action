package department_handler

import (
	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitDepartmentRoutes(router chi.Router, departmentHandler inbound.DepartmentPortHandler, authMiddleware middleware.AuthMiddleware,cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/departments", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: departmentHandler.GetAllDepartments,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: departmentHandler.CreateDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateDepartment)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: departmentHandler.UpdateDepartmentRequest,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateDepartment)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: departmentHandler.GetDepartmentByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method: http.MethodPatch,
				Path:  "/enable/{id}",
				Handler: departmentHandler.EnableDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateDepartment)),

			},
			
		},
	{
				Method: http.MethodPatch,
				Path:  "/disable/{id}",
				Handler: departmentHandler.DisableDepartment,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateDepartment)),

			},
			
		}}
		route.RegisterRoutes(r, routes)
	})
}
