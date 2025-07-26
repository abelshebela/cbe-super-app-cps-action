package accountlookup

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	accountLookUpRoutes "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_lookup"

	"github.com/go-chi/chi/v5"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func InitAccountLookUpRoutes(router chi.Router, accountLookUp accountLookUpRoutes.UserSearchAdapter, authMiddleware middleware.AuthMiddleware) {
	router.Route("/account_look_ups", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: accountLookUp.SearchUser,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
