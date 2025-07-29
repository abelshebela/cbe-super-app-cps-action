package permission_handler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
)

func InitPermissionRoutes(router chi.Router, permissionHandler inbound.PermissionPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/permissions", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: permissionHandler.CreatePermissionGroup,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: permissionHandler.GetPermissionGroups,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/categories",
				Handler: permissionHandler.GetAllPermissionCategoriesWithPermissions,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{group_name}",
				Handler: permissionHandler.GetPermissionGroup,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{group_name}",
				Handler: permissionHandler.UpdatePermissionGroup,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
