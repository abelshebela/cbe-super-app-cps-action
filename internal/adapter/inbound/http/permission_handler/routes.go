package permission_handler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"

	"github.com/go-chi/chi/v5"
)

func InitPermissionRoutes(router chi.Router, permissionHandler inbound.PermissionPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/permission", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: permissionHandler.CreatePermissionGroup,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/request/{action_code}",
				Handler: permissionHandler.ApprovePermissionGroup,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
