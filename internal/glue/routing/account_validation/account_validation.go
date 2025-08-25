package accountvalidation

import (
	accountvalidation "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler accountvalidation.AccountValidation, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/account_validation/{id}",
			Handler: handler.FetchAccountValidation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/account_validation",
			Handler: handler.FetchAllAccountValidation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/account_validation/update/{id}",
			Handler: handler.UpdateAccountValidation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
