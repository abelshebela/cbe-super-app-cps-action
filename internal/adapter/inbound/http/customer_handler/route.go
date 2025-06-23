package customerhandler

import (
	"net/http"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"

	"github.com/go-chi/chi/v5"
)

func InitCustomerRoutes(router chi.Router, customerHandler inbound.CustomerDetail,authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/customers", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: customerHandler.GetCustomerDetail,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: customerHandler.GetCustomerByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
