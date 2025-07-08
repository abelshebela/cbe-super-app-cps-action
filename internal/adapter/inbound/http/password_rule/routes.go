package passwordrule

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func RegisterPasswordRuleRoutes(router chi.Router, handler inbound.PasswordRuleInbound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/password-rule", func(r chi.Router) {
		routes := []sharedhttp.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: handler.RequestPasswordRuleUpdate,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve",
				Handler: handler.ApproveOrRejectPasswordRuleAction,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.GetPasswordRuleUpdateActionByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/check",
				Handler: handler.CheckPasswordRule,
			},
		}

		sharedhttp.RegisterRoutes(r, routes)
	})
}
