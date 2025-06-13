package faydaaccount

import (
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitFaydaRoutes(router chi.Router, faydaHandler inbound.FaydaAccount) {
	router.Route("/api/v1/cbesuperapp/cps_action/fayda_account_disable", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/initiate",
				Handler: faydaHandler.InitiateDisableFaydaAccount,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve",
				Handler: faydaHandler.AuthorizeFaydaAccountDisable,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/reject",
				Handler: faydaHandler.RejectFaydaAccountDisable,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
