package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/customer"
	"cbe-super-app-cps-action/internal/glue"
	middleware "cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler customer.CustomerDetail, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/customers",
			Handler: handler.GetCustomerDetail,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customers/{id}",
			Handler: handler.GetCustomerByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customers/blocked",
			Handler: handler.GetBlockedCustomer,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/customers/enable/{id}",
			Handler: handler.SetEnableCustomerSession,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customers/enable_otp_verify/{id}",
			Handler: handler.EnableCustomer,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/customers/disable/{id}",
			Handler: handler.DisableCustomer,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
