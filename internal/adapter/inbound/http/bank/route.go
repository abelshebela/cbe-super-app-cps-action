package bank

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	bankRoutes "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bank"

	"github.com/go-chi/chi/v5"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func InitBankRoutes(router chi.Router, bank bankRoutes.BankAdapter, authMiddleware middleware.AuthMiddleware) {
	router.Route("/banks", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: bank.CreateOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: bank.UpdateOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: bank.DeleteOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: bank.GetOneBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/list",
				Handler: bank.GetAllBank,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/{action_code}/approve",
				Handler: bank.Authorize,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/{action_code}/reject",
				Handler: bank.Reject,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/{id}/enable",
				Handler: bank.Enable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/{id}/disable",
				Handler: bank.Disable,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
