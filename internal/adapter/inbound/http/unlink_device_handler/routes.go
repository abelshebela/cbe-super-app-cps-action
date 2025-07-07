package unlink_device_handler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPUnlinkRoutes(router chi.Router, handler inbound.UnlinkPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/unlink", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/unlink_device",
				Handler: handler.UnlinkDevice,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			}, {
				Method:  http.MethodPost,
				Path:    "/approve_or_decline_unlink_device",
				Handler: handler.ApproveUnlinkDevice,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
