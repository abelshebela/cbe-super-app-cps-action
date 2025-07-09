package amount_based_auth_handler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/go-chi/chi/v5"
)

func InitAmountBasedAuthHandler(router chi.Router, handler inbound.AmountBasedAuthHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/amount_based_auth", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/update/{id}",
				Handler: handler.UpdateAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/approve/{id}",
				Handler: handler.ApproveAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/reject/{id}",
				Handler: handler.RejectAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker}),
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})

}
