package passwordrule

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
)

func RegisterPasswordRuleRoutes(r chi.Router, handler inbound.PasswordRuleInbound, authMiddleware middleware.AuthMiddleware) {
	routes := []sharedhttp.Route{
		{
			Method:  http.MethodPost,
			Path:    "/password_rule/update_request",
			Handler: handler.RequestPasswordRuleUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/password_rule/update_approve",
			Handler: handler.ApproveOrRejectPasswordRuleAction,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/password_rule/action_by_id",
			Handler: handler.GetPasswordRuleUpdateActionByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/password_rule/check",
			Handler: handler.CheckPasswordRule,
		},
	}

	sharedhttp.RegisterRoutes(r, routes)
}
