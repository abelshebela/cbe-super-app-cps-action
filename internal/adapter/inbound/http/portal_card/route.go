package service

import (
	"net/http"

	route "cbe-super-app-cps-action/internal/adapter/inbound/http"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/internal/port/inbound"

	"github.com/go-chi/chi/v5"
)

func InitPortalCardRoutes(router chi.Router, portalcardHandler inbound.PortalCardBound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/portalcard", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/fetch",
				Handler: portalcardHandler.GetAllPortalCard,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
