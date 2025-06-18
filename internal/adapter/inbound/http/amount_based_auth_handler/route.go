package amount_based_auth_handler

import (
	"github.com/go-chi/chi/v5"
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
	"net/http"
)

func InitAmountBasedAuthHandler(router chi.Router, handler inbound.AmountBasedAuthHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/amount-based-auth", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPut,
				Path:    "/update",
				Handler: handler.UpdateAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/approve/{id}",
				Handler: handler.ApproveAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER"}),
				},
			},
		}
		route.RegisterRoutes(r, routes)
	})

}
