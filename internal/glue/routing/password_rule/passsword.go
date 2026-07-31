package passwordrule

import (
	"net/http"

	password "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/password_rule"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler password.PasswordRule, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/password_rule/",
			Handler: handler.GetPasswordRule,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateTempToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/password_rule/internal",
			Handler: handler.GetPasswordRule,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/password_rule/{id}",
			Handler: handler.RequestPasswordRuleUpdate,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/password_rule/check",
			Handler: handler.CheckPasswordRule,
		},
	}

	glue.RegisterRoutes(router, routes)

}
