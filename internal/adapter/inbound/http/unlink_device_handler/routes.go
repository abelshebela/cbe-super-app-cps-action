package unlink_device_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/unlink"
)

func RegisterHTTPUnlinkRoutes(router chi.Router, handler inbound.UnlinkPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/unlink", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/unlink_device",
				Handler: handler.UnlinkDevice,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			}, {
				Method:  http.MethodPost,
				Path:    "/approve_or_decline_unlink_device",
				Handler: handler.ApproveUnlinkDevice,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
