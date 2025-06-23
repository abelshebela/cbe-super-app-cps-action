package amount_based_auth_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
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
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/approve/{id}",
				Handler: handler.ApproveAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker"}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/reject/{id}",
				Handler: handler.RejectAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER"}),
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})

}
