package faydaaccount

import (
	"net/http"

	route "cbe-super-app-cps-action/internal/adapter/inbound/http"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/internal/port/inbound"

	"github.com/go-chi/chi/v5"
)

func InitFaydaRoutes(router chi.Router, faydaHandler inbound.FaydaAccount, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/fayda_account_disable", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/initiate",
				Handler: faydaHandler.InitiateDisableFaydaAccount,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve",
				Handler: faydaHandler.AuthorizeFaydaAccountDisable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/reject",
				Handler: faydaHandler.RejectFaydaAccountDisable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
