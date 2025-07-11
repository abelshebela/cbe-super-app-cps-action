package service

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	"github.com/go-chi/chi/v5"
)

func InitPortalCardRoutes(router chi.Router, portalcardHandler inbound.PortalCardBound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/portalcard", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: portalcardHandler.GetAllPortalCard,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
