package customerhandler

import (
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound"
	"net/http"

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
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: customerHandler.GetCustomerByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
