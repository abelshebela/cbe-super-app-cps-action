package permission_handler

import (
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/permission"
	"net/http"

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
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/request/{action_code}",
				Handler: permissionHandler.ApprovePermissionGroup,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
