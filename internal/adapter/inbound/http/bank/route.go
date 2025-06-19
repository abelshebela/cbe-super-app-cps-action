package bank

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	bankRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/bank"
)

func InitBankRoutes(router chi.Router, bank bankRoutes.BankAdapter, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/bank", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: bank.CreateOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/update/{id}",
				Handler: bank.UpdateOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/delete/{id}",
				Handler: bank.DeleteOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: bank.GetOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/banks",
				Handler: bank.GetAllBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER", "CHECKER", "IFB-CHECKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/{action_code}",
				Handler: bank.Authorize,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},

			{
				Method:  http.MethodPost,
				Path:    "/reject/{action_code}",
				Handler: bank.Reject,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/enable/{id}",
				Handler: bank.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/disable/{id}",
				Handler: bank.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
