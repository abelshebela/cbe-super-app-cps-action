package permission

import (
	role "cbe-super-app-cps-action/internal/constants"
	permission "cbe-super-app-cps-action/internal/constants/interfaces/permission"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler permission.PermissionHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/permissions",
			Handler: handler.CreatePermissionGroup,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/permissions",
			Handler: handler.GetPermissionGroups,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/permissions/categories",
			Handler: handler.GetAllPermissionCategoriesWithPermissions,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/permissions/{group_name}",
			Handler: handler.GetPermissionGroup,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/permissions/by_id/{id}",
			Handler: handler.GetPermissionGroupById,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},

		{
			Method:  http.MethodPatch,
			Path:    "/permissions/{id}",
			Handler: handler.UpdatePermissionGroup,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/permissions/categories/{department_id}",
			Handler: handler.GetPermissionCategoriesByDepartment,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/permissions/groups/{department_id}",
			Handler: handler.GetPermissionGroupsByDepartment,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
